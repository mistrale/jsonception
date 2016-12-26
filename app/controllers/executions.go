package controllers

import (
	"fmt"

	"github.com/mistrale/jsonception/app/models"

	"github.com/revel/revel"
)

// References controller
type Executions struct {
	GorpController
}

// Index method to list all execution
func (c Executions) Index() revel.Result {
	var execs []models.Execution
	_, err := c.Txn.Select(&execs,
		`select * from Execution`)
	if err != nil {
		panic(err)
	}
	return c.Render(execs)
}

// Create method to add new execution in DB
func (c Executions) Create() revel.Result {
	name := c.Params.Get("execution.Name")
	script := []byte(c.Params.Get("execution.Script"))

	fmt.Printf("Params : name = %s\tscript = %s\n", name, script)

	c.Validation.Required(name).Message("Please enter an execution name")
	c.Validation.Required(script).Message("Script cannot be empty")

	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		fmt.Printf("error")

		c.Validation.Keep()
		c.FlashParams()
		return nil
	}
	return c.Render()
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
