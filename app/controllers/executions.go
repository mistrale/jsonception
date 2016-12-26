package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/mistrale/jsonception/app/models"

	"github.com/mistrale/jsonception/app/socket"

	"github.com/revel/revel"
	"github.com/satori/go.uuid"
)

// References controller
type Executions struct {
	GorpController
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

// Index method to list all execution
func (c Executions) Index() revel.Result {
	var execs []models.Execution
	_, err := c.Txn.Select(&execs,
		`select * from Execution`)
	if err != nil {
		panic(err)
	}
	return c.Render(execs)
}

// Create method to add new execution in DB
func (c Executions) Create() revel.Result {
	name := c.Params.Get("execution.Name")
	script := []byte(c.Params.Get("execution.Script"))

	fmt.Printf("Params : name = %s\tscript = %s\n", name, script)

	c.Validation.Required(name).Message("Please enter an execution name")
	c.Validation.Required(script).Message("Script cannot be empty")

	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		fmt.Printf("error")

		c.Validation.Keep()
		c.FlashParams()
		return nil
	}
	return c.Render()
}

// Run method to execute script
func (c Executions) Run(script string) revel.Result {
	obj := make(map[string]interface{})
	obj["success"] = false

	err := ioutil.WriteFile("/tmp/bash.sh", []byte(script), 0777)
	if err != nil {
		obj["error"] = err.Error()
		return c.RenderJson(obj)
	}
	uuid := uuid.NewV4()
	room := socket.CreateRoom(uuid.String())
	go func() {
		// exec command
		cmd := exec.Command("/tmp/test.sh")
		ch := make(chan string)
		out := newStream(ch)
		cmd.Stdout = out
		cmd.Stderr = out

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		go func(ch chan string) {
			for {
				msg := <-ch
				room.Chan <- msg
			}
		}(ch)
		cmd.Wait()
		room.Chan <- "end_" + uuid.String()
	}()
	obj["success"] = true
	obj["uuid"] = uuid
	return c.RenderJson(obj)
}
