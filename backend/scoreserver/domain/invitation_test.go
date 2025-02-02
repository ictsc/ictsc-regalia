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

func Test_ListInvitationCode(t *testing.T) {
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

	effect := invitationCodeListerFunc(func(_ context.Context, _ domain.InvitationCodeFilter) ([]*domain.InvitationCode, error) {
		return []*domain.InvitationCode{ic1, ic2}, nil
	})

	cases := map[string]struct {
		effect domain.InvitationCodeLister

		wants []*domain.InvitationCode
	}{
		"ok": {
			effect: effect,

			wants: []*domain.InvitationCode{ic1, ic2},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ics, err := domain.ListInvitationCodes(context.Background(), tt.effect)
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

type invitationCodeListerFunc func(ctx context.Context, filter domain.InvitationCodeFilter) ([]*domain.InvitationCode, error)

func (f invitationCodeListerFunc) ListInvitationCodes(ctx context.Context, filter domain.InvitationCodeFilter) ([]*domain.InvitationCode, error) {
	return f(ctx, filter)
}

func Test_CreateInvitationCode(t *testing.T) {
	t.Parallel()

	team1 := must(domain.NewTeam(domain.TeamInput{
		ID:           must(uuid.NewV4()),
		Code:         1,
		Name:         "team1",
		Organization: "org",
	}))

	now := must(time.Parse(time.RFC3339, "2025-01-01T00:00:00Z"))
	expiresAt := now.Add(24 * time.Hour)

	effect := &struct {
		invitationCodeCreatorFunc
		domain.ClockerFunc
	}{
		invitationCodeCreatorFunc: func(ctx context.Context, ic *domain.InvitationCode) error {
			return nil
		},
		ClockerFunc: func() time.Time {
			return now
		},
	}

	cases := map[string]struct {
		team      *domain.Team
		effect    domain.InvitationCodeCreateEffect
		expiresAt time.Time
		wantErr   domain.ErrType
		want      *domain.InvitationCode
	}{
		"ok": {
			team:      team1,
			effect:    effect,
			expiresAt: expiresAt,
			want: must(domain.NewInvitationCode(domain.InvitationCodeInput{
				ID:        must(uuid.NewV4()),
				Team:      team1,
				Code:      "dummy",
				ExpiresAt: expiresAt,
				CreatedAt: now,
			})),
		},
		"already expired": {
			team:      team1,
			effect:    effect,
			expiresAt: now.Add(-24 * time.Hour),
			wantErr:   domain.ErrTypeInvalidArgument,
		},
		"creation fails": {
			team: team1,
			effect: &struct {
				invitationCodeCreatorFunc
				domain.Clocker
			}{
				Clocker: effect,
				invitationCodeCreatorFunc: func(ctx context.Context, ic *domain.InvitationCode) error {
					return domain.NewError(domain.ErrTypeInternal, errors.New("dummy"))
				},
			},
			expiresAt: expiresAt,
			wantErr:   domain.ErrTypeInternal,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, err := tt.team.CreateInvitationCode(context.Background(), tt.effect, tt.expiresAt)
			if domain.ErrTypeFrom(err) != tt.wantErr {
				t.Errorf("want error type %v, got %v", tt.wantErr, err)
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
