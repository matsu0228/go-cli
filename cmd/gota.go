package cmd

import (
	"fmt"
	"github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
	"github.com/spf13/cobra"
	"os"
)

var FileName string

// gotaCmd represents the gota command
var gotaCmd = &cobra.Command{
	Use:   "gota",
	Short: "gota example that is library to handle data with dataframe",
	Long:  "gota exmaple",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gota called")
		// readbuffer()
		df := readcsv()
		// filter
		fil := df.Filter(
			dataframe.F{
				Colname:    "Age",
				Comparator: series.Greater,
				Comparando: 30,
			},
		)
		fmt.Println(fil)
	},
}

// https://godoc.org/github.com/kniren/gota/dataframe#ReadCSV
func readcsv() dataframe.DataFrame {
	fmt.Println("open file of '" + FileName + "' ")

	fp, err := os.Open(FileName)
	if err != nil {
		fmt.Println(err)
	}
	df := dataframe.ReadCSV(fp)
	defer fp.Close()
	return df
}

func init() {
	gotaCmd.Flags().StringVarP(&FileName, "filename", "f", "", "File name to read from")
	rootCmd.AddCommand(gotaCmd)
}
