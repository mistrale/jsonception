package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/mistrale/jsonception/app/dispatcher"
)

// Library : test container
type Library struct {
	gorm.Model `json:"library_id"`
	Name       string        `json:"name"`
	TestIDs    []uint        `json:"test_ids" sql:"-"`
	Tests      []Test        `json:"tests" gorm:"many2many:library_tests;"`
	Uuid       string        `json:"-" db:"-"`
	Orders     LibraryOrders `json:"test_orders" sql:"type:jsonb"`
}

// type Order to know when to run test
type Order struct {
	IdTest uint `json:"id_test"`
	Order  int  `json:"order"`
}

type LibraryOrders []Order

func (slice LibraryOrders) Value() (driver.Value, error) {
	valueString, err := json.Marshal(slice)
	return string(valueString), err
}

func (slice *LibraryOrders) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &slice); err != nil {
		return err
	}
	return nil
}

func (slice LibraryOrders) Len() int {
	return len(slice)
}

func (slice LibraryOrders) Less(i, j int) bool {
	return slice[i].Order < slice[j].Order
}

func (slice LibraryOrders) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (lib Library) findTest(order Order) *Test {
	for index, v := range lib.Tests {
		if v.ID == order.IdTest {
			lib.Tests = append(lib.Tests[:index], lib.Tests[index+1:]...)
			return &v
		}
	}
	return nil
}

func (lib Library) CheckOrder() error {
	for _, v := range lib.Orders {
		if v.Order > len(lib.Orders) {
			return errors.New(fmt.Sprintf("Error creating or updating library : Order %d is bigger than test number :%d\n", v.Order, len(lib.Orders)))
		}
	}
	return nil
}

func (lib Library) dealTestScript(test *Test, channel chan dispatcher.Event,
	end chan int, history *LibraryHistory, lib_room chan dispatcher.Event) {
	fmt.Printf("test id IN GO : %s\n", test.GetOrder())
	for {
		msg := <-channel
		body := make(map[string]interface{})
		body["test_body"] = msg.Body
		body["test_id"] = test.ID
		msg.Body = body
		fmt.Printf("msg : %s\n", msg.Body)
		lib_room <- msg
		if msg.Type == RESULT_TEST {
			hist := test.History
			history.Histories = append(history.Histories, hist)
			log.Printf("on rnetre ici %d\n", hist.TestID)
			if msg.Status == true {
				end <- 1
			} else {
				end <- 0
			}
			return
		}
		if msg.Status != true {
			hist := test.History
			history.Histories = append(history.Histories, hist)
			end <- 0
			fmt.Printf("ERROR : %s\n", msg.Errors)
			return
		}
	}
}

func (lib Library) Run(testsOrders map[int]int, end chan int, history *LibraryHistory,
	channel chan dispatcher.Event) {
	for _, o := range lib.Orders {
		test := lib.findTest(o)
		testsOrders[o.Order]++

		// if test needs to be runned in parallele
		if testsOrders[o.Order] > 1 {
			test.Order = "lib_" + strconv.Itoa(int(lib.ID)) + "_" + strconv.Itoa(testsOrders[o.Order])
		} else {
			test.Order = "lib_" + strconv.Itoa(int(lib.ID))
		}

		fmt.Printf("Order : %s\tfor test id :%d\tand size order : %d\n", test.Order, test.ID, len(lib.Orders))

		var runner dispatcher.IRunnable = test
		request := dispatcher.WorkRequest{Runner: &runner, Response: make(chan dispatcher.Event)}
		dispatcher.WorkQueue <- request

		go lib.dealTestScript(test, request.Response, end, history, channel)
	}
}
