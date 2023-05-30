/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (

	"github.com/dglitxh/tyrenz/spydey"
	"github.com/spf13/cobra"
)

// spydeyCmd represents the spydey command
var spydeyCmd = &cobra.Command{
	Use:   "spydey",
	Short: "A powerful file system crawler.",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		spydey.Search(args[0])
	},
}

func init() {
	rootCmd.AddCommand(spydeyCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// spydeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// spydeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
