package value

import (
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
	"github.com/oklog/ulid/v2"
)

// UserID ユーザーID
type UserID struct {
	value ulid.ULID
}

// NewUserID 文字列からユーザーIDを生成
func NewUserID(value string) (UserID, error) {
	id, err := ulid.Parse(value)
	if err != nil {
		return UserID{}, errors.Wrap(errors.ErrBadArgument, err)
	}

	return UserID{value: id}, nil
}

// Equals ユーザーIDが等しいか
func (id UserID) Equals(other UserID) bool {
	return id.value == other.value
}

// String 文字列に変換
func (id UserID) String() string {
	return id.value.String()
}

// UserName ユーザー名
type UserName struct {
	value string
}

// NewUserName ユーザー名を生成
func NewUserName(value string) (UserName, error) {
	if len(value) < 1 || len(value) > 20 {
		return UserName{}, errors.New(errors.ErrBadArgument, "Invalid value")
	}

	return UserName{value: value}, nil
}

// Equals ユーザー名が等しいか
func (name UserName) Equals(other UserName) bool {
	return name.value == other.value
}

// Value 値を取得
func (name UserName) Value() string {
	return name.value
}
