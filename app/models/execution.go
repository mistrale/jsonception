package models

import "github.com/revel/revel"

// Execution : script runned
type Execution struct {
	ExecutionID int    `json:"executionID"`
	Name        string `json:"name"`
	Script      string `json:"script"`
}

// func Run(script string, response chan map[string]interface{}) {
// 	uuid := uuid.NewV4()
// 	var file string
//
// 	if runtime.GOOS == "windows" {
// 		file = "/" + uuid.String() + ".sh"
// 	} else {
// 		file = "/tmp/bash_" + uuid.String() + ".sh"
// 	}
// 	fmt.Printf("script : %s\tfile : %s\n", script, file)
//
// 	if err := ioutil.WriteFile(file, []byte(script), 0777); err != nil {
// 		response <- utils.NewResponse(false, err.Error(), nil)
// 	}
// 	//defer os.Remove(file)
//
// 	go func() {
// 		var cmd *exec.Cmd
// 		cmd = exec.Command("bash", "-c", script)
//
// 		ch := make(chan string)
// 		out := newStream(ch)
// 		cmd.Stdout = out
// 		cmd.Stderr = out
// 		if err := cmd.Start(); err != nil {
// 			response <- utils.NewResponse(false, err.Error(), nil)
// 			//log.Fatal(err)
// 		}
// 		response <- utils.NewResponse(true, "", uuid.String())
// 		go func(ch chan string) {
// 			for {
// 				msg := <-ch
// 				response <- utils.NewResponse(true, "", msg)
// 				fmt.Printf("on push dansle chan : %s\n", msg)
// 				//room.Chan <- msg
// 			}
// 		}(ch)
// 		cmd.Wait()
// 		response <- utils.NewResponse(true, "", "end_"+uuid.String())
// 	}()
// }

// Validate Execution struct field for DB
func (exec *Execution) Validate(v *revel.Validation) {
	v.Required(exec.Name)
	v.Required(exec.Script)
}
