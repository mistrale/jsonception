package main

import (
	"fmt"
	"io/ioutil"
)


func main() {
	if _, err := ioutil.ReadFile("C:\\json_file\\click_set_return_debug.json"); err != nil {
		fmt.Printf("error : %s\n", err.Error())
	}
}
