package server

import (
	"database/sql"
	"net/http"
	"strconv"
	"sync"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/rs/xid"
)

type request struct {
	Preset      string `json:"preset" binding:"required"`
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

func createJobHandler(c *gin.Context) {
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	// Decode json.
	var json request
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Job and push the work to work queue.
	job := types.Job{
		GUID:        xid.New().String(),
		Preset:      json.Preset,
		Source:      json.Source,
		Destination: json.Destination,
		Status:      types.JobQueued, // Status queued.
	}

	// Send to work queue.
	_, err := enqueuer.Enqueue(config.Get().WorkerJobName, work.Q{
		"guid":        job.GUID,
		"preset":      job.Preset,
		"source":      job.Source,
		"destination": job.Destination,
	})
	if err != nil {
		log.Info(err)
	}

	db := data.New()
	created := db.Jobs.CreateJob(job)

	// Create the encode relationship.
	ed := types.Encode{
		JobID: created.ID,
		Progress: types.NullFloat64{
			NullFloat64: sql.NullFloat64{
				Float64: 0,
				Valid:   true,
			},
		},
		Probe: types.NullString{
			NullString: sql.NullString{
				String: "{}",
				Valid:  true,
			},
		},
		Options: types.NullString{
			NullString: sql.NullString{
				String: "{}",
				Valid:  true,
			},
		},
	}
	edCreated := db.Jobs.CreateEncode(ed)
	created.EncodeID = edCreated.EncodeID

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

	db := data.New()
	wg.Add(1)
	go func() {
		jobs = db.Jobs.GetJobs((pageInt-1)*countInt, countInt)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		jobsCount = db.Jobs.GetJobsCount()
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

	db := data.New()
	job, err := db.Jobs.GetJobByID(int64(jobInt))
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
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	// Decode json.
	var json updateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := data.New()
	job, err := db.Jobs.GetJobByID(int64(id))
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

	updatedJob := db.Jobs.UpdateJobByID(id, *job)
	c.JSON(http.StatusOK, updatedJob)
}

func getJobStatusByIDHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// Update status.
	db := data.New()
	status, _ := db.Jobs.GetJobStatusByID(int64(id))

	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"job_status": status,
	})
}

func cancelJobByIDHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	// Update status.
	db := data.New()
	db.Jobs.UpdateJobStatusByID(id, types.JobCancelled)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

func restartJobByIDHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	// Update status.
	db := data.New()
	db.Jobs.UpdateJobStatusByID(id, types.JobRestarting)

	job, _ := db.Jobs.GetJobByID(int64(id))

	// Send back to work queue.
	_, err := enqueuer.Enqueue(config.Get().WorkerJobName, work.Q{
		"guid":        job.GUID,
		"preset":      job.Preset,
		"source":      job.Source,
		"destination": job.Destination,
	})
	if err != nil {
		log.Info(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
