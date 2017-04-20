package controllers

import (
	"fmt"
	"html/template"

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

func initTemplate() {
	revel.TemplateFuncs["set_exec_uuid"] = func(exec *models.Script) template.JS {
		fmt.Printf("uuid : %s\n", exec.Uuid)
		if exec.Uuid == "" {
			exec.Uuid = uuid.NewV4().String()
		}
		fmt.Printf("uuid after : %s\n", exec.Uuid)

		return template.JS("")
	}

	revel.TemplateFuncs["set_test_uuid"] = func(test *models.Test) template.JS {
		test.Uuid = uuid.NewV4().String()
		test.Script.Uuid = test.Uuid
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
	revel.TemplateFuncs["newHistory"] = func(idTest uint) *models.TestHistory {
		if idTest != 0 {
			var test models.Test
			Dbm.First(&test, idTest)
			return &models.TestHistory{TestName: test.Name}
		}
		return &models.TestHistory{}
	}
	revel.TemplateFuncs["newUuid"] = func() string {
		return uuid.NewV4().String()
	}
}

func initScriptDB() {
	if !Dbm.HasTable(&models.Script{}) {
		Dbm.CreateTable(&models.Script{})
		params := models.Parameters{
			models.Parameter{Name: "run_uuid", Value: "", Type: "string"},
			models.Parameter{Name: "device_uuid", Value: "", Type: "string"},
			models.Parameter{Name: "json_file", Value: "", Type: "file"},
			models.Parameter{Name: "jesaispasquoicasert", Value: "", Type: "string"},
			models.Parameter{Name: "obsid", Value: "", Type: "string"},
		}

		execs := []*models.Script{
			&models.Script{Name: "ScenarioEngine", Content: `"C:\\Program Files\\Witbe\\applications\\witbe-scenario-engine\\bin\\ScenarioEngine.exe" -r $run_uuid -d $device_uuid -t $json_file -T $jesaispasquoicasert -o $obsid`, Params: params},
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
			&models.Test{
				Name: "CustomCode",
				PathRefFile: "C:\\json_file\\CustomCode.json",
				PathLogFile: "C:\\ProgramData\\Witbe\\storage\\data\\logs\\witbe-scenario-engine\\80a6f70e-c8d2-420f-7380-a881c48f3355_debug.json",
				ScriptID: 1,
				Config: `[{"ref_fields" : {},"config" : {"body" : {"data" : ["returncode", "status", "pad"]}}}]`, 
				Params: models.Parameters{
					models.Parameter{Name: "run_uuid", Value: "80a6f70e-c8d2-420f-7380-a881c48f3355", Type: "string"},
					models.Parameter{Name: "json_file", Value: "CustomCode", Type: "file"},					
					models.Parameter{Name: "device_uuid", Value: "b142a21e-b7c9-448a-9c57-37cb39d36530", Type: "string"},					
					models.Parameter{Name: "jesaispasquoicasert", Value: "5ce14f17-a44c-4497-bc52-01f01cb13677", Type: "string"},					
					models.Parameter{Name: "obsid", Value: "-2917977", Type: "string"},		
				},
			},
			
			&models.Test{
				Name: "AdministratifBlock",
				PathRefFile: "C:\\json_file\\AdministratifBlock.json",
				PathLogFile: "C:\\ProgramData\\Witbe\\storage\\data\\logs\\witbe-scenario-engine\\ab675c3c-ea5e-46a2-7a98-344881629649_debug.json",
				ScriptID: 1,
				Config: `[{"ref_fields" : {},"config" : {"body" : {"data" : ["returncode", "status", "pad"]}}}]`, 
				Params: models.Parameters{
					models.Parameter{Name: "run_uuid", Value: "ab675c3c-ea5e-46a2-7a98-344881629649", Type: "string"},
					models.Parameter{Name: "json_file", Value: "AdministratifBlock", Type: "file"},					
					models.Parameter{Name: "device_uuid", Value: "b142a21e-b7c9-448a-9c57-37cb39d36530", Type: "string"},					
					models.Parameter{Name: "jesaispasquoicasert", Value: "ef51accd-b992-43a9-bc89-a3bf122a1b23", Type: "string"},					
					models.Parameter{Name: "obsid", Value: "-1430563", Type: "string"},		
				},
			},
			
			
			&models.Test{
				Name: "ContainerDuration",
				PathRefFile: "C:\\json_file\\ContainerDuration.json",
				PathLogFile: "C:\\ProgramData\\Witbe\\storage\\data\\logs\\witbe-scenario-engine\\29aac857-61ad-44fb-6a5f-db6ac1d10b06_debug.json",
				ScriptID: 1,
				Config: `[{"ref_fields" : {},"config" : {"body" : {"data" : ["returncode", "status", "pad"]}}}]`, 
				Params: models.Parameters{
					models.Parameter{Name: "run_uuid", Value: "29aac857-61ad-44fb-6a5f-db6ac1d10b06", Type: "string"},
					models.Parameter{Name: "json_file", Value: "ContainerDuration", Type: "file"},					
					models.Parameter{Name: "device_uuid", Value: "b142a21e-b7c9-448a-9c57-37cb39d36530", Type: "string"},					
					models.Parameter{Name: "jesaispasquoicasert", Value: "7abba34f-c235-424f-9a88-c9812feb0d6", Type: "string"},					
					models.Parameter{Name: "obsid", Value: "-2261861", Type: "string"},		
				},
			},
			
			&models.Test{
				Name: "ContainerInAndOutParams",
				PathRefFile: "C:\\json_file\\ContainerInAndOutParams.json",
				PathLogFile: "C:\\ProgramData\\Witbe\\storage\\data\\logs\\witbe-scenario-engine\\30110e84-e876-4459-7ced-c065abd90389_debug.json",
				ScriptID: 1,
				Config: `[{"ref_fields" : {},"config" : {"body" : {"data" : ["returncode", "status", "pad"]}}}]`, 
				Params: models.Parameters{
					models.Parameter{Name: "run_uuid", Value: "30110e84-e876-4459-7ced-c065abd90389", Type: "string"},
					models.Parameter{Name: "json_file", Value: "ContainerInAndOutParams", Type: "file"},					
					models.Parameter{Name: "device_uuid", Value: "b142a21e-b7c9-448a-9c57-37cb39d36530", Type: "string"},					
					models.Parameter{Name: "jesaispasquoicasert", Value: "0f988986-758f-4dec-8f6b-2fc2a97cf56c", Type: "string"},					
					models.Parameter{Name: "obsid", Value: "-2001385", Type: "string"},		
				},
			},
			
			&models.Test{
				Name: "FunctionParams",
				PathRefFile: "C:\\json_file\\FunctionParams.json",
				PathLogFile: "C:\\ProgramData\\Witbe\\storage\\data\\logs\\witbe-scenario-engine\\a2fe0eba-575f-4f5e-69e8-c75ac6e6ea34_debug.json",
				ScriptID: 1,
				Config: `[{"ref_fields" : {},"config" : {"body" : {"data" : ["returncode", "status", "pad"]}}}]`, 
				Params: models.Parameters{
					models.Parameter{Name: "run_uuid", Value: "a2fe0eba-575f-4f5e-69e8-c75ac6e6ea34", Type: "string"},
					models.Parameter{Name: "json_file", Value: "FunctionParams", Type: "file"},					
					models.Parameter{Name: "device_uuid", Value: "b142a21e-b7c9-448a-9c57-37cb39d36530", Type: "string"},					
					models.Parameter{Name: "jesaispasquoicasert", Value: "c6ba999d-18b4-4a1e-a7b2-27fc888ef451", Type: "string"},					
					models.Parameter{Name: "obsid", Value: "-2126025", Type: "string"},		
				},
			},
			
			&models.Test{
				Name: "LoopV2",
				PathRefFile: "C:\\json_file\\LoopV2.json",
				PathLogFile: "C:\\ProgramData\\Witbe\\storage\\data\\logs\\witbe-scenario-engine\\0a3ae90c-daed-425d-7b9d-b22d880b3554_debug.json",
				ScriptID: 1,
				Config: `[{"ref_fields" : {},"config" : {"body" : {"data" : ["returncode", "status", "pad"]}}}]`, 
				Params: models.Parameters{
					models.Parameter{Name: "run_uuid", Value: "0a3ae90c-daed-425d-7b9d-b22d880b3554", Type: "string"},
					models.Parameter{Name: "json_file", Value: "LoopV2", Type: "file"},					
					models.Parameter{Name: "device_uuid", Value: "b142a21e-b7c9-448a-9c57-37cb39d36530", Type: "string"},					
					models.Parameter{Name: "jesaispasquoicasert", Value: "4f901d78-0272-4080-86c6-09f907db0088", Type: "string"},					
					models.Parameter{Name: "obsid", Value: "-1861469", Type: "string"},		
				},
			},
			
			&models.Test{
				Name: "ReturnCode",
				PathRefFile: "C:\\json_file\\ReturnCode.json",
				PathLogFile: "C:\\ProgramData\\Witbe\\storage\\data\\logs\\witbe-scenario-engine\\6a21f539-4b7d-4d08-4835-f39f4b18498e_debug.json",
				ScriptID: 1,
				Config: `[{"ref_fields" : {},"config" : {"body" : {"data" : ["returncode", "status", "pad"]}}}]`, 
				Params: models.Parameters{
					models.Parameter{Name: "run_uuid", Value: "6a21f539-4b7d-4d08-4835-f39f4b18498e", Type: "string"},
					models.Parameter{Name: "json_file", Value: "ReturnCode", Type: "file"},					
					models.Parameter{Name: "device_uuid", Value: "b142a21e-b7c9-448a-9c57-37cb39d36530", Type: "string"},					
					models.Parameter{Name: "jesaispasquoicasert", Value: "fe92a841-44bb-4dd2-9102-8e7cb3aaeefe", Type: "string"},					
					models.Parameter{Name: "obsid", Value: "-2626999", Type: "string"},		
				},
			},
			
			&models.Test{
				Name: "Timeout",
				PathRefFile: "C:\\json_file\\Timeout.json",
				PathLogFile: "C:\\ProgramData\\Witbe\\storage\\data\\logs\\witbe-scenario-engine\\7e495da1-5a71-43b4-5116-a00b3dc52cf4_debug.json",
				ScriptID: 1,
				Config: `[{"ref_fields" : {},"config" : {"body" : {"data" : ["returncode", "status", "pad"]}}}]`, 
				Params: models.Parameters{
					models.Parameter{Name: "run_uuid", Value: "7e495da1-5a71-43b4-5116-a00b3dc52cf4", Type: "string"},
					models.Parameter{Name: "json_file", Value: "Timeout", Type: "file"},					
					models.Parameter{Name: "device_uuid", Value: "b142a21e-b7c9-448a-9c57-37cb39d36530", Type: "string"},					
					models.Parameter{Name: "jesaispasquoicasert", Value: "a89dec01-0464-4b66-9394-7b882c670390", Type: "string"},					
					models.Parameter{Name: "obsid", Value: "-1710827", Type: "string"},		
				},
			},
			
			
			
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
		orders := models.LibraryOrders{
			models.Order{IdTest: 1, Order: 1},
			models.Order{IdTest: 2, Order: 1},
			models.Order{IdTest: 3, Order: 1},
			models.Order{IdTest: 4, Order: 1},
			models.Order{IdTest: 5, Order: 1},
			models.Order{IdTest: 6, Order: 1},
			models.Order{IdTest: 7, Order: 1},
			models.Order{IdTest: 8, Order: 1},
		}

		Dbm.Preload("Script").Find(&tests)

		lib := &models.Library{Name: "First lib", Tests: tests, Orders: orders}
		Dbm.Create(lib)
	}

	if !Dbm.HasTable(&models.LibraryHistory{}) {
		Dbm.CreateTable(&models.LibraryHistory{})
	}
}

func InitDB() {
	json_writer.Init()

	db, err := gorm.Open("postgres", "user=mistrale dbname=jsonception sslmode=disable password=admin77")
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
	initScriptDB()
	initTestDB()
	initLibraryDB()
	var exec_counts []models.Script
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
