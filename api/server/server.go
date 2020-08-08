package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alfg/openencoder/api/logging"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// Project constants.
const (
	ProjectName   = "openencoder"
	ProjectGithub = "https://github.com/alfg/openencoder"
	ProjectDocs   = "https://github.com/alfg/openencoder/blob/master/API.md"
	ProjectWiki   = "https://github.com/alfg/openencoder/wiki"
	OK            = "OK"
	NOK           = "NOK"

	// Machines.
	WorkerTag = "openencoder-worker"

	// JWT settings.
	JwtRealm       = "openencoder"
	JwtIdentityKey = "id"
	JwtRoleKey     = "role"
	JwtTimeout     = time.Hour // Duration a JWT is valid.
	JwtMaxRefresh  = time.Hour // Duration a JWT can be refreshed.

	// User role types.
	RoleAdmin    = "admin"
	RoleOperator = "operator"
	RoleGuest    = "guest"
)

var (
	// ProjectVersion gets the current version set by the CI build.
	ProjectVersion = os.Getenv("VERSION")

	// Server settings.
	redisPool *redis.Pool
	enqueuer  *work.Enqueuer
	log       = logging.Log
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
	r := gin.New()
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
