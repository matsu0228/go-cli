package cmd

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func moveLog(trgFilePath, dstPath string) { // targFilePath is included filename
	debug(fmt.Sprintf("move: %s to %s \n", trgFilePath, dstPath))
	targetName := getFilenameFromPath(trgFilePath)
	dstFilePath := strings.TrimRight(dstPath, "/") + "/" + targetName
	if err := os.Rename(trgFilePath, dstFilePath); err != nil {
		errorlog(fmt.Sprintln(err))
	}
}

func deleteFile(trgFilePath string) {
	if err := os.Remove(trgFilePath); err != nil {
		errorlog(fmt.Sprintln(err))
	}
}

func compressToGzip(trgFilePath, dstPath string) {
	level := gzip.DefaultCompression
	targetName := getFilenameFromPath(trgFilePath)
	dstFilePath := strings.TrimRight(dstPath, "/") + "/" + targetName + ".gz" // targetFileName
	file, err := os.Create(dstFilePath)
	if err != nil { // setup out file
		errorlog(fmt.Sprintln(err))
	}
	defer file.Close()
	writer, err := gzip.NewWriterLevel(file, level)
	if err != nil {
		errorlog(fmt.Sprintln(err))
	}
	defer writer.Close()
	// tw := tar.NewWriter(writer)
	// defer tw.Close() // end setup put file

	body, err := ioutil.ReadFile(trgFilePath)
	if err != nil { // get org body
		errorlog(fmt.Sprintln(err))
	}
	// if _, err := tw.Write(body); err != nil { // write gzip data
	if _, err := writer.Write(body); err != nil { // write gzip data
		errorlog(fmt.Sprintln(err))
	} else { // 圧縮が成功したときだけ元ファイル削除
		debug("success zip and delete:" + trgFilePath)
		deleteFile(trgFilePath)
	}
}
