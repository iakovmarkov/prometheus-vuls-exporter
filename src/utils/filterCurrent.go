package utils

import (
	"os"
)

func FilterCurrent(files []os.FileInfo) []os.FileInfo {
	var newFiles []os.FileInfo

	for _, file := range files {
		if file.Name() != "current" {
			newFiles = append(newFiles, file)
		}
	}

	return newFiles
}
