package server

import (
	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/net"
	"github.com/gin-gonic/gin"
)

type s3ListResponse struct {
	Folders []string `json:"folders"`
	Files   []file   `json:"files"`
}

type file struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func s3ListHandler(c *gin.Context) {
	prefix := c.DefaultQuery("prefix", "")

	db := data.New()
	ak := db.Settings.GetSetting("S3_ACCESS_KEY").Value
	sk := db.Settings.GetSetting("S3_SECRET_KEY").Value
	pv := db.Settings.GetSetting("S3_PROVIDER").Value
	rg := db.Settings.GetSetting("S3_INBOUND_BUCKET_REGION").Value
	ib := db.Settings.GetSetting("S3_INBOUND_BUCKET").Value
	ob := db.Settings.GetSetting("S3_OUTBOUND_BUCKET").Value

	s3 := net.NewS3(ak, sk, pv, rg, ib, ob)

	resp := s3ListResponse{}
	files, err := s3.S3ListFiles(prefix)
	if err != nil {
		c.JSON(200, gin.H{
			"data": resp,
		})
		return
	}

	// var prefixes &[]S3ListResponse.Folders
	for _, item := range files.CommonPrefixes {
		resp.Folders = append(resp.Folders, *item.Prefix)
	}

	for _, item := range files.Contents {
		var obj file
		obj.Name = *item.Key
		obj.Size = *item.Size
		resp.Files = append(resp.Files, obj)
	}

	c.JSON(200, gin.H{
		"data": resp,
	})
}
