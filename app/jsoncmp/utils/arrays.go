package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func GetJsonArray(filePath string, instance *[]interface{}) error {
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		return fmt.Errorf("Error opening file %s : %s", filePath, e.Error())
	}
	if err := json.Unmarshal(file, instance); err != nil {
		return fmt.Errorf("Error reading json file %s : %s", filePath, err.Error())

	}
	return nil
}

func RemoveFromArray(s []interface{}, e interface{}) []interface{} {
	for i, a := range s {
		if a == e {
			s = append(s[:i], s[i+1:]...)
			return s
		}
	}
	return s
}

func CopyArray(originalArray []interface{}) []interface{} {
	newArray := make([]interface{}, len(originalArray))
	for i, v := range originalArray {
		if mValue, ok := v.(map[string]interface{}); ok {
			newArray[i] = CopyMap(mValue)
		} else if aValue, ok := v.([]interface{}); ok {
			newArray[i] = CopyArray(aValue)
		} else {
			newArray[i] = v
		}
	}
	return newArray
}

func ArrayContains(s []interface{}, e interface{}) int {
	for i, a := range s {
		m1, ok := a.(map[string]interface{})
		m2, ok2 := e.(map[string]interface{})

		if ok && ok2 {
			for k, _ := range m1 {
				if _, ok := m2[k]; ok {
					return i
				}
			}
		} else if a == e {
			return i
		}
	}
	return -1
}

func MergeArrays(pArray, arrayMaps []interface{}) []interface{} {
	for _, vArray := range arrayMaps {

		// if vArray is not contained in params[k]
		if pos := ArrayContains(pArray, vArray); pos == -1 {
			/* if vArray is a string starting with '!' */
			if strArray, ok := vArray.(string); ok && strArray[0] == '!' {

				strArray = strArray[1:]

				if index := ArrayContains(pArray, strArray); index >= 0 {
					pArray = RemoveFromArray(pArray, strArray)
				}
			} else /* if vArray is a simple string */ {
				pArray = append(pArray, vArray)
			}
		} else /* else if vArray is contained and is an object */ if vobjectArray, ok := vArray.(map[string]interface{}); ok {

			MergeMaps(pArray[pos].(map[string]interface{}), vobjectArray)
		}
	}
	return pArray
}
