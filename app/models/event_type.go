package models

const (
	END_ROOM = "end_room"
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

	START_LIB  = "start_lib"
	EVENT_LIB  = "event_lib"
	RESULT_LIB = "result_lib"
)
