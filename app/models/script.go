package models

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/mistrale/jsonception/app/dispatcher"
)

// Script : script runned
type Script struct {
	gorm.Model
	Name    string     `json:"name"`
	Content string     `json:"content"`
	Uuid    string     `json:"-" sql:"-"`
	Order   string     `json:"-" sql:"-"`
	Params  Parameters `json:"parameters" sql:"type:jsonb"`
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
func (e *Script) GetOrder() string {
	return e.Order
}

func (e *Script) InitParams() {
	for _, v := range e.Params {
		v.Print()
		switch v.Type {
		case "bool":
			e.Content = strings.Replace(e.Content, "$"+v.Name, strconv.FormatBool(v.Value.(bool)), -1)
		case "int":
			if value, ok := v.Value.(float64); ok {
				e.Content = strings.Replace(e.Content, "$"+v.Name, strconv.FormatFloat(value, 'g', -1, 64), -1)
			} else {
				e.Content = strings.Replace(e.Content, "$"+v.Name, strconv.Itoa(v.Value.(int)), -1)
			}
		default:
			e.Content = strings.Replace(e.Content, "$"+v.Name, v.Value.(string), -1)
		}
	}
}

// Run method to exec script
func (e Script) Run(response chan dispatcher.Event) {
	e.InitParams()
	var cmd *exec.Cmd
	fmt.Printf("Content : %s\nWTTFFFF", e.Content)

	cmd = exec.Command("bash", "-c", e.Content)

	ch := make(chan string)
	out := newStream(ch)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Start(); err != nil {
		response <- dispatcher.Event{Type: RESULT_SCRIPT, Errors: []string{err.Error()}, Status: false, Body: nil}
	}
	response <- dispatcher.Event{Type: START_SCRIPT, Errors: nil, Status: true, Body: e.Uuid}
	go func(ch chan string) {
		for {
			response <- dispatcher.Event{Type: EVENT_SCRIPT, Errors: nil, Status: true, Body: <-ch}
		}
	}(ch)
	cmd.Wait()
	cmd.Process.Kill()
	response <- dispatcher.Event{Type: RESULT_SCRIPT, Errors: nil, Status: true, Body: "Script done !"}
}
