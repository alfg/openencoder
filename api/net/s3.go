package net

import (
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/alfg/openencoder/api/config"

	"github.com/alfg/openencoder/api/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

// S3Download downloads source files from S3.
func S3Download(job types.Job) error {
	log.Info("downloading from S3: ", job.Source)

	// Get local destination path.
	// tmpDir := "/tmp" + "/asdf/"
	// os.MkdirAll(tmpDir, 0700)

	// Open file for writing.
	// file, err := os.Create(tmpDir + path.Base(job.Source))
	file, err := os.Create(job.LocalSource)
	if err != nil {
		return err
	}

	// Create session.
	downloader := s3manager.NewDownloader(session.New(
		&aws.Config{
			Region: aws.String(config.C.S3InboundRegion),
		},
	))

	parsedURL, _ := url.Parse(job.Source)
	key := parsedURL.Path

	// Get object input details.
	objInput := s3.GetObjectInput{
		Bucket: aws.String(config.Get().S3InboundBucket),
		Key:    aws.String(key),
	}

	// Download file to local.
	_, err = downloader.Download(file, &objInput)
	file.Close()

	return err
}

// S3Upload uploads a file to S3.
func S3Upload(job types.Job) error {
	log.Info("uploading files to S3: ", job.Destination)
	defer log.Info("upload complete")

	// Get list of files in output dir.
	filelist := []string{}
	filepath.Walk(path.Dir(job.LocalSource)+"/dst", func(path string, f os.FileInfo, err error) error {
		if isDirectory(path) {
			return nil
		} else {
			filelist = append(filelist, path)
			return nil
		}
	})

	uploadDir(filelist, job)
	return nil
}

// S3ListFiles lists s3 objects for a given prefix.
func S3ListFiles(prefix string) (*s3.ListObjectsV2Output, error) {
	sess := session.New(
		&aws.Config{
			Region: aws.String(config.Get().S3InboundRegion),
		},
	)

	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(
		&s3.ListObjectsV2Input{
			Bucket:    aws.String(config.Get().S3InboundBucket),
			Delimiter: aws.String("/"),
			Prefix:    aws.String(prefix),
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func uploadDir(filelist []string, job types.Job) {
	for _, file := range filelist {
		uploadFile(file, job)
	}
}

func uploadFile(path string, job types.Job) error {
	log.Info("uploading file to S3.", job.Destination)

	// Open source path file.
	// tmpDir := "/tmp" + "/asdf/"
	// file, err := os.Open(tmpDir + path.Base(job.Source))
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	// Set key.
	parsedURL, _ := url.Parse(job.Destination)
	key := parsedURL.Path + filepath.Base(path)

	uploader := s3manager.NewUploader(session.New(&aws.Config{
		Region: aws.String(config.Get().S3OutboundRegion),
	}))
	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(config.Get().S3OutboundBucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	return nil
}

func isDirectory(path string) bool {
	fd, err := os.Stat(path)
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}
	switch mode := fd.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	}
	return false
}
