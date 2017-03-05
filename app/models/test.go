package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/mistrale/jsonception/app/jsoncmp"
	"github.com/mistrale/jsonception/app/utils"
)

// Test : script and logevent
type Test struct {
	TestID      int        `json:"test_id"  gorm:"primary_key"`
	Name        string     `json:"name"`
	Config      string     `json:"config"`
	PathRefFile string     `json:"path_test_log"`
	PathLogFile string     `json:"path_ref_log"`
	ScriptID    int        `json:"scriptID"`
	Script      Script     `json:"script" gorm:"ForeignKey:ScriptID;AssociationForeignKey:ID"`
	Params      Parameters `json:"parameters" sql:"type:jsonb"`
	Order       string     `json:"-" sql:"-"`
	Uuid        string     `json:"-" sql:"-"`
}

func (test Test) GetOrder() string {
	return test.Order
}

func (test Test) GetID() int {
	return test.TestID
}

// Run method to realise test
func (test Test) Compare(response chan map[string]interface{}) {
	fmt.Println("Running teeeeeest")
	errors := make([]string, 0)
	status := true
	b, err := ioutil.ReadFile(test.PathRefFile)
	if err != nil {
		response <- utils.NewResponse(false, "Error reading log reference file : "+err.Error(), nil)
		return
	}
	refJSon := make([]interface{}, 0)
	config := make([]interface{}, 0)

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
		params1 := utils.CopyArray(config)

		// get keys to compare for current event
		eventParams := jsoncmp.FindParameters(refEvent.(map[string]interface{}), params1)
		if eventParams == nil {
			response <- utils.NewResponse(false, "", nil)
		}
		// compare two events
		if err := jsoncmp.CompareEvent(refEvent.(map[string]interface{}), testJson[i].(map[string]interface{}), eventParams); err != nil {
			fmt.Printf("Error : %s\n", err.Error())
			errors = append(errors, err.Error()+"\n")
			status = false
			continue
		}
		jsonlogresp := refEvent.(map[string]interface{})
		jsonrefresp := testJson[i].(map[string]interface{})
		str1, err := json.Marshal(jsonlogresp)
		if err != nil {
			response <- utils.NewResponse(false, err.Error(), nil)
			return

		}
		str2, err2 := json.Marshal(jsonrefresp)

		if err2 != nil {
			response <- utils.NewResponse(false, err2.Error(), nil)
			return
		}
		resp := make(map[string]interface{})
		resp["event_type"] = EVENT_TEST
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
		resp["body"].(map[string]interface{})[TEST_LOG_EVENT] = string(prettyJSON1.Bytes())
		resp["body"].(map[string]interface{})[REF_LOG_EVENT] = string(prettyJSON2.Bytes())
		response <- utils.NewResponse(true, "", resp)
	}
	resp := make(map[string]interface{})
	resp["event_type"] = RESULT_TEST

	if status {
		resp["body"] = "files match !"
		response <- utils.NewResponse(true, "", resp)
	} else {

		resp["body"] = errors
		response <- utils.NewResponse(false, strings.Join(errors[:], ""), resp)
	}
	return
}

func (test Test) Run(room chan map[string]interface{}) {
	//	history := &TestHistory{TestID: test.TestID, RunUUID: test.Uuid}
	channel := make(chan map[string]interface{})
	output := ""

	//var runner IRunnable = test.Script
	fmt.Printf("on run test id : %d\n", test.TestID)
	// if there is an Script
	if test.ScriptID != 0 {
		test.Script.Params = test.Params
		go test.Script.Run(channel)

		response := <-channel
		if response["status"] != true {
			response["history"] = &TestHistory{TestID: test.TestID, RunUUID: test.Uuid, TestName: test.Name, Success: false, TimeRunned: time.Now().UnixNano()}
			room <- response
			return
		}

		response["response"] = make(map[string]interface{})
		response["response"].(map[string]interface{})["event_type"] = START_TEST
		response["response"].(map[string]interface{})["time_runned"] = time.Now().UnixNano()
		room <- response

		for {
			msg := <-channel
			if msg["status"] != true {
				response["history"] = &TestHistory{TestID: test.TestID, RunUUID: test.Uuid, TestName: test.Name, Success: false, OutputExec: output, OutputTest: msg["message"].(string)}
				room <- msg
				return
			}
			room <- msg
			if response2, ok := msg["response"].(map[string]interface{}); ok {
				if response2["event_type"] == RESULT_SCRIPT {
					break
				}
				output += msg["response"].(map[string]interface{})["body"].(string)
			}
		}
	}
	go test.Compare(channel)
	reflog := ""
	testlog := ""

	for {
		msg := <-channel
		if msg["status"] != true {
			msg["history"] = &TestHistory{TestID: test.TestID, RunUUID: test.Uuid, TestName: test.Name, Success: false,
				OutputExec: output, OutputTest: msg["message"].(string), Reflog: reflog, Testlog: testlog, TimeRunned: time.Now().UnixNano()}
			room <- msg
			break
		}
		response := msg["response"].(map[string]interface{})
		if response["event_type"] == EVENT_TEST {
			reflog += response["body"].(map[string]interface{})[REF_LOG_EVENT].(string)
			testlog += response["body"].(map[string]interface{})[TEST_LOG_EVENT].(string)
		} else if response["event_type"] == RESULT_TEST {
			fmt.Printf("on a fini le test : %d\n", test.TestID)
			response["history"] = &TestHistory{TestID: test.TestID, RunUUID: test.Uuid, TestName: test.Name, Success: true,
				OutputExec: output, OutputTest: response["body"].(string), Reflog: reflog, Testlog: testlog, TimeRunned: time.Now().UnixNano()}
			room <- msg
			break
		}
		room <- msg
	}
}
