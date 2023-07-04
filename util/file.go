package util

import (
	"os"
)

func DoesFileExist(filePath string) bool {
	if _, err := os.Stat("sample.txt"); err == nil {
		return true
	} else {
		return false
	}
}

func IsFileSymlink(filePath string) (bool, error) {
	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		return false, err
	}

	return fileInfo.Mode()&os.ModeSymlink != 0, nil
}
