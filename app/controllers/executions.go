package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

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

	exec := &models.Execution{ExecutionID: 0, Uuid: uuid.NewV4().String()}
	return c.Render(exec)
}

// Index method to list all execution
func (c Executions) All() revel.Result {
	var execs []models.Execution
	c.Txn.Find(&execs)

	// _, err := c.Txn.Select(&execs,
	// 	`select * from Execution`)
	// if err != nil {
	// 	return c.RenderJson(utils.NewResponse(false, "", err.Error()))
	// }
	testID := 0
	return c.Render(execs, testID)
}

// Index method to list all execution
func (c Executions) Get() revel.Result {
	var execs []models.Execution
	c.Txn.Find(&execs)

	return c.RenderJson(utils.NewResponse(true, "", execs))
}

// GetOne method to routes GET /execution/:id
func (c Executions) GetOne(id int) revel.Result {
	var exec models.Execution
	c.Txn.First(&exec, id)

	return c.RenderJson(utils.NewResponse(true, "", exec))
}

// GetOneTemplate method to routes GET /execution/:id throw template
func (c Executions) GetOneTemplate(id int) revel.Result {
	var exec models.Execution
	c.Txn.First(&exec, id)

	uuid := uuid.NewV4()
	c.Render(exec, uuid)
	return c.RenderTemplate("Executions/test.html")
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

	c.Txn.Create(exec)
	// if err != nil {
	// 	return c.RenderJson(utils.NewResponse(false, "", err.Error()))
	// }
	//dispatcher.AddWorker()
	return c.RenderJson(utils.NewResponse(true, "Execution successfully created", *exec))
}

// Run method to execute script
func (c Executions) Run(id_exec int, script string) revel.Result {
	uuid := uuid.NewV4()

	var request dispatcher.WorkRequest
	var exec models.Execution

	if id_exec != 0 {
		c.Txn.First(&exec, id_exec)
	} else {
		exec.Script = script
	}
	exec.Uuid = uuid.String()

	var runner models.IRunnable = exec
	request = dispatcher.WorkRequest{Runner: &runner, Response: make(chan map[string]interface{})}
	dispatcher.WorkQueue <- request

	response := <-request.Response
	if response["status"] != true {
		fmt.Printf("status : %s\n", response)
		return c.RenderJson(response)
	}
	room := socket.CreateRoom(uuid.String())

	go func(ch chan map[string]interface{}) {
		for {
			room.Chan <- <-request.Response
		}
	}(request.Response)
	fmt.Println("ALLELULIA")
	return c.RenderJson(response)
}
