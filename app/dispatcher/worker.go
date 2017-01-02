package dispatcher

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/mistrale/jsonception/app/utils"
)

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

type outstream struct {
	ch chan string
}

func newStream(ch chan string) *outstream {
	return &outstream{ch: ch}
}

func (out outstream) Write(p []byte) (int, error) {
	out.ch <- string(p)
	//	fmt.Println("wtf : " + string(p))
	return len(p), nil
}

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

// Run function for worker doing job
func (w *Worker) Run(Work WorkRequest) {
	fmt.Printf("script : %s\tuuid : %s\n", Work.Script, Work.Uuid)

	//defer os.Remove(file)

	var cmd *exec.Cmd
	cmd = exec.Command("bash", "-c", Work.Script)

	ch := make(chan string)
	out := newStream(ch)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Start(); err != nil {
		Work.Response <- utils.NewResponse(false, err.Error(), nil)
		//log.Fatal(err)
	}
	Work.Response <- utils.NewResponse(true, "ok", Work.Uuid)
	go func(ch chan string) {
		for {
			msg := <-ch
			response := make(map[string]interface{})
			response["type"] = "exec_event"
			response["body"] = msg
			Work.Response <- utils.NewResponse(true, "", response)
			fmt.Printf("on push dansle chan : %s\n", response)
			//room.Chan <- msg
		}
	}(ch)
	cmd.Wait()
	Work.Response <- utils.NewResponse(true, "", "end_"+Work.Uuid)
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
				w.Run(work)
				// Receive a work request.
				fmt.Printf("worker%d: Received work request, delaying for %f seconds\n", w.ID, work.Delay.Seconds())

				time.Sleep(work.Delay)
				fmt.Printf("worker%d: Hello, %s!\n", w.ID, work.Uuid)

			case <-w.QuitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.ID)
				return
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
