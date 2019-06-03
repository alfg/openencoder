package server

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// NewServer creates a new server
func NewServer(port string) {
	r := gin.Default()

	// Default.
	r.GET("/", indexHandler)

	// Web dashboard.
	r.Static("/dashboard", "./static")
	// r.StaticFile("/", "./static/index.html")

	// API.
	api := r.Group("/api")
	{
		api.GET("/", indexHandler)
		api.POST("/encode", encodeHandler)
		api.GET("/jobs", jobsHandler)
		api.GET("/workers", workersHandler)
	}

	log.Info("started server on port: ", port)
	r.Run()
}
