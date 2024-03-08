package cmd

import (
	"context"
	"log"

	"github.com/defval/di"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/repository"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/repository/bun"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/server"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/service"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/service/impl"
	"github.com/ictsc/ictsc-outlands/backend/internal/proto/anita/v1/v1connect"
	connectServer "github.com/ictsc/ictsc-outlands/backend/pkg/connect/server"
	"github.com/ictsc/ictsc-outlands/backend/pkg/db/rdb"
	rdbBun "github.com/ictsc/ictsc-outlands/backend/pkg/db/rdb/bun"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Connect server",
	Run: func(cmd *cobra.Command, args []string) {
		container, err := di.New(
			di.ProvideValue(cmd.Context(), di.As(new(context.Context))),

			di.ProvideValue(&config),
			di.Provide(provideRDBConfig),
			di.Provide(provideServerConfig),

			di.Provide(rdbBun.NewDB, di.As(new(rdb.Tx))),
			di.Provide(bun.NewUserRepository, di.As(new(repository.UserRepository))),
			di.Provide(bun.NewTeamRepository, di.As(new(repository.TeamRepository))),

			di.Provide(impl.NewUserService, di.As(new(service.UserService))),
			di.Provide(impl.NewTeamService, di.As(new(service.TeamService))),

			di.Provide(server.NewUserServiceHandler, di.As(new(v1connect.UserServiceHandler))),
			di.Provide(server.NewTeamServiceHandler, di.As(new(v1connect.TeamServiceHandler))),

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
