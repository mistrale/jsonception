package models

const (
	// REFLOGEVENT event log of reference
	REF_LOG_EVENT = "ref_log_event"

	// TESTLOGEVENT event log of current test uuid
	TEST_LOG_EVENT = "test_log_event"

	// RESULTEVENT for result of a test
	RESULT_TEST = "result_test"
	EVENT_TEST  = "event_test"
	START_TEST  = "start_test"

	START_SCRIPT  = "start_script"
	RESULT_SCRIPT = "result_script"
	EVENT_SCRIPT  = "event_script"
)

//
// IRunnable interface for all models implementing Run method
//
type IRunnable interface {
	GetOrder() string
	GetID() int
	Run(chan map[string]interface{})
}
