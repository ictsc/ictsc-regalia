package domain_test

import (
	"context"
	"iter"
	"slices"
	"strings"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

func TestUserNameValidation(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in string

		wantErr error
	}{
		"ok": {
			in: "ok_woman",
		},
		"no name": {
			in:      "",
			wantErr: domain.ErrInvalidUserName,
		},
		"too short name": {
			in:      "a",
			wantErr: domain.ErrInvalidUserName,
		},
		"too long name": {
			in:      "abcdefghijklmnopqrstuvwxyz123456abc",
			wantErr: domain.ErrInvalidUserName,
		},
		"invalid character name": {
			in:      "🙆",
			wantErr: domain.ErrInvalidUserName,
		},
		"repeated periods name": {
			in:      "a..b",
			wantErr: domain.ErrInvalidUserName,
		},
		"invalid character uppercase": {
			in:      "ABC",
			wantErr: domain.ErrInvalidUserName,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := domain.NewUserName(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want error type %v, but got %v", tt.wantErr, err)
			}
			if err != nil {
				return
			}
		})
	}
}

func TestGetUserByName(t *testing.T) {
	t.Parallel()

	effect := userListerFunc(func(context.Context, domain.UserListFilter) iter.Seq2[*domain.UserData, error] {
		return asErrIter(slices.Values([]*domain.UserData{
			{
				ID:   uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
				Name: "user1",
			},
			{
				ID:   uuid.FromStringOrNil("00000000-0000-0000-0000-000000000002"),
				Name: "user2",
			},
		}))
	})

	cases := map[string]struct {
		effect domain.UserLister
		name   domain.UserName

		wants   *domain.UserData
		wantErr error
	}{
		"ok": {
			effect: effect,
			name:   must(domain.NewUserName("user1")),

			wants: &domain.UserData{
				ID:   uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
				Name: "user1",
			},
		},
		"not found": {
			effect: effect,
			name:   must(domain.NewUserName("user3")),

			wantErr: domain.ErrNotFound,
		},
		"error": {
			effect: userListerFunc(func(context.Context, domain.UserListFilter) iter.Seq2[*domain.UserData, error] {
				return singleErrIter[*domain.UserData](nil, errors.New("some error"))
			}),
			name: must(domain.NewUserName("user1")),

			wantErr: domain.ErrInternal,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			user, err := tt.name.User(ctx, tt.effect)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want error %v, but got %v", tt.wantErr, err)
			}
			if err != nil {
				return
			}

			actual := user.Data()
			if diff := cmp.Diff(
				tt.wants, actual,
				cmp.AllowUnexported(domain.User{}),
			); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inName        string
		inDisplayName string

		want     *domain.UserProfileData
		wantErrs []error
	}{
		"ok": {
			inName:        "alice",
			inDisplayName: "Alice",

			want: &domain.UserProfileData{
				User: &domain.UserData{
					Name: "alice",
				},
				DisplayName: "Alice",
			},
		},
		"invalid name": {
			inName:        "アルファベットでない",
			inDisplayName: "Alice",

			wantErrs: []error{domain.ErrInvalidUserName},
		},
		"invalid display name": {
			inName:        "alice",
			inDisplayName: strings.Repeat("あ", 128),

			wantErrs: []error{domain.ErrInvalidDisplayName},
		},
		"composite error": {
			inName:        "アルファベットでない",
			inDisplayName: strings.Repeat("あ", 128),

			wantErrs: []error{domain.ErrInvalidUserName, domain.ErrInvalidDisplayName},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var createdProfile *domain.UserProfileData
			effect := userCreateFunc(func(_ context.Context, profile *domain.UserProfileData) error {
				createdProfile = profile
				return nil
			})

			got, err := domain.CreateUser(t.Context(), effect, tt.inName, tt.inDisplayName)
			if err == nil && len(tt.wantErrs) != 0 {
				t.Error("unexpected success")
			}
			for _, wantErr := range tt.wantErrs {
				if !errors.Is(err, wantErr) {
					t.Errorf("want error %v, but got %v", wantErr, err)
				}
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(createdProfile, got.Data()); diff != "" {
				t.Errorf("(-created, +got)\n%s", diff)
			}
			if diff := cmp.Diff(
				tt.want, got.Data(),
				cmpopts.IgnoreFields(domain.UserData{}, "ID"),
			); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

type userListerFunc func(context.Context, domain.UserListFilter) iter.Seq2[*domain.UserData, error]

func (f userListerFunc) ListUsers(ctx context.Context, filter domain.UserListFilter) iter.Seq2[*domain.UserData, error] {
	return f(ctx, filter)
}

type userCreateFunc func(context.Context, *domain.UserProfileData) error

func (f userCreateFunc) CreateUser(ctx context.Context, profile *domain.UserProfileData) error {
	return f(ctx, profile)
}
