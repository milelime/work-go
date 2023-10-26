/*
Copyright © 2023 Alexey Ayzin <alexey.ayzin@gmail.com>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var previousOffset int64 = 0

const defaultSessionLog = "session.log"

var SESSION_LOG string = defaultSessionLog

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "work-go",
	Short: "An hour logger with some display options",
	Long: `This hour logger is pretty simple. You can start and end sessions then
view them at your convenience. Funcionality includes viewing your work sessions
for this week and today. We'll also give you some totals and averages so you
can report them if you need to.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.work-go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	loadLogPath()
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func readLastLine(filename string) string {

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	lastLineSize := 0

	for {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		lastLineSize = len(line)
	}

	fileInfo, err := os.Stat(filename)

	buffer := make([]byte, lastLineSize)

	offset := fileInfo.Size() - int64(lastLineSize+1)
	numRead, err := file.ReadAt(buffer, offset)

	if previousOffset != offset {
		buffer = buffer[:numRead]
		previousOffset = offset
		return fmt.Sprintf("%s \n", buffer)
	}
	return fmt.Sprintf("%s \n", buffer)
}

// Check if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func handleMissingLog() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("The session.log file is missing.")
		fmt.Println("1. Create a new session.log file here.")
		fmt.Println("2. Locate the existing session.log file.")
		fmt.Print("Enter your choice (1/2): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			// Create a new log file
			_, err := os.Create(SESSION_LOG)
			if err != nil {
				log.Fatal("Error creating the session.log file:", err)
			}
			return
		case "2":
			// Locate the existing file
			fmt.Print("Please provide the full path to the session.log file: ")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSpace(path)
			if fileExists(path) {
				SESSION_LOG = path
				saveLogPath(path)
				return
			} else {
				fmt.Println("The provided file does not exist. Please try again.")
			}
		default:
			fmt.Println("Invalid choice. Please enter 1 or 2.")
		}
	}
}

// Save the file path internally
func saveLogPath(path string) {
	// In this example, we're saving the path to a file, but this could be
	// changed to a different method if needed, like a configuration file or a database.
	os.WriteFile("logpath.txt", []byte(path), 0644)
}

// Load the saved file path when the program starts
func loadLogPath() {
	if content, err := os.ReadFile("logpath.txt"); err == nil {
		SESSION_LOG = string(content)
	} else {
		SESSION_LOG = defaultSessionLog
	}
}
