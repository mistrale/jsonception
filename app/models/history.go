package models

// TestHistory model for all test history
type TestHistory struct {
	ID         int    `json:"id" gorm:"primary_key"`
	TestID     int    `json:"test_id"`
	OutputExec string `json:"output_exec"`
	Reflog     string `json:"ref_log"`
	Testlog    string `json:"test_log"`
	Status     string `json:"status"`
	Success    bool   `json:"sucess"`
	OutputTest string `json:"output_test"`
	RunUUID    string `json:"run_uuid"`
	Uuid       string `json:"-" db:"-"`
	TimeRunned int64  `json:"time_runned"`
}
