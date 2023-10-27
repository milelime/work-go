/*
Copyright © 2023 Alexey Ayzin <alexey.ayzin@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

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
		logEnd()
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

func logEnd() {
	if !fileExists(SESSION_LOG) {
		handleMissingLog()
	}
	buffer := readLastLine(SESSION_LOG)

	if buffer[0] == 'E' {
		log.Fatal("A session does not exist.\nPlease start a new session.")
	}

	now := time.Now()

	end := fmt.Sprintf(
		"END: %d-%02d-%02d %02d:%02d:%02d\n",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
	)

	f, err := os.OpenFile(SESSION_LOG, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(end); err != nil {
		log.Println(err)
	}

	fmt.Printf("Session ended at %s", end)
}
