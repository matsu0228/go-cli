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
	"io/ioutil"
	"os"
	"path/filepath"
)

// filehandlerCmd represents the filehandler command
var filehandlerCmd = &cobra.Command{
	Use:   "filehandler",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename := "test.txt"
		fmt.Println(dirwalk("./cmd"))
		readFile(filename)
		moveFile("test.txt", "test_dst.txt")
	},
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

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}
	return paths
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

func init() {
	rootCmd.AddCommand(filehandlerCmd)
}
