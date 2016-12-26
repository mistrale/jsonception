package models

import "github.com/revel/revel"

// Reference : script and logevent
type Reference struct {
	ReferenceID int
	Name        string
	Config      string
	Log         string

	ExecutionID string
	Execution   *Execution
}

// Validate Reference struct field for DB
func (ref *Reference) Validate(v *revel.Validation) {
	v.Check(ref.Name, revel.Required{})
	v.Check(ref.Log, revel.Required{})
	v.Check(ref.Config, revel.Required{})
	v.Required(ref.Execution)
}
