package server

import (
	"fmt"

	"github.com/alfg/openencoder/api/net"
	"github.com/gin-gonic/gin"
	"github.com/jlaffaye/ftp"
)

type ftpListResponse struct {
	Folders []string  `json:"folders"`
	Files   []ftpFile `json:"files"`
}

type ftpFile struct {
	Name string `json:"name"`
	Size uint64 `json:"size"`
}

func ftpListHandler(c *gin.Context) {
	prefix := c.DefaultQuery("prefix", "")

	host := "localhost"
	port := 21

	f := net.NewFTP(host, port)
	files, err := f.FTPListFiles(prefix)
	if err != nil {
		fmt.Println(err)
	}

	resp := ftpListResponse{}

	for _, item := range files {
		if item.Type != ftp.EntryTypeFolder {
			var obj ftpFile
			obj.Name = item.Name
			obj.Size = item.Size
			resp.Files = append(resp.Files, obj)
		} else {
			resp.Folders = append(resp.Folders, item.Name)
		}

		// fmt.Println(item.Size, item.Target, item.Type)
		fmt.Println("is folder? ", item.Name, item.Type == ftp.EntryTypeFolder)
	}

	c.JSON(200, gin.H{
		"data": resp,
	})
}
