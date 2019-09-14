package server

import (
	"net/http"

	"github.com/alfg/openencoder/api/data"
	"github.com/gin-gonic/gin"
)

func getStatsHandler(c *gin.Context) {
	db := data.New()
	stats, err := db.Jobs.GetJobsStats()
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
