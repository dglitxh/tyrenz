/*
Copyright © 2023 pomoCmd

*/
package cmd

import (
	"time"

	"github.com/dglitxh/tyrenz/pomodoro"
	"github.com/dglitxh/tyrenz/pomo_utils"
	"github.com/spf13/cobra"
)

var pomo time.Duration
var shortbrk time.Duration
var longbrk time.Duration

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
		app := pomodoro.App{}

		conf, err := pomo_utils.
		app.New()
		app.Run()
	},
}

func init() {
	pomodoroCmd.Flags().DurationVarP(&pomo, "pomo", "p", 25*time.Minute,
	"Pomodoro duration")
	pomodoroCmd.Flags().DurationVarP(&shortbrk, "short", "s", 5*time.Minute,
	"Short break duration")
	pomodoroCmd.Flags().DurationVarP(&longbrk, "long", "l", 15*time.Minute,
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
