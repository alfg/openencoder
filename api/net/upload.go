package net

import (
	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
)

// UploadFunc creates a upload.
type UploadFunc func(job types.Job) error

// GetUploader gets the upload function.
func GetUploader() *S3 {
	// Get credentials from settings.
	db := data.New()
	ak := db.Settings.GetSetting("AWS_ACCESS_KEY").Value
	sk := db.Settings.GetSetting("AWS_SECRET_KEY").Value
	rg := db.Settings.GetSetting("AWS_REGION").Value

	s3 := NewS3(ak, sk, rg)

	return s3
}
