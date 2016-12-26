package controllers

import (
	"database/sql"
	"fmt"

	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mistrale/jsonception/app/models"
	db "github.com/revel/modules/db/app"
	r "github.com/revel/revel"
)

var (
	Dbm *gorp.DbMap
)

func InitDB() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.SqliteDialect{}}

	Dbm.AddTable(models.Execution{}).SetKeys(true, "ExecutionID")

	Dbm.TraceOn("[gorp]", r.INFO)
	Dbm.CreateTables()

	execs := []*models.Execution{
		&models.Execution{Name: "Test 1 de l execution youlo", Script: "tata"},
		&models.Execution{Name: "Test 2 de l execution youlo", Script: "tata"},
	}
	fmt.Printf("test")
	for _, exec := range execs {
		if err := Dbm.Insert(exec); err != nil {
			panic(err)
		}
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
