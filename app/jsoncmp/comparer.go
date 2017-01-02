package jsoncmp

import (
	"fmt"
	"reflect"

	"github.com/mistrale/jsonception/app/jsoncmp/utils"
)

func findParameters(event, params map[string]interface{}) map[string]interface{} {
	newParams := params["fields"].(map[string]interface{})
	if typeParams, ok := params["event_type"].(map[string]interface{})[event["body"].(map[string]interface{})["event_type"].(string)]; ok {
		if !utils.MergeMaps(newParams, typeParams.(map[string]interface{})) {
			return nil
		}
	}
	if id, ok := event["body"].(map[string]interface{})["data"].(map[string]interface{})["block_uuid"].(string); ok {
		if blockParams, ok2 := params["block_uuid"].(map[string]interface{})[id]; ok2 {
			if !utils.MergeMaps(newParams, blockParams.(map[string]interface{})) {
				return nil
			}
		}

	}
	return newParams
}

func compareEvent(refEvent, testEvent map[string]interface{}, params interface{}) error {
	// if params is a string
	if sParams, ok := params.(string); ok {
		eq := reflect.DeepEqual(refEvent[sParams], testEvent[sParams])
		if !eq {
			return fmt.Errorf("Fields differ : Reference file has value = '%s' and test file has value = '%s'", refEvent[sParams], testEvent[sParams])
		}
		return nil
	} else /* else if params is an object */ if oParams, ok := params.(map[string]interface{}); ok {
		for k, v := range oParams {
			if err := compareEvent(refEvent[k].(map[string]interface{}), testEvent[k].(map[string]interface{}), v); err != nil {
				return fmt.Errorf("[%s] %s", k, err.Error())
			}
		}
	} else /* else if params is an array */ if aParams, ok := params.([]interface{}); ok {
		for _, v := range aParams {
			if err := compareEvent(refEvent, testEvent, v); err != nil {
				if _, vOk := v.(map[string]interface{}); vOk {
					return fmt.Errorf("[%s] ", err.Error())
				}
				return fmt.Errorf("[%s] %s", v, err.Error())
			}
		}
	}
	return nil
}

func CompareJSON(testFilePath string, refJson []interface{}, params map[string]interface{}) bool {

	// read json config for scenarios to be tested

	// get json event from two files to differ
	testJson := make([]interface{}, 0)
	if err := utils.GetJsonArray(testFilePath, &testJson); err != nil {
		fmt.Println(err.Error())
		return false
	}

	if len(refJson) != len(testJson) {
		fmt.Println("File diverg from different events number")
		return false
	}

	// iterate through reference file
	for i, refEvent := range refJson {
		params1 := utils.CopyMap(params)

		// get keys to compare for current event
		eventParams := findParameters(refEvent.(map[string]interface{}), params1)
		if eventParams == nil {
			return false
		}
		// compare two events
		if err := compareEvent(refEvent.(map[string]interface{}), testJson[i].(map[string]interface{}), eventParams); err != nil {
			fmt.Printf("Error : %s\n", err.Error())
			return false
		}

	}
	return true
}
