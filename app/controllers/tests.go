package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mistrale/jsonception/app/dispatcher"
	"github.com/mistrale/jsonception/app/models"
	"github.com/mistrale/jsonception/app/socket"
	"github.com/mistrale/jsonception/app/utils"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// References controller
type Tests struct {
	GorpController
}

func (c Tests) loadExecutionByID(id int) *models.Execution {
	h, err := c.Txn.Get(models.Execution{}, id)
	if err != nil {
		panic(err)
	}
	if h == nil {
		return nil
	}
	return h.(*models.Execution)
}

// Create method to add new execution in DB
func (c Tests) Create() revel.Result {
	test := &models.Test{}
	content, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(content, test); err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	fmt.Printf("Params : name = %s\tconfig = %s\tlogevent = %s\tpath_file : %s\t exeuc iD : %d\n", test.Name, test.Config, test.PathRefFile, test.PathLogFile, test.ExecutionID)

	// check params
	if test.Name == "" {
		return c.RenderJson(utils.NewResponse(false, "Test name cannot be empty.", nil))
	}

	if test.Config == "" {
		return c.RenderJson(utils.NewResponse(false, "Json Configuration name cannot be empty.", nil))
	}

	if test.PathLogFile == "" {
		return c.RenderJson(utils.NewResponse(false, "Path log file name cannot be empty.", nil))
	}

	if _, err := ioutil.ReadFile(test.PathLogFile); err != nil {
		return c.RenderJson(utils.NewResponse(false, "Error reading log file.", nil))
	}

	if _, err := ioutil.ReadFile(test.PathRefFile); err != nil {
		return c.RenderJson(utils.NewResponse(false, "Error reading log reference file.", nil))
	}

	var tests []models.Test

	if _, err := c.Txn.Select(&tests,
		`select * from Test where PathLogFile=?`, test.PathLogFile); err != nil {
		return c.RenderJson(utils.NewResponse(true, err.Error(), nil))
	}

	if len(tests) >= 1 {
		return c.RenderJson(utils.NewResponse(false, "A log file path of same name is already taken", nil))
	}
	// insert ref with ExecutionID
	if err := c.Txn.Insert(test); err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	return c.RenderJson(utils.NewResponse(true, "Successful test creation", *test))
}

// Refresh method to reset a reference
func (c Tests) Run(testID int) revel.Result {
	var test models.Test
	uuid := uuid.NewV4()
	room := socket.CreateRoom(uuid.String())
	var request dispatcher.WorkRequest

	if err := c.Txn.SelectOne(&test,
		`select * from Test where testID=?`, testID); err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	fmt.Printf("print result : %d and exec id : %d\n", testID, test.ExecutionID)

	// if there is an execution
	if test.ExecutionID != 0 {
		var exec models.Execution

		if err := c.Txn.SelectOne(&exec, "select * from Execution where ExecutionID=?", test.ExecutionID); err != nil {
			return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
		}
		request = dispatcher.WorkRequest{Uuid: uuid.String(), Script: exec.Script, Response: make(chan map[string]interface{})}
		dispatcher.WorkQueue[exec.ExecutionID] <- request
		response := <-request.Response
		if response["status"] != true {
			fmt.Printf("status : %s\n", response)
			return c.RenderJson(response)
		}
	}
	go func() {
		if test.ExecutionID != 0 {
			for {
				msg := <-request.Response
				if msg["status"] != true {
					fmt.Printf("status : %s\n", msg)
					room.Chan <- msg
					return
				}
				if msg["response"] == "end_"+uuid.String() {
					break
				}
				room.Chan <- msg
			}
		}
		ch := make(chan map[string]interface{})
		go test.Run(ch)
		go func(ch chan map[string]interface{}) {
			for {
				msg := <-ch
				response := make(map[string]interface{})
				response["type"] = "test_event"
				response["body"] = msg
				fmt.Printf("on push dansle chan : %s\n", msg)
				room.Chan <- msg
			}
		}(ch)
	}()
	//c.Executions.Run("lslsl")
	return c.RenderJson(utils.NewResponse(true, "", uuid.String()))
}

// Refresh method to reset a reference
func (c Tests) Show(testID int) revel.Result {
	return c.Render()
}

// Index method to page from reference index
func (c Tests) Index() revel.Result {
	test := models.Test{TestID: 0, Execution: &models.Execution{}}
	return c.Render(test)
}

// All method to get all reference
func (c Tests) All() revel.Result {
	var tests []models.Test
	_, err := c.Txn.Select(&tests,
		`select * from Test`)
	if err != nil {
		panic(err)
	}
	for i, v := range tests {
		exec := &models.Execution{}
		err = c.Txn.SelectOne(exec, "select * from Execution where ExecutionID=?", v.ExecutionID)
		fmt.Printf("executID name : %s\n", exec.Name)
		tests[i].Execution = exec
	}
	fmt.Printf("nm test : %d\n", len(tests))
	return c.Render(tests)
}

// Get method to get all reference
func (c Tests) Get() revel.Result {
	var tests []models.Test
	_, err := c.Txn.Select(&tests,
		`select * from Test`)
	if err != nil {
		panic(err)
	}
	for i, v := range tests {
		exec := &models.Execution{}
		err = c.Txn.SelectOne(exec, "select * from Execution where ExecutionID=?", v.ExecutionID)
		fmt.Printf("executID name : %s\n", exec.Name)
		tests[i].Execution = exec
	}
	return c.RenderJson(tests)
}
