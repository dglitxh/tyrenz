/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/dglitxh/trentz/spoofy"
	"github.com/spf13/cobra"
)
var ext string
var paste bool
// spoofCmd represents the spoof command
var spoofCmd = &cobra.Command{
	Use:   "spoof",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("spoofing.......")
		spoofy.Create(ext, paste)
	},
}

func init() {
	rootCmd.AddCommand(spoofCmd)
	spoofCmd.Flags().StringVarP(&ext, "extension", "e", "txt", "extension for the file to be spoofed.")
	spoofCmd.Flags().BoolVarP(&paste, "paste", "p", false, "are you pasting the text?")
}
