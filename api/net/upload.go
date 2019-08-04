package net

import "github.com/alfg/openencoder/api/types"

// UploadFunc creates a upload.
type UploadFunc func(job types.Job) error

// GetUploader gets the upload function.
func GetUploader() *S3 {
	return &S3{}
}
