package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/mistrale/jsonception/app/models"
	"github.com/mistrale/jsonception/app/socket"
	"github.com/mistrale/jsonception/app/utils"
	uuid "github.com/satori/go.uuid"

	"github.com/revel/revel"
)

const (
	CREATE = 1
	UPDATE = 2
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

	exec := &models.Execution{Uuid: uuid.NewV4().String()}
	return c.Render(exec)
}

// Index method to list all execution
func (c Executions) All() revel.Result {
	var execs []models.Execution
	c.Txn.Find(&execs)
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
func (c Executions) GetOneTemplate(executionID int, uuid string) revel.Result {
	exec := &models.Execution{}
	c.Txn.First(exec, executionID)
	exec.Uuid = uuid
	c.Render(exec)
	return c.RenderTemplate("Executions/one.html")
}

func (c Executions) InitExecutionModel(mode int) (*models.Execution, error) {
	exec := &models.Execution{}
	m := c.Request.MultipartForm
	name := c.Request.FormValue("name")
	script := c.Request.FormValue("script")
	params := c.Request.FormValue("parameters")
	fmt.Println(reflect.TypeOf(m))

	if name == "" {
		return nil, errors.New("Execution name cannot be empty.")
	}
	if mode == CREATE {
		var execs []models.Execution
		c.Txn.Find(&execs, "name = ?", name)
		if len(execs) > 0 {
			return nil, errors.New("Execution name already taken.")
		}
	}
	if script == "" {
		return nil, errors.New("Script name cannot be empty.")
	}
	parameters := &models.Parameters{}
	if err := json.Unmarshal([]byte(params), parameters); err != nil {
		return nil, err
	}
	if err := parameters.Check(); err != nil {
		return nil, err
	}
	if err := parameters.UploadFileFromParameters(m); err != nil {
		return nil, err
	}
	exec.Name = name
	exec.Script = script
	exec.Params = *parameters
	return exec, nil
}

// Create method to add new execution in DB
func (c Executions) Create() revel.Result {
	exec, err := c.InitExecutionModel(CREATE)
	if err != nil {
		c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	c.Txn.Create(exec)
	return c.RenderJson(utils.NewResponse(true, "Execution successfully created", *exec))
}

func (c Executions) Delete(id_exec int) revel.Result {
	if id_exec == 0 {
		return c.RenderJson(utils.NewResponse(false, "You need to provide id_exec", ""))
	}
	var exec models.Execution
	fmt.Printf("id : %d\n", id_exec)
	c.Txn.First(&exec, id_exec)
	c.Txn.Delete(&exec)
	return c.RenderJson(utils.NewResponse(true, "", "Execution deleted"))
}

func (c Executions) Update(id_exec int) revel.Result {
	fmt.Printf("id_exec : %d\n", id_exec)
	if id_exec == 0 {
		return c.RenderJson(utils.NewResponse(false, "You need to provide id_exec", ""))
	}
	exec := &models.Execution{}
	c.Txn.First(&exec, id_exec)

	new_exec, err := c.InitExecutionModel(UPDATE)
	if err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	exec.Name = new_exec.Name
	exec.Script = new_exec.Script
	exec.Params = new_exec.Params
	c.Txn.Save(exec)
	return c.RenderJson(utils.NewResponse(true, "", "Execution updated"))
}

// Run method to execute script
func (c Executions) Run(id_exec int, script string) revel.Result {
	uuid := uuid.NewV4()
	channel := make(chan map[string]interface{})
	var exec models.Execution

	if id_exec != 0 {
		c.Txn.First(&exec, id_exec)

	} else {
		exec.Script = script
	}
	exec.Uuid = uuid.String()

	go exec.Run(channel)
	response := <-channel
	if response["status"] != true {
		fmt.Printf("status : %s\n", response)
		return c.RenderJson(response)
	}
	room := socket.CreateRoom(uuid.String())

	go func(ch chan map[string]interface{}, exec_uuid string) {
		for {
			msg := <-channel
			room.Chan <- msg
			if msg["response"].(map[string]interface{})["event_type"] == models.RESULT_EXEC {
				room.Chan <- utils.NewResponse(true, "", "end_"+exec_uuid)
			}
		}
	}(channel, exec.Uuid)
	return c.RenderJson(response)
}
