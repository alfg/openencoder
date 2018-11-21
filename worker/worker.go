package worker

import (
	"time"

	"github.com/alfg/enc/types"
	log "github.com/sirupsen/logrus"
)

// Jobs contains a queue of jobs.
var Jobs chan types.Job

// NewWorker creates a new worker instance to listen and process jobs in the queue.
func NewWorker(maxWorkers int, maxQueueSize int) {
	// create job channel
	Jobs = make(chan types.Job, maxQueueSize)

	// create workers
	for i := 1; i <= maxWorkers; i++ {
		go func(i int) {
			for j := range Jobs {
				startJob(i, j)
			}
		}(i)
	}
}

func startJob(id int, j types.Job) {
	log.Infof("worker%d: started %s, delay at %f seconds\n", id, j.Task, j.Delay.Seconds())
	time.Sleep(j.Delay)

	// runWorkflow(j)
	runEncodeJob(j)
	log.Infof("worker%d: completed %s!\n", id, j.Task)
}
