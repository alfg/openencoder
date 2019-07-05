package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/alfg/enc/api/data"
	"github.com/alfg/enc/api/net"
	"github.com/alfg/enc/api/types"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/rs/xid"
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

	// Create Job and push the work to work queue.
	job := types.Job{
		GUID:        xid.New().String(),
		Profile:     json.Profile,
		Source:      json.Source,
		Destination: json.Destination,
	}

	// Send to work queue.
	_, err := enqueuer.Enqueue("encode", work.Q{
		"guid":        job.GUID,
		"profile":     job.Profile,
		"source":      job.Source,
		"destination": job.Destination,
	})
	if err != nil {
		log.Fatal(err)
	}

	created := data.CreateJob(job)
	fmt.Println(created)

	// Create response.
	resp := response{
		Message: "Job created",
		Status:  200,
	}
	c.JSON(http.StatusCreated, resp)
}

func jobsHandler(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	count := c.DefaultQuery("count", "10")
	pageInt, _ := strconv.Atoi(page)
	countInt, _ := strconv.Atoi(count)

	if page == "0" {
		pageInt = 1
	}

	var wg sync.WaitGroup
	var jobs *[]types.Job
	var jobsCount int

	wg.Add(1)
	go func() {
		jobs = data.GetJobs((pageInt-1)*countInt, countInt)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		jobsCount = data.GetJobsCount()
		wg.Done()
	}()
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{
		"count": jobsCount,
		"items": jobs,
	})
}

func workerQueuesHandler(c *gin.Context) {
	client := work.NewClient("enc", redisPool)

	queues, err := client.Queues()
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, queues)
}

func workerPoolsHandler(c *gin.Context) {
	client := work.NewClient("enc", redisPool)

	resp, err := client.WorkerPoolHeartbeats()
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, resp)
}

func workerBusyHandler(c *gin.Context) {
	client := work.NewClient("enc", redisPool)

	observations, err := client.WorkerObservations()
	if err != nil {
		fmt.Println(err)
	}

	var busyObservations []*work.WorkerObservation
	for _, ob := range observations {
		if ob.IsBusy {
			busyObservations = append(busyObservations, ob)
		}
	}
	c.JSON(200, busyObservations)
}

type S3ListResponse struct {
	Folders []string `json:"folders"`
	Files   []file   `json:"files"`
}

type file struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func s3ListHandler(c *gin.Context) {
	prefix := c.DefaultQuery("prefix", "")
	files, err := net.S3ListFiles(prefix)
	if err != nil {
		log.Fatal(err)
	}

	resp := S3ListResponse{}

	// var prefixes &[]S3ListResponse.Folders
	for _, item := range files.CommonPrefixes {
		resp.Folders = append(resp.Folders, *item.Prefix)
	}

	for _, item := range files.Contents {
		var obj file
		obj.Name = *item.Key
		obj.Size = *item.Size
		resp.Files = append(resp.Files, obj)
	}

	c.JSON(200, gin.H{
		"data": resp,
	})
}
