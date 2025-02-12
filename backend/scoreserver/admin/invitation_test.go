package admin_test

import (
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAdminInvitationService_Create(t *testing.T) {
	t.Parallel()

	now := time.Now()

	cases := map[string]struct {
		in       *adminv1.CreateInvitationCodeRequest
		wants    *adminv1.CreateInvitationCodeResponse
		wantCode connect.Code
	}{
		"ok": {
			in: &adminv1.CreateInvitationCodeRequest{
				InvitationCode: &adminv1.InvitationCode{
					TeamCode:  1,
					ExpiresAt: timestamppb.New(now.Add(24 * time.Hour)),
				},
			},
			wants: &adminv1.CreateInvitationCodeResponse{
				InvitationCode: &adminv1.InvitationCode{
					TeamCode:  1,
					ExpiresAt: timestamppb.New(now.Add(24 * time.Hour)),
				},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			enforcer := setupEnforcer(t)

			db := pgtest.SetupDB(t)

			mux := http.NewServeMux()
			mux.Handle(adminv1connect.NewInvitationServiceHandler(admin.NewInvitationServiceHandler(
				enforcer, pg.NewRepository(db),
			)))

			server := setupServer(t, mux)
			client := adminv1connect.NewInvitationServiceClient(server.Client(), server.URL)

			resp, err := client.CreateInvitationCode(ctx, connect.NewRequest(tt.in))
			assertCode(t, tt.wantCode, err)
			if err != nil {
				return
			}

			if diff := cmp.Diff(
				resp.Msg, tt.wants,
				cmpopts.IgnoreUnexported(
					adminv1.CreateInvitationCodeResponse{},
					adminv1.InvitationCode{},
					timestamppb.Timestamp{}),
				cmpopts.IgnoreFields(adminv1.CreateInvitationCodeResponse{}, "InvitationCode.Code"),
			); diff != "" {
				t.Errorf("(-got, +want)\n%s", diff)
			}
		})
	}
}
