package main

// import (
// 	"errors"
// 	"fmt"
//
// 	"github.com/Nyks06/abstract-godbus"
// 	"github.com/Nyks06/dbus"
// )
//

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/godbus/dbus"
)

func getDBusResp(ok *dbus.Call) map[string]interface{} {
	fmt.Printf("body %s\n", ok.Body[0])
	resp := make(map[string]interface{})
	err := json.Unmarshal([]byte(ok.Body[0].(string)), &resp)
	if err != nil {
		panic(err.Error())
	}
	if resp["status"] != true {
		panic(resp["message"])
	}
	return resp
}

func IsOnDbus(name string) (bool, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return false, errors.New(err.Error())
	}

	var s []string
	s = make([]string, 0)
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)
	if err != nil {
		return false, errors.New(err.Error())
	}
	for _, v := range s {
		fmt.Printf("name : %s\n", v)
		if v == name {
			return true, nil
		}
	}
	return false, errors.New(name + " is not present")
}

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	_, err = IsOnDbus("jsonception.JSonWriter")
	if err != nil {
		panic(err.Error())
	}
	ok := conn.Object("jsonception.JSonWriter", "/jsonception/JSonWriter").Call("jsonception.JSonWriter.Connect", 0, "TESTCREATION")
	token := getDBusResp(ok)["response"]
	ok = conn.Object("jsonception.JSonWriter", "/jsonception/JSonWriter").Call("jsonception.JSonWriter.Write", 0, token, `{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`)
	getDBusResp(ok)
	ok = conn.Object("jsonception.JSonWriter", "/jsonception/JSonWriter").Call("jsonception.JSonWriter.Write", 0, token, `{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`)
	getDBusResp(ok)

	//data, _ := json.MarshalIndent(node, "", "    ")
	//os.Stdout.Write(data)
}
