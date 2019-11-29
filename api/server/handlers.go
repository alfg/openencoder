package server

import (
	"os"

	"github.com/gin-gonic/gin"
)

func indexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"name":    "openencoder",
		"version": os.Getenv("VERSION"),
		"github":  "https://github.com/alfg/openencoder",
		"docs":    "https://github.com/alfg/openencoder/blob/master/API.md",
	})
}
