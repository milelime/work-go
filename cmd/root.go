/*
Copyright © 2023 Alexey Ayzin <alexey.ayzin@gmail.com>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var previousOffset int64 = 0

const SESSION_LOG = "session.log"

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
