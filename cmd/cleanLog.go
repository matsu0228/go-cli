// Copyright © 2017 matsuki
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
	// "compress/gzip"
	"encoding/csv"
	"fmt"
	"github.com/hashicorp/logutils"
	"github.com/spf13/cobra"
	"io"
	// "io/ioutil"
	"log"
	"os"
	// "path/filepath"
	"strconv"
	// "reflect"
	// "strings"
	"syscall"
	"time"
)

var ConfKeyMap map[string]int

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

// FileInfo: https://golang.org/pkg/syscall/#Stat_t
// isTargetLog は、ファイルが処理対象のログファイルかどうか(cond_*を満たすか)判定する
func isTargetLog(filePath string, conf []string) bool {
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
		// debug(fmt.Sprintf("size: %s > %s :>> %t\n", s.Size, condSize, isSize))
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
		// debug(fmt.Sprintf("time: %s before %s :>> %t\n", fileTime, logRimitTime, isTime))
	}
	if isSize && isTime {
		isTargetLogCond = true
	}
	return isTargetLogCond
}

// cleanLogCmd represents the cleanLogfile command
var cleanLogCmd = &cobra.Command{
	Use:   "cleanLog",
	Short: "clean log from conf",
	Long:  "clean loger. you put conf file of 'clean_log.conf. format is following'",
	Run: func(cmd *cobra.Command, args []string) {

		// log init  // TODO: define another function
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
