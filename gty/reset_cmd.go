package main

import (
	tickspot "github.com/nicored/gotickyourself"

	"github.com/spf13/cobra"
)

func getResetCmd(tick *tickspot.Tick) *cobra.Command {
	return &cobra.Command{
		Use:    "reset",
		Short:  "Resets all settings",
		Long:   ``,
		PreRun: initConfigFiles,
		Run:    runResetCmd,
	}
}

func runResetCmd(cmd *cobra.Command, args []string) {

}
