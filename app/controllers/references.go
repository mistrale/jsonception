package controllers

import (
	"fmt"

	"github.com/revel/revel"
)

// References controller
type References struct {
	GorpController
}

// Create method to add new execution in DB
func (c References) Create() revel.Result {
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

// Refresh method to reset a reference
func (c References) Refresh() revel.Result {
	return c.Render()
}

// Index method to page from reference index
func (c References) Index() revel.Result {
	return c.Render()
}

// All method to get all reference
func (c References) All() revel.Result {
	return c.Render()
}
