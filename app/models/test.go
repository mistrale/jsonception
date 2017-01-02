package models

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mistrale/jsonception/app/jsoncmp"
	"github.com/mistrale/jsonception/app/utils"
	"github.com/revel/revel"
)

// Test : script and logevent
type Test struct {
	TestID      int
	Name        string `json:"name"`
	Config      string `json:"config"`
	PathRefFile string `json:"log_events"`
	PathLogFile string `json:"path_log"`
	ExecutionID int    `json:"executionID"`
	Execution   *Execution
}

// Run method to realise test
func (test *Test) Run(response chan map[string]interface{}) {
	//	uuid := uuid.NewV4()
	b, err := ioutil.ReadFile(test.PathRefFile)
	if err != nil {
		response <- utils.NewResponse(false, "Error reading log reference file : "+err.Error(), nil)
		return
	}
	refJSon := make([]interface{}, 0)
	config := make(map[string]interface{})

	if err := json.Unmarshal([]byte(b), &refJSon); err != nil {
		response <- utils.NewResponse(false, err.Error(), nil)
		return
	}

	if err := json.Unmarshal([]byte(test.Config), &config); err != nil {
		response <- utils.NewResponse(false, err.Error(), nil)
		return
	}
	response <- utils.NewResponse(true, "matchs being", nil)

	//fmt.Printf(" test config : %s\n", test.Config)

	success := jsoncmp.CompareJSON(test.PathLogFile, refJSon, config)

	if success {
		response <- utils.NewResponse(true, "Files match.", nil)
	} else {
		response <- utils.NewResponse(true, "Files doesnt match.", nil)
	}
}

// Validate Reference struct field for DB
func (ref *Test) Validate(v *revel.Validation) {
	v.Check(ref.Name, revel.Required{})
	v.Check(ref.Config, revel.Required{})
	v.Check(ref.PathRefFile, revel.Required{})
	v.Check(ref.PathLogFile, revel.Required{})
	v.Required(ref.Execution)
}
