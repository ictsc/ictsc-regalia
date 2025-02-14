package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

func TestJoinTeam(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	cases := map[string]struct {
		inUser            *domain.User
		inCode            *domain.InvitationCode
		inTeamMemberCount uint

		wantUserID uuid.UUID
		wantCodeID uuid.UUID
		wantTeamID uuid.UUID
		wantErr    error
	}{
		"ok": {
			inUser: domain.FixUser1(t, &domain.UserData{
				ID: uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
			}),
			inCode: domain.FixInvitationCode1(t, &domain.InvitationCodeData{
				ID: uuid.FromStringOrNil("ad3f83d3-65be-4884-8a03-adb11a8127ef"),
				Team: &domain.TeamData{
					ID:         uuid.FromStringOrNil("69094ee0-70ce-4f07-8fd0-56bb8caf80a6"),
					MaxMembers: 3,
				},
			}),
			inTeamMemberCount: 0,

			wantUserID: uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
			wantCodeID: uuid.FromStringOrNil("ad3f83d3-65be-4884-8a03-adb11a8127ef"),
			wantTeamID: uuid.FromStringOrNil("69094ee0-70ce-4f07-8fd0-56bb8caf80a6"),
		},
		"expired": {
			inUser: domain.FixUser1(t, nil),
			inCode: domain.FixInvitationCode1(t, &domain.InvitationCodeData{
				ExpiresAt: now.Add(-24 * time.Hour),
				CreatedAt: now.Add(-48 * time.Hour),
				Team:      &domain.TeamData{MaxMembers: 3},
			}),
			inTeamMemberCount: 0,

			wantErr: domain.ErrInvitationCodeExpired,
		},
		"full member": {
			inUser: domain.FixUser1(t, nil),
			inCode: domain.FixInvitationCode1(t, &domain.InvitationCodeData{
				Team: &domain.TeamData{MaxMembers: 3},
			}),
			inTeamMemberCount: 3,

			wantErr: domain.ErrTeamIsFull,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var actualUserID, actualCodeID, actualTeamID uuid.UUID
			manager := teamMemberManagerFns{
				teamMemberGetterFns: teamMemberGetterFns{
					countTeamMembersFn: func(context.Context, uuid.UUID) (uint, error) {
						return tt.inTeamMemberCount, nil
					},
				},
				addTeamMemberFn: func(ctx context.Context, userID, invitationCodeID, teamID uuid.UUID) error {
					actualUserID, actualCodeID, actualTeamID = userID, invitationCodeID, teamID
					return nil
				},
			}

			err := tt.inUser.JoinTeam(t.Context(), manager, now, tt.inCode)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got error type %v, want %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			if actualUserID != tt.wantUserID {
				t.Errorf("got userID %v, want %v", actualUserID, tt.wantUserID)
			}
			if actualCodeID != tt.wantCodeID {
				t.Errorf("got codeID %v, want %v", actualCodeID, tt.wantCodeID)
			}
			if actualTeamID != tt.wantTeamID {
				t.Errorf("got teamID %v, want %v", actualTeamID, tt.wantTeamID)
			}
		})
	}
}

type (
	getTeamMemberByIDFn func(ctx context.Context, userID uuid.UUID) (*domain.TeamMemberData, error)
	countTeamMembersFn  func(ctx context.Context, teamID uuid.UUID) (uint, error)
	addTeamMemberFn     func(ctx context.Context, userID uuid.UUID, invitationCodeID uuid.UUID, teamID uuid.UUID) error
	teamMemberGetterFns struct {
		getTeamMemberByIDFn
		countTeamMembersFn
	}
	teamMemberManagerFns struct {
		teamMemberGetterFns
		addTeamMemberFn
	}
)

func (f getTeamMemberByIDFn) GetTeamMemberByID(ctx context.Context, userID uuid.UUID) (*domain.TeamMemberData, error) {
	return f(ctx, userID)
}

func (f countTeamMembersFn) CountTeamMembers(ctx context.Context, teamID uuid.UUID) (uint, error) {
	return f(ctx, teamID)
}

func (f addTeamMemberFn) AddTeamMember(ctx context.Context, userID uuid.UUID, invitationCodeID uuid.UUID, teamID uuid.UUID) error {
	return f(ctx, userID, invitationCodeID, teamID)
}
