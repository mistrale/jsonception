package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/mistrale/jsonception/app/jsoncmp"
	"github.com/mistrale/jsonception/app/utils"
)

// Test : script and logevent
type Test struct {
	TestID      int       `json:"test_id"  gorm:"primary_key"`
	Name        string    `json:"name"`
	Config      string    `json:"config"`
	PathRefFile string    `json:"log_events"`
	PathLogFile string    `json:"path_log"`
	ExecutionID int       `json:"executionID"`
	Execution   Execution `json:"execution" gorm:"ForeignKey:ExecutionID;AssociationForeignKey:ExecutionID"`

	Uuid string `json:"-" db:"-"`
}

func (test Test) GetID() int {
	return test.TestID
}

// Run method to realise test
func (test Test) Run(response chan map[string]interface{}) {
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
	// get json event from two files to differ
	var testJson []interface{}
	if err := utils.GetJsonArray(test.PathLogFile, &testJson); err != nil {
		fmt.Println(err.Error())
		response <- utils.NewResponse(false, err.Error(), nil)
		return
	}

	if len(refJSon) != len(testJson) {
		fmt.Println("File diverg from different events number")
		response <- utils.NewResponse(false, "File diverg from different events number", nil)
		return
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
		jsonlogresp := refEvent.(map[string]interface{})
		jsonrefresp := testJson[i].(map[string]interface{})

		//json.MarshalIndent(jsonresp, "", "    ")
		str1, err := json.Marshal(jsonlogresp)
		if err != nil {
			response <- utils.NewResponse(false, err.Error(), nil)
			return

		}
		str2, err2 := json.Marshal(jsonrefresp)
		if err2 != nil {
			response <- utils.NewResponse(false, err.Error(), nil)
			return
		}
		resp := make(map[string]interface{})
		resp["event_type"] = TESTEVENT
		var prettyJSON1 bytes.Buffer
		err = json.Indent(&prettyJSON1, str1, "", "\t")
		if err != nil {
			response <- utils.NewResponse(false, err.Error(), nil)
			return
		}
		var prettyJSON2 bytes.Buffer
		err = json.Indent(&prettyJSON2, str2, "", "\t")
		if err != nil {
			response <- utils.NewResponse(false, err.Error(), nil)
			return
		}
		resp["body"] = make(map[string]interface{})
		resp["body"].(map[string]interface{})[TESTLOGEVENT] = string(prettyJSON1.Bytes())
		resp["body"].(map[string]interface{})[REFLOGEVENT] = string(prettyJSON2.Bytes())

		response <- utils.NewResponse(true, "", resp)

	}
	//success := jsoncmp.CompareJSON(test.PathLogFile, refJSon, config)
	resp := make(map[string]interface{})
	resp["event_type"] = RESULTEVENT
	resp["body"] = "files match !"

	response <- utils.NewResponse(true, "", resp)
	log.Println("JOB DONE")
}
