package main

import (
	"tickspot"

	"regexp"

	"strconv"

	"strings"
	"time"

	"fmt"

	"log"
	"os"

	"github.com/jinzhu/now"
	"github.com/spf13/cobra"
)

func getLogCmd(tick *tickspot.Tick) *cobra.Command {
	logCmd := &cobra.Command{
		Use:    "log",
		Short:  "log new entries",
		Long:   ``,
		PreRun: initConfigFiles,
		Run:    runLogCmd,
	}

	logCmd.Flags().StringP("notes", "n", "", "Add notes to the log")
	logCmd.Flags().StringP("date", "d", "today", "Is the date for which the hours must be logged (default is today). (eg. '2017-07-02', 'yesterday', 'today'")

	return logCmd
}

//$ gty log # Log all day  (default task, default time)

func runLogCmd(cmd *cobra.Command, args []string) {
	const logSpecificHours = 1
	const logPeriod = 2

	logType := logPeriod
	task := projectsConfig.DefaultTask
	logHours := -1.0

	nArgs := 0
	period := ""

	if len(args) >= 1 && isSimplePeriod(args[0]) {
		period = args[0]
		nArgs += 1
	} else if len(args) >= 2 && isPeriodWithCount(strings.Join(args[0:1], " ")) {
		period = strings.Join(args[0:1], " ")
		nArgs += 2
	} else if len(args) > 0 {
		hours, isHours := getHours(args[0])
		if isHours && hours > 0 {
			logHours = hours
			nArgs += 1
			logType = logSpecificHours
		} else if hours == 0 {
			log.Println("Error: Hours must be between 0 and 24")
			os.Exit(1)
		}
	}

	if len(args) > nArgs {
		t, ok := Alias[args[nArgs]]
		errfOnMismatch(ok, true, "Alias %s does not exist\n", args[nArgs])
		task = t
	}

	if task == nil {
		log.Println("You must specify a task, or have a default task available for use")
		os.Exit(1)
	}

	if logType == logSpecificHours || period == "" {
		var err error
		date := cmd.Flag("date").Value.String()
		date, err = parseDate(date)

		if logHours < 0 {
			entries, err := tick.GetEntries(tickspot.DateRange{date, date})
			errfOnMismatch(err, nil, "Could not load entries")
			dayHours := getTotalEntriesHours(entries)

			expectedHours := getHoursPerDay()
			if dayHours < expectedHours {
				logHours = expectedHours - dayHours
			}
		}

		entry, err := tick.CreateEntry(date, logHours, cmd.Flag("notes").Value.String(), task, true)
		errfOnMismatch(err, nil, "An error occurred when creating the entry. %s\n", err)

		fmt.Printf("Logged successfully\n\n")
		entry.Print(tick)

		os.Exit(0)
	}
}

func getHoursPerDay() float64 {
	return settingsConfig.HoursPerWeek / float64(7-len(settingsConfig.NonWorkingDays))
}

func parseDate(arg string) (string, error) {
	valid := regexp.MustCompile("(?i)^((\\d{4}-\\d{1,2}-\\d{1,2})|yesterday|today)$").MatchString(arg)
	if !valid {
		return "", fmt.Errorf("%s is not a valid date\n", arg)
	}

	if arg != "today" && arg != "yesterday" {
		return arg, nil
	}

	var date time.Time
	if arg == "today" {
		date = time.Now()
	} else if arg == "yesterday" {
		date = time.Now().Add(-24 * time.Hour)
	}

	return fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day()), nil
}

func isHours(arg string) bool {
	return regexp.MustCompile("^\\d+(\\.\\d*)?$").MatchString(arg)
}

func getHours(arg string) (float64, bool) {
	if !isHours(arg) {
		return 0.0, false
	}

	hours, err := strconv.ParseFloat(arg, 64)
	errfOnMismatch(err, nil, "Error converting %s to number", arg)

	return hours, true
}

func getTimePeriodStart(arg string) (*time.Time, bool) {
	arg = strings.ToLower(arg)

	from, isTime := getFromTime(arg)
	if isTime == true {
		return from, isTime
	}

	return getFromNumberTime(arg)
}

func isSimplePeriod(period string) bool {
	return regexp.MustCompile("(?i)^\\s*(yesterday|today|week|fortnight|month)\\s*$").MatchString(period)
}

func getFromTime(arg string) (*time.Time, bool) {
	switch arg {
	case "today":
		bDay := now.BeginningOfDay()
		return &bDay, true
	case "yesterday":
		bYesterday := now.BeginningOfDay().Add(-24 * time.Hour)
		return &bYesterday, true
	case "week":
		bWeek := now.Monday()
		return &bWeek, true
	case "fortnight":
		bFortnight := now.Monday().Add(-14 * 24 * time.Hour)
		return &bFortnight, true
	case "month":
		bMonth := now.BeginningOfMonth()
		return &bMonth, true
	}

	return nil, false
}

func isPeriodWithCount(period string) bool {
	return regexp.MustCompile("(?i)^\\s*(\\d+\\s+)(days?|weeks?|months?)\\s*$").MatchString(period)
}

func getFromNumberTime(arg string) (*time.Time, bool) {
	if !isPeriodWithCount(arg) {
		return nil, false
	}

	sep := strings.Split(arg, " ")

	countStr := sep[0]
	count, err := strconv.Atoi(countStr)
	errfOnMismatch(err, nil, "Could not convert %s to integer", countStr)

	for pi, part := range sep {
		if pi == 0 {
			continue
		}

		if strings.TrimSpace(part) == "" {
			continue
		}

		switch part {
		case "day":
			fallthrough
		case "days":
			bDays := now.BeginningOfDay().Add(-time.Duration(count) * 24 * time.Hour)
			return &bDays, true
		case "week":
			fallthrough
		case "weeks":
			bWeeks := now.Monday().Add(-time.Duration(count) * 7 * 24 * time.Hour)
			return &bWeeks, true
		case "month":
			fallthrough
		case "months":
			timeNow := time.Now()
			bMonths := time.Date(timeNow.Year(), timeNow.Month()-time.Month(count), 0, 0, 0, 0, 0, time.Local)
			return &bMonths, true
		}

		break
	}

	return nil, false
}
