package main

import (
	tickspot "github.com/nicored/gotickyourself"

	"fmt"

	"strings"

	"github.com/spf13/cobra"
)

func getListCmd(tick *tickspot.Tick) *cobra.Command {
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

	for _, entry := range entries {
		projectName := "Unknown"
		taskName := "Unknown"

		if t, ok := Tasks[entry.TaskId]; ok {
			projectName = projectsConfig.Projects[t.ProjectId].Name
			taskName = t.Name
		}

		fmt.Println(entry.Date, entry.Hours, projectName, taskName, entry.TaskId)
	}
}
