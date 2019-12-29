package server

import (
	"fmt"
	"net/http"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/data"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

func indexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":    ProjectName,
		"version": ProjectVersion,
		"github":  ProjectGithub,
		"docs":    ProjectDocs,
		"wiki":    ProjectWiki,
	})
}

func healthHandler(c *gin.Context) {
	var (
		dbHealth    = OK
		redisHealth = OK
	)

	// Check database health.
	db, err := data.ConnectDB()
	if err != nil {
		dbHealth = NOK
	}
	err = db.Ping()
	if err != nil {
		dbHealth = NOK
	}
	defer db.Close()

	// Check Redis health.
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", config.Get().RedisHost, config.Get().RedisPort))
	if err != nil {
		redisHealth = NOK
	}
	defer conn.Close()

	// Check workers heartbeat.
	client := work.NewClient(config.Get().WorkerNamespace, redisPool)
	workerHeartbeats, _ := client.WorkerPoolHeartbeats()
	workers := len(workerHeartbeats)

	c.JSON(200, gin.H{
		"api":     OK,
		"db":      dbHealth,
		"redis":   redisHealth,
		"workers": workers,
	})
}
