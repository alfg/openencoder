package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/machine"
	"github.com/alfg/openencoder/api/net"
	"github.com/alfg/openencoder/api/types"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type request struct {
	Profile     string `json:"profile" binding:"required"`
	Source      string `json:"source" binding:"required"`
	Destination string `json:"dest" binding:"required"`
}

type updateRequest struct {
	Status string `json:"status"`
}

type response struct {
	Message string     `json:"message"`
	Status  int        `json:"status"`
	Job     *types.Job `json:"job"`
}

type machineRequest struct {
	Provider string `json:"provider" binding:"required"`
	Region   string `json:"region" binding:"required"`
	Size     string `json:"size" binding:"required"`
	Count    int    `json:"count" binding:"required,min=1,max=10"` // Max of 10 machines.
}

func indexHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)

	c.JSON(200, gin.H{
		"name":    "openencoder",
		"version": "0.0.1",
		"github":  "https://github.com/alfg/openencoder",
		"user_id": claims["id"],
		"user":    user.(*User).UserName,
	})
}

func createJobHandler(c *gin.Context) {

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
		Status:      types.JobQueued, // Status queued.
	}

	// Send to work queue.
	_, err := enqueuer.Enqueue(config.Get().WorkerJobName, work.Q{
		"guid":        job.GUID,
		"profile":     job.Profile,
		"source":      job.Source,
		"destination": job.Destination,
	})
	if err != nil {
		log.Fatal(err)
	}
	created := data.CreateJob(job)

	// Create the encode relationship.
	ed := types.EncodeData{
		JobID: created.ID,
		Progress: types.NullFloat64{
			NullFloat64: sql.NullFloat64{
				Float64: 0,
				Valid:   true,
			},
		},
		Data: types.NullString{
			NullString: sql.NullString{
				String: "{}",
				Valid:  true,
			},
		},
	}
	edCreated := data.CreateEncodeData(ed)
	created.EncodeDataID = edCreated.EncodeDataID

	// Create response.
	resp := response{
		Message: "Job created",
		Status:  200,
		Job:     created,
	}
	c.JSON(http.StatusCreated, resp)
}

func getJobsHandler(c *gin.Context) {
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

func getJobsByIDHandler(c *gin.Context) {
	id := c.Param("id")
	jobInt, _ := strconv.Atoi(id)

	job, err := data.GetJobByID(jobInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Job does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"job":    job,
	})
}

func updateJobByIDHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// Decode json.
	var json updateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := data.GetJobByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Job does not exist",
		})
		return
	}

	if json.Status != "" {
		job.Status = json.Status
	}

	updatedJob := data.UpdateJobByID(id, *job)
	c.JSON(http.StatusOK, updatedJob)
}

func getStatsHandler(c *gin.Context) {
	stats, err := data.GetJobsStats()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Job does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"stats": gin.H{
			"jobs": stats,
		},
	})
}

func workerQueueHandler(c *gin.Context) {
	client := work.NewClient(config.Get().WorkerNamespace, redisPool)

	queues, err := client.Queues()
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, queues)
}

func workerPoolsHandler(c *gin.Context) {
	client := work.NewClient(config.Get().WorkerNamespace, redisPool)

	resp, err := client.WorkerPoolHeartbeats()
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, resp)
}

func workerBusyHandler(c *gin.Context) {
	client := work.NewClient(config.Get().WorkerNamespace, redisPool)

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

type s3ListResponse struct {
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

	resp := s3ListResponse{}

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

func profilesHandler(c *gin.Context) {
	profiles := config.Get().Profiles

	c.JSON(200, gin.H{
		"profiles": profiles,
	})
}

func machinesHandler(c *gin.Context) {
	client, _ := machine.NewDigitalOceanClient()
	ctx := context.TODO()

	// Get list of machines from DO client.
	machines, err := client.ListDropletByTag(ctx, "openencoder")
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"machines": machines,
	})
}

func createMachineHandler(c *gin.Context) {
	// Decode json.
	var json machineRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, _ := machine.NewDigitalOceanClient()
	ctx := context.TODO()

	// Create machine.
	machine, err := client.CreateDroplets(ctx, json.Region, json.Size, json.Count)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"machine": machine,
	})
	return
}

func deleteMachineHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	client, _ := machine.NewDigitalOceanClient()
	ctx := context.TODO()

	// Create machine.
	machine, err := client.DeleteDropletByID(ctx, id)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"machine": machine,
	})
	return
}

func deleteMachineByTagHandler(c *gin.Context) {
	client, _ := machine.NewDigitalOceanClient()
	ctx := context.TODO()

	// Create machine.
	err := client.DeleteDropletByTag(ctx, "openencoder")
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"deleted": true,
	})
	return
}

func listMachineRegionsHandler(c *gin.Context) {
	client, _ := machine.NewDigitalOceanClient()
	ctx := context.TODO()

	// Get list of machine regions from DO client.
	regions, err := client.ListRegions(ctx)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"regions": regions,
	})
}

func listMachineSizesHandler(c *gin.Context) {
	client, _ := machine.NewDigitalOceanClient()
	ctx := context.TODO()

	// Get list of machine sizes from DO client.
	sizes, err := client.ListSizes(ctx)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"sizes": sizes,
	})
}

func registerHandler(c *gin.Context) {
	hash, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(hash))

	c.JSON(200, gin.H{})
}

func loginHandler(c *gin.Context) {
	err := bcrypt.CompareHashAndPassword([]byte("$2a$04$2wHmBSAneLjvdFddNzlxFevZ/LoL6ZV02CZ7q0DwmR0uRYCIj4vxu"), []byte("test"))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(err == nil)

	c.JSON(200, gin.H{})
}
