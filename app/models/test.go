package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mistrale/jsonception/app/dispatcher"
	"github.com/mistrale/jsonception/app/jsoncmp"
	"github.com/mistrale/jsonception/app/utils"
)

// Test : script and logevent
type Test struct {
	gorm.Model
	Name        string      `json:"name"`
	Config      string      `json:"config"`
	PathRefFile string      `json:"path_test_log"`
	PathLogFile string      `json:"path_ref_log"`
	ScriptID    int         `json:"scriptID"`
	Script      Script      `json:"script" gorm:"ForeignKey:ScriptID;AssociationForeignKey:ID"`
	Params      Parameters  `json:"parameters" sql:"type:jsonb"`
	Order       string      `json:"-" sql:"-"`
	Uuid        string      `json:"-" sql:"-"`
	History     TestHistory `json:"-" sql:"-"`
}

func (test *Test) Print() {
	fmt.Printf("Name : %d\tScriptID : %d\n", test.ID, test.ScriptID)
	for _, v := range test.Params {
		fmt.Printf("Params : name = %s\tvalue = %s\ttype = %s\n", v.Name, v.Value, v.Type)
	}
}

func (test Test) GetOrder() string {
	return test.Order
}

func (test Test) GetID() uint {
	return test.ID
}

// Run method to realise test
func (test *Test) Compare(response chan dispatcher.Event) {
	fmt.Println("Running teeeeeest")
	errors := make([]string, 0)
	status := true
	b, err := ioutil.ReadFile(test.PathRefFile)
	if err != nil {
		response <- dispatcher.Event{Status: false, Errors: []string{"Error reading log reference file : " + err.Error()}, Body: nil, Type: RESULT_TEST}
		return
	}
	var refJSon []interface{}
	var config []interface{}

	if err := json.Unmarshal([]byte(b), &refJSon); err != nil {
		response <- dispatcher.Event{Status: false, Errors: []string{err.Error()}, Type: RESULT_TEST}
		return
	}
	if err := json.Unmarshal([]byte(test.Config), &config); err != nil {
		response <- dispatcher.Event{Status: false, Errors: []string{err.Error()}, Type: RESULT_TEST}
		return
	}
	// get json event from two files to differ
	var testJson []interface{}
	if err := utils.GetJsonArray(test.PathLogFile, &testJson); err != nil {
		fmt.Println(err.Error())
		response <- dispatcher.Event{Status: false, Errors: []string{err.Error()}, Type: RESULT_TEST}
		return
	}

	if len(refJSon) != len(testJson) {
		fmt.Println("File diverg from different events number")
		response <- dispatcher.Event{Status: false, Errors: []string{"File diverg from different events number"}, Type: RESULT_TEST}
		return
	}
	// iterate through reference file
	for i, refEvent := range refJSon {
		params1 := utils.CopyArray(config)

		// get keys to compare for current event
		eventParams := jsoncmp.FindParameters(refEvent.(map[string]interface{}), params1)
		if eventParams == nil {
			continue
		}
		// compare two events
		if err := jsoncmp.CompareEvent(refEvent.(map[string]interface{}), testJson[i].(map[string]interface{}), eventParams); err != nil {
			fmt.Printf("Error : %s\n", err.Error())
			errors = append(errors, "Event id nb "+strconv.Itoa(i)+" : "+err.Error()+"\n")
			status = false
			continue
		}
		jsonlogresp := refEvent.(map[string]interface{})
		jsonrefresp := testJson[i].(map[string]interface{})
		str1, err := json.Marshal(jsonlogresp)
		if err != nil {
			response <- dispatcher.Event{Status: false, Errors: []string{"Event id nb " + strconv.Itoa(i) + " : " + err.Error()}, Type: RESULT_TEST}
			return

		}
		str2, err2 := json.Marshal(jsonrefresp)

		if err2 != nil {
			response <- dispatcher.Event{Status: false, Errors: []string{"Event id nb " + strconv.Itoa(i) + " : " + err.Error()}, Type: RESULT_TEST}
			return
		}

		var prettyJSON1 bytes.Buffer
		err = json.Indent(&prettyJSON1, str1, "", "\t")
		if err != nil {
			response <- dispatcher.Event{Status: false, Errors: []string{"Event id nb " + strconv.Itoa(i) + " : " + err.Error()}, Type: RESULT_TEST}
			return
		}
		var prettyJSON2 bytes.Buffer
		err = json.Indent(&prettyJSON2, str2, "", "\t")
		if err != nil {
			response <- dispatcher.Event{Status: false, Errors: []string{"Event id nb " + strconv.Itoa(i) + " : " + err.Error()}, Type: RESULT_TEST}
			return
		}
		body := make(map[string]interface{})
		body[TEST_LOG_EVENT] = string(prettyJSON1.Bytes())
		body[REF_LOG_EVENT] = string(prettyJSON2.Bytes())
		response <- dispatcher.Event{Status: true, Errors: nil, Type: EVENT_TEST, Body: body}
	}
	resp := make(map[string]interface{})
	resp["event_type"] = RESULT_TEST

	if status {
		response <- dispatcher.Event{Status: true, Errors: nil, Type: RESULT_TEST, Body: "files match !"}
	} else {
		response <- dispatcher.Event{Status: false, Errors: errors, Type: RESULT_TEST, Body: "files doesnt match !"}
	}
	return
}

func (test *Test) Run(room chan dispatcher.Event) {
	//	history := &TestHistory{TestID: test.TestID, RunUUID: test.Uuid}
	channel := make(chan dispatcher.Event)
	output := ""
	room <- dispatcher.Event{Type: START_TEST, Status: true, Errors: nil, Body: time.Now().UnixNano()}
	//var runner IRunnable = test.Script
	fmt.Printf("on run test id : %d\n", test.ID)
	// if there is an Script
	if test.ScriptID != 0 {
		test.Script.Params = test.Params
		go test.Script.Run(channel)

		response := <-channel
		if response.Status != true {
			test.History = TestHistory{TestID: test.ID, RunUUID: test.Uuid, TestName: test.Name, Success: false, TimeRunned: time.Now().UnixNano()}
			room <- response
			return
		}
		room <- response

		for {
			msg := <-channel
			if msg.Status != true {
				test.History = TestHistory{TestID: test.ID, RunUUID: test.Uuid, TestName: test.Name, Success: false, OutputExec: output, OutputTest: msg.Errors[0]}
				room <- msg
				return
			}
			room <- msg
			if msg.Type == RESULT_SCRIPT {
				break
			}
			output += msg.Body.(string)
		}
	}
	go test.Compare(channel)
	reflog := ""
	testlog := ""

	for {
		msg := <-channel
		room <- msg
		if msg.Status != true {
			test.History = TestHistory{TestID: test.ID, RunUUID: test.Uuid, TestName: test.Name, Success: false, Errors: msg.Errors,
				OutputExec: output, OutputTest: "Test failed !", Reflog: reflog, Testlog: testlog, TimeRunned: time.Now().UnixNano()}
			break
		}
		if msg.Type == EVENT_TEST {
			reflog += msg.Body.(map[string]interface{})[REF_LOG_EVENT].(string)
			testlog += msg.Body.(map[string]interface{})[TEST_LOG_EVENT].(string)
		} else if msg.Type == RESULT_TEST {
			fmt.Printf("on a fini le test : %d\n", test.ID)
			test.History = TestHistory{TestID: test.ID, RunUUID: test.Uuid, TestName: test.Name, Success: true,
				OutputExec: output, OutputTest: msg.Body.(string), Reflog: reflog, Testlog: testlog, TimeRunned: time.Now().UnixNano()}
			break
		}
	}
}
