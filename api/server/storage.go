package server

import (
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
	driver := db.Settings.GetSetting(types.StorageDriver).Value

	files := getFileList(driver, prefix)

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
	ak := db.Settings.GetSetting(types.S3AccessKey).Value
	sk := db.Settings.GetSetting(types.S3SecretKey).Value
	pv := db.Settings.GetSetting(types.S3Provider).Value
	rg := db.Settings.GetSetting(types.S3InboundBucketRegion).Value
	ib := db.Settings.GetSetting(types.S3InboundBucket).Value
	ob := db.Settings.GetSetting(types.S3OutboundBucket).Value

	s3 := net.NewS3(ak, sk, pv, rg, ib, ob)

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
	addr := db.Settings.GetSetting(types.FTPAddr).Value
	user := db.Settings.GetSetting(types.FTPUsername).Value
	pass := db.Settings.GetSetting(types.FTPPassword).Value

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
