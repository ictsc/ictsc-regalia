package bun

import (
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
)

func convertToDomainUser(user *User) (*domain.User, error) {
	id, err := value.NewUserID(user.ID)
	if err != nil {
		return nil, err
	}

	name, err := value.NewUserName(user.Name)
	if err != nil {
		return nil, err
	}

	teamID, err := value.NewTeamID(user.TeamID)
	if err != nil {
		return nil, err
	}

	return domain.NewUser(id, name, teamID), nil
}

func convertToDomainUsers(users []*User) ([]*domain.User, error) {
	domainUsers := make([]*domain.User, len(users))

	for i, user := range users {
		u, err := convertToDomainUser(user)
		if err != nil {
			return nil, err
		}

		domainUsers[i] = u
	}

	return domainUsers, nil
}

func convertFromDomainUser(user *domain.User) *User {
	return &User{ // nolint:exhaustruct
		ID:     user.ID().String(),
		Name:   user.Name().Value(),
		TeamID: user.TeamID().String(),
	}
}

func convertToDomainTeam(team *Team) (*domain.Team, error) {
	id, err := value.NewTeamID(team.ID)
	if err != nil {
		return nil, err
	}

	code, err := value.NewTeamCode(team.Code)
	if err != nil {
		return nil, err
	}

	name, err := value.NewTeamName(team.Name)
	if err != nil {
		return nil, err
	}

	organization, err := value.NewTeamOrganization(team.Organization)
	if err != nil {
		return nil, err
	}

	invitationCode, err := value.NewTeamInvitationCode(team.InvitationCode)
	if err != nil {
		return nil, err
	}

	codeRemaining, err := value.NewTeamCodeRemaining(team.CodeRemaining)
	if err != nil {
		return nil, err
	}

	domainTeam := domain.NewTeam(id, code, name, organization, invitationCode, codeRemaining)

	if team.Bastion != nil {
		bastion, err := value.NewBastion(team.Bastion.User, team.Bastion.Password, team.Bastion.Host, team.Bastion.Port)
		if err != nil {
			return nil, err
		}

		domainTeam.SetBastion(bastion)
	}

	members := make([]*domain.User, len(team.Members))

	for i, member := range team.Members {
		user, err := convertToDomainUser(member)
		if err != nil {
			return nil, err
		}

		members[i] = user
	}

	domainTeam.SetMembers(members)

	return domainTeam, nil
}

func convertToDomainTeams(teams []*Team) ([]*domain.Team, error) {
	domainTeams := make([]*domain.Team, len(teams))

	for i, team := range teams {
		t, err := convertToDomainTeam(team)
		if err != nil {
			return nil, err
		}

		domainTeams[i] = t
	}

	return domainTeams, nil
}

func convertFromDomainTeam(team *domain.Team) *Team {
	bastion := convertFromDomainBastion(team.ID(), team.Bastion().V)

	return &Team{ // nolint:exhaustruct
		ID:             team.ID().String(),
		Code:           team.Code().Value(),
		Name:           team.Name().Value(),
		Organization:   team.Organization().Value(),
		InvitationCode: team.InvitationCode().Value(),
		CodeRemaining:  team.CodeRemaining().Value(),
		Bastion:        &bastion,
	}
}

func convertFromDomainBastion(teamID value.TeamID, bastion value.Bastion) Bastion {
	return Bastion{ // nolint:exhaustruct
		TeamID:   teamID.String(),
		User:     bastion.User(),
		Password: bastion.Password(),
		Host:     bastion.Host(),
		Port:     bastion.Port(),
	}
}
