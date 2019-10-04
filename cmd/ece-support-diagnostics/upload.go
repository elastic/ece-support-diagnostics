package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/elastic/ece-support-diagnostics/uploader"
	"github.com/spf13/cobra"
)

var uploadID string

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&uploadID, "upload-id", "u", "", "The Upload ID (required)")
	uploadCmd.MarkFlagRequired("upload-id")
}

var uploadCmd = &cobra.Command{
	Use:   "upload {filename}",
	Short: "upload a single file",
	Long:  `Upload an already created diagnostic or log file to Elastic Support`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			printOptionsError("[!] Path to the file is missing")
			cmd.Usage()
			os.Exit(1)
		} else if len(args) > 1 {
			printOptionsError("[!] Please pass only a single file to the command")
			cmd.Usage()
			os.Exit(1)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// cobra.ExactArgs(1) validates that there should be only one file
		filepath := args[0]
		// 3. File does not exist
		if _, err := os.Stat(filepath); err != nil {
			if os.IsNotExist(err) {
				printOptionsError("[!] File %q does not exist", filepath)
				cmd.Usage()
				os.Exit(1)
			} else {
				printOptionsError("[!] Error reading file %q: %s", filepath, err)
				cmd.Usage()
				os.Exit(1)
			}
		}
		uploader.RunUpload(filepath, uploadID)
	},
}

func printOptionsError(format string, a ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf(format, a...)
	fmt.Println(strings.Repeat("-", 80))
}
