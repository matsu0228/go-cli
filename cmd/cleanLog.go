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
*/
func readConfCsv(confFile string) ([]string, [][]string) {
	var confData [][]string
	warn("open conf file of '" + confFile + "' ")
	fp, err := os.Open(confFile)
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(fp)
	defer fp.Close()
	header, _ := reader.Read()
	for {
		record, err := reader.Read()
		// fmt.Println(reflect.TypeOf(record))
		if len(record) != len(ConfKeyMap) {
			fmt.Printf("set %s colmuns. this line data is following\n", len(ConfKeyMap))
			fmt.Println(record)
		} else {
			confData = append(confData, record)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
	}
	return header, confData
}

func getChildFiles(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
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
	// var oldFiles []string
	isTargetLogCond := false
	isSize, isTime := false, false // 各条件の判定
	var s syscall.Stat_t
	syscall.Stat(filePath, &s)

	if _, err := os.Stat(filePath); err != nil { // isFile exist?
		fmt.Println(err)
		return isTargetLogCond
	}
	if errMsg := validateFileformat(filePath); errMsg != "" {
		fmt.Println(errMsg)
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

func moveLog(trgFilePath, dstPath string) { // targFilePath is included filename
	fmt.Printf("move: %s to %s \n", trgFilePath, dstPath)

	pos := strings.LastIndex(trgFilePath, "/")
	targetName := trgFilePath[pos+1:]
	dstFilePath := strings.TrimRight(dstPath, "/") + "/" + targetName
	if err := os.Rename(trgFilePath, dstFilePath); err != nil {
		fmt.Println(err)
	}
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

// cleanLogCmd represents the filehandler command
var cleanLogCmd = &cobra.Command{
	Use:   "cleanLog",
	Short: "clean log from conf",
	Long:  "clean loger. you put conf file of 'clean_log.conf. format is following'",
	Run: func(cmd *cobra.Command, args []string) {
		// readConf
		_, confData := readConfCsv("clean_log.conf")
		for _, conf := range confData {
			debug(fmt.Sprintln(conf))
			warn(fmt.Sprintln(conf))
			targetPath := conf[ConfKeyMap["TargetPath"]]
			if _, err := os.Stat(targetPath); err != nil {
				errorlog(fmt.Sprintln(err))
				break
			}
			logs := getChildFiles(targetPath)
			for _, logPath := range logs { // targetPath配下のファイルをすべてチェック
				// fmt.Println(logPath)
				if isTargetLog(logPath, conf) { // 条件を満たすファイルが存在した場合
					switch conf[ConfKeyMap["Mode"]] {
					case "MOVE":
						debug("called MOVE")
						moveLog(logPath, conf[ConfKeyMap["DstPath"]])
					case "DELETE":
						debug("DELETE")
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
}
