package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version will be populated by the compiler at build time
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of Tasks app",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Tasks app version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
