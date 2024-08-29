package impl

import (
	"context"

	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/repository"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/service"
	"github.com/ictsc/ictsc-outlands/backend/pkg/db/rdb"
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
)

// TeamService チームサービス
type TeamService struct {
	tx    rdb.Tx
	repo  repository.TeamRepository
	uRepo repository.UserRepository
}

var _ service.TeamService = (*TeamService)(nil)

// NewTeamService チームサービスを生成する
func NewTeamService(tx rdb.Tx, repo repository.TeamRepository, uRepo repository.UserRepository) *TeamService {
	return &TeamService{tx: tx, repo: repo, uRepo: uRepo}
}

func (s *TeamService) exists(ctx context.Context, id value.TeamID) bool {
	_, err := s.repo.SelectTeam(ctx, id)

	return err == nil
}

// ReadTeam チームを取得する
func (s *TeamService) ReadTeam(ctx context.Context, id value.TeamID) (*domain.Team, error) {
	return s.repo.SelectTeam(ctx, id)
}

// ReadTeams チームを取得する
func (s *TeamService) ReadTeams(ctx context.Context) ([]*domain.Team, error) {
	return s.repo.SelectTeams(ctx)
}

// ReadTeamByInvitationCode 招待コードからチームを取得する
func (s *TeamService) ReadTeamByInvitationCode(ctx context.Context, code value.TeamInvitationCode) (*domain.Team, error) {
	return s.repo.SelectTeamByInvitationCode(ctx, code)
}

// CreateTeam チームを作成する
func (s *TeamService) CreateTeam(ctx context.Context, args service.CreateTeamArgs) (*domain.Team, error) {
	invCode, err := value.NewRandTeamInvitationCode()
	if err != nil {
		return nil, err
	}

	team := domain.NewTeam(value.NewRandTeamID(), args.Code, args.Name, args.Org, invCode, args.CodeRemaining)

	err = s.tx.Do(ctx, nil, func(ctx context.Context) error {
		if s.exists(ctx, team.ID()) {
			return errors.New(errors.ErrAlreadyExists, "Team already exists")
		}

		return s.repo.UpsertTeam(ctx, team)
	})
	if err != nil {
		return nil, err
	}

	return team, nil
}

// UpdateTeam チームを更新する
func (s *TeamService) UpdateTeam(ctx context.Context, args service.UpdateTeamArgs) (*domain.Team, error) {
	var (
		team *domain.Team
		err  error
	)

	err = s.tx.Do(ctx, nil, func(ctx context.Context) error {
		team, err = s.repo.SelectTeam(ctx, args.ID)
		if err != nil {
			return err
		}

		if args.Code.Valid {
			team.SetCode(args.Code.V)
		}

		if args.Name.Valid {
			team.SetName(args.Name.V)
		}

		if args.Org.Valid {
			team.SetOrganization(args.Org.V)
		}

		if args.CodeRemaining.Valid {
			team.SetCodeRemaining(args.CodeRemaining.V)
		}

		if args.Bastion.Valid {
			team.SetBastion(args.Bastion.V)
		}

		return s.repo.UpsertTeam(ctx, team)
	})
	if err != nil {
		return nil, err
	}

	return team, nil
}

// DeleteTeam チームを削除する
func (s *TeamService) DeleteTeam(ctx context.Context, id value.TeamID) error {
	err := s.tx.Do(ctx, nil, func(ctx context.Context) error {
		if !s.exists(ctx, id) {
			return errors.New(errors.ErrNotFound, "Team not found")
		}

		return s.repo.DeleteTeam(ctx, id)
	})
	if err != nil {
		return err
	}

	return nil
}

// MoveMember メンバーを移動する
func (s *TeamService) MoveMember(ctx context.Context, to value.TeamID, memberID value.UserID) error {
	err := s.tx.Do(ctx, nil, func(ctx context.Context) error {
		member, err := s.uRepo.SelectUser(ctx, memberID)
		if err != nil {
			return err
		}

		fromTeam, err := s.repo.SelectTeam(ctx, member.TeamID())
		if err != nil {
			return err
		}

		toTeam, err := s.repo.SelectTeam(ctx, to)
		if err != nil {
			return err
		}

		err = domain.MoveMember(fromTeam, toTeam, member)
		if err != nil {
			return err
		}

		if err := s.uRepo.UpsertUser(ctx, member); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
