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
// goblinCmd represents the goblin command
var goblinCmd = &cobra.Command{
	Use:   "goblin",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var proc goblin.Process
		step, err := os.ReadFile("goblinConfig.json"); if err != nil {
			fmt.Println(err)
		}
		if err := json.Unmarshal([]byte(step), &proc); err != nil {
			fmt.Println(err)
		}
		proc.Run(timeout)
	},
}

func init() {
	goblinCmd.Flags().BoolVarP(&timeout, "timeout", "t", false, "sets a timeout for eache process")
	rootCmd.AddCommand(goblinCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// goblinCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// goblinCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
