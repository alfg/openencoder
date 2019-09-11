package server

import (
	"github.com/alfg/openencoder/api/config"
	"github.com/gin-gonic/gin"
)

func profilesHandler(c *gin.Context) {
	profiles := config.Get().Profiles

	c.JSON(200, gin.H{
		"profiles": profiles,
	})
}
