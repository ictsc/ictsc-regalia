package memory

import (
	"context"
	"sort"

	"github.com/ictsc/ictsc-regalia/backend/internal/domain"
)

// ListTeamsはメモリ上にあるチームをチーム番号順で返す
func (r *Repository) ListTeams(ctx context.Context) ([]*domain.Team, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	teams := make([]*domain.Team, 0, len(r.teams))
	for _, team := range r.teams {
		teams = append(teams, team)
	}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].TeamNumber() < teams[j].TeamNumber()
	})

	return teams, nil
}

// CreateTeamはチームをメモリ上に保存する
func (r *Repository) CreateTeam(ctx context.Context, team *domain.Team) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	code := team.TeamNumber()
	if _, exists := r.teams[code]; exists {
		return &domain.Error{
			Type: domain.ErrAlreadyExists,
			Msg:  "team already exists",
		}
	}

	r.teams[code] = team
	return nil
}
