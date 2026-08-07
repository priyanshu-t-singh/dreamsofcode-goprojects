package cmd

import (
	"fmt"
	"strconv"

	"github.com/mergestat/timediff"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01.5-todo-list-sqlite/internal/tasks"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01.5-todo-list-sqlite/internal/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long:  `List all tasks. Use -a or --all to list all tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showAll := cmd.Flags().Changed("all")
		tasksRepo := tasks.NewRepository(dbConn)
		tasks, err := tasksRepo.GetAll(showAll)
		if err != nil {
			return fmt.Errorf("error fetching tasks: %w", err)
		}

		var results [][]string

		// Append the header
		if showAll {
			results = append(results, []string{"ID", "Description", "CreatedAt", "Due", "Done"})
			results = append(results, []string{"--", "-----------", "---------", "---", "----"})
		} else {
			results = append(results, []string{"ID", "Description", "CreatedAt", "Due"})
			results = append(results, []string{"--", "-----------", "---------", "---"})
		}

		// Append the rows
		for _, task := range tasks {
			timeDiff := timediff.TimeDiff(task.CreatedAt)
			var dueTimeDiff string = "-"
			if task.DueDate != nil {
				dueTimeDiff = timediff.TimeDiff(*task.DueDate)
			}

			if showAll {
				results = append(results, []string{task.ID, task.Title, timeDiff, dueTimeDiff, strconv.FormatBool(task.CompletedAt != nil)})
			} else {
				results = append(results, []string{task.ID, task.Title, timeDiff, dueTimeDiff})
			}
		}

		utils.PrintTable(results)
		return nil
	},
}

func init() {
	listCmd.Flags().BoolP("all", "a", false, "list all tasks")
	rootCmd.AddCommand(listCmd)
}
