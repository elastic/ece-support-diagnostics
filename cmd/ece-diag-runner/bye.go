package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// byeCmd represents the bye command
var byeCmd = &cobra.Command{
	Use:   "bye",
	Short: "says bye",
	Long:  `This subcommand says bye`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bye called")
		fmt.Printf("%+v\n", cmd)
	},
}

func init() {
	// rootCmd.AddCommand(byeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// byeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// byeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
