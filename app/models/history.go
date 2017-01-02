package models

// TestHistory model for all test history
type TestHistory struct {
	TestHistory int    `json:"testhistory_id"`
	TestID      int    `json:"test_id"`
	OutputExec  string `json:"output_exec"`
	Status      string `json:"status"`
	Success     bool   `json:"sucess"`
	Response    string `json:"response"`
}
