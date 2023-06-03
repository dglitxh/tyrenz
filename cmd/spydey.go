/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/dglitxh/tyrenz/spydey"
	"github.com/spf13/cobra"
)

var allow bool
var dir string
// spydeyCmd represents the spydey command
var spydeyCmd = &cobra.Command{
	Use:   "spydey",
	Short: "A powerful file system crawler.",
	Long: `A cutting edge file system crawling cli applications that can 
	create, search, crawl and many more 
  **actions include [crawl, find]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			switch args[0] {
				case "find": 
					if err := checkArgs(args, "please add filename to search"); err == nil {
						spydey.Find(args[1], dir)
					} else {
						fmt.Println(err)
					}
				case "crawl":
					spydey.Crawl(allow, dir)
				default:
					fmt.Println(" Please add a valid action.")
		}
		} else {
			fmt.Println("Please add an action \n  **actions include [crawl, find]")
		}
	},
}

func init() {
	spydeyCmd.Flags().BoolVarP(&allow, "allow_hidden", "a", false, "crawl hidden files and dirs.")
	spydeyCmd.Flags().StringVarP(&dir, "directory", "d", ".", "add a root directory for operations.")
	rootCmd.AddCommand(spydeyCmd)
	
}
