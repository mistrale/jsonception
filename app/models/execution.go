package models

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"

	"github.com/mistrale/jsonception/app/socket"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// Execution : script runned
type Execution struct {
	ExecutionID int
	Name        string
	Script      string
}

type outstream struct {
	ch chan string
}

func newStream(ch chan string) *outstream {
	return &outstream{ch: ch}
}

func (out outstream) Write(p []byte) (int, error) {
	out.ch <- string(p)
	//	fmt.Println("wtf : " + string(p))
	return len(p), nil
}

type Response struct {
	Status int
	Data   string
}

func Run(script string, response chan Response) {
	uuid := uuid.NewV4()
	var file string

	if runtime.GOOS == "windows" {
		file = "/" + uuid.String() + ".sh"
	} else {
		file = "/tmp/bash_" + uuid.String() + ".sh"
	}
	fmt.Printf("script : %s\tfile : %s\n", script, file)

	if err := ioutil.WriteFile(file, []byte(script), 0777); err != nil {
		response <- Response{Status: 500, Data: err.Error()}
	}
	//defer os.Remove(file)
	room := socket.CreateRoom(uuid.String())

	go func() {
		var cmd *exec.Cmd
		cmd = exec.Command("bash", "-c", script)

		ch := make(chan string)
		out := newStream(ch)
		cmd.Stdout = out
		cmd.Stderr = out
		if err := cmd.Start(); err != nil {
			response <- Response{Status: 500, Data: err.Error()}
			//log.Fatal(err)
		}
		response <- Response{Status: 200, Data: uuid.String()}
		go func(ch chan string) {
			for {
				msg := <-ch
				fmt.Printf("on push dansle chan : %s\n", msg)
				room.Chan <- msg
			}
		}(ch)
		cmd.Wait()
		room.Chan <- "end_" + uuid.String()
	}()
}

// Validate Execution struct field for DB
func (exec *Execution) Validate(v *revel.Validation) {
	v.Required(exec.Name)
	v.Required(exec.Script)
}
