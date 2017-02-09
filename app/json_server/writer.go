package json_writer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/satori/go.uuid"
)

type FileInfo struct {
	name     string
	file     *os.File
	nb_event int
}

// type Event struct {
// 	Id             int           `json:"id"`
// 	Event_type     string        `json:"event_type"`
// 	Str_logs       string        `json:"str_logs"`
// 	Formatted_logs string        `json:"formatted_logs"`
// 	Params         []interface{} `json:"params"`
// }

type Writer struct {
	files map[string]*FileInfo
}

type responseWriter struct {
	Status   bool
	Message  string
	Response interface{}
}

const (
	DEBUG   = "debug"
	INFO    = "info"
	SUCCESS = "success"
	WARNING = "warning"
	DANGER  = "danger"
)

func generateResponse(status bool, message string, response interface{}) *responseWriter {
	return &responseWriter{
		Status:   status,
		Message:  message,
		Response: response,
	}
}

// Printf method
func (w *Writer) Write(token, params string) *responseWriter {
	if v, ok := w.files[token]; ok {
		fmt.Printf("Event : %s\n", params)
		event := make(map[string]interface{})
		//	event := Event{Id: v.nb_event, Event_type: eventType, Str_logs: format, Formatted_logs: fmt.Sprintf(format, a...), Params: a}
		w.files[token].nb_event++
		j, err := json.MarshalIndent(event, "", "  ")
		if err != nil {
			fmt.Printf("Error marshal indent : %s\n", err.Error())
			return generateResponse(false, err.Error(), nil)
		}
		v.file.Write(j)
	} else {
		fmt.Printf("Token not known : %s\n", token)
		return generateResponse(false, "Token not known "+token, nil)
	}
	return generateResponse(true, "", "ok")
}

func (w *Writer) CreateFile(fileName string) *responseWriter {
	if f, err := os.Create(fileName); err == nil {
		token := uuid.NewV4().String()
		w.files[token] = &FileInfo{name: fileName, file: f, nb_event: 0}
		return generateResponse(true, "", token)
	} else {
		fmt.Printf("Error for create file  :%s\n", err.Error())
		return generateResponse(false, err.Error(), nil)
	}
}

func (w *Writer) CloseFile(token string) *responseWriter {
	w.files[token].file.Close()
	delete(w.files, token)
	return generateResponse(true, "", nil)
}
