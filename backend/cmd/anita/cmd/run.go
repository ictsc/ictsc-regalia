package cmd

import (
	"log"

	"github.com/defval/di"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/repository"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/repository/bun"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/server"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/service"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/service/impl"
	connectServer "github.com/ictsc/ictsc-outlands/backend/pkg/connect/server"
	"github.com/ictsc/ictsc-outlands/backend/pkg/db/rdb"
	rdbBun "github.com/ictsc/ictsc-outlands/backend/pkg/db/rdb/bun"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the connect server",
	Run: func(_ *cobra.Command, _ []string) {
		container, err := di.New(
			di.Provide(rdbBun.NewDB, di.As(new(rdb.Tx))),
			di.Provide(bun.NewUserRepository, di.As(new(repository.UserRepository))),
			di.Provide(bun.NewTeamRepository, di.As(new(repository.TeamRepository))),

			di.Provide(impl.NewUserService, di.As(new(service.UserService))),
			di.Provide(impl.NewTeamService, di.As(new(service.TeamService))),

			di.Provide(server.NewServer),
		)
		if err != nil {
			log.Panic(err)
		}

		if err = container.Invoke(connectServer.Start); err != nil {
			log.Panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
