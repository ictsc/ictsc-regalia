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
	UserID      = userID
	UserName    = userName
	User        = user
	UserProfile = userProfile
)

func NewUserName(name string) (UserName, error) {
	return newUserName(name)
}

func (n UserName) User(ctx context.Context, eff UserLister) (*User, error) {
	return getUser(ctx, eff, UserListFilter{Name: string(n)})
}

func (u *User) ID() UserID {
	return u.userID
}

func (u *User) Name() UserName {
	return u.name
}

// ユーザーに関する操作集合
type (
	UserData struct {
		ID   uuid.UUID
		Name string
	}
	UserListFilter struct {
		Name string
	}
	UserLister interface {
		// ListUsers - ユーザを取得する
		//
		// `filter`の条件を満たさないユーザを返してもよいが，満たすユーザは全て返す
		ListUsers(ctx context.Context, filter UserListFilter) iter.Seq2[*UserData, error]
	}
	UserProfileData struct {
		User        *UserData
		DisplayName string
	}
	UserCreator interface {
		CreateUser(ctx context.Context, user *UserProfileData) error
	}
)

func (u *User) Data() *UserData {
	return &UserData{
		ID:   uuid.UUID(u.userID),
		Name: string(u.name),
	}
}

func (p *UserProfile) Data() *UserProfileData {
	return &UserProfileData{
		User:        p.user.Data(),
		DisplayName: p.displayName,
	}
}

type (
	userID   uuid.UUID
	userName string
	user     struct {
		userID
		name userName
	}
	userProfile struct {
		*user
		displayName string
	}
)

var validUserName = regexp.MustCompile(`^[a-z0-9_.]+$`)

func newUserName(name string) (userName, error) {
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

func (d *UserData) parse() (*user, error) {
	name, err := newUserName(d.Name)
	if err != nil {
		return nil, err
	}
	return &User{
		userID: userID(d.ID),
		name:   name,
	}, nil
}

//nolint:unused
func (d *UserProfileData) _parse() (*UserProfile, error) {
	user, err := d.User.parse()
	if err != nil {
		return nil, err
	}
	//nolint:mnd
	if len(d.DisplayName) > 128 {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("display name length must be less than 128"))
	}
	return &UserProfile{
		user:        user,
		displayName: d.DisplayName,
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

			if !yield(userIn.parse()) {
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
