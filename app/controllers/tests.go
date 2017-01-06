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
	output := ""

	if err := c.Txn.SelectOne(&test,
		`select * from Test where testID=?`, testID); err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}

	// create history
	history := &models.TestHistory{TestID: testID, RunUUID: uuid.String(), Status: "running"}
	if err := c.Txn.Insert(history); err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}

	// if there is an execution
	if test.Execution != nil {
		request = dispatcher.WorkRequest{Uuid: uuid.String(), Script: test.Execution.Script, Response: make(chan map[string]interface{})}
		dispatcher.WorkQueue[test.ExecutionID-1] <- request
		response := <-request.Response
		if response["status"] != true {
			fmt.Printf("status : %s\n", response)
			updateEndHistory(history, "", "", "", response["message"].(string), false)
			return c.RenderJson(response)
		}
	}
	go func() {

		if test.Execution != nil {
			for {
				msg := <-request.Response
				if msg["status"] != true {
					fmt.Printf("status : %s\n", msg)
					updateEndHistory(history, output, "", "", msg["message"].(string), false)
					room.Chan <- msg
					return
				}
				if msg["response"] == "end_"+uuid.String() {
					break
				}
				output += msg["response"].(map[string]interface{})["body"].(string)
				room.Chan <- msg
			}
		}
		ch := make(chan map[string]interface{})
		go test.Run(ch)

		reflog := ""
		testlog := ""

		for {
			msg := <-ch
			room.Chan <- msg
			if msg["status"] != true {
				updateEndHistory(history, output, reflog, testlog, msg["message"].(string), false)
				break
			}
			fmt.Printf("response : %s\n", msg)
			response := msg["response"].(map[string]interface{})
			if response["type"] == models.TESTEVENT {
				reflog += response["body"].(map[string]interface{})[models.REFLOGEVENT].(string)
				testlog += response["body"].(map[string]interface{})[models.TESTLOGEVENT].(string)
			} else if response["type"] == models.RESULTEVENT {
				updateEndHistory(history, output, reflog, testlog, response["body"].(string), true)
				break
			}
		}
	}()
	return c.RenderJson(utils.NewResponse(true, "", uuid.String()))
}

// GetHistory method to get all history from one test
func (c Tests) GetHistory(testID int) revel.Result {
	var history []models.TestHistory

	if _, err := c.Txn.Select(&history,
		`select * from TestHistory where TestID=?`, testID); err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	return c.RenderJson(history)
}

// GetHistory method to get all history from one test and template html
func (c Tests) GetHistoryTemplate(testID int) revel.Result {
	var history []models.TestHistory

	if _, err := c.Txn.Select(&history,
		`select * from TestHistory where TestID=?`, testID); err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	c.Render(history)
	return c.RenderTemplate("TestHistory/all.html")
}

// Refresh method to reset a reference
func (c Tests) Show(testID int) revel.Result {
	return c.Render()
}

// Index method to page from reference index
func (c Tests) Index() revel.Result {
	test := models.Test{TestID: 0, Execution: &models.Execution{}, Uuid: uuid.NewV4().String()}
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
