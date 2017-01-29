package controllers

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
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
	c.Txn.Find(&hist)
	return c.RenderJson(hist)
}

func updateEndHistory(db *gorm.DB, history *models.TestHistory, outputexec, reflog, testlog, outputTest string, success bool) {
	log.Printf("run uuid  : %s\n", history.RunUUID)
	history.OutputExec = outputexec
	history.Reflog = reflog
	history.Testlog = testlog
	history.OutputTest = outputTest
	history.Success = success
	history.TimeRunned = time.Now().UnixNano()
	Dbm.Create(history)
}
