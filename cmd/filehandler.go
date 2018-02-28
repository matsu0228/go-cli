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
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var dstFileName, trgFileName string

func debugMoveFileForWindows(trgFileName, dstFileName string) {
	fmt.Println("move from '" + trgFileName + "' to '" + dstFileName)
	trgFile, err := os.Open(trgFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if stat, err := os.Stat(dstFileName); err == nil {
		fmt.Println(dstFileName+" is exist.", stat)
		return
	}
	dstFile, err := os.Create(dstFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, trgFile) // windowsでは、異なるdrive間の移動ができないためコピーする
	if err != nil {
		fmt.Println(err)
		return
	}
	trgFile.Close()
	if err := os.Remove(trgFileName); err != nil {
		fmt.Println(err)
	}
}

func moveFile(srcName, dstName string) {
	if _, err := os.Stat(srcName); err != nil {
		fmt.Println(err)
		return
	}
	if stat, err := os.Stat(dstName); err == nil {
		fmt.Println(dstName+" is exist.", stat)
		return
	}
	if err := os.Rename(srcName, dstName); err != nil {
		fmt.Println(err)
	}
}

// tips
// https://qiita.com/kamol/items/fae07e8533b36f553714
// --------------------------------------------------------------
// 追記
// //引数: ファイルのパス, フラグ, パーミッション(わからなければ0666でおっけーです)
// file, err := os.OpenFile("test.txt", os.O_WRONLY|os.O_APPEND, 0666)
//
// 存在しなければ新規作成
//     file, err := os.OpenFile("test.txt", os.O_WRONLY|os.O_CREATE, 0666)

func readFile(filename string) {
	fp, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

// filehandlerCmd represents the filehandler command
var filehandlerCmd = &cobra.Command{
	Use:   "filehandler",
	Short: "file handler",
	Long:  "file handler. log cleaner",
	Run: func(cmd *cobra.Command, args []string) {
		// logs := getChildFiles("log/")
		// fileHandler ---------------------
		fmt.Println("move from '" + trgFileName + "' to '" + dstFileName)
		debugMoveFileForWindows(trgFileName, dstFileName)
		// fmt.Println("move from '" + trgFileName + "'")
		// filename := "test.txt"
		// readFile(FileName)
		// moveFile(FileName)
		// moveFile("test.txt", "test_dst.txt")
	},
}

func init() {
	filehandlerCmd.Flags().StringVarP(&trgFileName, "targget FileName", "t", "", "File name to move from")
	filehandlerCmd.Flags().StringVarP(&dstFileName, "dstFileName", "d", "", "File name to move to")
	rootCmd.AddCommand(filehandlerCmd)
}
