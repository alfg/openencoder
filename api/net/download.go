package net

import (
	"errors"

	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
)

// Download downloads a job source based on the driver setting.
func Download(job types.Job) error {
	db := data.New()
	driver := db.Settings.GetSetting(types.StorageDriver).Value

	if driver == "s3" {
		if err := s3Download(job); err != nil {
			return err
		}
		return nil
	} else if driver == "ftp" {
		if err := ftpDownload(job); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no driver set")
}

// GetPresignedURL gets a presigned URL from S3.
func GetPresignedURL(job types.Job) (string, error) {
	db := data.New()
	ak := db.Settings.GetSetting(types.S3AccessKey).Value
	sk := db.Settings.GetSetting(types.S3SecretKey).Value
	pv := db.Settings.GetSetting(types.S3Provider).Value
	rg := db.Settings.GetSetting(types.S3OutboundBucketRegion).Value
	ib := db.Settings.GetSetting(types.S3InboundBucket).Value
	ob := db.Settings.GetSetting(types.S3OutboundBucket).Value

	s3 := NewS3(ak, sk, pv, rg, ib, ob)
	str, err := s3.GetPresignedURL(job)
	if err != nil {
		return str, err
	}
	return str, nil
}

// S3Download sets the download function.
func s3Download(job types.Job) error {
	db := data.New()
	ak := db.Settings.GetSetting(types.S3AccessKey).Value
	sk := db.Settings.GetSetting(types.S3SecretKey).Value
	pv := db.Settings.GetSetting(types.S3Provider).Value
	rg := db.Settings.GetSetting(types.S3OutboundBucketRegion).Value
	ib := db.Settings.GetSetting(types.S3InboundBucket).Value
	ob := db.Settings.GetSetting(types.S3OutboundBucket).Value

	// Get job data.
	j, err := db.Jobs.GetJobByGUID(job.GUID)
	if err != nil {
		log.Error(err)
		return err
	}
	encodeID := j.EncodeID

	s3 := NewS3(ak, sk, pv, rg, ib, ob)

	// Download with progress updates.
	go trackTransferProgress(encodeID, s3)
	err = s3.Download(job)
	close(progressCh)

	return err
}

// FTPDownload sets the FTP download function.
func ftpDownload(job types.Job) error {
	db := data.New()
	addr := db.Settings.GetSetting(types.FTPAddr).Value
	user := db.Settings.GetSetting(types.FTPUsername).Value
	pass := db.Settings.GetSetting(types.FTPPassword).Value

	f := NewFTP(addr, user, pass)
	err := f.Download(job)
	return err
}
