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
	ak := db.Settings.GetSetting("S3_ACCESS_KEY").Value
	sk := db.Settings.GetSetting("S3_SECRET_KEY").Value
	pv := db.Settings.GetSetting("S3_PROVIDER").Value
	rg := db.Settings.GetSetting("S3_INBOUND_BUCKET_REGION").Value
	ib := db.Settings.GetSetting("S3_INBOUND_BUCKET").Value
	ob := db.Settings.GetSetting("S3_OUTBOUND_BUCKET").Value

	s3 := NewS3(ak, sk, pv, rg, ib, ob)

	return s3
}
