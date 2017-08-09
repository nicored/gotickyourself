package main

import (
	"github.com/spf13/cobra"
)

func getSettingsCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "settings",
		Short:  "Shows all settings",
		Long:   ``,
		PreRun: initConfigFiles,
		Run:    runSettingsCmd,
	}
}

func runSettingsCmd(cmd *cobra.Command, args []string) {

}
