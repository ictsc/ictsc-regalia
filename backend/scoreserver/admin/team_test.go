package admin_test

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestAdminTeamService_Create(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in *adminv1.CreateTeamRequest

		wants    *adminv1.CreateTeamResponse
		wantCode connect.Code
	}{
		"ok": {
			in: &adminv1.CreateTeamRequest{
				Team: &adminv1.Team{Code: 3, Name: "トラブルバスターズ", Organization: "ICTSC Committee", MemberLimit: 1},
			},
			wants: &adminv1.CreateTeamResponse{
				Team: &adminv1.Team{Code: 3, Name: "トラブルバスターズ", Organization: "ICTSC Committee", MemberLimit: 1},
			},
		},

		"invalid code": {
			in: &adminv1.CreateTeamRequest{
				Team: &adminv1.Team{Code: -1, Name: "below zero", Organization: "ICTSC Committee"},
			},
			wantCode: connect.CodeInvalidArgument,
		},
		"duplicate code": {
			in: &adminv1.CreateTeamRequest{
				Team: &adminv1.Team{Code: 1, Name: "duplicator", Organization: "ICTSC Committee", MemberLimit: 2},
			},
			wantCode: connect.CodeAlreadyExists,
		},

		"empty name": {
			in: &adminv1.CreateTeamRequest{
				Team: &adminv1.Team{Code: 3, Name: "", Organization: "ICTSC Committee"},
			},
			wantCode: connect.CodeInvalidArgument,
		},
		"duplicate name": {
			in: &adminv1.CreateTeamRequest{
				Team: &adminv1.Team{Code: 3, Name: "トラブルシューターズ", Organization: "ICTSC Committee", MemberLimit: 2},
			},
			wantCode: connect.CodeAlreadyExists,
		},

		"empty organization": {
			in: &adminv1.CreateTeamRequest{
				Team: &adminv1.Team{Code: 3, Name: "empty org", Organization: ""},
			},
			wantCode: connect.CodeInvalidArgument,
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := setupTeamService(t)

			resp, err := client.CreateTeam(ctx, connect.NewRequest(tt.in))
			assertCode(t, tt.wantCode, err)
			if err != nil {
				return
			}

			if diff := cmp.Diff(
				resp.Msg, tt.wants,
				cmpopts.IgnoreUnexported(adminv1.CreateTeamResponse{}, adminv1.Team{}),
			); diff != "" {
				t.Errorf("(-got, +want)\n%s", diff)
			}
		})
	}
}

func TestAdminTeamService_List(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in *adminv1.ListTeamsRequest

		wants    *adminv1.ListTeamsResponse
		wantCode connect.Code
	}{
		"ok": {
			in: &adminv1.ListTeamsRequest{},

			wants: &adminv1.ListTeamsResponse{
				Teams: []*adminv1.Team{
					{Code: 1, Name: "トラブルシューターズ", Organization: "ICTSC Association", MemberLimit: 6},
					{Code: 2, Name: "トラブルメイカーズ", Organization: "ICTSC Association", MemberLimit: 2},
				},
			},
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := setupTeamService(t)

			resp, err := client.ListTeams(ctx, connect.NewRequest(tt.in))
			assertCode(t, tt.wantCode, err)
			if err != nil {
				return
			}

			if diff := cmp.Diff(
				resp.Msg, tt.wants,
				cmpopts.IgnoreUnexported(adminv1.ListTeamsResponse{}, adminv1.Team{}),
			); diff != "" {
				t.Errorf("(-got, +want)\n%s", diff)
			}
		})
	}
}

func TestAdminTeamService_Get(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in *adminv1.GetTeamRequest

		wants    *adminv1.GetTeamResponse
		wantCode connect.Code
	}{
		"by code": {
			in: &adminv1.GetTeamRequest{Code: 1},
			wants: &adminv1.GetTeamResponse{
				Team: &adminv1.Team{Code: 1, Name: "トラブルシューターズ", Organization: "ICTSC Association", MemberLimit: 6},
			},
		},
		"no code": {
			in:       &adminv1.GetTeamRequest{Code: 10},
			wantCode: connect.CodeNotFound,
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := setupTeamService(t)

			resp, err := client.GetTeam(ctx, connect.NewRequest(tt.in))
			assertCode(t, tt.wantCode, err)
			if err != nil {
				return
			}

			if diff := cmp.Diff(
				resp.Msg, tt.wants,
				cmpopts.IgnoreUnexported(adminv1.GetTeamResponse{}, adminv1.Team{}),
			); diff != "" {
				t.Errorf("(-got, +want)\n%s", diff)
			}
		})
	}
}

func TestAdminTeamService_Update(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in *adminv1.UpdateTeamRequest

		wants    *adminv1.UpdateTeamResponse
		wantCode connect.Code
	}{
		"update org": {
			in: &adminv1.UpdateTeamRequest{
				Team: &adminv1.Team{Code: 1, Organization: "ICTSC Committee"},
			},
			wants: &adminv1.UpdateTeamResponse{
				Team: &adminv1.Team{Code: 1, Name: "トラブルシューターズ", Organization: "ICTSC Committee", MemberLimit: 6},
			},
		},
		"cannot update name": {
			in: &adminv1.UpdateTeamRequest{
				Team: &adminv1.Team{Code: 1, Name: "new name"},
			},
			wantCode: connect.CodeInvalidArgument,
		},
		"no team": {
			in:       &adminv1.UpdateTeamRequest{Team: &adminv1.Team{Code: 10}},
			wantCode: connect.CodeNotFound,
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := setupTeamService(t)

			resp, err := client.UpdateTeam(ctx, connect.NewRequest(tt.in))
			assertCode(t, tt.wantCode, err)
			if err != nil {
				return
			}

			if diff := cmp.Diff(
				resp.Msg, tt.wants,
				cmpopts.IgnoreUnexported(adminv1.UpdateTeamResponse{}, adminv1.Team{}),
			); diff != "" {
				t.Errorf("(-got, +want)\n%s", diff)
			}
		})
	}
}

func TestAdminTeamService_Delete(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in       *adminv1.DeleteTeamRequest
		wants    *adminv1.DeleteTeamResponse
		wantCode connect.Code
	}{
		"ok": {
			in:    &adminv1.DeleteTeamRequest{Code: 1},
			wants: &adminv1.DeleteTeamResponse{},
		},
		"no team": {
			in:       &adminv1.DeleteTeamRequest{Code: 10},
			wantCode: connect.CodeNotFound,
		},
	}

	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := setupTeamService(t)

			resp, err := client.DeleteTeam(ctx, connect.NewRequest(tt.in))
			assertCode(t, tt.wantCode, err)
			if err != nil {
				return
			}

			if diff := cmp.Diff(
				resp.Msg, tt.wants,
				cmpopts.IgnoreUnexported(adminv1.DeleteTeamResponse{}),
			); diff != "" {
				t.Errorf("(-got, +want)\n%s", diff)
			}
		})
	}
}

func setupTeamService(t *testing.T) adminv1connect.TeamServiceClient {
	t.Helper()

	db := pgtest.SetupDB(t)

	enforcer := setupEnforcer(t)
	repo := pg.NewRepository(db)

	handler := admin.NewTeamServiceHandler(enforcer, repo)
	mux := http.NewServeMux()
	mux.Handle(adminv1connect.NewTeamServiceHandler(handler))

	server := setupServer(t, mux)

	client := adminv1connect.NewTeamServiceClient(http.DefaultClient, server.URL)
	return client
}
