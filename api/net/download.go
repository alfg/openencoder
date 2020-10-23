package net

import (
	"errors"

	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
)

// Download downloads a job source based on the driver setting.
func Download(job types.Job) error {
	db := data.New()
	setting, err := db.Settings.GetSetting(types.StorageDriver)
	if err != nil {
		return err
	}
	driver := setting.Value

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
	settings := db.Settings.GetSettings()

	config := S3Config{
		AccessKey:      types.GetSetting(types.S3AccessKey, settings),
		SecretKey:      types.GetSetting(types.S3SecretKey, settings),
		Provider:       types.GetSetting(types.S3Provider, settings),
		Region:         types.GetSetting(types.S3OutboundBucketRegion, settings),
		InboundBucket:  types.GetSetting(types.S3InboundBucket, settings),
		OutboundBucket: types.GetSetting(types.S3OutboundBucket, settings),
	}
	s3 := NewS3(config)
	str, err := s3.GetPresignedURL(job)
	if err != nil {
		return str, err
	}
	return str, nil
}

// S3Download sets the download function.
func s3Download(job types.Job) error {
	db := data.New()
	settings := db.Settings.GetSettings()

	// Get job data.
	j, err := db.Jobs.GetJobByGUID(job.GUID)
	if err != nil {
		log.Error(err)
		return err
	}
	encodeID := j.EncodeID

	config := S3Config{
		AccessKey:      types.GetSetting(types.S3AccessKey, settings),
		SecretKey:      types.GetSetting(types.S3SecretKey, settings),
		Provider:       types.GetSetting(types.S3Provider, settings),
		Region:         types.GetSetting(types.S3OutboundBucketRegion, settings),
		InboundBucket:  types.GetSetting(types.S3InboundBucket, settings),
		OutboundBucket: types.GetSetting(types.S3OutboundBucket, settings),
	}
	s3 := NewS3(config)

	// Download with progress updates.
	go trackTransferProgress(encodeID, s3)
	err = s3.Download(job)
	close(progressCh)

	return err
}

// FTPDownload sets the FTP download function.
func ftpDownload(job types.Job) error {
	db := data.New()
	settings := db.Settings.GetSettings()

	addr := types.GetSetting(types.FTPAddr, settings)
	user := types.GetSetting(types.FTPUsername, settings)
	pass := types.GetSetting(types.FTPPassword, settings)

	f := NewFTP(addr, user, pass)
	err := f.Download(job)
	return err
}
