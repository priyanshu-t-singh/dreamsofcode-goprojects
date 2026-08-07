package cmd

import (
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01.5-todo-list-sqlite/internal/tasks"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <taskid>",
	Short: "Delete a task with taskid",
	Long:  "This command delete the given task with taskid.",
	RunE: func(cmd *cobra.Command, args []string) error {
		taskId := args[0]
		tasksRepo := tasks.NewRepository(dbConn)
		return tasksRepo.Delete(taskId)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
