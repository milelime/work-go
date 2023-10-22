/*
Copyright © 2023 Alexey Ayzin <alexey.ayzin@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// endCmd represents the end command
var endCmd = &cobra.Command{
	Use:   "end",
	Short: "End an existing session",
	Long: `Log the end of an existing session. If a session has not been started,
you will need to start a new one.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("end called")
	},
}

func init() {
	rootCmd.AddCommand(endCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// endCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// endCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
