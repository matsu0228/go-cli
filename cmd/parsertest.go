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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var echoTimes int

// parsertestCmd represents the parsertest command
var parsertestCmd = &cobra.Command{
	Use:   "parsertest",
	Short: "sample of argment parser",
	Long: `example of argment parser:
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < echoTimes; i++ {
			fmt.Println("Print: " + strings.Join(args, " "))
		}
	},
}

func init() {
	parsertestCmd.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to excute print")
	rootCmd.AddCommand(parsertestCmd)
}
