/*
Copyright © 2023 ydzly <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dglitxh/tyrenz/goblin"
	"github.com/spf13/cobra"
)

var timeout bool
var input bool

// goblinCmd represents the goblin command
var goblinCmd = &cobra.Command{
	Use:   "goblin",
	Short: "A tool to run external programs and tasks",
	Long: `This tool runs tasks and programs listed in a user created file with the name 
	"goblinConfig.json" with an array of json map in the following format
	 {
		"name": "make directory",
		"msg": "make a directory",
		"cmd": "mkdir",
		"dir": ".",
		"args": ["temp"]
  	}`,
	Run: func(cmd *cobra.Command, args []string) {
		var proc goblin.Process
		f, err := os.OpenFile("goblinConfig.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			os.Exit(1)
		}
		proc.ReadConfig()
		if input {
			proc.ScanInput()
			j, err := json.MarshalIndent(proc, " ", " ")
			if err != nil {
				fmt.Println(err)
			}
			os.WriteFile(f.Name(), j, 0664)
		}
		proc.Run(timeout)
	},
}

func init() {
	goblinCmd.Flags().BoolVarP(&timeout, "timeout", "t", false, "sets a timeout for eache process")
	goblinCmd.Flags().BoolVarP(&input, "input", "i", false, "get input from user and skip config read.")
	rootCmd.AddCommand(goblinCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// goblinCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// goblinCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
