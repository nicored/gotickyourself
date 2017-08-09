package main

import (
	"github.com/spf13/cobra"
)

func getRolesCmd() *cobra.Command {
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
