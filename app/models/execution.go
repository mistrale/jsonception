package models

import (
	"fmt"
	"os/exec"

	"github.com/mistrale/jsonception/app/utils"
)

// Execution : script runned
type Execution struct {
	ExecutionID int    `json:"executionID" gorm:"primary_key"`
	Name        string `json:"name" sql:"unique"`
	Script      string `json:"script"`
	Uuid        string `json:"-" sql:"-"`
	Order       string `json:"-" sql:"-"`
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

// GetID method to retrieve model's id
func (e Execution) GetOrder() string {
	return e.Order
}

// Run method to start script
func (e Execution) Run(response chan map[string]interface{}) {
	fmt.Println("ON RUN")
	var cmd *exec.Cmd
	cmd = exec.Command("bash", "-c", e.Script)

	ch := make(chan string)
	out := newStream(ch)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Start(); err != nil {
		response <- utils.NewResponse(false, err.Error(), nil)
		//log.Fatal(err)
	}
	response <- utils.NewResponse(true, "ok", e.Uuid)
	go func(ch chan string) {
		for {
			msg := <-ch
			resp := make(map[string]interface{})
			resp["event_type"] = EXEC_EVENT
			resp["body"] = msg

			response <- utils.NewResponse(true, "", resp)
			//room.Chan <- msg
		}
	}(ch)
	cmd.Wait()
	cmd.Process.Kill()
	resp := make(map[string]interface{})
	resp["event_type"] = RESULT_EXEC
	resp["body"] = "Execution done !"
	response <- utils.NewResponse(true, "", resp)

}
