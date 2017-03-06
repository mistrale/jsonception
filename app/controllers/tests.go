package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

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

func (c Tests) InitTestModel(mode int) (*models.Test, error) {
	test := &models.Test{}
	m := c.Request.MultipartForm
	name := c.Request.FormValue("name")
	config := c.Request.FormValue("config")
	log_events := c.Request.FormValue("path_test_log")
	path_log := c.Request.FormValue("path_ref_log")
	params := c.Request.FormValue("parameters")
	ScriptID, err := strconv.Atoi(c.Request.FormValue("scriptID"))

	fmt.Printf("name : %s\tconfig : %s\tlog_events : %s\tpath_log : %s\tparams : %s\tScript ID : %d\n", name, config, log_events, path_log, params, ScriptID)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, errors.New("Test name cannot be empty.")
	}
	if mode == CREATE {
		var tests []models.Test
		c.Txn.Find(&tests, "name = ?", name)
		if len(tests) > 0 {
			return nil, errors.New("Test name already taken.")
		}
	}
	if path_log == "" {
		return nil, errors.New("path_log path name cannot be empty.")
	}
	if log_events == "" {
		return nil, errors.New("Logs_event path name cannot be empty.")
	}
	if config == "" {
		return nil, errors.New("Config name cannot be empty.")
	}
	var check_config []map[string]interface{}
	if err := json.Unmarshal([]byte(config), &check_config); err != nil {
		return nil, err
	}
	parameters := &models.Parameters{}
	if err := json.Unmarshal([]byte(params), parameters); err != nil {
		fmt.Printf("ici")
		return nil, err
	}
	if err := parameters.Check(); err != nil {
		return nil, err
	}
	if err := parameters.UploadFileFromParameters(m); err != nil {
		return nil, err
	}
	test.Name = name
	test.Config = config
	test.Params = *parameters
	test.PathLogFile = log_events
	test.ScriptID = ScriptID
	test.PathRefFile = path_log
	return test, nil
}

// Create method to add new Script in DB
func (c Tests) Create() revel.Result {
	test, err := c.InitTestModel(CREATE)
	if err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	exec := &models.Script{}
	c.Txn.First(exec, test.ScriptID)
	if err := test.Params.CheckTestParamsWithExecParams(&exec.Params); err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
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
	c.Txn.Preload("Script").First(&test, id_test)

	new_test, err := c.InitTestModel(UPDATE)
	if err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	c.Txn.First(&new_test.Script, new_test.ScriptID)

	if err := new_test.Params.CheckTestParamsWithExecParams(&new_test.Script.Params); err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	new_test.ID = test.ID
	c.Txn.Save(new_test)
	return c.RenderJson(utils.NewResponse(true, "Test updated", *test))
}

// Run method to start a test
func (c Tests) Run(testID int) revel.Result {
	var test models.Test
	channel := make(chan dispatcher.Event)
	test_uuid := uuid.NewV4().String()
	room := socket.CreateRoom(test_uuid)

	c.Txn.Preload("Script").First(&test, testID)
	test.Order = "order_test_" + strconv.Itoa(int(test.ID))
	test.Uuid = test_uuid

	go func(channel chan dispatcher.Event, test_uuid string) {
		go test.Run(channel)
		for {
			msg := <-channel
			room.Chan <- msg
			if msg.Status == false {
				room.Chan <- dispatcher.Event{Status: true, Body: "end_" + test_uuid, Type: models.END_ROOM}
				fmt.Printf("Hstory : %d\t%d\t%s\n", test.History.ID, test.History.TestID, test.History.OutputExec)
				Dbm.Create(&test.History)
				break
			}
			if msg.Type == models.RESULT_TEST {
				fmt.Printf("Hstory : %d\t%d\t%s\n", test.History.ID, test.History.TestID, test.History.OutputExec)

				room.Chan <- dispatcher.Event{Status: true, Body: "end_" + test_uuid, Type: models.END_ROOM}
				Dbm.Create(&test.History)
				break
			}
		}
	}(channel, test_uuid)
	return c.RenderJson(utils.NewResponse(true, "", test_uuid))
}

// GetHistory method to get all history from one test
func (c Tests) GetHistory(testID string) revel.Result {
	var history []models.TestHistory
	c.Txn.Where("test_id = ?", testID).Find(&history)
	return c.RenderJson(history)
}

// GetHistory method to get all history from one test and template html
func (c Tests) GetHistoryTemplate(testID int) revel.Result {
	var history []models.TestHistory

	c.Txn.Where("test_id = ?", testID).Find(&history)
	fmt.Println(len(history))

	c.Render(history)
	return c.RenderTemplate("TestHistory/all.html")
}

// Refresh method to reset a reference
func (c Tests) Show(testID int) revel.Result {
	return c.Render()
}

// Index method to page from reference index
func (c Tests) Index() revel.Result {
	test := &models.Test{Script: models.Script{}, Uuid: uuid.NewV4().String()}
	return c.Render(test)
}

// All method to get all reference
func (c Tests) All() revel.Result {
	var tests []models.Test
	c.Txn.Preload("Script").Find(&tests)
	return c.Render(tests)
}

// Get method to get all reference
func (c Tests) Get() revel.Result {
	var tests []models.Test
	c.Txn.Preload("Script").Find(&tests)
	return c.RenderJson(tests)
}

// Get method to get all reference
func (c Tests) GetOneTemplate(testID int) revel.Result {
	test := &models.Test{}
	c.Txn.Preload("Script").First(test, testID)
	c.Render(test)
	return c.RenderTemplate("Tests/TestTemplate.html")
}
