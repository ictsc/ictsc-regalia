package value

import (
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
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

// NewTeamIDFrom チームIDを生成
func NewTeamIDFrom(value string) (TeamID, error) {
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
