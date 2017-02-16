package json_writer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/Nyks06/dbus"
)

var (
	dbusConn *dbus.Conn
	logger   *Writer
)

const (
	OBJECT_NAME      = "jsonception.JSonWriter"
	OBJECT_PATH      = "/jsonception/JSonWriter"
	OBJECT_INTERFACE = "jsonception.JSonWriter"
)

type DbusExporter struct {
}

func (exp *DbusExporter) Disconnect(token string) (string, *dbus.Error) {
	fmt.Printf("on se disconnect avec params : %s\n", token)
	resp := logger.CloseFile(token)
	b, _ := json.Marshal(resp)
	return string(b), nil
}

func (exp *DbusExporter) Connect(file string) (string, *dbus.Error) {
	fmt.Printf("on se connect avec params : %s\n", file)
	resp := logger.CreateFile(file)
	b, _ := json.Marshal(resp)
	return string(b), nil
}

func (exp *DbusExporter) Write(token, params string) (string, *dbus.Error) {
	fmt.Println("test write")
	fmt.Printf("params : %s\n", params)
	//w.Printf(token, event_type, format, a...)
	resp := logger.Write(token, params)
	b, _ := json.Marshal(resp)
	return string(b), nil
}

func Init() {
	fmt.Println("On init dbus json writer")
	logger = &Writer{files: make(map[string]*FileInfo)}

	var err error
	if runtime.GOOS == "windows" {
		if dbusConn, err = dbus.SystemBus(); err != nil {
			panic(err.Error())
		}
	} else {
		if dbusConn, err = dbus.SessionBus(); err != nil {
			panic(err.Error())
		}
	}

	e := DbusExporter{}
	dbusConn.RequestName(OBJECT_INTERFACE, dbus.NameFlagDoNotQueue)
	if err = dbusConn.Export(&e, OBJECT_PATH, OBJECT_INTERFACE); err != nil {
		panic(err.Error())
	}
	os.Setenv("LOG_DEBUG", "1")
	log.Println("DBus json writer well exported")
}
