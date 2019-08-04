package net

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ProgressWriter tracks the download progress.
type ProgressWriter struct {
	written int64
	writer  io.WriterAt
	size    int64
}

func (pw *ProgressWriter) WriteAt(p []byte, off int64) (int, error) {
	atomic.AddInt64(&pw.written, int64(len(p)))
	// percentageDownloaded := float32(pw.written*100) / float32(pw.size)
	// fmt.Printf("File size:%d downloaded:%d percentage:%.2f%%\r", pw.size, pw.written, percentageDownloaded)
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

// ProgressReader for uploading progress.
type ProgressReader struct {
	fp       *os.File
	size     int64
	read     int64
	Progress int
}

func (r *ProgressReader) Read(p []byte) (int, error) {
	return r.fp.Read(p)
}

func (r *ProgressReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := r.fp.ReadAt(p, off)
	if err != nil {
		return n, err
	}
	atomic.AddInt64(&r.read, int64(n))
	// fmt.Printf("total read:%d progress:%d%%\r", r.read/2, int(float32(r.read*100/2)/float32(r.size)))
	return n, err
}

func (r *ProgressReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}
