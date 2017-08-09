package main

import (
	"tickspot"

	"github.com/spf13/cobra"
)

func getRolesCmd(tick *tickspot.Tick) *cobra.Command {
	return &cobra.Command{
		Use:    "roles",
		Short:  "Resets all settings",
		Long:   ``,
		PreRun: initConfigFiles,
		Run:    runRolesCmd,
	}
}

func runRolesCmd(cmd *cobra.Command, args []string) {

}
