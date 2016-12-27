package controllers

import (
	"fmt"
	"log"

	"github.com/mistrale/jsonception/app/models"

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
	return c.Render()
}

// Index method to list all execution
func (c Executions) All() revel.Result {
	var execs []models.Execution
	_, err := c.Txn.Select(&execs,
		`select * from Execution`)
	if err != nil {
		panic(err)
	}
	return c.Render(execs)
}

// Create method to add new execution in DB
func (c Executions) Create(name, script string) revel.Result {
	fmt.Printf("name : %s\tscript : %s\n", name, script)
	c.Validation.Required(name).Message("Please enter an execution name")
	c.Validation.Required(script).Message("Script cannot be empty")

	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		fmt.Printf("error")

		c.Validation.Keep()
		c.FlashParams()
		return nil
	}
	err := c.Txn.Insert(&models.Execution{Name: name, Script: script})
	if err != nil {
		log.Println(err.Error())
	}
	return c.Redirect(Executions.All)
}

// Run method to execute script
func (c Executions) Run(script string) revel.Result {

	fmt.Printf("script LOL : %x\n", script)
	obj := make(map[string]interface{})
	ch := make(chan models.Response)
	models.Run(script, ch)
	response := <-ch
	if response.Status == 200 {
		obj["success"] = true
		obj["uuid"] = response.Data
	} else {
		obj["success"] = false
		obj["error"] = response.Data
		fmt.Printf("error : %s\n", response.Data)

	}
	return c.RenderJson(obj)
}
