package value

import (
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
	"github.com/ictsc/ictsc-outlands/backend/pkg/random"
	"github.com/oklog/ulid/v2"
)

// TeamID チームID
type TeamID struct {
	value ulid.ULID
}

// NewRandTeamID ランダムなチームIDを生成
func NewRandTeamID() TeamID {
	return TeamID{value: ulid.Make()}
}

// NewTeamID チームIDを生成
func NewTeamID(value string) (TeamID, error) {
	id, err := ulid.Parse(value)
	if err != nil {
		return TeamID{}, errors.Wrap(errors.ErrBadArgument, err)
	}

	return TeamID{value: id}, nil
}

// Equals チームIDが等しいか
func (id TeamID) Equals(other TeamID) bool {
	return id.value == other.value
}

// String 文字列に変換
func (id TeamID) String() string {
	return id.value.String()
}

// TeamCode チームコード
type TeamCode struct {
	value int
}

// NewTeamCode チームコードを生成
func NewTeamCode(value int) (TeamCode, error) {
	if value < 1 || value > 100 {
		return TeamCode{}, errors.New(errors.ErrBadArgument, "Invalid value")
	}

	return TeamCode{value: value}, nil
}

// Equals チームコードが等しいか
func (code TeamCode) Equals(other TeamCode) bool {
	return code.value == other.value
}

// Value 値を取得
func (code TeamCode) Value() int {
	return code.value
}

// TeamName チーム名
type TeamName struct {
	value string
}

// NewTeamName チーム名を生成
func NewTeamName(value string) (TeamName, error) {
	if len(value) < 1 || len(value) > 20 {
		return TeamName{}, errors.New(errors.ErrBadArgument, "Invalid value")
	}

	return TeamName{value: value}, nil
}

// Equals チーム名が等しいか
func (name TeamName) Equals(other TeamName) bool {
	return name.value == other.value
}

// Value 値を取得
func (name TeamName) Value() string {
	return name.value
}

// TeamOrganization チームの所属組織
type TeamOrganization struct {
	value string
}

// NewTeamOrganization チームの所属組織を生成
func NewTeamOrganization(value string) (TeamOrganization, error) {
	if len(value) < 1 || len(value) > 50 {
		return TeamOrganization{}, errors.New(errors.ErrBadArgument, "Invalid value")
	}

	return TeamOrganization{value: value}, nil
}

// Equals チームの所属組織が等しいか
func (org TeamOrganization) Equals(other TeamOrganization) bool {
	return org.value == other.value
}

// Value 値を取得
func (org TeamOrganization) Value() string {
	return org.value
}

const invitationCodeDigit = 32

// TeamInvitationCode チームの招待コード
type TeamInvitationCode struct {
	value string
}

// NewRandTeamInvitationCode ランダムなチームの招待コードを生成
func NewRandTeamInvitationCode() (TeamInvitationCode, error) {
	code, err := random.NewString(invitationCodeDigit)
	if err != nil {
		return TeamInvitationCode{}, err
	}

	return TeamInvitationCode{value: code}, nil
}

// NewTeamInvitationCode チームの招待コードを生成
func NewTeamInvitationCode(value string) (TeamInvitationCode, error) {
	if len(value) != invitationCodeDigit {
		return TeamInvitationCode{}, errors.New(errors.ErrBadArgument, "Invalid code")
	}

	return TeamInvitationCode{value: value}, nil
}

// Equals チームの招待コードが等しいか
func (code TeamInvitationCode) Equals(other TeamInvitationCode) bool {
	return code.value == other.value
}

// Value 値を取得
func (code TeamInvitationCode) Value() string {
	return code.value
}

// TeamCodeRemaining 招待コードの残り回数
type TeamCodeRemaining struct {
	value int
}

// NewTeamCodeRemaining 招待コードの残り回数を生成
func NewTeamCodeRemaining(value int) (TeamCodeRemaining, error) {
	if value < 0 || value > 5 {
		return TeamCodeRemaining{}, errors.New(errors.ErrBadArgument, "Invalid value")
	}

	return TeamCodeRemaining{value: value}, nil
}

// Equals 招待コードの残り回数が等しいか
func (codeRemaining TeamCodeRemaining) Equals(other TeamCodeRemaining) bool {
	return codeRemaining.value == other.value
}

// Value 値を取得
func (codeRemaining TeamCodeRemaining) Value() int {
	return codeRemaining.value
}
