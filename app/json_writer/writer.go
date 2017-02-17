package json_writer

import (
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"

	"github.com/satori/go.uuid"
)

type FileInfo struct {
	name     string
	nb_event int
	events   []map[string]interface{}
}

type Writer struct {
	files map[string]*FileInfo
}

type responseWriter struct {
	Status   bool        `json:"status"`
	Message  string      `json:"message"`
	Response interface{} `json:"response"`
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
	if _, ok := w.files[token]; ok {
		fmt.Printf("Event : %s\n", params)
		event := make(map[string]interface{})
		if err := json.Unmarshal([]byte(params), &event); err != nil {
			fmt.Printf("Bad json format : %s\n", params)
			return generateResponse(false, "Bad json format : "+params, nil)
		}
		event["event_id"] = w.files[token].nb_event
		w.files[token].nb_event++
		w.files[token].events = append(w.files[token].events, event)

		b, err := json.Marshal(w.files[token].events)
		if err != nil {
			return generateResponse(false, err.Error(), nil)
		}
		ioutil.WriteFile(w.files[token].name, b, 0644)
		
	} else {
		fmt.Printf("Token not known : %s\n", token)
		return generateResponse(false, "Token not known "+token, nil)
	}
	return generateResponse(true, "", "ok")
}

func (w *Writer) CreateFile(fileName string) *responseWriter {
		token := uuid.NewV4().String()
		w.files[token] = &FileInfo{name: fileName, nb_event: 0}
		return generateResponse(true, "", token)
}

func (w *Writer) CloseFile(token string) *responseWriter {
	delete(w.files, token)
	return generateResponse(true, "", nil)
}
