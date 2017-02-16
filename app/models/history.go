package models

// TestHistory model for all test history
type TestHistory struct {
	ID         int    `json:"id" gorm:"primary_key"`
	TestID     int    `json:"test_id"`
	TestName   string `json:"test_name"`
	OutputExec string `json:"output_exec"`
	Reflog     string `json:"ref_log"`
	Testlog    string `json:"test_log"`
	Success    bool   `json:"sucess"`
	OutputTest string `json:"output_test"`
	RunUUID    string `json:"run_uuid"`
	Uuid       string `json:"-" sql:"-"`
	TimeRunned int64  `json:"time_runned"`
}

// LibraryHistory model for all lib history
type LibraryHistory struct {
	ID         int           `json:"id" gorm:"primary_key"`
	LibID      int           `json:"library_id"`
	LibName    string        `json:"lib_name"`
	Result     string        `json:"result"`
	Success    bool          `json:"success"`
	RunUUID    string        `json:"run_uuid"`
	TimeRunned int64         `json:"time_runned"`
	Histories  []TestHistory `json:"test_histories" gorm:"many2many:tests_hist;"`
}
