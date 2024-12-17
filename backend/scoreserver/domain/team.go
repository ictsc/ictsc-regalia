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
	if code < 0 || 100 < code {
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

type TeamCreateInput struct {
	Code         int
	Name         string
	Organization string
}

func CreateTeam(input TeamCreateInput) (*Team, error) {
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

type TeamUpdateInput = TeamCreateInput

func (t Team) Updated(input TeamUpdateInput) (*Team, error) {
	if input.Code != 0 && input.Code != int(t.code) {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("cannot update team code"))
	}
	if input.Name != "" && input.Name != t.name {
		// Discord のラベル名がチーム名に紐付くため，チーム名を変更可能にするには Discord のラベル ID とチーム ID の結びつきを保存して管理する必要がある
		// これは実装を複雑にするため，現状はチーム名の変更を許可しない
		// オペレーション上必要になった場合は実装を検討する
		return nil, NewError(ErrTypeInvalidArgument, errors.New("cannot update team name"))
	}
	if len(input.Organization) > 0 {
		t.organization = input.Organization
	}
	return &t, nil
}

// チームに関するワークフロー

type TeamListWorkflow struct {
	Lister TeamsLister
}

func (w *TeamListWorkflow) Run(ctx context.Context) ([]*Team, error) {
	teams, err := w.Lister.ListTeams(ctx)
	return teams, NewError(ErrTypeInternal, err)
}

type TeamGetWorkflow struct {
	Getter TeamGetter
}
type TeamGetInput struct {
	Code int
}

func (w *TeamGetWorkflow) Run(ctx context.Context, input TeamGetInput) (*Team, error) {
	code, err := NewTeamCode(input.Code)
	if err != nil {
		return nil, err
	}

	team, err := w.Getter.GetTeamByCode(ctx, code)
	if err != nil {
		if ErrTypeFrom(err) == ErrTypeNotFound {
			return nil, err
		}
		return nil, NewError(ErrTypeInternal, err)
	}

	return team, nil
}

type (
	TeamCreateWorkflow struct {
		RunTx TxFunc[TeamCreateTxEffect]
	}
	TeamCreateTxEffect interface {
		TeamCreator
	}
)

func (w *TeamCreateWorkflow) Run(ctx context.Context, input TeamCreateInput) (*Team, error) {
	team, err := CreateTeam(input)
	if err != nil {
		return nil, err
	}

	if err := w.RunTx(ctx, func(effect TeamCreateTxEffect) error {
		if err := effect.CreateTeam(ctx, team); err != nil {
			if ErrTypeFrom(err) == ErrTypeAlreadyExists {
				return NewError(ErrTypeAlreadyExists, errors.New("team already exists"))
			}
			return NewError(ErrTypeInternal, err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return team, nil
}

// チームを更新するワークフロー
type (
	TeamUpdateWorkflow struct {
		RunTx TxFunc[TeamUpdateTxEffect]
	}
	TeamUpdateTxEffect interface {
		TeamGetter
		TeamUpdater
	}
)

func (w *TeamUpdateWorkflow) Run(ctx context.Context, input TeamUpdateInput) (*Team, error) {
	var teamResult *Team
	if err := w.RunTx(ctx, func(effect TeamUpdateTxEffect) error {
		getWf := TeamGetWorkflow{Getter: effect}
		team, err := getWf.Run(ctx, TeamGetInput{Code: input.Code})
		if err != nil {
			return err
		}

		updated, err := team.Updated(input)
		if err != nil {
			return err
		}

		if err := effect.UpdateTeam(ctx, updated); err != nil {
			if ErrTypeFrom(err) == ErrTypeAlreadyExists {
				return err
			}
			return NewError(ErrTypeInternal, err)
		}

		teamResult = updated

		return nil
	}); err != nil {
		return nil, err
	}
	return teamResult, nil
}

type (
	TeamDeleteWorkflow struct {
		RunTx TxFunc[TeamDeleteTxEffect]
	}
	TeamDeleteTxEffect interface {
		TeamGetter
		TeamDeleter
	}
	TeamDeleteInput = TeamGetInput
)

func (w *TeamDeleteWorkflow) Run(ctx context.Context, input TeamDeleteInput) error {
	if err := w.RunTx(ctx, func(effect TeamDeleteTxEffect) error {
		getWf := TeamGetWorkflow{Getter: effect}
		team, err := getWf.Run(ctx, input)
		if err != nil {
			return err
		}

		if err := effect.DeleteTeam(ctx, team); err != nil {
			return NewError(ErrTypeInternal, err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
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
