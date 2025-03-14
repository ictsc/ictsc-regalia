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

func (u UserID) Profile(ctx context.Context, eff UserProfileReader) (*UserProfile, error) {
	profileData, err := eff.GetUserProfileByID(ctx, uuid.UUID(u))
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get user profile")
	}
	return profileData.parse()
}

func (p *UserProfile) User() *User {
	return p.user
}

func (p *UserProfile) DisplayName() string {
	return p.profile.displayName
}

func CreateUser(ctx context.Context, eff UserCreator, name, displayName string) (*UserProfile, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, WrapAsInternal(err, "failed to generate user ID")
	}

	profile, err := (&UserProfileData{
		User:    &UserData{ID: id, Name: name},
		Profile: &ProfileData{DisplayName: displayName},
	}).parse()
	if err != nil {
		return nil, err
	}

	if err := eff.CreateUser(ctx, profile.Data()); err != nil {
		return nil, WrapAsInternal(err, "failed to create user")
	}
	return profile, nil
}

var (
	ErrInvalidUserName    = NewInvalidArgumentError("invalid user name", nil)
	ErrDuplicateUserName  = NewAlreadyExistsError("user name", nil)
	ErrInvalidDisplayName = NewInvalidArgumentError("invalid display name", nil)
)

// ユーザーに関する操作集合
type (
	UserData struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
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
		User    *UserData    `json:"user"`
		Profile *ProfileData `json:"profile"`
	}
	ProfileData struct {
		DisplayName string `json:"display_name"`
	}
	UserCreator interface {
		CreateUser(ctx context.Context, user *UserProfileData) error
	}
	UserProfileReader interface {
		GetUserProfileByID(ctx context.Context, userID uuid.UUID) (*UserProfileData, error)
	}
)

func (u *User) Data() *UserData {
	return &UserData{
		ID:   uuid.UUID(u.userID),
		Name: string(u.name),
	}
}

func (p *profile) data() *ProfileData {
	return &ProfileData{DisplayName: p.displayName}
}

func (p *UserProfile) Data() *UserProfileData {
	return &UserProfileData{
		User:    p.user.Data(),
		Profile: p.profile.data(),
	}
}

type (
	userID   uuid.UUID
	userName string
	user     struct {
		userID
		name userName
	}
	profile struct {
		displayName string
	}
	userProfile struct {
		*user
		*profile
	}
)

var validUserName = regexp.MustCompile(`^[a-z0-9_.]+$`)

func newUserName(name string) (userName, error) {
	// Discordのユーザー名の制限に合わせる
	if name == "" {
		return "", errors.Join(ErrInvalidUserName, errors.New("name is required"))
	}
	if len(name) < 2 || len(name) > 32 {
		return "", errors.Join(ErrInvalidUserName, errors.New("name length must be between 2 and 32"))
	}
	if !validUserName.MatchString(name) {
		return "", errors.Join(ErrInvalidUserName, errors.New("name contains invalid characters"))
	}
	if strings.Contains(name, "..") {
		return "", errors.Join(ErrInvalidUserName, errors.New("name contains repeated periods"))
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

func (d *ProfileData) parse() (*profile, error) {
	//nolint:mnd
	if len(d.DisplayName) > 64 {
		return nil, errors.Join(ErrInvalidDisplayName, errors.New("display name length must be less than 64"))
	}
	return &profile{displayName: d.DisplayName}, nil
}

func (d *UserProfileData) parse() (*UserProfile, error) {
	var errs []error
	user, err := d.User.parse()
	if err != nil {
		errs = append(errs, err)
	}

	profile, err := d.Profile.parse()
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return &UserProfile{
		user:    user,
		profile: profile,
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
			return nil, newInternalError(errors.New("multiple users found"))
		}
		result = user
	}
	if result == nil {
		return nil, NewNotFoundError("user", nil)
	}
	return result, nil
}
