package cmd

import (
	"fmt"
	"os"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01-todo-list/constants"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01-todo-list/internal/csv_handler"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01-todo-list/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete <taskid>",
	Short: "Delete a task with taskid",
	Long:  "This command delete the given task with taskid.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskid := args[0]

		csv := csv_handler.CSVHandler{}
		tasks, err := csv.ReadFromCSV(constants.CsvFilePath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var index = findTaskIndex(tasks, taskid)
		if index == -1 {
			fmt.Printf("Task not found of id: %v\n", taskid)
			os.Exit(1)
		}

		tasks = removeItem(tasks, index)
		err = csv.WriteToCSV(constants.CsvFilePath, tasks)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

	},
}

func removeItem(arr []model.Task, index int) []model.Task {
	if index < 0 || index >= len(arr) {
		return arr // Return original array if index is out of bounds
	}

	newArr := make([]model.Task, len(arr)-1)
	copy(newArr, arr[:index]) // Copy elements before the index

	// Only copy if there is something after `index`
	if index < len(arr)-1 {
		copy(newArr[index:], arr[index+1:]) // Copy elements after the index
	}

	return newArr
}
