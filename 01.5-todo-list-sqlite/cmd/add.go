package cmd

import (
	"fmt"
	"time"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01.5-todo-list-sqlite/internal/tasks"
	"github.com/spf13/cobra"
)

var dueDateStr string

var addCmd = &cobra.Command{
	Use:     "add <title>",
	Short:   "Add a new task",
	Long:    "This `add` command adds a new task to the task list",
	Aliases: []string{"create"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var due *string
		if dueDateStr != "" {
			_, err := time.Parse("2006-01-02", dueDateStr)
			if err != nil {
				return fmt.Errorf("invalid due date, expected YYYY-MM-DD: %w", err)
			}
			due = &dueDateStr
		}

		description := args[0]
		tasksRepo := tasks.NewRepository(dbConn)
		_, err := tasksRepo.Create(description, due)

		return err
	},
}

func init() {
	addCmd.Flags().StringVar(&dueDateStr, "due", "", "due date (YYYY-MM-DD), optional")
	rootCmd.AddCommand(addCmd)
}
