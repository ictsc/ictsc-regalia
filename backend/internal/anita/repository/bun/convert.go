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
