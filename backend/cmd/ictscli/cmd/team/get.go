package team

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-regalia/backend/cmd/ictscli/client"
	v1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/spf13/cobra"
)

var getCode int64

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a team by code",
	Long: `Get detailed information about a team by its code.
Example: ictscli team get --code 1`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		req := connect.NewRequest(&v1.GetTeamRequest{
			Code: getCode,
		})

		res, err := client.GetTeamClient().GetTeam(ctx, req)
		if err != nil {
			fmt.Printf("Failed to get team: %v\n", err)
			return
		}

		team := res.Msg.GetTeam()
		client.PrintJSON(team, func() {
			fmt.Printf("Team details:\n")
			fmt.Printf("  Code: %d\n", team.GetCode())
			fmt.Printf("  Name: %s\n", team.GetName())
			fmt.Printf("  Organization: %s\n", team.GetOrganization())
		})
	},
}

func init() {
	getCmd.Flags().Int64Var(&getCode, "code", 0, "Team code (required)")
	getCmd.MarkFlagRequired("code")
}
