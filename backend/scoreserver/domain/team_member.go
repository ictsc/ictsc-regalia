package domain

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
)

type (
	TeamMember struct {
		*user
		team *team
	}
	teamMember        = TeamMember
	TeamMemberProfile struct {
		*teamMember
		profile       *profile
		discordUserID DiscordUserID
	}
)

func (u userID) TeamMember(ctx context.Context, eff TeamMemberGetter) (*TeamMember, error) {
	member, err := eff.GetTeamMemberByID(ctx, uuid.UUID(u))
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get team member")
	}
	return member.parse()
}

func (u *user) JoinTeam(ctx context.Context, eff TeamMemberManager, now time.Time, invitationCode *InvitationCode) error {
	if invitationCode.Expired(now) {
		return errors.WithStack(ErrInvitationCodeExpired)
	}

	memberCount, err := eff.CountTeamMembers(ctx, uuid.UUID(invitationCode.team.teamID))
	if err != nil {
		return WrapAsInternal(err, "failed to count team members")
	}
	if memberCount >= invitationCode.team.maxMembers {
		return errors.WithStack(ErrTeamIsFull)
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

func ListTeamMembers(ctx context.Context, eff TeamMemberProfileReader) ([]*TeamMemberProfile, error) {
	list, err := eff.ListTeamMembers(ctx)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to list team members")
	}

	members := make([]*TeamMemberProfile, 0, len(list))
	for _, tmData := range list {
		member, err := tmData.parse()
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (t *Team) Members(ctx context.Context, eff TeamMemberProfileReader) ([]*TeamMemberProfile, error) {
	list, err := eff.ListTeamMembersByTeamID(ctx, uuid.UUID(t.teamID))
	if err != nil {
		return nil, WrapAsInternal(err, "failed to list team members")
	}

	members := make([]*TeamMemberProfile, 0, len(list))
	for _, tmData := range list {
		member, err := tmData.parse()
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (mp *TeamMemberProfile) UserProfile() *UserProfile {
	return &UserProfile{user: mp.user, profile: mp.profile}
}

func (mp *TeamMemberProfile) DiscordUserID() DiscordUserID {
	return mp.discordUserID
}

var (
	ErrTeamIsFull = NewInvalidArgumentError("team is full", nil)
)

type (
	TeamMemberData struct {
		User *UserData `json:"user"`
		Team *TeamData `json:"team"`
	}
	TeamMemberProfileData struct {
		User          *UserData    `json:"user"`
		Team          *TeamData    `json:"team"`
		Profile       *ProfileData `json:"profile"`
		DiscordUserID int64        `json:"discord_user_id"`
	}
	TeamMemberGetter interface {
		GetTeamMemberByID(ctx context.Context, userID uuid.UUID) (*TeamMemberData, error)
		CountTeamMembers(ctx context.Context, teamID uuid.UUID) (uint, error)
	}
	TeamMemberProfileReader interface {
		ListTeamMembers(ctx context.Context) ([]*TeamMemberProfileData, error)
		ListTeamMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*TeamMemberProfileData, error)
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

func (d *TeamMemberProfileData) parse() (*TeamMemberProfile, error) {
	user, err := d.User.parse()
	if err != nil {
		return nil, err
	}

	team, err := d.Team.parse()
	if err != nil {
		return nil, err
	}

	profile, err := d.Profile.parse()
	if err != nil {
		return nil, err
	}

	return &TeamMemberProfile{
		teamMember:    &teamMember{user: user, team: team},
		profile:       profile,
		discordUserID: DiscordUserID(d.DiscordUserID),
	}, nil
}
