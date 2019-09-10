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

	ak := data.GetSetting("AWS_ACCESS_KEY").Value
	sk := data.GetSetting("AWS_SECRET_KEY").Value
	rg := data.GetSetting("AWS_REGION").Value

	s3 := net.NewS3(ak, sk, rg)

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
