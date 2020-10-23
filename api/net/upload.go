package net

import (
	"errors"

	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
)

// Upload uploads a job based on the driver setting.
func Upload(job types.Job) error {
	db := data.New()
	driver, err := db.Settings.GetSetting(types.StorageDriver)
	if err != nil {
		return errors.New("no driver set")
	}

	if driver.Value == "s3" {
		if err := s3Upload(job); err != nil {
			return err
		}
		return nil
	} else if driver.Value == "ftp" {
		if err := ftpUpload(job); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no driver set")
}

// GetUploader gets the upload function.
func s3Upload(job types.Job) error {
	// Get credentials from settings.
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

	// Get job data.
	j, err := db.Jobs.GetJobByGUID(job.GUID)
	if err != nil {
		log.Error(err)
		return err
	}
	encodeID := j.EncodeID

	s3 := NewS3(config)
	go trackTransferProgress(encodeID, s3)
	err = s3.Upload(job)
	close(progressCh)

	return err
}

// GetFTPUploader sets the FTP upload function.
func ftpUpload(job types.Job) error {
	db := data.New()
	settings := db.Settings.GetSettings()

	addr := types.GetSetting(types.FTPAddr, settings)
	user := types.GetSetting(types.FTPUsername, settings)
	pass := types.GetSetting(types.FTPPassword, settings)

	f := NewFTP(addr, user, pass)
	err := f.Upload(job)
	return err
}
