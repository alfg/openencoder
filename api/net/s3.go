package net

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

// S3Download downloads source files from S3.
// AWS_REGION, AWS_ACCESS_KEY, and AWS_SECRET_KEY envvars must be set!
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
	writer := &progressWriter{writer: file, size: size, written: 0}
	objInput := s3.GetObjectInput{
		Bucket: aws.String(config.Get().S3InboundBucket),
		Key:    aws.String(key),
	}

	// Download file to local.
	if _, err = downloader.Download(writer, &objInput); err != nil {
		log.Printf("download failed! deleting file: %s", file.Name())
		os.Remove(file.Name())
		panic(err)
	}
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
		}
		filelist = append(filelist, path)
		return nil
	})

	uploadDir(filelist, job)
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

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("upload error: ", err)
		return err
	}

	// Set key.
	parsedURL, _ := url.Parse(job.Destination)
	key := parsedURL.Path + filepath.Base(path)

	reader := &progressReader{
		fp:   file,
		size: fileInfo.Size(),
	}

	uploader := s3manager.NewUploader(session.New(), func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.LeavePartsOnError = true
	})

	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
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

// progressWriter tracks the download progress.
type progressWriter struct {
	written int64
	writer  io.WriterAt
	size    int64
}

func (pw *progressWriter) WriteAt(p []byte, off int64) (int, error) {
	atomic.AddInt64(&pw.written, int64(len(p)))
	percentageDownloaded := float32(pw.written*100) / float32(pw.size)
	fmt.Printf("File size:%d downloaded:%d percentage:%.2f%%\r", pw.size, pw.written, percentageDownloaded)
	return pw.writer.WriteAt(p, off)
}

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func getFileSize(svc *s3.S3, bucket, prefix string) (filesize int64, error error) {
	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(prefix),
	}

	resp, err := svc.HeadObject(params)
	if err != nil {
		return 0, err
	}
	return *resp.ContentLength, nil
}

// progressReader for uploading progress.
type progressReader struct {
	fp   *os.File
	size int64
	read int64
}

func (r *progressReader) Read(p []byte) (int, error) {
	return r.fp.Read(p)
}

func (r *progressReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := r.fp.ReadAt(p, off)
	if err != nil {
		return n, err
	}
	atomic.AddInt64(&r.read, int64(n))
	fmt.Printf("total read:%d progress:%d%%\r", r.read/2, int(float32(r.read*100/2)/float32(r.size)))
	return n, err
}

func (r *progressReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}
