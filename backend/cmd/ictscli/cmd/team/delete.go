package team

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-regalia/backend/cmd/ictscli/client"
	v1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/spf13/cobra"
)

var deleteCode int64

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a team",
	Long: `Delete a team by its code.
Example: ictscli team delete --code 1`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		req := connect.NewRequest(&v1.DeleteTeamRequest{
			Code: deleteCode,
		})

		_, err := client.GetTeamClient().DeleteTeam(ctx, req)
		if err != nil {
			fmt.Printf("Failed to delete team: %v\n", err)
			return
		}

		client.PrintJSON(map[string]interface{}{
			"success": true,
			"code":    deleteCode,
			"message": "Team deleted successfully",
		}, func() {
			fmt.Printf("Successfully deleted team with code: %d\n", deleteCode)
		})
	},
}

func init() {
	deleteCmd.Flags().Int64Var(&deleteCode, "code", 0, "Team code (required)")
	deleteCmd.MarkFlagRequired("code")
}
