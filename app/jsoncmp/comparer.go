package jsoncmp

import (
	"fmt"
	"reflect"

	"github.com/mistrale/jsonception/app/utils"
)

func FindParameters(event, params map[string]interface{}) map[string]interface{} {
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

func CompareEvent(refEvent, testEvent map[string]interface{}, params interface{}) error {
	// if params is a string
	if sParams, ok := params.(string); ok {
		eq := reflect.DeepEqual(refEvent[sParams], testEvent[sParams])
		if !eq {
			return fmt.Errorf("Fields differ : Reference file has value = '%s' and test file has value = '%s'", refEvent[sParams], testEvent[sParams])
		}
		return nil
	} else /* else if params is an object */ if oParams, ok := params.(map[string]interface{}); ok {
		for k, v := range oParams {
			if err := CompareEvent(refEvent[k].(map[string]interface{}), testEvent[k].(map[string]interface{}), v); err != nil {
				return fmt.Errorf("[%s] %s", k, err.Error())
			}
		}
	} else /* else if params is an array */ if aParams, ok := params.([]interface{}); ok {
		for _, v := range aParams {
			if err := CompareEvent(refEvent, testEvent, v); err != nil {
				if _, vOk := v.(map[string]interface{}); vOk {
					return fmt.Errorf("[%s] ", err.Error())
				}
				return fmt.Errorf("[%s] %s", v, err.Error())
			}
		}
	}
	return nil
}
