package helpers

import (
	"os"
	"path"
)

func CreateLocalSourcePath(workDir string, src string, ID string) string {
	// Get local destination path.
	tmpDir := workDir + "/" + ID + "/"
	os.MkdirAll(tmpDir, 0700)
	os.MkdirAll(tmpDir+"src", 0700)
	os.MkdirAll(tmpDir+"dst", 0700)
	return tmpDir + path.Base(src)
}

func GetTmpPath(workDir string, ID string) string {
	tmpDir := workDir + "/" + ID + "/"
	return tmpDir
}
