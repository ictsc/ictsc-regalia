package team

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-regalia/backend/cmd/ictscli/client"
	v1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var (
	updateCode         int64
	updateName         string
	updateOrganization string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a team",
	Long: `Update a team's information. Only specified fields will be updated.
Example: ictscli team update --code 1 --name "新トラブルシューターズ"`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		paths := []string{}
		if cmd.Flags().Changed("name") {
			paths = append(paths, "name")
		}
		if cmd.Flags().Changed("organization") {
			paths = append(paths, "organization")
		}

		req := connect.NewRequest(&v1.UpdateTeamRequest{
			Team: &v1.Team{
				Code:         updateCode,
				Name:         updateName,
				Organization: updateOrganization,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: paths,
			},
		})

		res, err := client.GetTeamClient().UpdateTeam(ctx, req)
		if err != nil {
			fmt.Printf("Failed to update team: %v\n", err)
			return
		}

		team := res.Msg.GetTeam()
		client.PrintJSON(team, func() {
			fmt.Printf("Successfully updated team:\n")
			fmt.Printf("  Code: %d\n", team.GetCode())
			fmt.Printf("  Name: %s\n", team.GetName())
			fmt.Printf("  Organization: %s\n", team.GetOrganization())
		})
	},
}

func init() {
	updateCmd.Flags().Int64Var(&updateCode, "code", 0, "Team code (required)")
	updateCmd.Flags().StringVar(&updateName, "name", "", "New team name")
	updateCmd.Flags().StringVar(&updateOrganization, "organization", "", "New team organization")
	updateCmd.MarkFlagRequired("code")
}
