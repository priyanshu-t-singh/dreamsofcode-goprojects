package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mergestat/timediff"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01-todo-list/constants"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01-todo-list/internal/csv_handler"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01-todo-list/internal/model"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01-todo-list/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.Flags().BoolP("all", "a", false, "list all tasks")
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long:  `List all tasks. Use -a or --all to list all tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		showAll := cmd.Flags().Changed("all")

		csv := csv_handler.CSVHandler{}

		tasks, err := csv.ReadFromCSV(constants.CsvFilePath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		filteredTasks := filterTasks(tasks, showAll)

		var results [][]string

		// Append the header
		if showAll {
			results = append(results, []string{"ID", "Description", "CreatedAt", "Done"})
			results = append(results, []string{"--", "-----------", "---------", "----"})
		} else {
			results = append(results, []string{"ID", "Description", "CreatedAt"})
			results = append(results, []string{"--", "-----------", "---------"})
		}

		// Append the rows
		for _, task := range filteredTasks {
			parsedTime, err := time.Parse(time.RFC3339, task.CreatedAt)
			if err != nil {
				fmt.Println("Error parsing time:", err)
				os.Exit(1)
			}

			timeDiff := timediff.TimeDiff(parsedTime)

			if showAll {
				results = append(results, []string{task.ID, task.Description, timeDiff, strconv.FormatBool(task.IsCompleted)})
			} else {
				results = append(results, []string{task.ID, task.Description, timeDiff})
			}
		}

		utils.PrintTable(results)
	},
}

func filterTasks(tasks []model.Task, showAll bool) []model.Task {
	var filteredTasks []model.Task
	for _, task := range tasks {
		if !task.IsCompleted || showAll {
			filteredTasks = append(filteredTasks, task)
		}
	}
	return filteredTasks
}
