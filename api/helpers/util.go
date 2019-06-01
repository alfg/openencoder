package helpers

import (
	"os"
	"path"
)

func GetLocalSourcePath(src string, ID string) string {
	// Get local destination path.
	tmpDir := "/tmp/" + ID + "/"
	os.MkdirAll(tmpDir, 0700)
	os.MkdirAll(tmpDir+"src", 0700)
	os.MkdirAll(tmpDir+"dst", 0700)
	return tmpDir + path.Base(src)
}
