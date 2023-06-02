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
// spydeyCmd represents the spydey command
var spydeyCmd = &cobra.Command{
	Use:   "spydey",
	Short: "A powerful file system crawler.",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			switch args[0] {
				case "find": 
					if err := checkArgs(args, "please add filename to search"); err == nil {
						spydey.Find(args[1])
					} else {
						fmt.Println(err)
					}
				case "crawl":
					spydey.Crawl(allow)
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
	rootCmd.AddCommand(spydeyCmd)
	
}
