package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

var (
	redisPool *redis.Pool
	enqueuer  *work.Enqueuer
)

// Config defines configuration for creating a NewServer.
type Config struct {
	ServerPort  string
	RedisHost   string
	RedisPort   int
	Namespace   string
	JobName     string
	Concurrency uint
}

// NewServer creates a new server
func NewServer(serverCfg Config) {
	// Setup redis queue.
	redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",
				fmt.Sprintf("%s:%d", serverCfg.RedisHost, serverCfg.RedisPort))
		},
	}
	enqueuer = work.NewEnqueuer(serverCfg.Namespace, redisPool)

	// Setup server.
	r := gin.Default()

	// Default redirect to dashboard.
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	})

	// Web dashboard.
	r.Static("/dashboard", "./web/dist")

	// API routes.
	api := r.Group("/api")
	{
		// Index.
		api.GET("/", indexHandler)

		// S3.
		api.GET("/s3/list", s3ListHandler)

		// Profiles.
		api.GET("/profiles", profilesHandler)

		// Jobs.
		api.POST("/jobs", createJobHandler)
		api.GET("/jobs", getJobsHandler)
		api.GET("/jobs/:id", getJobsByIDHandler)
		api.PUT("/jobs/:id", updateJobByIDHandler)

		// Stats.
		api.GET("/stats", getStatsHandler)

		// Worker info.
		api.GET("/worker/queues", workerQueuesHandler)
		api.GET("/worker/pools", workerPoolsHandler)
		api.GET("/worker/busy", workerBusyHandler)
	}

	log.Info("started server on port: ", serverCfg.ServerPort)
	r.Run()
}
