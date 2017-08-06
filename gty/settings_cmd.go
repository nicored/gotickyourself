package main

import (
	tickspot "github.com/nicored/gotickyourself"

	"github.com/spf13/cobra"
)

func getSettingsCmd(tick *tickspot.Tick) *cobra.Command {
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
