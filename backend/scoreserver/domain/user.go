package domain

import (
	"context"
	"iter"
	"regexp"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
)

type (
	UserName string
	User     struct {
		id   uuid.UUID
		name UserName
	}
	UserData struct {
		ID   uuid.UUID
		Name string
	}
	UserProfile struct {
		user        *User
		displayName string
	}
	UserProfileData struct {
		User        *UserData
		DisplayName string
	}
)

func NewUserName(name string) (UserName, error) {
	return newUserName(name)
}

func (n UserName) User(ctx context.Context, eff UserLister) (*User, error) {
	return getUser(ctx, eff, UserListFilter{Name: string(n)})
}

func NewUser(input *UserData) (*User, error) {
	return newUser(input)
}

func (u *User) Name() UserName {
	return u.name
}

func NewUserProfile(input *UserProfileData) (*UserProfile, error) {
	return newUserProfile(input)
}

func (u *UserProfile) User() *User {
	return u.user
}

func (u *UserProfile) DisplayName() string {
	return u.displayName
}

// ユーザーに関する操作集合
type (
	UserListFilter struct {
		Name string
	}
	UserLister interface {
		// ListUsers - ユーザを取得する
		//
		// `filter`の条件を満たさないユーザを返してもよいが，満たすユーザは全て返す
		ListUsers(ctx context.Context, filter UserListFilter) iter.Seq2[*UserData, error]
	}
	UserCreator interface {
		CreateUser(ctx context.Context, user *UserProfileData) error
	}
)

var validUserName = regexp.MustCompile(`^[a-z0-9_.]+$`)

func newUserName(name string) (UserName, error) {
	// Discordのユーザー名の制限に合わせる
	if name == "" {
		return "", NewError(ErrTypeInvalidArgument, errors.New("name is required"))
	}
	if len(name) < 2 || len(name) > 32 {
		return "", NewError(ErrTypeInvalidArgument, errors.New("name length must be between 2 and 32"))
	}
	if !validUserName.MatchString(name) {
		return "", NewError(ErrTypeInvalidArgument, errors.New("name contains invalid characters"))
	}
	if strings.Contains(name, "..") {
		return "", NewError(ErrTypeInvalidArgument, errors.New("name contains repeated periods"))
	}
	return UserName(name), nil
}

func newUser(input *UserData) (*User, error) {
	name, err := NewUserName(input.Name)
	if err != nil {
		return nil, err
	}
	return &User{
		id:   input.ID,
		name: name,
	}, nil
}

func newUserProfile(input *UserProfileData) (*UserProfile, error) {
	user, err := newUser(input.User)
	if err != nil {
		return nil, err
	}

	//nolint:mnd
	if len(input.DisplayName) > 128 {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("display name length must be less than 128"))
	}
	return &UserProfile{
		user:        user,
		displayName: input.DisplayName,
	}, nil
}

func listUsers(ctx context.Context, lister UserLister, filter UserListFilter) iter.Seq2[*User, error] {
	return func(yield func(*User, error) bool) {
		for userIn, err := range lister.ListUsers(ctx, filter) {
			if err != nil {
				if !yield(nil, WrapAsInternal(err, "failed to list users")) {
					return
				}
			}

			if filter.Name != "" && userIn.Name != filter.Name {
				continue
			}

			if !yield(NewUser(userIn)) {
				return
			}
		}
	}
}

func getUser(ctx context.Context, lister UserLister, filter UserListFilter) (*User, error) {
	var result *User
	for user, err := range listUsers(ctx, lister, filter) {
		if err != nil {
			return nil, err
		}
		if result != nil {
			return nil, NewError(ErrTypeInternal, errors.New("multiple users found"))
		}
		result = user
	}
	if result == nil {
		return nil, NewError(ErrTypeNotFound, errors.New("user not found"))
	}
	return result, nil
}
