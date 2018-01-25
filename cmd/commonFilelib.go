package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func getChildFiles(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		errorlog(fmt.Sprintln(err))
	}
	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, getChildFiles(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}
	return paths
}

func getFilenameWithoutExtention(fileName string) string {
	pos := strings.LastIndex(fileName, ".")
	return fileName[:pos]
}

func validateFileformat(fileName string) string {
	errMsg := ""
	pos := strings.LastIndex(fileName, ".")
	if fileName[pos:] != ".log" {
		// err = errors.New("this file extention is NOT log")
		errMsg = "this file's extention is NOT log"
	}
	return errMsg
}

func getFilenameFromPath(targetFilePath string) string {
	pos := strings.LastIndex(targetFilePath, "/")
	targetName := targetFilePath[pos+1:]
	return targetName
}
