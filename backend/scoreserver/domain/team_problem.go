package domain

import (
	"context"
	"slices"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
)

type (
	TeamProblem struct {
		*problem
		team  *Team
		score *Score
	}
	TeamProblemDetail struct {
		team          *Team
		problemDetail *DescriptiveProblem
		score         *Score
	}
)

type TeamProblemLister interface {
	TeamsLister
	ProblemReader
	TeamProblemScoreReader
}

func ListTeamProblemsForAdmin(ctx context.Context, eff TeamProblemLister) ([]*TeamProblem, error) {
	return listTeamProblems(ctx, false, eff)
}

func listTeamProblems(ctx context.Context, isPublic bool, eff TeamProblemLister) ([]*TeamProblem, error) {
	teams, err := ListTeams(ctx, eff)
	if err != nil {
		return nil, err
	}

	problems, err := ListProblems(ctx, eff)
	if err != nil {
		return nil, err
	}

	scores, err := eff.ListTeamProblemScores(ctx, isPublic)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to list scores")
	}

	teamProblems := make([]*TeamProblem, 0)
	for _, team := range teams {
		for _, problem := range problems {
			teamProblem := &TeamProblem{team: team, problem: problem}

			idx := slices.IndexFunc(scores, func(score *TeamProblemScoreData) bool {
				return score.TeamID == uuid.UUID(team.teamID) && score.ProblemID == uuid.UUID(problem.problemID)
			})
			if idx >= 0 {
				score, err := scores[idx].Score.parse(problem)
				if err != nil {
					return nil, err
				}
				teamProblem.score = score
			}

			teamProblems = append(teamProblems, teamProblem)
		}
	}

	return teamProblems, nil
}

type TeamProblemReader interface {
	ProblemReader
	TeamProblemScoreReader
}

func (t *Team) ProblemsForPublic(ctx context.Context, eff TeamProblemReader) ([]*TeamProblem, error) {
	return t.problems(ctx, true, eff)
}

func (t *Team) problems(ctx context.Context, isPublic bool, eff TeamProblemReader) ([]*TeamProblem, error) {
	problems, err := ListProblems(ctx, eff)
	if err != nil {
		return nil, err
	}

	scores, err := eff.ListTeamProblemScoresByTeamID(ctx, isPublic, uuid.UUID(t.teamID))
	if err != nil {
		return nil, WrapAsInternal(err, "failed to list scores")
	}

	teamProblems := make([]*TeamProblem, 0, len(problems))
	for _, problem := range problems {
		teamProblem := &TeamProblem{team: t, problem: problem}

		idx := slices.IndexFunc(scores, func(score *TeamProblemScoreData) bool {
			return score.ProblemID == uuid.UUID(problem.problemID)
		})
		if idx >= 0 {
			score, err := scores[idx].Score.parse(problem)
			if err != nil {
				return nil, err
			}
			teamProblem.score = score
		}

		teamProblems = append(teamProblems, teamProblem)
	}
	return teamProblems, nil
}

func (tp *TeamProblem) Team() *Team {
	return tp.team
}

func (tp *TeamProblem) Problem() *Problem {
	return tp.problem
}

func (tp *TeamProblem) ProblemID() uuid.UUID {
	return uuid.UUID(tp.problemID)
}

func (tp *TeamProblem) Score() *Score {
	return tp.score
}

func (tp *TeamProblem) Details(ctx context.Context, eff ProblemReader) (*TeamProblemDetail, error) {
	detail, err := tp.problem.DescriptiveProblem(ctx, eff)
	if err != nil {
		return nil, err
	}
	return &TeamProblemDetail{
		team:          tp.team,
		problemDetail: detail,
		score:         tp.score,
	}, nil
}

func (t *Team) ProblemDetailByCodeForPublic(ctx context.Context, eff TeamProblemReader, code ProblemCode) (*TeamProblemDetail, error) {
	problem, err := code.Problem(ctx, eff)
	if err != nil {
		return nil, err
	}

	detail, err := problem.DescriptiveProblem(ctx, eff)
	if err != nil {
		return nil, err
	}

	teamProblem := &TeamProblemDetail{team: t, problemDetail: detail}

	if scoreData, err := eff.GetTeamProblemScore(
		ctx, true, uuid.UUID(t.teamID), uuid.UUID(problem.problemID),
	); err == nil {
		score, err := scoreData.parse(problem)
		if err != nil {
			return nil, err
		}
		teamProblem.score = score
	} else if !errors.Is(err, ErrNotFound) {
		return nil, WrapAsInternal(err, "failed to get score")
	}

	return teamProblem, nil
}

func (t *Team) ProblemByCodeForPublic(
	ctx context.Context, eff TeamProblemReader, code ProblemCode,
) (*TeamProblem, error) {
	return t.problemByCode(ctx, eff, true, code)
}

func (t *Team) ProblemByCodeForAdmin(
	ctx context.Context, eff TeamProblemReader, code ProblemCode,
) (*TeamProblem, error) {
	return t.problemByCode(ctx, eff, false, code)
}

func (t *Team) problemByCode(ctx context.Context, eff TeamProblemReader, isPublic bool, code ProblemCode) (*TeamProblem, error) {
	problem, err := code.Problem(ctx, eff)
	if err != nil {
		return nil, err
	}

	teamProblem := &TeamProblem{team: t, problem: problem}

	if scoreData, err := eff.GetTeamProblemScore(
		ctx, isPublic, uuid.UUID(t.teamID), uuid.UUID(problem.problemID),
	); err == nil {
		score, err := scoreData.parse(problem)
		if err != nil {
			return nil, err
		}
		teamProblem.score = score
	} else if !errors.Is(err, ErrNotFound) {
		return nil, WrapAsInternal(err, "failed to get score")
	}

	return teamProblem, nil
}

func (tp *TeamProblemDetail) Team() *Team {
	return tp.team
}

func (tp *TeamProblemDetail) TeamProblem() *TeamProblem {
	return &TeamProblem{
		team:    tp.team,
		problem: tp.problemDetail.problem,
		score:   tp.score,
	}
}

func (tp *TeamProblemDetail) ProblemDetail() *DescriptiveProblem {
	return tp.problemDetail
}

type (
	TeamProblemScoreData struct {
		TeamID    uuid.UUID `json:"team_id"`
		ProblemID uuid.UUID `json:"problem_id"`
		Score     ScoreData `json:"score"`
	}
	TeamProblemScoreReader interface {
		GetTeamProblemScore(ctx context.Context, isPublic bool, teamID, problemID uuid.UUID) (*ScoreData, error)
		ListTeamProblemScoresByTeamID(ctx context.Context, isPublic bool, teamID uuid.UUID) ([]*TeamProblemScoreData, error)
		ListTeamProblemScores(ctx context.Context, isPublic bool) ([]*TeamProblemScoreData, error)
	}
)
