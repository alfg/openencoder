package net

import (
	"fmt"
	"time"

	"github.com/alfg/openencoder/api/types"
	"github.com/jlaffaye/ftp"
)

// FTP connection details.
type FTP struct {
	Host     string
	Port     int
	Username string
	Password string
	Timeout  time.Duration
}

// NewFTP creates a new S3 instance.
func NewFTP(host string, port int) *FTP {
	return &FTP{
		Host:    host,
		Port:    port,
		Timeout: 5,
	}
}

// FTPDownload download a file from an FTP connection.
func (f *FTP) FTPDownload(job types.Job) error {
	addr := fmt.Sprintf("%s:%d", f.Host, f.Port)
	c, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Error(err)
		return err
	}

	err = c.Login("anonymous", "anonymous")
	if err != nil {
		log.Error(err)
		return err
	}

	// Do something with the FTP conn

	if err := c.Quit(); err != nil {
		log.Error(err)
		return err
	}
	return err
}

// FTPListFiles lists s3 objects for a given prefix.
func (f *FTP) FTPListFiles(prefix string) ([]*ftp.Entry, error) {
	addr := fmt.Sprintf("%s:%d", f.Host, f.Port)
	c, err := ftp.Dial(addr, ftp.DialWithTimeout(f.Timeout*time.Second))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = c.Login("username", "mypass")
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Do something with the FTP conn
	entries, err := c.List("/")
	for _, e := range entries {
		fmt.Println(e)
	}

	if err := c.Quit(); err != nil {
		log.Error(err)
		return nil, err
	}
	return entries, nil
}
