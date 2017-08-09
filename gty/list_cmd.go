package main

import (
	"strings"

	"github.com/spf13/cobra"
)

func getListCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "ls",
		Short:  "List logs",
		Long:   ``,
		PreRun: initConfigFiles,
		Run:    runListCmd,
	}
}

func runListCmd(cmd *cobra.Command, args []string) {
	argsStr := strings.ToLower(strings.TrimSpace(strings.Join(args, " ")))
	if argsStr == "" {
		argsStr = "today"
	}

	entries, err := tick.GetEntries(getDateRange(argsStr))
	errfOnMismatch(err, nil, "Could not get entries. %s", err)

	tick.PrintEntries(entries)
}
