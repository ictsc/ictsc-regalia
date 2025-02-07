package domain

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
)

type (
	TeamCode = teamCode
	Team     = team
)

func NewTeamCode(code int64) (TeamCode, error) {
	if code < 1 || 99 < code {
		return 0, NewError(ErrTypeInvalidArgument, errors.New("invalid team code"))
	}
	return TeamCode(code), nil
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

func (t *Team) MaxMembers() uint {
	return t.maxMembers
}

type TeamGetEffect = TeamGetter

func (tc TeamCode) Team(ctx context.Context, effect TeamGetEffect) (*Team, error) {
	data, err := effect.GetTeamByCode(ctx, int64(tc))
	if err != nil {
		return nil, err
	}
	return data.parse()
}

type TeamListEffect = TeamsLister

func ListTeams(ctx context.Context, effect TeamListEffect) ([]*Team, error) {
	data, err := effect.ListTeams(ctx)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to list teams")
	}

	teams := make([]*Team, 0, len(data))
	for _, d := range data {
		t, err := d.parse()
		if err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}

	return teams, nil
}

type (
	TeamCreateInput struct {
		Code         int
		Name         string
		Organization string
		MaxMembers   uint
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
		return effect.CreateTeam(ctx, team.Data())
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

	team, err := (&TeamData{
		ID:           id,
		Code:         int64(input.Code),
		Name:         input.Name,
		Organization: input.Organization,
		MaxMembers:   input.MaxMembers,
	}).parse()
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
		return effect.UpdateTeam(ctx, updated.Data())
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
	return effect.DeleteTeam(ctx, uuid.UUID(t.teamID))
}

// チームの操作のためのインターフェース
type (
	TeamData struct {
		ID           uuid.UUID
		Code         int64
		Name         string
		Organization string
		MaxMembers   uint
	}
	TeamsLister interface {
		ListTeams(ctx context.Context) ([]*TeamData, error)
	}
	TeamGetter interface {
		GetTeamByCode(ctx context.Context, code int64) (*TeamData, error)
	}
	TeamCreator interface {
		CreateTeam(ctx context.Context, team *TeamData) error
	}
	TeamUpdater interface {
		UpdateTeam(ctx context.Context, team *TeamData) error
	}
	TeamDeleter interface {
		DeleteTeam(ctx context.Context, teamID uuid.UUID) error
	}
)

func (t *TeamData) merge(data *TeamData) *TeamData {
	if !data.ID.IsNil() {
		t.ID = data.ID
	}
	if data.Code != 0 {
		t.Code = data.Code
	}
	if data.Name != "" {
		t.Name = data.Name
	}
	if data.Organization != "" {
		t.Organization = data.Organization
	}
	if data.MaxMembers != 0 {
		t.MaxMembers = data.MaxMembers
	}
	return t
}

func (t *Team) Data() *TeamData {
	return &TeamData{
		ID:           uuid.UUID(t.teamID),
		Code:         int64(t.code),
		Name:         t.name,
		Organization: t.organization,
		MaxMembers:   t.maxMembers,
	}
}

type (
	teamID   uuid.UUID
	teamCode int64
	team     struct {
		teamID
		code         teamCode
		name         string
		organization string
		maxMembers   uint
	}
)

func (t *TeamData) parse() (*team, error) {
	code, err := NewTeamCode(t.Code)
	if err != nil {
		return nil, err
	}

	if len(t.Name) == 0 {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("team name must not be empty"))
	}

	if len(t.Organization) == 0 {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("team organization must not be empty"))
	}

	if t.MaxMembers < 1 {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("team max members must be greater than 0"))
	}

	return &Team{
		teamID:       teamID(t.ID),
		code:         code,
		name:         t.Name,
		organization: t.Organization,
		maxMembers:   t.MaxMembers,
	}, nil
}
