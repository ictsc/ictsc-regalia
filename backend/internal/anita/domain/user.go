package domain

import "github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"

// User ユーザー
type User struct {
	id     value.UserID
	name   value.UserName
	teamID value.TeamID
}

// NewUser ユーザーを作成する
func NewUser(id value.UserID, name value.UserName, teamID value.TeamID) *User {
	return &User{
		id:     id,
		name:   name,
		teamID: teamID,
	}
}

// ID ユーザーIDを取得する
func (u *User) ID() value.UserID {
	return u.id
}

// Name ユーザー名を取得する
func (u *User) Name() value.UserName {
	return u.name
}

// TeamID チームIDを取得する
func (u *User) TeamID() value.TeamID {
	return u.teamID
}

// SetName ユーザー名を設定する
func (u *User) SetName(name value.UserName) {
	u.name = name
}

// Equals ユーザーが等しいか
func (u *User) Equals(other *User) bool {
	return u.id.Equals(other.id)
}
