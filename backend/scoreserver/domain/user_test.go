package domain_test

import (
	"context"
	"iter"
	"slices"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

func TestUserValidation(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in *domain.UserData

		wantErr domain.ErrType
	}{
		"ok": {
			in: &domain.UserData{
				ID:   uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
				Name: "ok_woman",
			},
		},
		"no name": {
			in: &domain.UserData{
				ID: uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
			},
			wantErr: domain.ErrTypeInvalidArgument,
		},
		"too short name": {
			in: &domain.UserData{
				ID:   uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
				Name: "a",
			},
			wantErr: domain.ErrTypeInvalidArgument,
		},
		"invalid character name": {
			in: &domain.UserData{
				ID:   uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
				Name: "🙆",
			},
			wantErr: domain.ErrTypeInvalidArgument,
		},
		"repeated periods name": {
			in: &domain.UserData{
				ID:   uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
				Name: "a..b",
			},
			wantErr: domain.ErrTypeInvalidArgument,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := domain.NewUser(tt.in)
			if typ := domain.ErrTypeFrom(err); typ != tt.wantErr {
				t.Errorf("want error type %v, but got %v", tt.wantErr, typ)
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

		wants   *domain.User
		wantErr domain.ErrType
	}{
		"ok": {
			effect: effect,
			name:   must(domain.NewUserName("user1")),

			wants: must(domain.NewUser(&domain.UserData{
				ID:   uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
				Name: "user1",
			})),
		},
		"not found": {
			effect: effect,
			name:   must(domain.NewUserName("user3")),

			wantErr: domain.ErrTypeNotFound,
		},
		"error": {
			effect: userListerFunc(func(context.Context, domain.UserListFilter) iter.Seq2[*domain.UserData, error] {
				return singleErrIter[*domain.UserData](nil, errors.New("some error"))
			}),
			name: must(domain.NewUserName("user1")),

			wantErr: domain.ErrTypeInternal,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			actual, err := tt.name.User(ctx, tt.effect)
			if typ := domain.ErrTypeFrom(err); typ != tt.wantErr {
				t.Errorf("want error type %v, but got %v", tt.wantErr, typ)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(
				tt.wants, actual,
				cmp.AllowUnexported(domain.User{}),
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
