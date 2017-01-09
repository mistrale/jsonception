package controllers

import (
	"html/template"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"

	_ "github.com/lib/pq"

	"github.com/mistrale/jsonception/app/dispatcher"
	"github.com/mistrale/jsonception/app/models"

	r "github.com/revel/revel"
)

var (
	Dbm *gorm.DB
)

func InitDB() {
	db, err := gorm.Open("sqlite3", "/tmp/post_db.bin")
	if err != nil {
		panic(err)
	}

	// Ping function checks the database connectivity
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}

	Dbm = db

	Dbm.CreateTable(&models.Execution{})
	Dbm.CreateTable(&models.Test{})
	Dbm.CreateTable(&models.Library{})
	Dbm.CreateTable(&models.TestHistory{})

	// 		db, err := sql.Open("sqlite3", "/tmp/post_db.bin")
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// //	Dbm = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	//
	// 	Dbm.AddTable(models.Execution{}).SetKeys(true, "ExecutionID")
	// 	Dbm.AddTable(models.Test{}).SetKeys(true, "TestID")
	// 	Dbm.AddTable(models.TestHistory{}).SetKeys(true, "ID")
	// 	Dbm.AddTable(models.Library{}).SetKeys(true, "LibraryID")
	//
	// 	//Dbm.AddTable(models.TestHistory{}).SetKeys(true, "TestHistoryID")
	//
	// 	Dbm.TraceOn("[gorp]", r.INFO)
	// 	Dbm.CreateTablesIfNotExists()
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
	Dbm.Find(&exec_counts)

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

	revel.TemplateFuncs["set_library_uuid"] = func(lib *models.Library) template.JS {
		lib.Uuid = uuid.NewV4().String()
		return template.JS("")
	}
}

type GorpController struct {
	*r.Controller
	Txn *gorm.DB
}

func (c *GorpController) Begin() r.Result {
	txn := Dbm.Begin()
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	c.Txn.Commit()
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	c.Txn.Rollback()
	c.Txn = nil
	return nil
}
