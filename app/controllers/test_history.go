package controllers

import (
	"log"
	"time"

	"github.com/mistrale/jsonception/app/models"
	"github.com/revel/revel"
)

// TestHistory controller
type TestHistory struct {
	GorpController
}

// Get method to get all history
func (c TestHistory) Get() revel.Result {
	var hist []models.TestHistory
	// _, err := c.Txn.Select(&hist,
	// 	`select * from TestHistory`)
	// if err != nil {
	// 	panic(err)
	// }
	c.Txn.Find(&hist)
	return c.RenderJson(hist)
}

func updateEndHistory(history *models.TestHistory, outputexec, reflog, testlog, outputTest string, success bool) {
	log.Printf("unix nano : %d\n", time.Now().UnixNano())
	history.OutputExec = outputexec
	history.Reflog = reflog
	history.Testlog = testlog
	history.OutputTest = outputTest
	history.Success = success
	history.Status = "finished"
	history.TimeRunned = time.Now().UnixNano()
	Dbm.Save(history)
	// _, err := Dbm.Exec(`update TestHistory set OutputExec = ?,
	// 																							RefLog = ?,
	// 																							TestLog = ?,
	// 																							OutputTest = ?,
	// 																							Success = ?,
	// 																							Status = ?,
	//                                               TimeRunned = ? where ID = ?`,
	// 	outputexec, reflog, testlog, outputTest, success, "finished", time.Now().UnixNano(), history.ID)
	// if err != nil {
	// 	fmt.Printf("History id : %d\n", history.ID)
	// 	log.Printf("error : %s\n", err.Error())
	//
	// 	panic(err)
	// }
}
