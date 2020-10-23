package net

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Config describes a configuration for setting up S3.
type S3Config struct {
	Provider       string
	Endpoint       string
	AccessKey      string
	SecretKey      string
	Region         string
	InboundBucket  string
	OutboundBucket string
}

// S3 creates a new S3 instance.
type S3 struct {
	Progress progress
	Writer   *ProgressWriter
	Reader   *ProgressReader

	Config S3Config
}

type progress struct {
	quit     chan struct{}
	Progress float32
}

// NewS3 creates a new S3 instance.
func NewS3(config S3Config) *S3 {
	config.Endpoint = getEndpoint(config.Provider, config.Region)

	return &S3{
		Config: config,
	}
}

// Download downloads source files from S3.
func (s *S3) Download(job types.Job) error {
	log.Info("downloading from S3: ", job.Source)

	// Open file for writing.
	file, err := os.Create(job.LocalSource)
	if err != nil {
		return err
	}

	// Create session and client.
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(s.Config.Endpoint),
		Region:      aws.String(s.Config.Region),
		Credentials: credentials.NewStaticCredentials(s.Config.AccessKey, s.Config.SecretKey, ""),
	})
	if err != nil {
		panic(err)
	}
	s3Client := s3.New(sess)
	downloader := s3manager.NewDownloader(sess)

	parsedURL, _ := url.Parse(job.Source)
	key := parsedURL.Path

	size, err := getFileSize(s3Client, s.Config.InboundBucket, key)
	if err != nil {
		panic(err)
	}
	log.Println("starting download, size: ", byteCountDecimal(size))

	// Get object input details.
	s.Writer = &ProgressWriter{writer: file, size: size, written: 0}
	objInput := s3.GetObjectInput{
		Bucket: aws.String(s.Config.InboundBucket),
		Key:    aws.String(key),
	}

	// Download file to local.
	go s.trackProgress("download")
	if _, err = downloader.Download(s.Writer, &objInput); err != nil {
		log.Printf("download failed! deleting file: %s", file.Name())
		os.Remove(file.Name())
		panic(err)
	}
	file.Close()

	s.finish()
	return err
}

func (s *S3) trackProgress(t string) {
	s.Progress.quit = make(chan struct{})
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-s.Progress.quit:
			ticker.Stop()
			return
		case <-ticker.C:
			if t == "download" {
				// Download progress.
				s.Progress.Progress = float32(s.Writer.written*100) / float32(s.Writer.size)
			} else if t == "upload" {
				// Upload progress.
				s.Progress.Progress = float32(s.Reader.read*100/2) / float32(s.Reader.size) // Upload.
			}

		}
	}
}

func (s *S3) finish() {
	close(s.Progress.quit)
}

// Upload uploads a file to S3.
func (s *S3) Upload(job types.Job) error {
	log.Info("uploading files to S3: ", job.Destination)
	defer log.Info("upload complete")

	// Get list of files in output dir.
	filelist := []string{}
	filepath.Walk(path.Dir(job.LocalSource)+"/dst", func(path string, f os.FileInfo, err error) error {
		if isDirectory(path) {
			return nil
		}
		filelist = append(filelist, path)
		return nil
	})

	s.uploadDir(filelist, job)
	return nil
}

func (s *S3) uploadDir(filelist []string, job types.Job) {
	for _, file := range filelist {
		s.uploadFile(file, job)
	}
}

func (s *S3) uploadFile(path string, job types.Job) error {
	log.Info("uploading file to S3.", job.Destination)

	// Open source path file.
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("upload error: ", err)
		return err
	}

	// Set key.
	parsedURL, _ := url.Parse(job.Destination)
	key := parsedURL.Path + filepath.Base(path)

	s.Reader = &ProgressReader{
		fp:   file,
		size: fileInfo.Size(),
	}

	go s.trackProgress("upload")

	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(s.Config.Endpoint),
		Region:      aws.String(s.Config.Region),
		Credentials: credentials.NewStaticCredentials(s.Config.AccessKey, s.Config.SecretKey, ""),
	})
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.LeavePartsOnError = true
	})

	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:   s.Reader,
		Bucket: aws.String(s.Config.OutboundBucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

// S3ListFiles lists s3 objects for a given prefix.
func (s *S3) S3ListFiles(prefix string) (*s3.ListObjectsV2Output, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(s.Config.Endpoint),
		Region:      aws.String(s.Config.Region),
		Credentials: credentials.NewStaticCredentials(s.Config.AccessKey, s.Config.SecretKey, ""),
	})
	svc := s3.New(sess)

	resp, err := svc.ListObjectsV2(
		&s3.ListObjectsV2Input{
			Bucket:    aws.String(s.Config.InboundBucket),
			Delimiter: aws.String("/"),
			Prefix:    aws.String(prefix),
		},
	)
	return resp, err
}

// GetPresignedURL generates a presigned URL from S3.
func (s *S3) GetPresignedURL(job types.Job) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(s.Config.Endpoint),
		Region:      aws.String(s.Config.Region),
		Credentials: credentials.NewStaticCredentials(s.Config.AccessKey, s.Config.SecretKey, ""),
	})
	svc := s3.New(sess)

	parsedURL, _ := url.Parse(job.Source)
	key := parsedURL.Path

	objInput := s3.GetObjectInput{
		Bucket: aws.String(s.Config.InboundBucket),
		Key:    aws.String(key),
	}

	req, _ := svc.GetObjectRequest(&objInput)
	urlStr, err := req.Presign(PresignedDuration)

	return urlStr, err
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

func getEndpoint(provider, region string) string {
	if strings.ToUpper(provider) == types.Custom {
		db := data.New()
		endpoint, err := db.Settings.GetSetting(types.S3Endpoint)
		if err != nil {
			return ""
		}
		return endpoint.Value
	} else if strings.ToUpper(provider) == types.DigitalOceanSpaces {
		return EndpointDigitalOceanSpacesRegion(region)
	}
	return EndpointAmazonAWSRegion(region)
}

func trackTransferProgress(encodeID int64, s3 *S3) {
	db := data.New()
	progressCh = make(chan struct{})
	ticker := time.NewTicker(ProgressInterval)

	for {
		select {
		case <-progressCh:
			ticker.Stop()
			return
		case <-ticker.C:
			log.Info("transfer progress: ", s3.Progress.Progress)
			err := db.Jobs.UpdateTransferProgressByID(encodeID, float64(s3.Progress.Progress))
			fmt.Println(float64(s3.Progress.Progress))
			if err != nil {
				log.Error(err)
			}
		}
	}
}
