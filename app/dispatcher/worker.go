package dispatcher

import (
	"fmt"
	"log"
)

type Worker struct {
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(workerQueue chan chan WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

// Run function for worker doing job
func (w *Worker) Run(Work WorkRequest) {
	log.Println("run starting")

	(*Work.Runner).Run(Work.Response)

	log.Println("run done")
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.

			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				fmt.Println("on va wrok")
				w.Run(work)
				fmt.Println("on a wrok")

				// Receive a work request.
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
