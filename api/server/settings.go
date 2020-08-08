package server

import (
	"net/http"

	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
	"github.com/gin-gonic/gin"
)

type settingsUpdateRequest struct {
	StorageDriver string `json:"STORAGE_DRIVER" binding:"eq=s3|eq=ftp"`

	S3AccessKey            string `json:"S3_ACCESS_KEY"`
	S3SecretKey            string `json:"S3_SECRET_KEY"`
	S3InboundBucket        string `json:"S3_INBOUND_BUCKET"`
	S3InboundBucketRegion  string `json:"S3_INBOUND_BUCKET_REGION"`
	S3OutboundBucket       string `json:"S3_OUTBOUND_BUCKET"`
	S3OutboundBucketRegion string `json:"S3_OUTBOUND_BUCKET_REGION"`
	S3Provider             string `json:"S3_PROVIDER" binding:"eq=digitaloceanspaces|eq=amazonaws|eq=custom|eq="`
	S3Streaming            string `json:"S3_STREAMING" binding:"eq=enabled|eq=disabled"`
	S3Endpoint             string `json:"S3_ENDPOINT"`

	FTPAddr     string `json:"FTP_ADDR"`
	FTPUsername string `json:"FTP_USERNAME"`
	FTPPassword string `json:"FTP_PASSWORD"`

	DigitalOceanEnabled     string `json:"DIGITAL_OCEAN_ENABLED" binding:"eq=enabled|eq=disabled"`
	DigitalOceanAccessToken string `json:"DIGITAL_OCEAN_ACCESS_TOKEN"`
	DigitalOceanRegion      string `json:"DIGITAL_OCEAN_REGION"`
	DigitalOceanVPC         string `json:"DIGITAL_OCEAN_VPC"`
	SlackWebhook            string `json:"SLACK_WEBHOOK"`
}

func settingsHandler(c *gin.Context) {
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdmin(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
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
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdmin(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	// Decode json.
	var json settingsUpdateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s := map[string]string{
		types.StorageDriver: json.StorageDriver,

		types.S3AccessKey:            json.S3AccessKey,
		types.S3SecretKey:            json.S3SecretKey,
		types.S3InboundBucket:        json.S3InboundBucket,
		types.S3InboundBucketRegion:  json.S3InboundBucketRegion,
		types.S3OutboundBucket:       json.S3OutboundBucket,
		types.S3OutboundBucketRegion: json.S3OutboundBucketRegion,
		types.S3Provider:             json.S3Provider,
		types.S3Streaming:            json.S3Streaming,
		types.S3Endpoint:             json.S3Endpoint,

		types.FTPAddr:     json.FTPAddr,
		types.FTPUsername: json.FTPUsername,
		types.FTPPassword: json.FTPPassword,

		types.DigitalOceanEnabled:     json.DigitalOceanEnabled,
		types.DigitalOceanAccessToken: json.DigitalOceanAccessToken,
		types.DigitalOceanRegion:      json.DigitalOceanRegion,
		types.DigitalOceanVPC:         json.DigitalOceanVPC,
		types.SlackWebhook:            json.SlackWebhook,
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
