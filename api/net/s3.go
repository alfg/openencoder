package net

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

// S3 creates a new S3 instance.
type S3 struct {
	Progress progress
	Writer   *ProgressWriter
	Reader   *ProgressReader
}

type progress struct {
	quit     chan struct{}
	Progress float64
}

// S3Download downloads source files from S3.
// AWS_REGION, AWS_ACCESS_KEY, and AWS_SECRET_KEY envvars must be set!
func (s *S3) S3Download(job types.Job) error {
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

	// Create session and client.
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	s3Client := s3.New(sess)
	downloader := s3manager.NewDownloader(sess)

	parsedURL, _ := url.Parse(job.Source)
	key := parsedURL.Path

	size, err := getFileSize(s3Client, config.Get().S3InboundBucket, key)
	if err != nil {
		panic(err)
	}
	log.Println("starting download, size: ", byteCountDecimal(size))

	// Get object input details.
	s.Writer = &ProgressWriter{writer: file, size: size, written: 0}
	objInput := s3.GetObjectInput{
		Bucket: aws.String(config.Get().S3InboundBucket),
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
				s.Progress.Progress = float64(s.Writer.written*100) / float64(s.Writer.size)
			} else if t == "upload" {
				// Upload progress.
				s.Progress.Progress = float64(s.Reader.read*100/2) / float64(s.Reader.size) // Upload.
			}

		}
	}
}

func (s *S3) finish() {
	close(s.Progress.quit)
}

// S3Upload uploads a file to S3.
func (s *S3) S3Upload(job types.Job) error {
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
	// tmpDir := "/tmp" + "/asdf/"
	// file, err := os.Open(tmpDir + path.Base(job.Source))
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
	uploader := s3manager.NewUploader(session.New(), func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.LeavePartsOnError = true
	})

	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:   s.Reader,
		Bucket: aws.String(config.Get().S3OutboundBucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

// S3ListFiles lists s3 objects for a given prefix.
func S3ListFiles(prefix string) (*s3.ListObjectsV2Output, error) {
	sess := session.New()

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
