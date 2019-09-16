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
	DigitalOceanAccessToken string `json:"DIGITAL_OCEAN_ACCESS_TOKEN"`
	SlackWebhook            string `json:"SLACK_WEBHOOK"`
	S3InboundBucket         string `json:"S3_INBOUND_BUCKET"`
	S3InboundBucketRegion   string `json:"S3_INBOUND_BUCKET_REGION"`
	S3OutboundBucket        string `json:"S3_OUTBOUND_BUCKET"`
	S3OutboundBucketRegion  string `json:"S3_OUTBOUND_BUCKET_REGION"`
}

func settingsHandler(c *gin.Context) {
	user, _ := c.Get(identityKey)
	username := user.(*types.User).Username

	d := data.New()
	settings := d.Settings.GetSettingsByUsername(username)
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
		types.DigitalOceanAccessToken: json.DigitalOceanAccessToken,
		types.SlackWebhook:            json.SlackWebhook,
		types.S3InboundBucket:         json.S3InboundBucket,
		types.S3InboundBucketRegion:   json.S3InboundBucketRegion,
		types.S3OutboundBucket:        json.S3OutboundBucket,
		types.S3OutboundBucketRegion:  json.S3OutboundBucketRegion,
	}

	db := data.New()
	userID := db.Users.GetUserID(username)

	d := data.New()
	err := d.Settings.UpdateSettingsByUserID(userID, s)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "error updating settings",
		})
	}

	c.JSON(200, gin.H{
		"settings": "updated",
	})
}
