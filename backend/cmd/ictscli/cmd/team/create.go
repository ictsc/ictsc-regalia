package team

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-regalia/backend/cmd/ictscli/client"
	v1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/spf13/cobra"
)

var (
	createCode         int64
	createName         string
	createOrganization string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new team",
	Long: `Create a new team with the specified code, name and organization.
Example: ictscli team create --code 1 --name "トラブルシューターズ" --organization "ICTSC"`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		req := connect.NewRequest(&v1.CreateTeamRequest{
			Team: &v1.Team{
				Code:         createCode,
				Name:         createName,
				Organization: createOrganization,
			},
		})

		res, err := client.GetTeamClient().CreateTeam(ctx, req)
		if err != nil {
			fmt.Printf("Failed to create team: %v\n", err)
			return
		}

		team := res.Msg.GetTeam()
		client.PrintJSON(team, func() {
			fmt.Printf("Successfully created team:\n")
			fmt.Printf("  Code: %d\n", team.GetCode())
			fmt.Printf("  Name: %s\n", team.GetName())
			fmt.Printf("  Organization: %s\n", team.GetOrganization())
		})
	},
}

func init() {
	createCmd.Flags().Int64Var(&createCode, "code", 0, "Team code (required)")
	createCmd.Flags().StringVar(&createName, "name", "", "Team name (required)")
	createCmd.Flags().StringVar(&createOrganization, "organization", "", "Team organization (required)")
	createCmd.MarkFlagRequired("code")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("organization")
}
