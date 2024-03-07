package bun

import (
	"context"

	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/repository"
	"github.com/ictsc/ictsc-outlands/backend/pkg/db/rdb/bun"
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
)

// TeamRepository チームリポジトリ
type TeamRepository struct {
	db *bun.DB
}

var _ repository.TeamRepository = (*TeamRepository)(nil)

// NewTeamRepository チームリポジトリを生成する
func NewTeamRepository(db *bun.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

// SelectTeam チームを取得する
func (repo *TeamRepository) SelectTeam(ctx context.Context, id value.TeamID) (*domain.Team, error) {
	db := repo.db.GetIDB(ctx)
	team := new(Team)

	exists, err := db.NewSelect().Model(team).Where("id = ?", id.String()).Exists(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnknown, err)
	}

	if !exists {
		return nil, errors.Wrap(errors.ErrNotFound, nil)
	}

	err = db.NewSelect().Model(team).
		Relation("Members").
		Join("LEFT OUTER JOIN bastion").JoinOn("bastion.team_id=team.id").
		Where("id = ?", id.String()).
		Scan(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnknown, err)
	}

	return convertToDomainTeam(team)
}

// SelectTeams チームを取得する
func (repo *TeamRepository) SelectTeams(ctx context.Context) ([]*domain.Team, error) {
	db := repo.db.GetIDB(ctx)
	teams := make([]*Team, 0)

	err := db.NewSelect().Model(&teams).Relation("Members").Relation("Bastion").Scan(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnknown, err)
	}

	return convertToDomainTeams(teams)
}

// SelectTeamByInvitationCode 招待コードからチームを取得する
func (repo *TeamRepository) SelectTeamByInvitationCode(ctx context.Context, code value.TeamInvitationCode) (*domain.Team, error) {
	db := repo.db.GetIDB(ctx)
	team := new(Team)

	exists, err := db.NewSelect().Model(team).Where("invitation_code = ?", code.Value()).Exists(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnknown, err)
	}

	if !exists {
		return nil, errors.Wrap(errors.ErrNotFound, nil)
	}

	err = db.NewSelect().Model(team).Relation("Members").Relation("Bastion").Where("invitation_code = ?", code.Value()).Scan(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnknown, err)
	}

	return convertToDomainTeam(team)
}

// UpsertTeam チームを挿入・更新する
func (repo *TeamRepository) UpsertTeam(ctx context.Context, team *domain.Team) error {
	db := repo.db.GetIDB(ctx)
	bunTeam := convertFromDomainTeam(team)

	_, err := db.NewInsert().Model(bunTeam).On("DUPLICATE KEY UPDATE").Exec(ctx)
	if err != nil {
		return errors.Wrap(errors.ErrUnknown, err)
	}

	if bunTeam.Bastion != nil {
		_, err = db.NewInsert().Model(bunTeam.Bastion).On("DUPLICATE KEY UPDATE").Exec(ctx)
		if err != nil {
			return errors.Wrap(errors.ErrUnknown, err)
		}
	}

	return nil
}

// DeleteTeam チームを削除する
func (repo *TeamRepository) DeleteTeam(ctx context.Context, id value.TeamID) error {
	db := repo.db.GetIDB(ctx)

	_, err := db.NewDelete().Model((*Team)(nil)).Where("id = ?", id.String()).Exec(ctx)
	if err != nil {
		return errors.Wrap(errors.ErrUnknown, err)
	}

	_, err = db.NewDelete().Model((*Bastion)(nil)).Where("team_id = ?", id.String()).Exec(ctx)
	if err != nil {
		return errors.Wrap(errors.ErrUnknown, err)
	}

	return nil
}
