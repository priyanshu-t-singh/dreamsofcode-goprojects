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
	rootCmd.AddCommand(completeCmd)
}

var completeCmd = &cobra.Command{
	Use:   "complete <taskid>",
	Short: "Mark the task as complete",
	Long:  "This command marks the task to be completed.",
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
			fmt.Printf("Task not found of id: %v\n", index)
			os.Exit(1)
		}

		tasks[index].IsCompleted = true
		err = csv.WriteToCSV(constants.CsvFilePath, tasks)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

	},
}

func findTaskIndex(tasks []model.Task, taskid string) int {
	for i, task := range tasks {
		if task.ID == taskid {
			return i
		}
	}
	return -1
}
