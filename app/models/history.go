package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type Errors []string

func (j Errors) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *Errors) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

// TestHistory model for all test history
type TestHistory struct {
	gorm.Model `json:"id"`
	TestID     uint   `json:"test_id"`
	TestName   string `json:"test_name"`
	OutputExec string `json:"output_exec"`
	Reflog     string `json:"ref_log"`
	Testlog    string `json:"test_log"`
	Success    bool   `json:"sucess"`
	OutputTest string `json:"output_test"`
	RunUUID    string `json:"run_uuid"`
	Errors     Errors `json:"errors" sql:"type:jsonb"`
	Uuid       string `json:"-" sql:"-"`
	TimeRunned int64  `json:"time_runned"`
}

// LibraryHistory model for all lib history
type LibraryHistory struct {
	gorm.Model `json:"id"`
	LibID      uint          `json:"library_id"`
	LibName    string        `json:"lib_name"`
	Result     string        `json:"result"`
	Success    bool          `json:"success"`
	RunUUID    string        `json:"run_uuid"`
	TimeRunned int64         `json:"time_runned"`
	Histories  []TestHistory `json:"test_histories" gorm:"many2many:tests_hist;"`
}
