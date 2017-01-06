package controllers

import (
	"database/sql"
	"html/template"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"

	"github.com/mistrale/jsonception/app/dispatcher"
	"github.com/mistrale/jsonception/app/models"

	r "github.com/revel/revel"
)

var (
	Dbm *gorp.DbMap
)

func InitDB() {
	db, err := sql.Open("sqlite3", "/tmp/post_db.bin")
	if err != nil {
		panic(err)
	}
	Dbm = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	Dbm.AddTable(models.Execution{}).SetKeys(true, "ExecutionID")
	Dbm.AddTable(models.Test{}).SetKeys(true, "TestID")
	Dbm.AddTable(models.TestHistory{}).SetKeys(true, "ID")

	//Dbm.AddTable(models.TestHistory{}).SetKeys(true, "TestHistoryID")

	Dbm.TraceOn("[gorp]", r.INFO)
	Dbm.CreateTablesIfNotExists()
	//
	// test := &models.Execution{ExecutionID: 0, Name: "No execution", Script: ""}
	//
	// // execs := []*models.Execution{
	// 	test,
	// 	&models.Execution{Name: "Test 2 de l execution youlo", Script: "tata"},
	// }

	// refs := []*models.Test{
	// 	&models.Test{Name: "CATASDWADAS", Config: "", PathRefFile: "", PathLogFile: "", ExecutionID: 1, Execution: nil},
	// 	&models.Test{Name: "CATASDWADAS", Config: "", PathRefFile: "", PathLogFile: "", ExecutionID: 0, Execution: nil},
	// }
	// fmt.Printf("test")
	// for _, ref := range refs {
	// 	if err := Dbm.Insert(ref); err != nil {
	// 		panic(err)
	// 	}
	// }
	//
	// for _, exec := range execs {
	// 	if err := Dbm.Insert(exec); err != nil {
	// 		panic(err)
	// 	}
	// }

	var exec_counts []models.Execution
	_, err = Dbm.Select(&exec_counts,
		`select * from Execution`)
	if err != nil {
		panic(err)
	}

	dispatcher.StartDispatcher(len(exec_counts))
	revel.TemplateFuncs["set_exec_uuid"] = func(exec *models.Execution) template.JS {
		exec.Uuid = uuid.NewV4().String()
		return template.JS("")
	}

	revel.TemplateFuncs["set_test_uuid"] = func(test *models.Test) template.JS {
		test.Uuid = uuid.NewV4().String()
		return template.JS("")
	}

	revel.TemplateFuncs["set_test_history_uuid"] = func(test *models.TestHistory) template.JS {
		test.Uuid = uuid.NewV4().String()
		return template.JS("")
	}
}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
