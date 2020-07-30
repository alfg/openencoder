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
	ak := db.Settings.GetSetting(types.S3AccessKey).Value
	sk := db.Settings.GetSetting(types.S3SecretKey).Value
	pv := db.Settings.GetSetting(types.S3Provider).Value
	rg := db.Settings.GetSetting(types.S3OutboundBucketRegion).Value
	ib := db.Settings.GetSetting(types.S3InboundBucket).Value
	ob := db.Settings.GetSetting(types.S3OutboundBucket).Value

	s3 := NewS3(ak, sk, pv, rg, ib, ob)

	return s3
}

// GetFTPUploader sets the FTP upload function.
func GetFTPUploader() *FTP {
	db := data.New()
	addr := db.Settings.GetSetting(types.FTPAddr).Value
	user := db.Settings.GetSetting(types.FTPUsername).Value
	pass := db.Settings.GetSetting(types.FTPPassword).Value

	f := NewFTP(addr, user, pass)
	return f
}
