package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information for ece-support-diagnostics",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ece-support-diagnostics")
		fmt.Println(VERSION)
	},
}
