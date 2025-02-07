package domain

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
)

type (
	TeamMember = teamMember
	teamMember struct {
		*user
		team *team
	}
)

func (u *user) TeamMember(ctx context.Context, eff TeamMemberGetter) (*TeamMember, error) {
	member, err := eff.GetTeamMemberByID(ctx, uuid.UUID(u.userID))
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get team member")
	}
	return member.parse()
}

func (u *user) JoinTeam(ctx context.Context, eff TeamMemberManager, now time.Time, invitationCode *InvitationCode) error {
	if invitationCode.Expired(now) {
		return NewError(ErrTypeInvalidArgument, errors.New("invitation code is expired"))
	}

	memberCount, err := eff.CountTeamMembers(ctx, uuid.UUID(invitationCode.team.teamID))
	if err != nil {
		return WrapAsInternal(err, "failed to count team members")
	}
	if memberCount >= invitationCode.team.maxMembers {
		return NewError(ErrTypeInvalidArgument, errors.New("team is full"))
	}

	if err := eff.AddTeamMember(ctx,
		uuid.UUID(u.userID), invitationCode.id, uuid.UUID(invitationCode.team.teamID),
	); err != nil {
		return WrapAsInternal(err, "failed to add team member")
	}
	return nil
}

func (m *teamMember) Team() *Team {
	return m.team
}

type (
	TeamMemberData struct {
		User *UserData
		Team *TeamData
	}
	TeamMemberGetter interface {
		GetTeamMemberByID(ctx context.Context, userID uuid.UUID) (*TeamMemberData, error)
		CountTeamMembers(ctx context.Context, teamID uuid.UUID) (uint, error)
	}
	TeamMemberManager interface {
		TeamMemberGetter
		AddTeamMember(ctx context.Context, userID uuid.UUID, invitationCodeID uuid.UUID, teamID uuid.UUID) error
	}
)

func (d *TeamMemberData) parse() (*teamMember, error) {
	user, err := d.User.parse()
	if err != nil {
		return nil, err
	}

	team, err := d.Team.parse()
	if err != nil {
		return nil, err
	}

	return &teamMember{
		user: user,
		team: team,
	}, nil
}

func (m *teamMember) Data() *TeamMemberData {
	return &TeamMemberData{
		User: m.user.Data(),
		Team: m.team.Data(),
	}
}
