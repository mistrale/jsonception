package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/mistrale/jsonception/app/models"
	"github.com/mistrale/jsonception/app/socket"
	"github.com/mistrale/jsonception/app/utils"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// Libraries controller
type Libraries struct {
	GorpController
}

// Create method to add new execution in DB
func (c Libraries) Create() revel.Result {
	lib := &models.Library{}
	content, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(content, lib); err != nil {
		fmt.Printf("err : %s\n", err.Error())
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	if lib.Name == "" {
		return c.RenderJson(utils.NewResponse(false, "Lib name cannot be empty.", nil))
	}
	for _, v := range lib.TestIDs {
		test := &models.Test{}
		c.Txn.First(test, v)
		lib.Tests = append(lib.Tests, *test)
	}
	//fmt.Printf("name : %s\tids : %d\n", lib.Name, lib.TestIDs[0])
	c.Txn.Create(lib)
	// 	return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	// }
	return c.RenderJson(utils.NewResponse(true, "Successful lib creation", *lib))
}

func (c Libraries) Run(libID int) revel.Result {
	lib := &models.Library{}
	lib_uuid := uuid.NewV4()
	lib_room := socket.CreateRoom(lib_uuid.String())
	end := make(chan bool, len(lib.Tests))

	history := &models.LibraryHistory{LibID: libID, TimeRunned: time.Now().UnixNano(), Histories: make([]models.TestHistory, 0),
		RunUUID: lib_uuid.String()}

	c.Txn.Preload("Tests.Execution").Preload("Tests").First(&lib, libID)
	for i, v := range lib.Tests {
		channel := make(chan map[string]interface{})
		test_uuid := uuid.NewV4()
		//	room := socket.CreateRoom(test_uuid.String())

		go func(test models.Test, ite int) {
			RunTest(&test, test_uuid.String(), channel, c.Txn)

			fmt.Printf("test id IN GO : %d\n", test.GetID())
			for {
				msg := <-channel
				if response, ok := msg["response"].(map[string]interface{}); ok {
					if response["type"] == models.RESULTEVENT {
						var hist models.TestHistory

						Dbm.Where("run_uuid = ?", test_uuid.String()).First(&hist)
						history.Histories = append(history.Histories, hist)
						end <- true
						return
					}
				}

				libmsg := make(map[string]interface{})
				libmsg["test_id"] = test.GetID()
				libmsg["event"] = msg
				lib_room.Chan <- libmsg
				if msg["status"] != true {
					fmt.Printf("ERROR : %s\n", msg["message"])
					return
				}
			}
		}(v, i)
	}
	go func(end chan bool) {
		nb_test := 0
		for {
			<-end
			nb_test++
			if nb_test == len(lib.Tests) {
				Dbm.Create(history)
				return
			}
		}
	}(end)
	return c.RenderJson(utils.NewResponse(true, "", lib_uuid.String()))
}

func (c Libraries) GetHistory(libID int) revel.Result {
	var history []models.LibraryHistory
	c.Txn.Preload("Histories").Where("lib_id = ?", libID).Find(&history)
	return c.RenderJson(history)
}

// Get method to get all library in json
func (c Libraries) GetOne(idLib int) revel.Result {
	lib := &models.Library{}
	c.Txn.Preload("Tests.Execution").Preload("Tests").First(&lib, idLib)
	return c.RenderJson(lib)
}

// Get method to get all library in json
func (c Libraries) Get() revel.Result {
	var libs []models.Library
	c.Txn.Preload("Tests.Execution").Preload("Tests").Find(&libs)
	return c.RenderJson(libs)
}

// Index method to page from library index
func (c Libraries) Index() revel.Result {
	library := &models.Library{LibraryID: 0}
	return c.Render(library)
}

// All method to get all library in index
func (c Libraries) All() revel.Result {
	var libs []models.Library
	c.Txn.Preload("Tests.Execution").Preload("Tests").Find(&libs)
	//spew.Dump(libs)

	return c.Render(libs)
}
