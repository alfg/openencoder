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
	ak := db.Settings.GetSetting(types.S3AccessKey).Value
	sk := db.Settings.GetSetting(types.S3SecretKey).Value
	pv := db.Settings.GetSetting(types.S3Provider).Value
	rg := db.Settings.GetSetting(types.S3OutboundBucketRegion).Value
	ib := db.Settings.GetSetting(types.S3InboundBucket).Value
	ob := db.Settings.GetSetting(types.S3OutboundBucket).Value

	s3 := NewS3(ak, sk, pv, rg, ib, ob)

	return s3
}

// GetFTPDownloader sets the FTP download function.
func GetFTPDownloader() *FTP {
	db := data.New()
	addr := db.Settings.GetSetting(types.FTPAddr).Value
	user := db.Settings.GetSetting(types.FTPUsername).Value
	pass := db.Settings.GetSetting(types.FTPPassword).Value

	f := NewFTP(addr, user, pass)
	return f
}
