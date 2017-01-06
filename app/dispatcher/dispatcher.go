package dispatcher

import (
	"fmt"
	"log"
	"time"
)

var (
	nb_work = 0
)

type WorkRequest struct {
	Uuid     string
	Script   string
	Delay    time.Duration
	Response chan map[string]interface{}
}

var workerQueue []chan chan WorkRequest

// working stuff
var WorkQueue = make([]chan WorkRequest, 0)

// AddWorker when a new execution is created
func AddWorker() {
	workerQueue = append(workerQueue, make(chan chan WorkRequest))
	WorkQueue = append(WorkQueue, make(chan WorkRequest))

	worker := NewWorker(nb_work+1, workerQueue[len(workerQueue)-1])
	worker.Start()
	startWorkerQueue(len(workerQueue) - 1)
	nb_work++
}

func StartDispatcher(nb_workers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	workerQueue = make([]chan chan WorkRequest, nb_workers)
	nb_work = nb_workers
	// Now, create all of our workers.
	for i := 0; i < nb_workers; i++ {
		workerQueue[i] = make(chan chan WorkRequest)
		WorkQueue = append(WorkQueue, make(chan WorkRequest))

		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, workerQueue[i])
		worker.Start()
		startWorkerQueue(i)
	}
	log.Printf("size de worker : %d\n", nb_workers)
}

func startWorkerQueue(id int) {
	go func() {
		for {
			select {
			case work := <-WorkQueue[id]:
				fmt.Println("Received work requeust")
				go func() {
					worker := <-workerQueue[id]

					fmt.Println("Dispatching work request")
					worker <- work
				}()
			}
		}
	}()
}
