package server

import (
	"net/http"

	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/net"
	"github.com/alfg/openencoder/api/types"
	"github.com/gin-gonic/gin"
	"github.com/jlaffaye/ftp"
)

// Storage drivers available.
const (
	StorageS3  = "s3"
	StorageFTP = "ftp"
)

type storageListResponse struct {
	Folders []string `json:"folders"`
	Files   []file   `json:"files"`
}

type file struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type s3ListResponse struct {
	Folders []string `json:"folders"`
	Files   []file   `json:"files"`
}

func storageListHandler(c *gin.Context) {
	prefix := c.DefaultQuery("prefix", "")

	db := data.New()
	driver, err := db.Settings.GetSetting(types.StorageDriver)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "storage not configured",
		})
	}

	files := getFileList(driver.Value, prefix)

	c.JSON(200, gin.H{
		"data": files,
	})
}

func getFileList(driver string, prefix string) *storageListResponse {
	resp := &storageListResponse{}
	if driver == StorageS3 {
		resp, _ = getS3FileList(prefix)
	} else if driver == StorageFTP {
		resp, _ = getFTPFileList(prefix)
	}
	return resp
}

func getS3FileList(prefix string) (*storageListResponse, error) {
	db := data.New()
	settings := db.Settings.GetSettings()

	config := net.S3Config{
		AccessKey:      types.GetSetting(types.S3AccessKey, settings),
		SecretKey:      types.GetSetting(types.S3SecretKey, settings),
		Provider:       types.GetSetting(types.S3Provider, settings),
		Region:         types.GetSetting(types.S3OutboundBucketRegion, settings),
		InboundBucket:  types.GetSetting(types.S3InboundBucket, settings),
		OutboundBucket: types.GetSetting(types.S3OutboundBucket, settings),
	}

	s3 := net.NewS3(config)

	resp := &storageListResponse{}
	files, err := s3.S3ListFiles(prefix)
	if err != nil {
		return nil, err
	}

	for _, item := range files.CommonPrefixes {
		resp.Folders = append(resp.Folders, *item.Prefix)
	}

	for _, item := range files.Contents {
		var obj file
		obj.Name = *item.Key
		obj.Size = *item.Size
		resp.Files = append(resp.Files, obj)
	}
	return resp, nil
}

func getFTPFileList(prefix string) (*storageListResponse, error) {
	db := data.New()
	settings := db.Settings.GetSettings()

	addr := types.GetSetting(types.FTPAddr, settings)
	user := types.GetSetting(types.FTPUsername, settings)
	pass := types.GetSetting(types.FTPPassword, settings)

	f := net.NewFTP(addr, user, pass)
	files, err := f.ListFiles(prefix)
	if err != nil {
		return nil, err
	}

	resp := &storageListResponse{}

	for _, item := range files {
		if item.Type != ftp.EntryTypeFolder {
			var obj file
			obj.Name = item.Name
			obj.Size = int64(item.Size)
			resp.Files = append(resp.Files, obj)
		} else {
			resp.Folders = append(resp.Folders, item.Name+"/")
		}
	}
	return resp, nil
}
