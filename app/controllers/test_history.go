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
}
