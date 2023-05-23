/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/dglitxh/trentz/todo"
	"github.com/spf13/cobra"
)
var complete bool
var description string
var tl todo.TodoList
var todoFn string

func checkArgs (args []string, err string) error {
	if len(args) < 2 {
		return errors.New(err)
	}
	return nil
}
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "A todo cli with multiple actions",
	Long: ` This todo command has multiple actions to manage tasks
these actions include....
	- add (requires a todo name/title)
	- list 
	- done (requires an id)
	- get  (requires an id)
	- delete (requires an id)`,	

	Run: func(cmd *cobra.Command, args []string) {
		tl.ReadTodo()
		if len(args) > 0 {
			switch args[0] {
				case "add": 
					tl.AddTodo(complete)
				case "list":
					tl.ListTodo()
				case "get":
					if err := checkArgs(args, "please add an index."); err == nil {
						tl.GetTodo(args[1])
					} else {
						fmt.Println(err)
					}
				case "delete":
					if err := checkArgs(args, "please add an index."); err == nil {
						tl.DeleteTodo(args[1])
					} else {
						fmt.Println(err)
					}
				case "complete":
					if err := checkArgs(args, "please add an index."); err == nil {
						tl.ToggleComplete(args[1])
					}else {
						fmt.Println(err)
					}
				case "edit":
					if err := checkArgs(args, "please add an index."); err == nil {
						tl.EditTodo(args[1])
					}else {
						fmt.Println(err)
					}
				default:
					fmt.Println(" Please add a valid action.")
		}
		} else {
			fmt.Println("Please add an action \n  **actions include [add, list, delete, get]")
		}
		
	},
}


func init() {
	
	rootCmd.AddCommand(todoCmd)
	todoCmd.Flags().BoolVarP(&complete, "complete", "c", false, "task is done?")
	todoCmd.Flags().StringVarP(&description, "description", "d", "", "describes the task")
	todoCmd.Flags().StringVarP(&todoFn, "name", "n", "todo", "name for the file to be spoofed.")


}
