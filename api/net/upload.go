package net

import "github.com/alfg/openencoder/api/types"

// UploadFunc creates a upload.
type UploadFunc func(job types.Job) error

// GetUploadFunc sets the upload function.
func GetUploadFunc() UploadFunc {
	return S3Upload
}
