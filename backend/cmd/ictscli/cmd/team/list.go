package team

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-regalia/backend/cmd/ictscli/client"
	v1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all teams",
	Long: `List all registered teams with their details.
Example: ictscli team list`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		req := connect.NewRequest(&v1.ListTeamsRequest{})

		res, err := client.GetTeamClient().ListTeams(ctx, req)
		if err != nil {
			fmt.Printf("Failed to list teams: %v\n", err)
			return
		}

		teams := res.Msg.GetTeams()
		client.PrintJSON(teams, func() {
			if len(teams) == 0 {
				fmt.Println("No teams found")
				return
			}

			fmt.Println("Teams:")
			for _, team := range teams {
				fmt.Printf("  Code: %d\n", team.GetCode())
				fmt.Printf("  Name: %s\n", team.GetName())
				fmt.Printf("  Organization: %s\n", team.GetOrganization())
				fmt.Println()
			}
		})
	},
}

func init() {
}
