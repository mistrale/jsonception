package dispatcher

import (
	"fmt"
	"log"
	"time"
)

type Event struct {
	Status bool        `json:"status"`
	Type   string      `json:"event_type"`
	Errors []string    `json:"errors"`
	Body   interface{} `json:"body"`
}

type IRunnable interface {
	GetOrder() string
	GetID() uint
	Run(chan Event)
}

type WorkRequest struct {
	Runner   *IRunnable
	Delay    time.Duration
	Response chan Event
}

//var resources map[string]*sync.Mutex

var workerQueue map[string]chan chan WorkRequest

// working stuff
var WorkQueue = make(chan WorkRequest)

func StartDispatcher(nb_workers int) {
	workerQueue = make(map[string]chan chan WorkRequest)
	log.Printf("size de worker : %d\n", nb_workers)

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				id := (*work.Runner).GetOrder()
				fmt.Printf("receiv word :%s\nand runner id : %d\n", id, (*work.Runner).GetID())

				_, ok := workerQueue[id]
				if !ok {
					workerQueue[id] = make(chan chan WorkRequest)
					worker := NewWorker(workerQueue[id])
					go worker.Start()
				}
				fmt.Println("On recupere un worker")
				//		go func() {
				worker := <-workerQueue[id]
				fmt.Println("On push")

				//workerQueue[id] <- work
				worker <- work
				fmt.Println("On a push")
				//		}()
			}
		}
	}()
}
