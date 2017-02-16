package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mistrale/jsonception/app/dispatcher"
	"github.com/mistrale/jsonception/app/models"
	"github.com/mistrale/jsonception/app/socket"
	"github.com/mistrale/jsonception/app/utils"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// Tests controller
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

	// insert ref with ExecutionID
	c.Txn.Create(test)
	return c.RenderJson(utils.NewResponse(true, "Successful test creation", *test))
}

func (c Tests) Delete(id_test int) revel.Result {
	var test models.Test
	c.Txn.First(&test, id_test)
	c.Txn.Delete(&test)
	return c.RenderJson(utils.NewResponse(true, "", "Test deleted"))
}

func (c Tests) Update(id_test int) revel.Result {
	if id_test == 0 {
		return c.RenderJson(utils.NewResponse(false, "You need to provide id_test", ""))
	}
	test := &models.Test{}
	fmt.Printf("id test : %d\n", id_test)
	content, _ := ioutil.ReadAll(c.Request.Body)
	c.Txn.First(&test, id_test)
	if err := json.Unmarshal(content, test); err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	fmt.Printf("content : %s\n", content)
	c.Txn.Save(&test)
	return c.RenderJson(utils.NewResponse(true, "", "Test updated"))
}

// Generic Run test method
func RunTest(test *models.Test, n_uuid string, room chan map[string]interface{}, db *gorm.DB) bool {
	var request dispatcher.WorkRequest
	output := ""

	test.Uuid = n_uuid
	test.Execution.Uuid = n_uuid
	// create history
	history := &models.TestHistory{TestID: test.TestID, RunUUID: n_uuid}

	var runner models.IRunnable = test.Execution

	// if there is an execution
	if test.ExecutionID != 0 {
		request = dispatcher.WorkRequest{Runner: &runner, Response: make(chan map[string]interface{})}
		dispatcher.WorkQueue <- request
		response := <-request.Response
		if response["status"] != true {
			updateEndHistory(db, history, "", "", "", response["message"].(string), test.Name, false)
			room <- response
			return false
		}
		response["response"] = make(map[string]interface{})
		response["response"].(map[string]interface{})["event_type"] = models.START_TEST
		response["response"].(map[string]interface{})["time_runned"] = time.Now().UnixNano()
		room <- response
	}

	go func() {

		if test.ExecutionID != 0 {
			for {
				fmt.Printf("chaann addr 2 : %p\n", room)
				msg := <-request.Response
				room <- msg

				if msg["status"] != true {
					updateEndHistory(db, history, output, "", "", msg["message"].(string), test.Name, false)
					return
				}
				if response, ok := msg["response"].(map[string]interface{}); ok {
					if response["event_type"] == models.RESULT_EXEC {
						break
					}
					output += msg["response"].(map[string]interface{})["body"].(string)
				}
				//				fmt.Printf("OUTPUT : %s\n", msg)
			}
		}
		runner = test
		request := dispatcher.WorkRequest{Runner: &runner, Response: make(chan map[string]interface{})}
		dispatcher.WorkQueue <- request
		ch := request.Response

		//go test.Run(ch)

		reflog := ""
		testlog := ""

		for {

			msg := <-ch
			//fmt.Printf("response : %s\n", msg)

			if msg["status"] != true {
				updateEndHistory(db, history, output, reflog, testlog, msg["message"].(string), test.Name, false)
				room <- msg
				return
			}
			response := msg["response"].(map[string]interface{})
			if response["event_type"] == models.TESTEVENT {
				reflog += response["body"].(map[string]interface{})[models.REFLOGEVENT].(string)
				testlog += response["body"].(map[string]interface{})[models.TESTLOGEVENT].(string)
			} else if response["event_type"] == models.RESULTEVENT {
				updateEndHistory(db, history, output, reflog, testlog, response["body"].(string), test.Name, true)
				room <- msg

				return
			}
			room <- msg
		}
	}()
	return true
}

// Run method to start a test
func (c Tests) Run(testID int) revel.Result {
	var test models.Test

	test_uuid := uuid.NewV4()
	room := socket.CreateRoom(test_uuid.String())
	c.Txn.Preload("Execution").First(&test, testID)

	// run test
	if !RunTest(&test, test_uuid.String(), room.Chan, c.Txn) {
		return c.RenderJson(<-room.Chan)
	}
	return c.RenderJson(utils.NewResponse(true, "", test_uuid.String()))
}

// GetHistory method to get all history from one test
func (c Tests) GetHistory(testID string) revel.Result {
	var history []models.TestHistory
	c.Txn.Where("test_id = ?", testID).Find(&history)
	fmt.Printf("format uuid : %s\n", testID)
	return c.RenderJson(history)
}

// GetHistory method to get all history from one test and template html
func (c Tests) GetHistoryTemplate(testID int) revel.Result {
	var history []models.TestHistory
	c.Txn.Where("test_id = ?", testID).Find(&history)
	c.Render(history)
	return c.RenderTemplate("TestHistory/all.html")
}

// Refresh method to reset a reference
func (c Tests) Show(testID int) revel.Result {
	return c.Render()
}

// Index method to page from reference index
func (c Tests) Index() revel.Result {
	test := &models.Test{TestID: 0, Execution: models.Execution{}, Uuid: uuid.NewV4().String()}
	return c.Render(test)
}

// All method to get all reference
func (c Tests) All() revel.Result {
	var tests []models.Test
	c.Txn.Preload("Execution").Find(&tests)
	return c.Render(tests)
}

// Get method to get all reference
func (c Tests) Get() revel.Result {
	var tests []models.Test
	c.Txn.Preload("Execution").Find(&tests)
	return c.RenderJson(tests)
}

// Get method to get all reference
func (c Tests) GetOneTemplate(testID int) revel.Result {
	test := &models.Test{}
	c.Txn.Preload("Execution").First(test, testID)
	c.Render(test)
	return c.RenderTemplate("Tests/TestTemplate.html")
}
