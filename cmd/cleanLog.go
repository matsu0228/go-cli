// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"github.com/hashicorp/logutils"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	// "reflect"
	"strings"
	"syscall"
	"time"
)

// type ConfCsv struct {
var ConfKeyMap map[string]int

/*
0. logging
1. parse config
 - conf: src_dir / dst_dir / cmd / condition
 - access remote server (use local nw)
2. loop conf
2-1. check condition
 - condition= days / KB (and)
 - getChildFiles
 - check days / KB
2-2. select mode
cmd= delete or move or zipmove
2-3. exec mode function

3. logging success

4. test
*/
func readConfCsv(confFile string) ([]string, [][]string) {
	var confData [][]string
	warn("open conf file of '" + confFile + "' ")
	fp, err := os.Open(confFile)
	if err != nil {
		errorlog(fmt.Sprintln(err))
	}
	reader := csv.NewReader(fp)
	defer fp.Close()
	header, _ := reader.Read()
	for {
		record, err := reader.Read()
		// fmt.Println(reflect.TypeOf(record))
		if len(record) != len(ConfKeyMap) {
			warn(fmt.Sprintf("set %s colmuns. this line data is following\n", len(ConfKeyMap)))
			warn(fmt.Sprintln(record))
		} else {
			confData = append(confData, record)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			errorlog(fmt.Sprintln(err))
		}
	}
	return header, confData
}

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

func isTargetLog(filePath string, conf []string) bool {
	// https://golang.org/pkg/syscall/#Stat_t
	isTargetLogCond := false
	isSize, isTime := false, false // 各条件の判定
	var s syscall.Stat_t
	syscall.Stat(filePath, &s)

	if _, err := os.Stat(filePath); err != nil { // isFile exist?
		fmt.Println(err)
		return isTargetLogCond
	}
	if errMsg := validateFileformat(filePath); errMsg != "" {
		errorlog(errMsg)
		return isTargetLogCond
	}

	if conf[ConfKeyMap["ConfSize"]] == "*" { // check Size
		isSize = true
	} else {
		condSize, _ := strconv.ParseInt(conf[ConfKeyMap["CondSize"]], 10, 64)
		if s.Size > condSize { // KB
			isSize = true
		}
		// fmt.Printf("size: %s > %s :>> %t\n", s.Size, condSize, isSize)
	}
	if conf[ConfKeyMap["ConfTime"]] == "*" { // check Time
		isTime = true
	} else {
		condTime, _ := strconv.Atoi(conf[ConfKeyMap["CondTime"]])
		fileTime := time.Unix(s.Mtim.Unix())
		logRimitTime := time.Now().AddDate(0, 0, condTime) // https://ashitani.jp/golangtips/tips_time.html#time_Duration
		if fileTime.Before(logRimitTime) {                 // file>=logRimit --> !file<logRimit
			isTime = true
		}
		// fmt.Printf("time: %s before %s :>> %t\n", fileTime, logRimitTime, isTime)
	}
	if isSize && isTime {
		isTargetLogCond = true
	}
	return isTargetLogCond
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

func deleteLog(trgFilePath string) { // targFilePath is included filename
	debug(fmt.Sprintf("delete: %s \n", trgFilePath))
}
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

// cleanLogCmd represents the cleanLogfile command
var cleanLogCmd = &cobra.Command{
	Use:   "cleanLog",
	Short: "clean log from conf",
	Long:  "clean loger. you put conf file of 'clean_log.conf. format is following'",
	Run: func(cmd *cobra.Command, args []string) {

		// log init
		logfile, err := os.OpenFile("./cleanLog.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic("cannnot open test.log:" + err.Error())
		}
		defer logfile.Close()
		filter := &logutils.LevelFilter{
			Levels: []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
			// MinLevel: logutils.LogLevel("DEBUG"), // for debug
			MinLevel: logutils.LogLevel("WARN"),
			Writer:   logfile, //os.Stderr,
		}
		log.SetOutput(filter)
		// end log int ------------------

		// readConf
		_, confData := readConfCsv("clean_log.conf")
		warn(fmt.Sprintln("start batch"))
		for _, conf := range confData {
			debug(fmt.Sprintln(conf))
			targetPath := conf[ConfKeyMap["TargetPath"]]
			if _, err := os.Stat(targetPath); err != nil {
				errorlog(fmt.Sprintln(err))
				break
			}
			logs := getChildFiles(targetPath)
			for _, logPath := range logs { // targetPath配下のファイルをすべてチェック
				if isTargetLog(logPath, conf) { // 条件を満たすファイルが存在した場合
					switch conf[ConfKeyMap["Mode"]] {
					case "MOVE":
						debug("MOVE: " + logPath)
						moveLog(logPath, conf[ConfKeyMap["DstPath"]])
					case "ZIPMOVE":
						debug("ZIPMOVE: " + logPath)
						compressToGzip(logPath, conf[ConfKeyMap["DstPath"]])
					case "DELETE":
						debug("DELETE: " + logPath)
						deleteFile(logPath)
					}
				}
			}
		} // end loop conf
	},
}

func debug(msg string) {
	log.Print("[DEBUG] " + msg)
}
func warn(msg string) {
	log.Print("[WARN] " + msg)
}

// errorであってもerrであっても被ることが多いのでここだけはerrorlogにする
func errorlog(msg string) {
	log.Print("[ERROR] " + msg)
}
func init() {
	rootCmd.AddCommand(cleanLogCmd)
	ConfKeyMap = map[string]int{
		"TargetPath": 0,
		"DstPath":    1,
		"Mode":       2,
		"CondTime":   3,
		"CondSize":   4,
	}
}
