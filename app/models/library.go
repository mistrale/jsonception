package models

import (
	"fmt"
	"log"
	"strconv"

	"github.com/mistrale/jsonception/app/dispatcher"
)

// Library : test container
type Library struct {
	LibraryID   int           `json:"library_id" gorm:"primary_key"`
	Name        string        `json:"name"`
	TestIDs     []int         `json:"test_ids" sql:"-"`
	Tests       []Test        `json:"tests" gorm:"many2many:library_tests;"`
	Uuid        string        `json:"-" db:"-"`
	Orders      LibraryOrders `json:"test_orders" db:"-"`
	OrderString string        `json:"orders"`
}

type LibraryOrders []Order

func (slice LibraryOrders) Len() int {
	return len(slice)
}

func (slice LibraryOrders) Less(i, j int) bool {
	return slice[i].Order < slice[j].Order
}

func (slice LibraryOrders) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// type Order to know when to run test
type Order struct {
	IdTest int `json:"id_test"`
	Order  int `json:"order"`
}

func (lib Library) findTest(order Order) *Test {
	for index, v := range lib.Tests {
		if v.TestID == order.IdTest {
			lib.Tests = append(lib.Tests[:index], lib.Tests[index+1:]...)
			return &v
		}
	}
	return nil
}

func (lib Library) dealTestExecution(test *Test, channel chan map[string]interface{},
	end chan int, history *LibraryHistory, lib_room chan map[string]interface{}) {
	fmt.Printf("test id IN GO : %s\n", test.GetOrder())
	for {
		msg := <-channel
		msg["test_id"] = test.TestID
		lib_room <- msg
		if response, ok := msg["response"].(map[string]interface{}); ok {
			if response["event_type"] == RESULTEVENT {
				hist := response["history"].(*TestHistory)
				history.Histories = append(history.Histories, *hist)
				log.Printf("on rnetre ici %d\n", hist.TestID)
				if msg["status"] == true {
					end <- 1
				} else {
					end <- 0
				}
				return
			}
		}
		if msg["status"] != true {
			hist := msg["history"].(*TestHistory)
			history.Histories = append(history.Histories, *hist)
			end <- 0
			fmt.Printf("ERROR : %s\n", msg["message"])
			return
		}
	}
}

func (lib Library) Run(testsOrders map[int]int, end chan int, history *LibraryHistory,
	channel chan map[string]interface{}) {
	for _, o := range lib.Orders {
		test := lib.findTest(o)
		testsOrders[o.Order]++

		// if test needs to be runned in parallele
		if testsOrders[o.Order] > 1 {
			test.Order = "lib_" + strconv.Itoa(lib.LibraryID) + "_" + strconv.Itoa(testsOrders[o.Order])
		} else {
			test.Order = "lib_" + strconv.Itoa(lib.LibraryID)
		}

		fmt.Printf("Order : %s\tfor test id :%d\tand size order : %d\n", test.Order, test.TestID, len(lib.Orders))

		var runner dispatcher.IRunnable = test
		request := dispatcher.WorkRequest{Runner: &runner, Response: make(chan map[string]interface{})}
		dispatcher.WorkQueue <- request

		go lib.dealTestExecution(test, request.Response, end, history, channel)
	}
}
