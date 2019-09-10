package server

import (
	"net/http"

	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
	"github.com/gin-gonic/gin"
)

type settingsUpdateRequest struct {
	AWSAccessKey            string `json:"AWS_ACCESS_KEY"`
	AWSSecretKey            string `json:"AWS_SECRET_KEY"`
	AWSRegion               string `json:"AWS_REGION"`
	DigitalOceanAccessToken string `json:"DIGITAL_OCEAN_ACCESS_TOKEN"`
	SlackWebhook            string `json:"SLACK_WEBHOOK"`
}

func settingsHandler(c *gin.Context) {
	user, _ := c.Get(identityKey)
	username := user.(*types.User).Username

	settings := data.GetSettingsByUsername(username)
	settingOptions := data.GetSettingsOptions()

	// Get all settings for response and set blank defaults.
	var resp []types.SettingsForm
	for _, v := range settingOptions {
		s := types.SettingsForm{
			Title:       v.Title,
			Name:        v.Name,
			Description: v.Description,
			Secure:      v.Secure,
		}
		for _, j := range settings {
			if j.Name == v.Name {
				s.Value = j.Value
			}
		}
		resp = append(resp, s)
	}

	c.JSON(200, gin.H{
		"settings": resp,
	})
}

func updateSettingsHandler(c *gin.Context) {
	user, _ := c.Get(identityKey)
	username := user.(*types.User).Username

	// Decode json.
	var json settingsUpdateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s := map[string]string{
		types.AWSAccessKey:            json.AWSAccessKey,
		types.AWSSecretKey:            json.AWSSecretKey,
		types.AWSRegion:               json.AWSRegion,
		types.DigitalOceanAccessToken: json.DigitalOceanAccessToken,
		types.SlackWebhook:            json.SlackWebhook,
	}

	// userID is 0 for some reason.
	userID := data.GetUserID(username)
	err := data.UpdateSettingsByUserID(userID, s)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "error updating settings",
		})
	}

	c.JSON(200, gin.H{
		"settings": "updated",
	})
}
