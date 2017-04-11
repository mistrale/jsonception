package main

import (
	"encoding/json"
	"fmt"

	"github.com/Nyks06/dbus"
)

type test struct {
	Status   bool        `json:"status"`
	Message  string      `json:"message"`
	Response interface{} `json:"response"`
}

func main() {
	answer := &test{}
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Println(err.Error())
	}
	call := conn.Object("jsonception.JSonWriter", "/jsonception/JSonWriter").Call("jsonception.JSonWriter"+".Connect", 0, "/Users/Sikorav/tata.test")
	if call.Err != nil {
	}
	if err := json.Unmarshal([]byte(call.Body[0].(string)), answer); err != nil {

	}
	conn.Object("jsonception.JSonWriter", "/jsonception/JSonWriter").Call("jsonception.JSonWriter"+".Write", 0, answer.Response, `{"event_id" : 1, "name" : "titi", "obj" : {"id" : 1}}`)
	conn.Object("jsonception.JSonWriter", "/jsonception/JSonWriter").Call("jsonception.JSonWriter"+".Write", 0, answer.Response, `{"event_id" : 2, "name" : "titi", "obj" : {"id" : 1}}`)

	conn.Object("jsonception.JSonWriter", "/jsonception/JSonWriter").Call("jsonception.JSonWriter"+".Disconnect", 0, answer.Response)

	fmt.Printf("response %s\n", answer.Response)

}
