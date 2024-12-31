package team

import (
	"github.com/spf13/cobra"
)

var TeamCmd = &cobra.Command{
	Use:   "team",
	Short: "Manage teams",
	Long: `Manage teams in the ICTSC Regalia system.

This command provides functionality to create, list, get, update, and delete teams.
Each team has a unique code, name, and organization.`,
}

func init() {
	TeamCmd.AddCommand(createCmd)
	TeamCmd.AddCommand(listCmd)
	TeamCmd.AddCommand(getCmd)
	TeamCmd.AddCommand(updateCmd)
	TeamCmd.AddCommand(deleteCmd)
}
