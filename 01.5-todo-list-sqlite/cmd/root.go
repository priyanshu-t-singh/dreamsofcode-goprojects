package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01.5-todo-list-sqlite/internal/database"
	"github.com/spf13/cobra"
)

var dbConn *sql.DB
var dbPath string = "./tasks.sqlite"

var rootCmd = &cobra.Command{
	Use:   "tasks",
	Short: "tasks is a todo list app using sqlite",
	Long:  "An cli application for managing tasks in the terminal. Using sqlite as the database.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		conn, err := database.Open(dbPath)
		if err != nil {
			return err
		}
		dbConn = conn
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if dbConn != nil {
			return dbConn.Close()
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
