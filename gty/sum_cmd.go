package main

import (
	"strings"

	tickspot "github.com/nicored/gotickyourself"

	"fmt"

	"github.com/spf13/cobra"
)

func getSumCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "sum",
		Short:  "Log summary",
		Long:   ``,
		PreRun: initConfigFiles,
		Run:    runSumCmd,
	}
}

func runSumCmd(cmd *cobra.Command, args []string) {
	nArgs := 0
	period := ""
	if len(args) >= 1 && isSimplePeriod(args[0]) {
		period = args[0]
		nArgs += 1
	} else if len(args) >= 2 && isPeriodWithCount(strings.Join(args[0:2], " ")) {
		period = strings.Join(args[0:2], " ")
		nArgs += 2
	} else {
		period = "today"
	}

	entries, err := tick.GetEntries(getDateRange(period))
	errfOnMismatch(err, nil, "Could not get entries. %s", err)

	totalHours := getTotalEntriesHours(entries)

	fmt.Println("Total number of hours: ", totalHours)
}

func getTotalEntriesHours(entries []*tickspot.Entry) float64 {
	countHours := 0.0
	for _, entry := range entries {
		countHours += entry.Hours
	}

	return countHours
}
