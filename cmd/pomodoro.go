/*
Copyright © 2023 pomoCmd

*/
package cmd

import (
	"fmt"
	"time"

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
	Short: "A Simple yet powerful pomodoro timer",
	Long: `This app is a simple pomodoro timer with a simple interface. you can specify time intervals 
	      pomodoro, long and short breaks. Although there is a default duration which is widely used.`,
	Run: func(cmd *cobra.Command, args []string) {
		i := &pomodoro.Instance{}
		s := pomodoro.UserSpecs{
		LongBreak: time.Duration(longbrk),
		ShortBreak: time.Duration(shortbrk),
		Interval: time.Duration(pomo),
	}
		i.Specs = s
		inst := pomodoro.NewInstance(i, pomodoro.CatPomodoro, int(s.Interval), int(s.LongBreak), int(s.ShortBreak))
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
