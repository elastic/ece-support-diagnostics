package cmd

import (
	"fmt"

	"github.com/elastic/ece-support-diagnostics/pkg/release"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of the ece-support-diagnostics",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf(
			// streams.Out,
			"ece-support-diagnostics\n\tVersion:\t%s\n\tGit Commit:\t%s\n\tBuilt:\t\t%s\n\tGo Version:\t%s\n",
			release.Version(),
			release.Commit(),
			release.BuildTime().UTC().Format("Mon Jan 02 15:04:05 2006 MST"),
			release.GoVersion(),
		)
	},
}
