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
	ak := data.GetSetting("AWS_ACCESS_KEY").Value
	sk := data.GetSetting("AWS_SECRET_KEY").Value
	rg := data.GetSetting("AWS_REGION").Value

	s3 := NewS3(ak, sk, rg)

	return s3
}
