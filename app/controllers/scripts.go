package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/mistrale/jsonception/app/dispatcher"
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
type Scripts struct {
	GorpController
}

func init() {
	revel.OnAppStart(InitDB)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}

func (c Scripts) Index() revel.Result {

	exec := &models.Script{Uuid: uuid.NewV4().String()}
	return c.Render(exec)
}

// Index method to list all Script
func (c Scripts) All() revel.Result {
	var execs []models.Script
	c.Txn.Find(&execs)
	testID := 0
	return c.Render(execs, testID)
}

// Index method to list all Script
func (c Scripts) Get() revel.Result {
	var execs []models.Script
	c.Txn.Find(&execs)
	return c.RenderJson(utils.NewResponse(true, "", execs))
}

// GetOne method to routes GET /Script/:id
func (c Scripts) GetOne(id int) revel.Result {
	var exec models.Script
	c.Txn.First(&exec, id)

	return c.RenderJson(utils.NewResponse(true, "", exec))
}

// GetOneTemplate method to routes GET /Script/:id throw template
func (c Scripts) GetOneTemplate(scriptID int, uuid string) revel.Result {
	exec := &models.Script{}
	c.Txn.First(exec, scriptID)
	exec.Uuid = uuid
	c.Render(exec)
	return c.RenderTemplate("Scripts/one.html")
}

func (c Scripts) InitScriptModel(mode int) (*models.Script, error) {
	exec := &models.Script{}
	m := c.Request.MultipartForm
	name := c.Request.FormValue("name")
	content := c.Request.FormValue("content")
	params := c.Request.FormValue("parameters")
	fmt.Println(reflect.TypeOf(m))

	if name == "" {
		return nil, errors.New("Script name cannot be empty.")
	}
	if mode == CREATE {
		var execs []models.Script
		c.Txn.Find(&execs, "name = ?", name)
		if len(execs) > 0 {
			return nil, errors.New("Script name already taken.")
		}
	}
	if content == "" {
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
	exec.Content = content
	exec.Params = *parameters
	return exec, nil
}

// Create method to add new Script in DB
func (c Scripts) Create() revel.Result {
	exec, err := c.InitScriptModel(CREATE)
	if err != nil {
		c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	for _, v := range exec.Params {
		v.Print()
	}
	c.Txn.Create(exec)
	return c.RenderJson(utils.NewResponse(true, "Script successfully created", *exec))
}

func (c Scripts) Delete(scriptID int) revel.Result {
	if scriptID == 0 {
		return c.RenderJson(utils.NewResponse(false, "You need to provide id_exec", ""))
	}
	var exec models.Script
	fmt.Printf("id : %d\n", scriptID)
	c.Txn.First(&exec, scriptID)
	c.Txn.Delete(&exec)
	return c.RenderJson(utils.NewResponse(true, "", "Script deleted"))
}

func (c Scripts) Update(scriptID int) revel.Result {
	fmt.Printf("id_exec : %d\n", scriptID)
	if scriptID == 0 {
		return c.RenderJson(utils.NewResponse(false, "You need to provide id_exec", ""))
	}
	exec := &models.Script{}
	c.Txn.First(&exec, scriptID)

	new_exec, err := c.InitScriptModel(UPDATE)
	if err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	exec.Name = new_exec.Name
	exec.Content = new_exec.Content
	exec.Params = new_exec.Params
	c.Txn.Save(exec)

	// check test params
	var tests []models.Test
	c.Txn.Find(&tests, "script_id = ?", exec.ID)
	for i_test, test := range tests {

		// for each params in test
		for i_test_params, p_test := range test.Params {

			isHere := false
			// find params in exec params
			for _, p_exec := range exec.Params {
				if p_exec.Name == p_test.Name {
					isHere = true
				}
			}
			if !isHere {
				tests[i_test].Params = append(tests[i_test].Params[:i_test_params], tests[i_test].Params[i_test_params+1:]...)
			}
		}
		c.Txn.Save(&tests[i_test])
	}
	fmt.Printf("size test : %d\n", len(tests))
	return c.RenderJson(utils.NewResponse(true, "", "Script updated"))
}

// Run method to execute script
func (c Scripts) Run(scriptID int, content string, params []models.Parameters) revel.Result {
	uuid := uuid.NewV4()
	//channel := make(chan map[string]interface{})
	channel := make(chan dispatcher.Event)
	var exec models.Script

	if scriptID != 0 {
		c.Txn.First(&exec, scriptID)

	} else {
		exec.Content = content
	}
	exec.Uuid = uuid.String()

	go exec.Run(channel)
	response := <-channel
	if response.Status != true {
		fmt.Printf("status : %v\n", response.Status)
		return c.RenderJson(response)
	}
	room := socket.CreateRoom(uuid.String())

	go func(ch chan dispatcher.Event, exec_uuid string) {
		for {
			msg := <-channel
			room.Chan <- msg
			if msg.Type == models.RESULT_SCRIPT {
				room.Chan <- dispatcher.Event{Type: models.END_ROOM, Status: true, Body: "end_" + exec_uuid}
			}
		}
	}(channel, exec.Uuid)
	return c.RenderJson(utils.NewResponse(true, "", exec.Uuid))
}
