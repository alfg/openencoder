package server

import (
	"fmt"
	"os"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/data"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

func indexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"name":    "openencoder",
		"version": os.Getenv("VERSION"),
		"github":  "https://github.com/alfg/openencoder",
		"docs":    "https://github.com/alfg/openencoder/blob/master/API.md",
		"wiki":    "https://github.com/alfg/openencoder/wiki",
	})
}

func healthHandler(c *gin.Context) {
	var (
		dbHealth    = "OK"
		redisHealth = "OK"
	)

	// Check database health.
	db, err := data.ConnectDB()
	if err != nil {
		dbHealth = "NOTOK"
	}
	err = db.Ping()
	if err != nil {
		dbHealth = "NOTOK"
	}
	defer db.Close()

	// Check Redis health.
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", config.Get().RedisHost, config.Get().RedisPort))
	if err != nil {
		redisHealth = "NOTOK"
	}
	defer conn.Close()

	// Check workers heartbeat.
	client := work.NewClient(config.Get().WorkerNamespace, redisPool)
	workerHeartbeats, _ := client.WorkerPoolHeartbeats()
	workers := len(workerHeartbeats)

	c.JSON(200, gin.H{
		"api":     "OK",
		"db":      dbHealth,
		"redis":   redisHealth,
		"workers": workers,
	})
}
