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

type presetUpdateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Data        string `json:"data" binding:"required"`
	Active      *bool  `json:"active" binding:"exists"`
}

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

func getPresetByIDHandler(c *gin.Context) {
	id := c.Param("id")
	presetInt, _ := strconv.Atoi(id)

	db := data.New()
	preset, err := db.Presets.GetPresetByID(presetInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Preset does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"preset": preset,
	})
}

func updatePresetByIDHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, _ := c.Get(identityKey)
	role := user.(*types.User).Role

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
	}

	// Decode json.
	var json presetUpdateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := data.New()
	preset, err := db.Presets.GetPresetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Preset does not exist",
		})
		return
	}

	// Update struct with new data if provided.
	if json.Name != "" {
		preset.Name = json.Name
	}

	if json.Description != "" {
		preset.Description = json.Description
	}

	if json.Data != "" {
		preset.Data = json.Data
	}

	preset.Active = json.Active

	updatedPreset := db.Presets.UpdatePresetByID(id, *preset)
	c.JSON(http.StatusOK, updatedPreset)
}
