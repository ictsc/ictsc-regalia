package server

import (
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
	v1 "github.com/ictsc/ictsc-outlands/backend/internal/proto/anita/v1"
)

func fromDomainUser(user *domain.User) *v1.User {
	return &v1.User{
		Id:     user.ID().String(),
		Name:   user.Name().Value(),
		TeamId: user.TeamID().String(),
	}
}

func fromDomainUsers(users []*domain.User) []*v1.User {
	res := make([]*v1.User, len(users))

	for i, user := range users {
		res[i] = fromDomainUser(user)
	}

	return res
}

func fromDomainTeam(team *domain.Team) *v1.Team {
	return &v1.Team{
		Id:             team.ID().String(),
		Code:           int32(team.Code().Value()),
		Name:           team.Name().Value(),
		Organization:   team.Organization().Value(),
		InvitationCode: team.InvitationCode().Value(),
		CodeRemaining:  int32(team.CodeRemaining().Value()),
	}
}

func fromDomainTeams(teams []*domain.Team) []*v1.Team {
	res := make([]*v1.Team, len(teams))

	for i, team := range teams {
		res[i] = fromDomainTeam(team)
	}

	return res
}

func fromDomainBastion(bastion value.Bastion) *v1.Bastion {
	return &v1.Bastion{
		User:     bastion.User(),
		Password: bastion.Password(),
		Host:     bastion.Host(),
		Port:     int32(bastion.Port()),
	}
}
