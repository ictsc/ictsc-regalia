package domain

import (
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
	"github.com/ictsc/ictsc-outlands/backend/pkg/optional"
)

const maxTeamMembers = 5

// Team チーム
type Team struct {
	id             value.TeamID
	code           value.TeamCode
	name           value.TeamName
	organization   value.TeamOrganization
	invitationCode value.TeamInvitationCode
	codeRemaining  value.TeamCodeRemaining

	bastion optional.Of[value.Bastion]

	members []*User
}

// NewTeam チームを作成する
func NewTeam(
	id value.TeamID,
	code value.TeamCode,
	name value.TeamName,
	organization value.TeamOrganization,
	invitationCode value.TeamInvitationCode,
	codeRemaining value.TeamCodeRemaining,
) *Team {
	return &Team{
		id:             id,
		code:           code,
		name:           name,
		organization:   organization,
		invitationCode: invitationCode,
		codeRemaining:  codeRemaining,
		bastion:        optional.New(value.Bastion{}, false),
		members:        make([]*User, 0, maxTeamMembers),
	}
}

// ID チームIDを取得する
func (t *Team) ID() value.TeamID {
	return t.id
}

// Code チームコードを取得する
func (t *Team) Code() value.TeamCode {
	return t.code
}

// Name チーム名を取得する
func (t *Team) Name() value.TeamName {
	return t.name
}

// Organization 組織名を取得する
func (t *Team) Organization() value.TeamOrganization {
	return t.organization
}

// InvitationCode 招待コードを取得する
func (t *Team) InvitationCode() value.TeamInvitationCode {
	return t.invitationCode
}

// CodeRemaining 招待コードの残り回数を取得する
func (t *Team) CodeRemaining() value.TeamCodeRemaining {
	return t.codeRemaining
}

// Bastion 踏み台サーバーを取得する
func (t *Team) Bastion() optional.Of[value.Bastion] {
	return t.bastion
}

// Members メンバーを取得する
func (t *Team) Members() []*User {
	return t.members
}

// SetCode チームコードを設定する
func (t *Team) SetCode(code value.TeamCode) {
	t.code = code
}

// SetName チーム名を設定する
func (t *Team) SetName(name value.TeamName) {
	t.name = name
}

// SetOrganization 組織名を設定する
func (t *Team) SetOrganization(organization value.TeamOrganization) {
	t.organization = organization
}

// SetCodeRemaining 招待コードの残り回数を設定する
func (t *Team) SetCodeRemaining(codeRemaining value.TeamCodeRemaining) {
	t.codeRemaining = codeRemaining
}

// SetBastion 踏み台サーバーを設定する
func (t *Team) SetBastion(bastion value.Bastion) {
	t.bastion = optional.NewValid(bastion)
}

// SetMembers メンバーを設定する
//
//	型変換以外での使用は禁止
func (t *Team) SetMembers(members []*User) {
	t.members = members
}

// DecrementCodeRemaining 招待コードの残り回数を減らす
func (t *Team) DecrementCodeRemaining() error {
	codeRemaining, err := value.NewTeamCodeRemaining(t.codeRemaining.Value() - 1)
	if err != nil {
		return err
	}

	t.codeRemaining = codeRemaining

	return nil
}
