package models

import "github.com/revel/revel"

// Execution : script runned
type Execution struct {
	ExecutionID int
	Name        string
	Script      string
}

// Validate Execution struct field for DB
func (exec *Execution) Validate(v *revel.Validation) {
	v.Required(exec.Name)
	v.Required(exec.Script)
}
