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

	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	ic1 := &domain.InvitationCodeData{
		ID: must(uuid.NewV4()),
		Team: &domain.TeamData{
			ID:           must(uuid.NewV4()),
			Code:         1,
			Name:         "team1",
			Organization: "org1",
			MaxMembers:   5,
		},
		Code:      "ABCD1234EFGH5678",
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}

	ic2 := &domain.InvitationCodeData{
		ID: must(uuid.NewV4()),
		Team: &domain.TeamData{
			ID:           must(uuid.NewV4()),
			Code:         2,
			Name:         "team2",
			Organization: "org2",
			MaxMembers:   3,
		},
		Code:      "WXYZ9876MNPQ5432",
		ExpiresAt: now.Add(48 * time.Hour),
		CreatedAt: now,
	}

	effect := invitationCodeListerFunc(func(_ context.Context, _ domain.InvitationCodeFilter) ([]*domain.InvitationCodeData, error) {
		return []*domain.InvitationCodeData{ic1, ic2}, nil
	})

	cases := map[string]struct {
		effect domain.InvitationCodeLister

		wants []*domain.InvitationCodeData
	}{
		"ok": {
			effect: effect,

			wants: []*domain.InvitationCodeData{ic1, ic2},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ics, err := domain.ListInvitationCodes(t.Context(), tt.effect)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			actual := make([]*domain.InvitationCodeData, 0, len(ics))
			for _, ic := range ics {
				actual = append(actual, ic.Data())
			}
			if diff := cmp.Diff(
				tt.wants, actual,
				cmp.AllowUnexported(domain.Team{}),
			); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type invitationCodeListerFunc func(ctx context.Context, filter domain.InvitationCodeFilter) ([]*domain.InvitationCodeData, error)

func (f invitationCodeListerFunc) ListInvitationCodes(ctx context.Context, filter domain.InvitationCodeFilter) ([]*domain.InvitationCodeData, error) {
	return f(ctx, filter)
}

func Test_CreateInvitationCode(t *testing.T) {
	t.Parallel()

	team1 := domain.FixTeam1(t, nil)

	now := must(time.Parse(time.RFC3339, "2025-01-01T00:00:00Z"))
	expiresAt := now.Add(24 * time.Hour)

	effect := invitationCodeCreatorFunc(func(context.Context, *domain.InvitationCodeData) error {
		return nil
	})

	cases := map[string]struct {
		team      *domain.Team
		effect    domain.InvitationCodeCreator
		expiresAt time.Time

		want    *domain.InvitationCodeData
		wantErr error
	}{
		"ok": {
			team:      team1,
			effect:    effect,
			expiresAt: expiresAt,
			want: &domain.InvitationCodeData{
				Team:      team1.Data(),
				ExpiresAt: expiresAt,
				CreatedAt: now,
			},
		},
		"already expired": {
			team:      team1,
			effect:    effect,
			expiresAt: now.Add(-24 * time.Hour),
			wantErr:   domain.ErrInvalidArgument,
		},
		"creation fails": {
			team: team1,
			effect: invitationCodeCreatorFunc(func(context.Context, *domain.InvitationCodeData) error {
				return domain.WrapAsInternal(errors.New("dummy"), "failed to create invitation code")
			}),
			expiresAt: expiresAt,
			wantErr:   domain.ErrInternal,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			code, err := tt.team.CreateInvitationCode(t.Context(), tt.effect, now, tt.expiresAt)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want error type %v, got %v", tt.wantErr, err)
			}
			if err != nil {
				return
			}

			actual := code.Data()
			if diff := cmp.Diff(
				tt.want, actual,
				cmp.AllowUnexported(domain.Team{}),
				cmpopts.IgnoreFields(domain.InvitationCodeData{}, "ID", "Code"),
			); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type invitationCodeCreatorFunc func(ctx context.Context, ic *domain.InvitationCodeData) error

func (f invitationCodeCreatorFunc) CreateInvitationCode(ctx context.Context, ic *domain.InvitationCodeData) error {
	return f(ctx, ic)
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
