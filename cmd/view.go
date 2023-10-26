/*
Copyright © 2023 Alexey Ayzin <alexey.ayzin@gmail.com>
*/package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View your sessions",
	Long: `View your sessions in either weekly or daily formats. You can also
view the important metrics like total hours, total sessions, hours per session,
hours per day, and some other metrics.`,
	Run: func(cmd *cobra.Command, args []string) {
		todayFlag, _ := cmd.Flags().GetBool("today")
		weekFlag, _ := cmd.Flags().GetBool("week")

		if todayFlag {
			viewLog("today")
		} else if weekFlag {
			viewLog("week")
		} else {
			fmt.Println("Please specify a flag, either --today (-t) or --week (-w).")
		}
	},
}

func init() {
	viewCmd.Flags().BoolP("today", "t", false, "View logs for today")
	viewCmd.Flags().BoolP("week", "w", false, "View logs for the past week")
	rootCmd.AddCommand(viewCmd)
}

type SessionMetrics struct {
	TotalSessions int
	TotalHours    int
	AvgHours      int
}

func viewLog(mode string) {
	if !fileExists(SESSION_LOG) {
		handleMissingLog()
	}
	buffer := readLastLine(SESSION_LOG)
	if strings.HasPrefix(buffer, "S") {
		fmt.Println("There is a session ongoing.")
		return
	}

	switch mode {
	case "today":
		viewTodayLogs()
	case "week":
		viewWeekLogs()
	default:
		fmt.Println("Invalid mode.")
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
