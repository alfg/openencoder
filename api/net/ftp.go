package net

import (
	"bufio"
	"io"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/alfg/openencoder/api/types"
	"github.com/jlaffaye/ftp"
)

const (
	// ErrorFileExists error return from FTP client.
	ErrorFileExists = "Can't create directory: File exists"
)

// FTP connection details.
type FTP struct {
	Addr     string
	Username string
	Password string
	Timeout  time.Duration
}

// NewFTP creates a new FTP instance.
func NewFTP(addr string, username string, password string) *FTP {
	return &FTP{
		Addr:     addr,
		Username: username,
		Password: password,
		Timeout:  5,
	}
}

// Download download a file from an FTP connection.
func (f *FTP) Download(job types.Job) error {
	log.Info("downloading from FTP: ", job.Source)

	// Create FTP connection.
	c, err := ftp.Dial(f.Addr, ftp.DialWithTimeout(f.Timeout*time.Second))
	if err != nil {
		log.Error(err)
		return err
	}

	// Login.
	err = c.Login(f.Username, f.Password)
	if err != nil {
		log.Error(err)
		return err
	}

	resp, err := c.Retr(job.Source)
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Close()

	outputFile, _ := os.OpenFile(job.LocalSource, os.O_WRONLY|os.O_CREATE, 0644)
	defer outputFile.Close()

	reader := bufio.NewReader(resp)
	p := make([]byte, 1024*4)

	for {
		n, err := reader.Read(p)
		if err == io.EOF {
			break
		}
		outputFile.Write(p[:n])
	}

	// Quit connection.
	if err := c.Quit(); err != nil {
		log.Error(err)
		return err
	}
	return err
}

// Upload uploads a file to FTP.
func (f *FTP) Upload(job types.Job) error {
	log.Info("uploading files to FTP: ", job.Destination)
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

	f.uploadDir(filelist, job)
	return nil
}

func (f *FTP) uploadDir(filelist []string, job types.Job) {
	for _, file := range filelist {
		f.uploadFile(file, job)
	}
}

// UploadFile uploads a file from an FTP connection.
func (f *FTP) uploadFile(path string, job types.Job) error {
	// Create FTP connection.
	c, err := ftp.Dial(f.Addr, ftp.DialWithTimeout(f.Timeout*time.Second))
	if err != nil {
		log.Error(err)
		return err
	}

	// Login.
	err = c.Login(f.Username, f.Password)
	if err != nil {
		log.Error(err)
		return err
	}

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)

	// Set destination path.
	parsedURL, _ := url.Parse(job.Destination)
	key := parsedURL.Path + filepath.Base(path)

	// Create directory.
	err = c.MakeDir(parsedURL.Path)
	if err != nil && err.(*textproto.Error).Msg != ErrorFileExists {
		log.Error(err)
		return err
	}

	err = c.Stor(key, reader)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// ListFiles lists FTP files for a given prefix.
func (f *FTP) ListFiles(prefix string) ([]*ftp.Entry, error) {
	c, err := ftp.Dial(f.Addr, ftp.DialWithTimeout(f.Timeout*time.Second))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = c.Login(f.Username, f.Password)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	entries, err := c.List(prefix)

	if err := c.Quit(); err != nil {
		log.Error(err)
		return nil, err
	}
	return entries, nil
}
