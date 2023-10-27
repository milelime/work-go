/*
Copyright © 2023 Alexey Ayzin <alexey.ayzin@gmail.com>
*/package cmd

import (
	"fmt"
	"strings"

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
