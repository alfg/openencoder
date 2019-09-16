package net

import (
	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
)

// DownloadFunc creates a download.
type DownloadFunc func(job types.Job) error

// GetDownloader sets the download function.
func GetDownloader() *S3 {

	// Get credentials from settings.
	db := data.New()
	ak := db.Settings.GetSetting("AWS_ACCESS_KEY").Value
	sk := db.Settings.GetSetting("AWS_SECRET_KEY").Value
	rg := db.Settings.GetSetting("S3_INBOUND_BUCKET_REGION").Value
	ib := db.Settings.GetSetting("S3_INBOUND_BUCKET").Value
	ob := db.Settings.GetSetting("S3_OUTBOUND_BUCKET").Value

	s3 := NewS3(ak, sk, rg, ib, ob)

	return s3
}
