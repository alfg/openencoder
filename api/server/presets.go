package server

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
	"github.com/gin-gonic/gin"
)

type createPresetRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Output      string `json:"output" binding:"required"`
	Data        string `json:"data" binding:"required"`
	Active      *bool  `json:"active" binding:"required"`
}

type presetUpdateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Output      string `json:"output" binding:"required"`
	Data        string `json:"data" binding:"required"`
	Active      *bool  `json:"active" binding:"required"`
}

func createPresetHandler(c *gin.Context) {
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	// Decode json.
	var json createPresetRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Preset.
	preset := types.Preset{
		Name:        json.Name,
		Description: json.Description,
		Output:      json.Output,
		Data:        json.Data,
		Active:      json.Active,
	}

	db := data.New()
	created, err := db.Presets.CreatePreset(preset)
	if err != nil {
		log.Error(err)
	}

	// Create response.
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"preset": created,
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
		"count":   presetsCount,
		"presets": presets,
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
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
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

	if json.Output != "" {
		preset.Output = json.Output
	}

	preset.Active = json.Active

	updatedPreset := db.Presets.UpdatePresetByID(id, *preset)
	c.JSON(http.StatusOK, updatedPreset)
}
