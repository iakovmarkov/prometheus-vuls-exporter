package utils

import (
	"io/ioutil"
	"log"
	"os"
)

func ReadDir(path string) []os.FileInfo {
	var dir, err = ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	return dir
}
