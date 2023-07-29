/*
Copyright © 2023 pomoCmd

*/
package cmd

import (
	"fmt"

	"github.com/dglitxh/tyrenz/helpers"
	"github.com/dglitxh/tyrenz/pomodoro"
	"github.com/spf13/cobra"
)

var pomo int
var shortbrk int
var longbrk int

// pomodoroCmd represents the pomodoro command
var pomodoroCmd = &cobra.Command{
	Use:   "pomodoro",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		i := &pomodoro.Instance{}
		inst := pomodoro.NewInstance(i, pomodoro.CatPomodoro, int(pomo), int(longbrk), int(shortbrk))
		
		app, err := inst.New()
		if err != nil {
			helpers.Logger("default state for start error please help me.")
			fmt.Println(err)
		}
		app.Run()
	},
}

func init() {
	pomodoroCmd.Flags().IntVarP(&pomo, "pomo", "p", 25,
	"Pomodoro duration")
	pomodoroCmd.Flags().IntVarP(&shortbrk, "short", "s", 5,
	"Short break duration")
	pomodoroCmd.Flags().IntVarP(&longbrk, "long", "l", 15,
	"Long break duration")
	rootCmd.AddCommand(pomodoroCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pomodoroCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pomodoroCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
