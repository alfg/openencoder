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
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Default redirect to dashboard.
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	})

	// Web dashboard.
	r.Static("/dashboard", "./web/dist")

	// Catch all fallback for HTML5 History Mode.
	// https://router.vuejs.org/guide/essentials/history-mode.html
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})

	registerRoutes(r)

	log.Info("started server on port: ", serverCfg.ServerPort)
	r.Run()
}
