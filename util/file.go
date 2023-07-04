package util

import (
	"os"
)

func DoesFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func IsSymlink(fileInfo os.FileInfo) bool {
	return fileInfo.Mode()&os.ModeSymlink != 0
}
