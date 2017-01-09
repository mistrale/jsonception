package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mistrale/jsonception/app/models"
	"github.com/mistrale/jsonception/app/utils"
	"github.com/revel/revel"
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

// Get method to get all library in json
func (c Libraries) Get() revel.Result {
	var library []models.Library
	c.Txn.Preload("Tests").Find(&library)
	return c.RenderJson(library)
}

// Index method to page from library index
func (c Libraries) Index() revel.Result {
	library := &models.Library{LibraryID: 0}
	return c.Render(library)
}

// All method to get all library in index
func (c Libraries) All() revel.Result {
	var libs []models.Library
	c.Txn.Preload("Tests").Find(&libs)
	return c.Render(libs)
}
