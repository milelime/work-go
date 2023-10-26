/*
Copyright © 2023 Alexey Ayzin <alexey.ayzin@gmail.com>
*/package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new session",
	Long: `Start logging a new session. If one has already started, you will
need to end it.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")
		logStart()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func logStart() {
	buffer := readLastLine(SESSION_LOG)

	if buffer[0] == 'S' {
		log.Fatal("A session has already been started.\nPlease end it before creating a new one.")
	}

	now := time.Now()

	start := fmt.Sprintf(
		"START: %d-%02d-%02d %02d:%02d:%02d\n",
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
	if _, err := f.WriteString(start); err != nil {
		log.Println(err)
	}
}
