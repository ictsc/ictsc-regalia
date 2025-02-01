package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

func Test_InvitationCodeListWorkflow(t *testing.T) {
	t.Parallel()

	team1 := must(domain.NewTeam(domain.TeamInput{
		ID:           must(uuid.NewV4()),
		Code:         1,
		Name:         "team1",
		Organization: "org1",
	}))

	team2 := must(domain.NewTeam(domain.TeamInput{
		ID:           must(uuid.NewV4()),
		Code:         2,
		Name:         "team2",
		Organization: "org2",
	}))

	now := must(time.Parse(time.RFC3339, "2025-01-01T00:00:00Z"))
	ic1 := must(domain.NewInvitationCode(domain.InvitationCodeInput{
		ID:        must(uuid.NewV4()),
		Team:      team1,
		Code:      "ABCD1234EFGH5678",
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}))

	ic2 := must(domain.NewInvitationCode(domain.InvitationCodeInput{
		ID:        must(uuid.NewV4()),
		Team:      team2,
		Code:      "WXYZ9876MNPQ5432",
		ExpiresAt: now.Add(48 * time.Hour),
		CreatedAt: now,
	}))

	workflow := &domain.InvitationCodeListWorkflow{
		Lister: invitationCodeListerFunc(func(_ context.Context, _ domain.InvitationCodeFilter) ([]*domain.InvitationCode, error) {
			return []*domain.InvitationCode{ic1, ic2}, nil
		}),
	}

	cases := map[string]struct {
		w *domain.InvitationCodeListWorkflow

		wants []*domain.InvitationCode
	}{
		"ok": {
			w: workflow,

			wants: []*domain.InvitationCode{ic1, ic2},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ics, err := tt.w.Run(context.Background())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(
				tt.wants, ics,
				cmp.AllowUnexported(domain.InvitationCode{}, domain.Team{}),
			); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_InvitationCodeCreateWorkflow(t *testing.T) {
	t.Parallel()

	team1 := must(domain.NewTeam(domain.TeamInput{
		ID:           must(uuid.NewV4()),
		Code:         1,
		Name:         "team1",
		Organization: "org",
	}))
	icCreator := invitationCodeCreatorFunc(
		func(context.Context, *domain.InvitationCode) error { return nil })
	workflow := &domain.InvitationCodeCreateWorkflow{
		TeamGetter: teamCodeGetterFunc(
			func(_ context.Context, code domain.TeamCode) (*domain.Team, error) {
				if code != 1 {
					return nil, domain.NewError(domain.ErrTypeNotFound, errors.New("team not found"))
				}
				return team1, nil
			}),

		RunTx: func(_ context.Context, f func(eff domain.InvitationCodeCreator) error) error {
			return f(icCreator)
		},
		Clock: func() time.Time {
			return must(time.Parse(time.RFC3339, "2025-01-01T00:00:00Z"))
		},
	}

	cases := map[string]struct {
		w  *domain.InvitationCodeCreateWorkflow
		in domain.InvitationCodeCreateInput

		want    *domain.InvitationCode
		wantErr domain.ErrType
	}{
		"ok": {
			w: workflow,
			in: domain.InvitationCodeCreateInput{
				TeamCode:  1,
				ExpiresAt: must(time.Parse(time.RFC3339, "2025-04-03T09:00:00Z")),
			},

			want: must(domain.NewInvitationCode(domain.InvitationCodeInput{
				ID:        must(uuid.NewV4()),
				Code:      "dummy",
				Team:      team1,
				ExpiresAt: must(time.Parse(time.RFC3339, "2025-04-03T09:00:00Z")),
				CreatedAt: must(time.Parse(time.RFC3339, "2025-01-01T00:00:00Z")),
			})),
		},
		"no team code": {
			w: workflow,
			in: domain.InvitationCodeCreateInput{
				ExpiresAt: must(time.Parse(time.RFC3339, "2025-04-03T09:00:00Z")),
			},

			wantErr: domain.ErrTypeInvalidArgument,
		},
		"no expires at": {
			w: workflow,
			in: domain.InvitationCodeCreateInput{
				TeamCode: 1,
			},

			wantErr: domain.ErrTypeInvalidArgument,
		},
		"team not found": {
			w: workflow,
			in: domain.InvitationCodeCreateInput{
				TeamCode:  2,
				ExpiresAt: must(time.Parse(time.RFC3339, "2025-04-03T09:00:00Z")),
			},

			wantErr: domain.ErrTypeNotFound,
		},
		"already expired": {
			w: workflow,
			in: domain.InvitationCodeCreateInput{
				TeamCode:  1,
				ExpiresAt: must(time.Parse(time.RFC3339, "2024-04-01T00:00:00Z")),
			},

			wantErr: domain.ErrTypeInvalidArgument,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, err := tt.w.Run(context.Background(), tt.in)
			if domain.ErrTypeFrom(err) != tt.wantErr {
				t.Errorf("want error typ %v, got %v", tt.wantErr, err)
			}
			if err != nil {
				t.Logf("error: %v", err)
				return
			}
			if diff := cmp.Diff(
				tt.want, actual,
				cmp.AllowUnexported(domain.InvitationCode{}, domain.Team{}),
				cmpopts.IgnoreFields(domain.InvitationCode{}, "id", "code"),
			); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type invitationCodeListerFunc func(ctx context.Context, filter domain.InvitationCodeFilter) ([]*domain.InvitationCode, error)

func (f invitationCodeListerFunc) ListInvitationCodes(ctx context.Context, filter domain.InvitationCodeFilter) ([]*domain.InvitationCode, error) {
	return f(ctx, filter)
}

type invitationCodeCreatorFunc func(ctx context.Context, ic *domain.InvitationCode) error

func (f invitationCodeCreatorFunc) CreateInvitationCode(ctx context.Context, ic *domain.InvitationCode) error {
	return f(ctx, ic)
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
