package server

import (
	"fmt"
	"net/http"

	"github.com/alfg/enc/api/data"
	"github.com/alfg/enc/api/types"
	"github.com/alfg/enc/api/worker"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

type request struct {
	Profile     string `json:"profile" binding:"required"`
	Source      string `json:"source" binding:"required"`
	Destination string `json:"dest" binding:"required"`
	Delay       string `json:"delay" binding:"required"`
}

type response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type index struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Docs    string `json:"docs"`
	Github  string `json:"github"`
}

// type job struct {
// 	ID int `json:"id"`
// 	Profile string `json:"profile"`
// }

func indexHandler(c *gin.Context) {
	resp := index{
		Name:    "enc",
		Version: "0.0.1",
		Docs:    "http://localhost/",
		Github:  "https://github.com/alfg/enc",
	}
	c.JSON(200, resp)
}

func encodeHandler(c *gin.Context) {

	// Decode json.
	var json request
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Job and push the work to nsq.
	job := types.Job{
		GUID:        xid.New().String(),
		Profile:     json.Profile,
		Source:      json.Source,
		Destination: json.Destination,
	}

	created := data.CreateJob(job)
	fmt.Println(created)

	go func() {
		log.Info("added: ", job.Profile)

		worker.Jobs <- job

		// // Encode message to bytes.
		// buf := new(bytes.Buffer)
		// enc := gob.NewEncoder(buf)
		// enc.Encode(job)

		// // Send to nsq.
		// config := nsq.NewConfig()
		// p, err := nsq.NewProducer("127.0.0.1:4150", config)
		// if err != nil {
		// 	log.Panic(err)
		// }
		// err = p.Publish("encode", buf.Bytes())
		// if err != nil {
		// 	log.Panic(err)
		// }
	}()

	// Create response.
	resp := response{
		Message: "Job created",
		Status:  200,
	}
	c.JSON(http.StatusCreated, resp)
}

func jobsHandler(c *gin.Context) {
	jobs := data.GetJobs()
	c.JSON(http.StatusOK, jobs)
}

type workerResponse struct {
	ID int
}

func workersHandler(c *gin.Context) {
	// config := nsq.NewConfig()
	// p, err := nsq.NewProducer("127.0.0.1:4150", config)

	// q, _ := nsq.NewConsumer("encode", "encode", config)
	// q.ConnectToNSQLookupd("127.0.0.1:4150")
	// q.ConnectToNSQD("127.0.0.1:4150")
	// stats := q.Stats()
	// fmt.Println(stats)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// opts := nsqlookupd.NewOptions()

	// Make my own queue manager?

	resp := workerResponse{
		ID: 1,
	}
	c.JSON(200, resp)
}
