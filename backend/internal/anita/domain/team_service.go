package domain

import "github.com/ictsc/ictsc-outlands/backend/pkg/errors"

// MoveMember メンバーのチームを移動
func MoveMember(from, to *Team, member *User) error {
	exists := false

	for _, v := range from.members {
		if v.Equals(member) {
			exists = true

			break
		}
	}

	if !exists {
		return errors.New(errors.ErrNotFound, "Member not found in from")
	}

	for _, v := range to.members {
		if v.Equals(member) {
			return errors.New(errors.ErrAlreadyExists, "Member already exists in to")
		}
	}

	if len(to.members)-1 >= maxTeamMembers {
		return errors.New(errors.ErrBadArgument, "Team is full")
	}

	for i, v := range from.members {
		if v.Equals(member) {
			from.members = append(from.members[:i], from.members[i+1:]...)
		}
	}

	to.members = append(to.members, member)
	member.teamID = to.id

	return nil
}
