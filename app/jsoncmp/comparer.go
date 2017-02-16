package jsoncmp

import (
	"fmt"
	"reflect"

	"github.com/mistrale/jsonception/app/utils"
)

func isConcerned(event interface{}, obj interface{}) bool {
	eventObj, isEventObj := event.(map[string]interface{})
	if isEventObj == false {
		fmt.Printf("Warning : event is not an object and ref fields ask for an object in params :%s\n ", obj)
		return false
	}
	if newObj, isObj := obj.(map[string]interface{}); isObj == true {
		// if it is an object iterate on all fields
		for newObj_key, newObj_value := range newObj {

			// if iteration is an object call recursively isConcerned() function
			if _, isObj = newObj_value.(map[string]interface{}); isObj == true {
				// check if key in event is also an object

				if !isConcerned(eventObj[newObj_key], newObj_value) {
					return false
				}
			}

			// if iteration is a string then check fields to check wheter it matchs or not
			if newString, isString := newObj_value.(string); isString == true {

				// check if field in event is also a string
				if eventString, eventIsString := eventObj[newObj_key].(string); eventIsString == true {

					if eventString == newString {
						fmt.Printf("Success : event field  match ! event = :%s\t field = %s\n ", eventString, newString)
						return true
					} else /* if event field doesnt match*/ {
						fmt.Printf("Warning : event field doesnt match : event = :%s\t field = %s\n ", eventString, newString)
						return false
					}
				} else /* if event is not a string return false */ {
					fmt.Printf("Warning : event is not an object and ref fields ask for an object in params :%s\n ", newObj_value)
					return false
				}
			}

		}
	} else {
		fmt.Printf("Warning : wrong config ref_fields in map :%s\n ", obj)
		return false
	}
	return true
}

func FindParameters(event map[string]interface{}, params []interface{}) map[string]interface{} {
	newParams := make(map[string]interface{})
	for _, v := range params {
		if isConcerned(event, v.(map[string]interface{})["ref_fields"]) {
			findParams := v.(map[string]interface{})["config"]
			if !utils.MergeMaps(newParams, findParams.(map[string]interface{})) {
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
