package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/mistrale/jsonception/app/dispatcher"
	"github.com/mistrale/jsonception/app/models"
	"github.com/mistrale/jsonception/app/socket"
	"github.com/mistrale/jsonception/app/utils"
	uuid "github.com/satori/go.uuid"

	"github.com/revel/revel"
)

// References controller
type Executions struct {
	GorpController
}

func init() {
	revel.OnAppStart(InitDB)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}

func (c Executions) Index() revel.Result {
	exec := models.Execution{ExecutionID: 0}
	return c.Render(exec)
}

// Index method to list all execution
func (c Executions) All() revel.Result {
	var execs []models.Execution
	_, err := c.Txn.Select(&execs,
		`select * from Execution`)
	if err != nil {
		panic(err)
	}
	testID := 0
	return c.Render(execs, testID)
}

// Index method to list all execution
func (c Executions) Get() revel.Result {
	var execs []models.Execution
	_, err := c.Txn.Select(&execs,
		`select * from Execution`)
	if err != nil {
		panic(err)
	}
	return c.RenderJson(utils.NewResponse(true, "", execs))
}

// Create method to add new execution in DB
func (c Executions) Create() revel.Result {
	exec := &models.Execution{}
	content, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(content, exec); err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	// check params
	if exec.Name == "" {
		return c.RenderJson(utils.NewResponse(false, "Execution name cannot be empty.", nil))
	}

	if exec.Script == "" {
		return c.RenderJson(utils.NewResponse(false, "Script name cannot be empty.", nil))
	}

	fmt.Printf("name : %s\tscript : %s\n", exec.Name, exec.Script)

	err := c.Txn.Insert(exec)
	if err != nil {
		log.Println(err.Error())
	}
	//	dispatcher.AddWorker(work_ID)
	return c.RenderJson(utils.NewResponse(true, "Execution successfully created", *exec))
}

// Run method to execute script

func (c Executions) Run(id_exec int, script string) revel.Result {
	uuid := uuid.NewV4()

	var request dispatcher.WorkRequest
	if id_exec != 0 {
		var exec models.Execution
		if err := c.Txn.SelectOne(&exec, "select * from Execution where ExecutionID=?", id_exec); err != nil {
			return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
		}
		request = dispatcher.WorkRequest{Uuid: uuid.String(), Script: exec.Script, Response: make(chan map[string]interface{})}
		dispatcher.WorkQueue[exec.ExecutionID] <- request
	} else {
		request = dispatcher.WorkRequest{Uuid: uuid.String(), Script: script, Response: make(chan map[string]interface{})}
		work := dispatcher.Worker{}
		work.Run(request)
	}

	response := <-request.Response
	if response["status"] != true {
		fmt.Printf("status : %s\n", response)
		return c.RenderJson(response)
	}
	room := socket.CreateRoom(uuid.String())

	go func(ch chan map[string]interface{}) {
		for {
			msg := <-request.Response
			fmt.Printf("on push dansle chan : %s\n", msg)
			room.Chan <- msg
		}
	}(request.Response)
	return c.RenderJson(response)
}
