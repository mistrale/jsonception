package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

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
	ip := "10.177.17.101"

	return c.Render(exec, ip)
}

// Index method to list all execution
func (c Executions) All() revel.Result {
	var execs []models.Execution
	c.Txn.Find(&execs)
	testID := 0
	ip := "10.177.17.101"
	return c.Render(execs, testID, ip)
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
	if id_exec == 0 {
		return c.RenderJson(utils.NewResponse(false, "You need to provide id_exec", ""))
	}
	exec := &models.Execution{}
	c.Txn.First(&exec, id_exec)
	content, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(content, exec); err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	c.Txn.Save(&exec)
	//	fmt.Printf("id exec : %d\tname : %s\n", id_exec, name)
	//  db.Model(&exec).Update("Name", name, )
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
