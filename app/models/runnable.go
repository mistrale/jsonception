package models

const (
	// REFLOGEVENT event log of reference
	REFLOGEVENT = "ref_log_event"

	// TESTLOGEVENT event log of current test uuid
	TESTLOGEVENT = "test_log_event"

	// RESULTEVENT for result of a test
	RESULTEVENT = "result_event"

	// TESTEVENT for event when test is running
	TESTEVENT = "test_event"

	EXEC_EVENT = "exec_event"

	RESULT_EXEC = "result_exec"

	START_TEST = "start_test"
)

//
// IRunnable interface for all models implementing Run method
//
type IRunnable interface {
	GetOrder() string
	GetID() int
	Run(chan map[string]interface{})
}
