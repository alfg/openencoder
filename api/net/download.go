package net

import "github.com/alfg/openencoder/api/types"

// DownloadFunc creates a download.
type DownloadFunc func(job types.Job) error

// GetDownloader sets the download function.
func GetDownloader() *S3 {
	return &S3{}
}
