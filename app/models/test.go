package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-gorp/gorp"
	"github.com/mistrale/jsonception/app/jsoncmp"
	"github.com/mistrale/jsonception/app/utils"
	"github.com/revel/revel"
)

// Test : script and logevent
type Test struct {
	TestID      int    `json:"test_id"`
	Name        string `json:"name"`
	Config      string `json:"config"`
	PathRefFile string `json:"log_events"`
	PathLogFile string `json:"path_log"`
	ExecutionID int    `json:"executionID"`
	Execution   *Execution
	Uuid        string `json:"-" db:"-"`
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
	// get json event from two files to differ
	var testJson []interface{}
	if err := utils.GetJsonArray(test.PathLogFile, &testJson); err != nil {
		fmt.Println(err.Error())
		response <- utils.NewResponse(false, err.Error(), nil)
	}

	if len(refJSon) != len(testJson) {
		fmt.Println("File diverg from different events number")
		response <- utils.NewResponse(false, "File diverg from different events number", nil)
	}

	// iterate through reference file
	for i, refEvent := range refJSon {
		params1 := utils.CopyMap(config)

		// get keys to compare for current event
		eventParams := jsoncmp.FindParameters(refEvent.(map[string]interface{}), params1)
		if eventParams == nil {
			response <- utils.NewResponse(false, "", nil)
		}
		// compare two events
		if err := jsoncmp.CompareEvent(refEvent.(map[string]interface{}), testJson[i].(map[string]interface{}), eventParams); err != nil {
			fmt.Printf("Error : %s\n", err.Error())
			response <- utils.NewResponse(false, err.Error(), nil)
		}
	}
	//success := jsoncmp.CompareJSON(test.PathLogFile, refJSon, config)
	response <- utils.NewResponse(true, "", "Files match !")
}

// Validate Reference struct field for DB
func (ref *Test) Validate(v *revel.Validation) {
	v.Check(ref.Name, revel.Required{})
	v.Check(ref.Config, revel.Required{})
	v.Check(ref.PathRefFile, revel.Required{})
	v.Check(ref.PathLogFile, revel.Required{})
	v.Required(ref.Execution)
}

func (t *Test) PostGet(exe gorp.SqlExecutor) error {
	var (
		obj interface{}
		err error
	)

	obj, err = exe.Get(Execution{}, t.ExecutionID)
	if err != nil {
		return fmt.Errorf("Error loading a test's execution (%d): %s", t.ExecutionID, err)
	}
	t.Execution = obj.(*Execution)
	return nil
}
