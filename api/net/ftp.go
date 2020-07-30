package net

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"time"

	"github.com/alfg/openencoder/api/types"
	"github.com/jlaffaye/ftp"
)

// FTP connection details.
type FTP struct {
	Addr     string
	Username string
	Password string
	Timeout  time.Duration
}

// NewFTP creates a new S3 instance.
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

	// Open file for writing.
	// file, err := os.Create(job.LocalSource)
	// if err != nil {
	// 	return err
	// }

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

	// resp, err := c.Retr("tears-of-steel-2s.mp4")
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

	// buf, err := ioutil.ReadAll(resp)

	// Quit connection.
	if err := c.Quit(); err != nil {
		log.Error(err)
		return err
	}
	return err
}

func (f *FTP) Upload(job types.Job) error {
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

	data := bytes.NewBufferString("testing")
	err = c.Stor("test-file.txt", data)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// ListFiles lists s3 objects for a given prefix.
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
	// for _, e := range entries {
	// 	fmt.Println(e)
	// }

	if err := c.Quit(); err != nil {
		log.Error(err)
		return nil, err
	}
	return entries, nil
}
