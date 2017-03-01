package controllers

import (
	"fmt"
	"html/template"
	"net"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

func initTemplate() {
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

func initExecutionDB() {
	if !Dbm.HasTable(&models.Execution{}) {
		Dbm.CreateTable(&models.Execution{})
		execs := []*models.Execution{
 		&models.Execution{Name: "Click_element_set_return", Script: "ScenarioEngine.exe -r click_set_return -d b142a21e-b7c9-448a-9c57-37cb39d36530 -t 5d84fccb-3836-42a8-a0c8-56bc93518cb4 -T 63432494-c134-4ebd-9dc5-de6f6150060f -o -1748335"},
 	 	&models.Execution{Name: "gla func", Script: "ScenarioEngine.exe -r gla_func -d b142a21e-b7c9-448a-9c57-37cb39d36530 -t 03b77937-aacb-4e32-8933-1b44c22e77ea -T a8ac7bc1-561f-4657-a94b-118c498756f3 -o -2608857"},
 		&models.Execution{Name: "test_amazone", Script: "ScenarioEngine.exe -r test_amazon -d b142a21e-b7c9-448a-9c57-37cb39d36530 -t 68ea0181-b048-4375-91bd-cd0cd6d7434c -T b9d19b1e-14df-410f-b97a-f871bf6094d6 -o -1160361"},
 		}
		for _, exec := range execs {
			Dbm.Create(exec)
		}
	}
}

func initTestDB() {
	if !Dbm.HasTable(&models.Test{}) {
		Dbm.CreateTable(&models.Test{})
		tests := []*models.Test{
		&models.Test{Name: "test_click_element_set_return", PathRefFile:"C:\\json_file\\click_set_return_debug.json",
			PathLogFile:"C:\\ProgramData\\Witbe\\storage\\data\\logs\\witbe-scenario-engine\\click_set_return_debug.json", ExecutionID:1,
			Config:`"[{"ref_fields" : {},"config" : {"body" : {"data" : ["returncode", "status", "pad"]}}}]"`},
		&models.Test{Name: "test_gla_func", PathRefFile:"C:\\json_file\\gla_func_debug.json",
			PathLogFile:"C:\\ProgramData\\Witbe\\storage\\data\\logs\\witbe-scenario-engine\\gla_func_debug.json", ExecutionID:2,
			Config:`"[{"ref_fields" : {},"config" : {"body" : {"data" : ["returncode", "status", "pad"]}}}]"`},
		&models.Test{Name: "test_amazon", PathRefFile:"C:\\json_file\test_amazon_debug.json",
			PathLogFile:"C:\\ProgramData\\Witbe\\storage\\data\\logs\\witbe-scenario-engine\test_amazon_debug.json", ExecutionID:3,
			Config:`"[{"ref_fields" : {},"config" : {"body" : {"data" : ["returncode", "status", "pad"]}}}]"`},
		}
		for _, test := range tests {
			Dbm.Create(test)
		}
	}
	if !Dbm.HasTable(&models.TestHistory{}) {
		Dbm.CreateTable(&models.TestHistory{})
	}
}

func initLibraryDB() {
	if !Dbm.HasTable(&models.Library{}) {
		Dbm.CreateTable(&models.Library{})
		var tests []models.Test

		Dbm.Preload("Execution").Find(&tests)

		lib := &models.Library{Name:"First lib", Tests: tests}
		Dbm.Create(lib)
	}

	if !Dbm.HasTable(&models.LibraryHistory{}) {
		Dbm.CreateTable(&models.LibraryHistory{})
	}
}

func InitDB() {
	json_writer.Init()

	//db, err := gorm.Open("postgres", "user=postgres dbname=jsonception sslmode=disable password=witbe77")
	 db, err := gorm.Open("sqlite3", "/tmp/post_db.bin")
	if err != nil {
		panic(err)
	}
	//
	// // Ping function checks the database connectivity
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}
	//
	Dbm = db

	var exec_counts []models.Execution
	Dbm.Find(&exec_counts)

	dispatcher.StartDispatcher(len(exec_counts))
	initTemplate()
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
