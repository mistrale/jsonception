package controllers

import (
	"fmt"
	"html/template"
	"net"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"

	_ "github.com/lib/pq"

	"github.com/mistrale/jsonception/app/dispatcher"
	"github.com/mistrale/jsonception/app/json_writer"
	"github.com/mistrale/jsonception/app/models"

	r "github.com/revel/revel"
)

var (
	Dbm *gorm.DB
)

func InitDB() {
	json_writer.Init()

	db, err := gorm.Open("sqlite3", "/tmp/post_db.bin")

	//db, err := gorm.Open("postgres", "user=Sikorav dbname=testata sslmode=disable")
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
	Dbm.CreateTable(&models.LibraryHistory{})

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
		fmt.Printf("uuid : %s\n", exec.Uuid)
		if exec.Uuid == "" {
			exec.Uuid = uuid.NewV4().String()
		}
		fmt.Printf("uuid after : %s\n", exec.Uuid)

		return template.JS("")
	}

	revel.TemplateFuncs["set_test_uuid"] = func(test *models.Test) template.JS {
		test.Uuid = uuid.NewV4().String()
		test.Execution.Uuid = test.Uuid
		return template.JS("")
	}

	revel.TemplateFuncs["set_test_history_uuid"] = func(test *models.TestHistory) template.JS {
		if test.Uuid == "" {
			test.Uuid = uuid.NewV4().String()
		}
		return template.JS("")
	}

	revel.TemplateFuncs["set_library_uuid"] = func(lib *models.Library) template.JS {
		lib.Uuid = uuid.NewV4().String()
		return template.JS("")
	}
	revel.TemplateFuncs["newHistory"] = func(idTest int) *models.TestHistory {
		if idTest != 0 {
			var test models.Test
			Dbm.First(&test, idTest)
			return &models.TestHistory{ID: -1, TestName: test.Name}
		}
		return &models.TestHistory{ID: -1}
	}
	revel.TemplateFuncs["newUuid"] = func() string {
		return uuid.NewV4().String()
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	var ip string
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				os.Stdout.WriteString("IP NET : " + ipnet.IP.String() + "\n")
				ip = ipnet.IP.String()
			}
		}
	}
	revel.TemplateFuncs["getIP"] = func() string {
		return ip
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
