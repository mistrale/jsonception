package models

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/mistrale/jsonception/app/utils"
)

// Execution : script runned
type Execution struct {
	gorm.Model
	Name   string     `json:"name"`
	Script string     `json:"script"`
	Uuid   string     `json:"-" sql:"-"`
	Order  string     `json:"-" sql:"-"`
	Params Parameters `json:"parameters" sql:"type:jsonb"`
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
func (e *Execution) GetOrder() string {
	return e.Order
}

func (e *Execution) InitParams() {
	for _, v := range e.Params {
		v.Print()
		switch v.Type {
		case "bool":
			e.Script = strings.Replace(e.Script, "$"+v.Name, strconv.FormatBool(v.Value.(bool)), -1)
		case "int":
			if value, ok := v.Value.(float64); ok {
				e.Script = strings.Replace(e.Script, "$"+v.Name, strconv.FormatFloat(value, 'g', -1, 64), -1)
			} else {
				e.Script = strings.Replace(e.Script, "$"+v.Name, strconv.Itoa(v.Value.(int)), -1)
			}
		default:
			e.Script = strings.Replace(e.Script, "$"+v.Name, v.Value.(string), -1)
		}
	}
	fmt.Printf("script : %s\n", e.Script)
}

// Run method to start script
func (e Execution) Run(response chan map[string]interface{}) {
	fmt.Println("ON RUN")
	e.InitParams()
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
