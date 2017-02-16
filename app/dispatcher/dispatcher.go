package dispatcher

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/mistrale/jsonception/app/models"
)

type WorkRequest struct {
	Runner   *models.IRunnable
	Delay    time.Duration
	Response chan map[string]interface{}
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
				id := reflect.TypeOf(*work.Runner).String() + "_" + strconv.Itoa((*work.Runner).GetID())
				fmt.Printf("receiv word :%s\n", id)

				_, ok := workerQueue[id]
				if !ok {
					workerQueue[id] = make(chan chan WorkRequest)
					worker := NewWorker(workerQueue[id])
					worker.Start()
				}
				fmt.Println("On recupere un worker")
				worker := <-workerQueue[id]
				fmt.Println("On push")

				//workerQueue[id] <- work
				worker <- work
				fmt.Println("On a push")

			}
		}
	}()
}
