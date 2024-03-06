package server

import (
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain"
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
