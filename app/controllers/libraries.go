package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"time"

	"github.com/mistrale/jsonception/app/models"
	"github.com/mistrale/jsonception/app/socket"
	"github.com/mistrale/jsonception/app/utils"
	"github.com/revel/revel"
	"github.com/satori/go.uuid"
)

// Libraries controller
type Libraries struct {
	GorpController
}

// Create method to add new Script in DB
func (c Libraries) Create() revel.Result {
	lib := &models.Library{}
	content, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Printf("content : %s\n", content)
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
	for _, v := range lib.Orders {
		fmt.Printf("test id : %d\torder : %d\n", v.IdTest, v.Order)
	}
	c.Txn.Create(lib)
	return c.RenderJson(utils.NewResponse(true, "Successful lib creation", *lib))
}

func (c Libraries) Delete(id_lib int) revel.Result {
	var lib models.Library
	c.Txn.First(&lib, id_lib)
	c.Txn.Delete(&lib)
	return c.RenderJson(utils.NewResponse(true, "", "Library deleted"))
}

func (c Libraries) Update(id_lib int) revel.Result {
	lib := &models.Library{}
	fmt.Printf("id lib : %d\n", id_lib)
	c.Txn.First(&lib, id_lib)

	content, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(content, lib); err != nil {
		return c.RenderJson(utils.NewResponse(false, err.Error(), nil))
	}
	for _, v := range lib.TestIDs {
		test := &models.Test{}
		c.Txn.First(test, v)
		lib.Tests = append(lib.Tests, *test)
	}
	c.Txn.Save(&lib)
	return c.RenderJson(utils.NewResponse(true, "", "Library updated"))
}

func (c Libraries) initRun(lib_uuid string, lib *models.Library, idLib int) error {
	c.Txn.Preload("Tests.Script").Preload("Tests").First(&lib, idLib)

	// check if all test are present
	tests := make(map[int]int)
	for _, v := range lib.Orders {
		tests[v.IdTest] += 1
	}
	for i, v := range lib.Tests {
		tests[v.TestID] -= 1
		lib.Tests[i].Uuid = uuid.NewV4().String()
		if tests[v.TestID] == 0 {
			delete(tests, v.TestID)
		}
	}
	if len(tests) > 0 {
		return errors.New("Ordering for running lib is wrong")
	}
	sort.Sort(lib.Orders)

	fmt.Println("\nSorted")
	for _, c := range lib.Orders {
		fmt.Printf("test order : %d\tid : %d\n", c.Order, c.IdTest)
	}
	return nil
}

func (c Libraries) Start(lib_uuid string, lib *models.Library) {
	fmt.Println("Starting lib")
	lib_room := socket.CreateRoom(lib_uuid)
	testsOrders := make(map[int]int)
	nb_test := len(lib.Tests)
	end := make(chan int, len(lib.Tests))

	history := &models.LibraryHistory{LibID: lib.LibraryID, TimeRunned: time.Now().UnixNano(), Histories: make([]models.TestHistory, 0),
		RunUUID: lib_uuid}

	go lib.Run(testsOrders, end, history, lib_room.Chan)

	// check end of test
	go func(end chan int) {
		history.Success = true
		for {
			v := <-end
			if v == 0 {
				history.Success = false
			}
			nb_test--
			fmt.Printf("NB TEST : %d\n", nb_test)
			if nb_test == 0 {
				history.LibName = lib.Name
				Dbm.Create(history)

				lib_room.Chan <- utils.NewResponse(true, "", "end_"+lib_uuid)
				return
			}
		}
	}(end)
}

func (c Libraries) Run(idLib int) revel.Result {
	fmt.Printf("idlib : %d\n", idLib)
	lib_uuid := uuid.NewV4()
	lib := &models.Library{}
	if idLib == 0 {
		return c.RenderJson(utils.NewResponse(true, "Empty idLib for run", nil))
	}
	if err := c.initRun(lib_uuid.String(), lib, idLib); err != nil {
		return c.RenderJson(utils.NewResponse(false, "", err.Error()))
	}
	c.Start(lib_uuid.String(), lib)
	return c.RenderJson(utils.NewResponse(true, "", lib_uuid.String()))
}

func (c Libraries) GetHistory(libID int) revel.Result {
	var history []models.LibraryHistory
	c.Txn.Preload("Histories").Where("lib_id = ?", libID).Find(&history)
	return c.RenderJson(history)
}

// DeleteHistory method to delete all history
func (c Libraries) DeleteHistory(id_lib int) revel.Result {
	var history []models.LibraryHistory
	c.Txn.Preload("Histories").Where("lib_id = ?", id_lib).Find(&history)
	for _, v := range history {
		for _, testHist := range v.Histories {
			c.Txn.Delete(&testHist)
		}
	}
	c.Txn.Delete(&history)
	return c.RenderJson(utils.NewResponse(true, "", "Library history deleted"))
}

// GetHistory method to get all history from one librairie and template html
func (c Libraries) GetHistoryTemplate(libID int) revel.Result {
	var history []models.LibraryHistory
	c.Txn.Preload("Histories").Where("lib_id = ?", libID).Find(&history)
	c.Render(history)
	return c.RenderTemplate("LibraryHistory/all.html")
}

// Get method to get all library in json
func (c Libraries) GetOne(libID int) revel.Result {
	lib := &models.Library{}
	c.Txn.Preload("Tests.Script").Preload("Tests").First(&lib, libID)
	return c.RenderJson(lib)
}

// Get method to get all library in json
func (c Libraries) Get() revel.Result {
	var libs []models.Library
	c.Txn.Preload("Tests.Script").Preload("Tests").Find(&libs)
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
	c.Txn.Preload("Tests.Script").Preload("Tests").Find(&libs)
	return c.Render(libs)
}
