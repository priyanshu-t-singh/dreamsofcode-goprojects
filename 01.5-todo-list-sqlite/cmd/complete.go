package cmd

import (
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01.5-todo-list-sqlite/internal/tasks"
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete <taskid>",
	Short: "Mark the task as complete",
	Long:  "This command marks the task to be completed.",
	RunE: func(cmd *cobra.Command, args []string) error {
		taskId := args[0]
		tasksRepo := tasks.NewRepository(dbConn)
		return tasksRepo.MarkAsCompleted(taskId)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
