package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func getFileInfo(filePath string) []string {
	info, err := os.Stat(filePath)
	if err != nil { // isFile exist?
		fmt.Println(err)
	}
	fSize := info.Size() // KB
	// fTime := time.Unix(info.ModTime())
	fTime := info.ModTime()
	aFileInfo := []string{
		fmt.Sprintf("%v", fSize),
		fmt.Sprintf("%v", fTime),
	}
	return aFileInfo
}

func getFileInfoOld(filePath string) []string {
	var s syscall.Stat_t
	syscall.Stat(filePath, &s)
	if _, err := os.Stat(filePath); err != nil { // isFile exist?
		fmt.Println(err)
	}
	fSize := s.Size // KB
	fTime := time.Unix(s.Mtim.Unix())
	aFileInfo := []string{
		fmt.Sprintf("%v", fSize),
		fmt.Sprintf("%v", fTime),
	}
	return aFileInfo
}

func main() {
	aFileInfo := getFileInfo("test.txt")
	fmt.Printf("info: %v", aFileInfo)
}
