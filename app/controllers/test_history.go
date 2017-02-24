package controllers

import (
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

// Get method to get one history
func (c TestHistory) GetOne(uuid, testName string) revel.Result {
	hist := &models.TestHistory{}
	hist.Uuid = uuid
	hist.TestName = testName
	c.Render(hist)
	return c.RenderTemplate("TestHistory/one.html")
}

func updateEndHistory(db *gorm.DB, history *models.TestHistory, outputexec, reflog, testlog, outputTest, testName string, success bool) {
	history.OutputExec = outputexec
	history.Reflog = reflog
	history.Testlog = testlog
	history.OutputTest = outputTest
	history.Success = success
	history.TimeRunned = time.Now().UnixNano()
	history.TestName = testName
	Dbm.Create(history)
}
