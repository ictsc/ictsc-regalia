package domain

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
)

type (
	TeamCode int
	Team     struct {
		id           uuid.UUID
		code         TeamCode
		name         string
		organization string
	}
	TeamInput struct {
		ID           uuid.UUID
		Code         int
		Name         string
		Organization string
	}
)

func NewTeam(input TeamInput) (*Team, error) {
	code, err := NewTeamCode(input.Code)
	if err != nil {
		return nil, err
	}

	if len(input.Name) == 0 {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("team name must not be empty"))
	}

	if len(input.Organization) == 0 {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("team organization must not be empty"))
	}

	return &Team{
		id:           input.ID,
		code:         code,
		name:         input.Name,
		organization: input.Organization,
	}, nil
}

func NewTeamCode(code int) (TeamCode, error) {
	if code < 1 || 99 < code {
		return 0, NewError(ErrTypeInvalidArgument, errors.New("invalid team code"))
	}
	return TeamCode(code), nil
}

func (t *Team) ID() uuid.UUID {
	return t.id
}

func (t *Team) Code() TeamCode {
	return t.code
}

func (t *Team) Name() string {
	return t.name
}

func (t *Team) Organization() string {
	return t.organization
}

type TeamGetEffect = TeamGetter

func (tc TeamCode) Team(ctx context.Context, effect TeamGetEffect) (*Team, error) {
	team, err := effect.GetTeamByCode(ctx, tc)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (tc TeamCode) Team(ctx context.Context, effect TeamGetEffect) (*Team, error) {
	return GetTeamByCode(ctx, effect, tc)
}

type TeamListEffect = TeamsLister

func ListTeams(ctx context.Context, effect TeamListEffect) ([]*Team, error) {
	teams, err := effect.ListTeams(ctx)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

type (
	TeamCreateInput struct {
		Code         int
		Name         string
		Organization string
	}
	TeamCreateEffect   = Tx[TeamCreateTxEffect]
	TeamCreateTxEffect = TeamCreator
)

func CreateTeam(ctx context.Context, effect TeamCreateEffect, input TeamCreateInput) (*Team, error) {
	team, err := createTeam(input)
	if err != nil {
		return nil, err
	}

	if err := effect.RunInTx(ctx, func(effect TeamCreateTxEffect) error {
		return effect.CreateTeam(ctx, team)
	}); err != nil {
		return nil, err
	}
	return team, nil
}

func createTeam(input TeamCreateInput) (*Team, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, NewError(ErrTypeInternal, errors.Wrap(err, "failed to generate uuid"))
	}

	team, err := NewTeam(TeamInput{
		ID:           id,
		Code:         input.Code,
		Name:         input.Name,
		Organization: input.Organization,
	})
	if err != nil {
		return nil, err
	}
	return team, nil
}

type (
	TeamUpdateInput struct {
		// Discord のラベル名がチーム名に紐付くため，チーム名を変更可能にするには Discord のラベル ID とチーム ID の結びつきを保存して管理する必要がある
		// これは実装を複雑にするため，現状はチーム名の変更を許可しない
		// オペレーション上必要になった場合は実装を検討する
		// Name string

		Organization string
	}
	TeamUpdateEffect   = Tx[TeamUpdateTxEffect]
	TeamUpdateTxEffect = TeamUpdater
)

func (t *Team) Update(ctx context.Context, effect TeamUpdateEffect, input TeamUpdateInput) error {
	updated := t.update(input)

	if err := effect.RunInTx(ctx, func(effect TeamUpdateTxEffect) error {
		return effect.UpdateTeam(ctx, updated)
	}); err != nil {
		return err
	}

	*t = *updated
	return nil
}

func (t *Team) update(input TeamUpdateInput) *Team {
	updated := *t
	if len(input.Organization) > 0 {
		updated.organization = input.Organization
	}
	return &updated
}

func (t *Team) Delete(ctx context.Context, effect TeamDeleter) error {
	return effect.DeleteTeam(ctx, t)
}

// チームの操作のためのインターフェース
type (
	TeamsLister interface {
		ListTeams(ctx context.Context) ([]*Team, error)
	}
	TeamGetter interface {
		GetTeamByCode(ctx context.Context, code TeamCode) (*Team, error)
	}
	TeamCreator interface {
		CreateTeam(ctx context.Context, team *Team) error
	}
	TeamUpdater interface {
		UpdateTeam(ctx context.Context, team *Team) error
	}
	TeamDeleter interface {
		DeleteTeam(ctx context.Context, team *Team) error
	}
)
