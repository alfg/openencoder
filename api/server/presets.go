package server

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
	"github.com/gin-gonic/gin"
)

func profilesHandler(c *gin.Context) {
	profiles := config.Get().Profiles

	c.JSON(200, gin.H{
		"profiles": profiles,
	})
}

func getPresetsHandler(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	count := c.DefaultQuery("count", "10")
	pageInt, _ := strconv.Atoi(page)
	countInt, _ := strconv.Atoi(count)

	if page == "0" {
		pageInt = 1
	}

	var wg sync.WaitGroup
	var presets *[]types.Preset
	var presetsCount int

	db := data.New()
	wg.Add(1)
	go func() {
		presets = db.Presets.GetPresets((pageInt-1)*countInt, countInt)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		presetsCount = db.Presets.GetPresetsCount()
		wg.Done()
	}()
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{
		"count": presetsCount,
		"items": presets,
	})
}
