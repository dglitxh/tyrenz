/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/dglitxh/tyrenz/spoofy"
	"github.com/spf13/cobra"
)
var ext string
var paste bool
var fn string
// spoofCmd represents the spoof command
var spoofCmd = &cobra.Command{
	Use:   "spoof",
	Short: "Creates (spoofs) text files",
	Long: `A cli tool to create text files on the go.
  it is simple and allows a **paste mode**; where up to 15 consecutive blank 
  lines can be allowed.
  *NOTE: 1. normal mode does not allow blank lines
	 2. to exit paste mode you can either type "end..." (on a new line) or press the enter key 16 times.
		                       *the choice is yours*`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("spoofing.......")
		spoofy.Spoof(ext, fn, paste)
	},
}

func init() {
	rootCmd.AddCommand(spoofCmd)
	spoofCmd.Flags().StringVarP(&ext, "extension", "e", "txt", "extension for the file to be spoofed.")
	spoofCmd.Flags().BoolVarP(&paste, "paste", "p", false, "activates paste mode where up to 15 blank lines can be allowed")
	spoofCmd.Flags().StringVarP(&fn, "name", "n", "file", "name for the file to be spoofed.")
	
}
