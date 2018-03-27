package cmd

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/spf13/cobra"
	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

// gormCmd represents the gorm command
var gormCmd = &cobra.Command{
	Use:   "gorm",
	Short: "gorm example that is library to handle database",
	Long:  "gorm exmaple",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gorm called")
		loadDbInfo()
	},
}

func loadDbInfo() {
	cfg, err := ini.Load("db_conf.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	// Classic read of values, default section can be represented as empty string
	fmt.Println("Data Path:", cfg.Section("database").Key("db.host").String())
}

// var Gdb *gorm.DB
//
// func InitDB() {
// 	var err error
// 	connectionString := getConnectionString()
// 	Gdb, err = gorm.Open("mysql", connectionString)
// 	Gdb.LogMode(true)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// }
// func (c *GormController) ConnectDB() r.Result {
// 	c.DB = Gdb
// 	return nil
// }

func init() {
	rootCmd.AddCommand(gormCmd)
}
