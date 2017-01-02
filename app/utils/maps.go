package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func GetJsonMap(filePath string, instance map[string]interface{}) error {
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		return fmt.Errorf("Error opening file %s : %s", filePath, e.Error())
	}
	if err := json.Unmarshal(file, &instance); err != nil {
		return fmt.Errorf("Error reading json file %s : %s", filePath, err.Error())

	}
	return nil
}

func CopyMap(originalMap map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range originalMap {
		if mObj, ok := v.(map[string]interface{}); ok {
			newMap[k] = CopyMap(mObj)
		} else if aObj, ok := v.([]interface{}); ok {
			newMap[k] = CopyArray(aObj)
		} else {
			newMap[k] = v
		}
	}
	return newMap
}

func MergeMaps(params, newParams map[string]interface{}) bool {
	//fmt.Printf("On rentre dans merParams avec params %s\nnewParams %s\n\n", params, newParams)
	for k, v := range newParams {

		// if keys doesnt exist in current params
		if _, ok := params[k]; !ok {
			params[k] = v
			continue
		}

		// if value is an object iterate throught it
		if maps, ok := v.(map[string]interface{}); ok {
			pMaps, ok := params[k].(map[string]interface{})
			if !ok {
				fmt.Printf("Key '%s' matchs but types differ\n", k)
				return false
			}

			for kMaps, vMaps := range maps {

				// if vMaps is not contained in params[k]
				if _, ok := pMaps[kMaps]; !ok {
					pMaps[kMaps] = vMaps
				} else /* if vMaps is contained in params[k] */ {

					// if vMaps is an object
					if objectMaps, ok := vMaps.(map[string]interface{}); ok {

						MergeMaps(pMaps[kMaps].(map[string]interface{}), objectMaps)
					} else /* if vMaps is an array */ if arrayMaps, ok := vMaps.([]interface{}); ok {
						if pArray, ok := pMaps[kMaps].([]interface{}); !ok {
							fmt.Printf("Key '%s' matchs but types differ for array\n", k)
							return false
						} else {
							pMaps[kMaps] = MergeArrays(pArray, arrayMaps)
						}
					}
				}
			}
		}
		// if value is an array iterate throught it
		if arrayMaps, ok := v.([]interface{}); ok {
			if pArray, ok := params[k].([]interface{}); !ok {
				fmt.Printf("Key '%s' matchs but types differ for array\n", k)
				return false
			} else {
				params[k] = MergeArrays(pArray, arrayMaps)
			}
		}

	}
	return true
}

func MapEquals(obj1, obj2 map[string]interface{}) error {

	return nil
}
