package utils

import (
	"os"
	"strings"
)

func AcceptJSON(files []os.FileInfo) []os.FileInfo {
	var newFiles []os.FileInfo

	for _, file := range files {
		parts := strings.Split(file.Name(), ".")
		lastPart := parts[len(parts)-1]
		if lastPart == "json" {
			newFiles = append(newFiles, file)
		}
	}

	return newFiles
}
