package server

import (
	"net/http"

	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
	"github.com/gin-gonic/gin"
)

type settingsUpdateRequest struct {
	S3AccessKey             string `json:"S3_ACCESS_KEY"`
	S3SecretKey             string `json:"S3_SECRET_KEY"`
	DigitalOceanAccessToken string `json:"DIGITAL_OCEAN_ACCESS_TOKEN"`
	SlackWebhook            string `json:"SLACK_WEBHOOK"`
	S3InboundBucket         string `json:"S3_INBOUND_BUCKET"`
	S3InboundBucketRegion   string `json:"S3_INBOUND_BUCKET_REGION"`
	S3OutboundBucket        string `json:"S3_OUTBOUND_BUCKET"`
	S3OutboundBucketRegion  string `json:"S3_OUTBOUND_BUCKET_REGION"`
	S3Provider              string `json:"S3_PROVIDER" binding:"eq=digitalocean|eq=aws|eq="`
}

func settingsHandler(c *gin.Context) {
	user, _ := c.Get(identityKey)
	role := user.(*types.User).Role

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
	}

	d := data.New()
	settings := d.Settings.GetSettings()
	settingOptions := d.Settings.GetSettingsOptions()

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

	c.JSON(http.StatusOK, gin.H{
		"settings": resp,
	})
}

func updateSettingsHandler(c *gin.Context) {
	user, _ := c.Get(identityKey)
	role := user.(*types.User).Role

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
	}

	// Decode json.
	var json settingsUpdateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s := map[string]string{
		types.S3AccessKey:             json.S3AccessKey,
		types.S3SecretKey:             json.S3SecretKey,
		types.DigitalOceanAccessToken: json.DigitalOceanAccessToken,
		types.SlackWebhook:            json.SlackWebhook,
		types.S3InboundBucket:         json.S3InboundBucket,
		types.S3InboundBucketRegion:   json.S3InboundBucketRegion,
		types.S3OutboundBucket:        json.S3OutboundBucket,
		types.S3OutboundBucketRegion:  json.S3OutboundBucketRegion,
		types.S3Provider:              json.S3Provider,
	}

	db := data.New()
	// userID := db.Users.GetUserID(username)

	err := db.Settings.UpdateSettings(s)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error updating settings",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"settings": "updated",
	})
}
