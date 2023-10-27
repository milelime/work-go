/*
Copyright © 2023 Alexey Ayzin <alexey.ayzin@gmail.com>
*/package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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

func viewTodayLogs() {
	date := time.Now().Format("2006-01-02")
	metrics := getMetricsForDate(date)

	if metrics.TotalSessions > 0 {
		totalHours, totalMinutes, totalSeconds := secondsToTime(metrics.TotalHours)
		averageHours, averageMinutes, averageSeconds := secondsToTime(metrics.AvgHours)

		fmt.Printf(
			"You completed %d sessions today.\nYou totalled %s time.\nYou averaged %s time per session.\n",
			metrics.TotalSessions,
			formatTime(totalHours, totalMinutes, totalSeconds),
			formatTime(averageHours, averageMinutes, averageSeconds),
		)
	}
}

func viewWeekLogs() {
	now := time.Now()

	fmt.Println("Date       | Total Sessions | Total Hours  | Average Hours")
	fmt.Println("----------------------------------------------------------")

	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		metrics := getMetricsForDate(date)
		if metrics.TotalSessions > 0 {
			totalHours, totalMinutes, totalSeconds := secondsToTime(metrics.TotalHours)
			avgHours, avgMinutes, avgSeconds := secondsToTime(metrics.AvgHours)
			fmt.Printf(
				"%s | %14d | %12s | %13s\n",
				date,
				metrics.TotalSessions,
				formatTime(totalHours, totalMinutes, totalSeconds),
				formatTime(avgHours, avgMinutes, avgSeconds),
			)
		}
	}
}

func secondsToTime(seconds int) (int, int, int) {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	sec := seconds % 60
	return hours, minutes, sec
}

func formatTime(hours, minutes, seconds int) string {
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func getMetricsForDate(targetDate string) SessionMetrics {
	lines := readLogFile(SESSION_LOG)

	metrics := SessionMetrics{}
	var sessionStart int

	for _, line := range lines {
		action, date, timeStr := parseLogData(line)
		timeInSeconds := getTimeInSeconds(timeStr)

		if action == "START:" && date == targetDate {
			sessionStart = timeInSeconds
		}

		if action == "END:" && date == targetDate {
			metrics.TotalHours += timeInSeconds - sessionStart
			metrics.TotalSessions++
		}
	}

	if metrics.TotalSessions > 0 {
		metrics.AvgHours = metrics.TotalHours / metrics.TotalSessions
	}

	return metrics
}

func parseLogData(line string) (action, date, timeStr string) {
	data := strings.Split(line, " ")
	return data[0], data[1], data[2]
}

func getTimeInSeconds(timeStr string) int {
	timeData := strings.Split(timeStr, ":")
	hour, _ := strconv.Atoi(timeData[0])
	minute, _ := strconv.Atoi(timeData[1])
	second, _ := strconv.Atoi(timeData[2])
	return (3600 * hour) + (60 * minute) + second
}

func readLogFile(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
