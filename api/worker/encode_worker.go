package worker

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/alfg/enc/api/types"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

// Make a redis pool
var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	},
}

// Make an enqueuer with a particular namespace
var enqueuer = work.NewEnqueuer("enc", redisPool)

type Context struct {
	GUID        string
	Profile     string
	Source      string
	Destination string
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}

func (c *Context) FindJob(job *work.Job, next work.NextMiddlewareFunc) error {
	if _, ok := job.Args["guid"]; ok {
		c.GUID = job.ArgString("guid")
		if err := job.ArgError(); err != nil {
			return err
		}
	}
	if _, ok := job.Args["profile"]; ok {
		c.Profile = job.ArgString("profile")
		if err := job.ArgError(); err != nil {
			return err
		}
	}
	if _, ok := job.Args["source"]; ok {
		c.Source = job.ArgString("source")
		if err := job.ArgError(); err != nil {
			return err
		}
	}
	if _, ok := job.Args["destination"]; ok {
		c.Destination = job.ArgString("destination")
		if err := job.ArgError(); err != nil {
			return err
		}
	}
	return next()
}

func (c *Context) SendJob(job *work.Job) error {
	guid := job.ArgString("guid")
	profile := job.ArgString("profile")
	source := job.ArgString("source")
	destination := job.ArgString("destination")
	startJob(0, types.Job{
		GUID:        guid,
		Profile:     profile,
		Source:      source,
		Destination: destination,
	})
	return nil
}

func (c *Context) Export(job *work.Job) error {
	return nil
}

// NewWorker creates a new worker instance to listen and process jobs in the queue.
func NewWorker(maxWorkers int, maxQueueSize int) {

	// Make a new pool.
	pool := work.NewWorkerPool(Context{}, 10, "enc", redisPool)

	// Add middleware that will be executed for each job
	pool.Middleware((*Context).Log)
	pool.Middleware((*Context).FindJob)

	// Map the name of jobs to handler functions
	pool.Job("encode", (*Context).SendJob)

	// Customize options:
	pool.JobWithOptions("export", work.JobOptions{Priority: 10, MaxFails: 1}, (*Context).Export)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}

func startJob(id int, j types.Job) {
	log.Infof("worker: started %s\n", j.Profile)

	// runWorkflow(j)
	runEncodeJob(j)
	log.Infof("worker%d: completed %s!\n", j.Profile)
}
