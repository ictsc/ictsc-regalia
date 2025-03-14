package domain

import (
	"context"
	"log/slog"
	"slices"

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
	TeamProblemScoreData struct {
		ProblemID   uuid.UUID `db:"problem_id"`
		MarkedScore uint32    `db:"marked_score"`
		Penalty     uint32    `db:"penalty"`
		TotalScore  uint32    `db:"total_score"`
		MaxScore    uint32    `db:"max_score"`
	}
)

type TeamProblemReader interface {
	TeamsLister
	ProblemReader
}

func ListTeamProblems(ctx context.Context, eff TeamProblemReader) ([]*TeamProblem, error) {
	teams, err := ListTeams(ctx, eff)
	if err != nil {
		return nil, err
	}
	slog.DebugContext(ctx, "teams listed", "count", len(teams))

	teamProblems := make([]*TeamProblem, 0)
	for _, team := range teams {
		tps, err := team.Problems(ctx, eff)
		if err != nil {
			return nil, err
		}
		slog.DebugContext(ctx, "problems listed", "team", team.Code(), "count", len(tps))
		teamProblems = append(teamProblems, tps...)
	}

	return teamProblems, nil
}

func (t *Team) Problems(ctx context.Context, eff ProblemReader) ([]*TeamProblem, error) {
	problems, err := ListProblems(ctx, eff)
	if err != nil {
		return nil, err
	}

	scores, err := t.listProblemsScore(ctx, eff)
	if err != nil {
		return nil, err
	}

	teamProblems := make([]*TeamProblem, 0, len(problems))
	for _, problem := range problems {
		teamProblem := &TeamProblem{team: t, problem: problem}

		idx := slices.IndexFunc(scores, func(score *TeamProblemScoreData) bool {
			return score.ProblemID == uuid.UUID(problem.problemID)
		})
		if idx >= 0 {
			score, err := scores[idx].parse(problem)
			if err != nil {
				return nil, err
			}
			teamProblem.score = score
		}

		teamProblems = append(teamProblems, teamProblem)
	}
	return teamProblems, nil
}

func (t *Team) listProblemsScore(ctx context.Context, eff ProblemReader) ([]*TeamProblemScoreData, error) {
	score, err := eff.ListProblemsScoreByTeamID(ctx, uuid.UUID(t.teamID))
	if err != nil {
		return nil, err
	}
	return score, nil
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

func (tp *TeamProblemScoreData) parse(problem *Problem) (*Score, error) {
	if tp.ProblemID != uuid.UUID(problem.problemID) {
		return nil, NewInvalidArgumentError("problem_id does not match", nil)
	}
	return (&ScoreData{
		MarkedScore: tp.MarkedScore,
	}).parse(problem)
}

func (tp *TeamProblem) Details(ctx context.Context, eff ProblemReader) (*TeamProblemDetail, error) {
	detail, err := tp.problem.DescriptiveProblem(ctx, eff)
	if err != nil {
		return nil, err
	}
	return &TeamProblemDetail{
		team:          tp.team,
		problemDetail: detail,
	}, nil
}

func (t *Team) ProblemDetailByCode(ctx context.Context, eff ProblemReader, code ProblemCode) (*TeamProblemDetail, error) {
	problem, err := code.Problem(ctx, eff)
	if err != nil {
		return nil, err
	}
	detail, err := problem.DescriptiveProblem(ctx, eff)
	if err != nil {
		return nil, err
	}
	return &TeamProblemDetail{
		team:          t,
		problemDetail: detail,
	}, nil
}

func (t *Team) ProblemByCode(ctx context.Context, eff ProblemReader, code ProblemCode) (*TeamProblem, error) {
	problem, err := code.Problem(ctx, eff)
	if err != nil {
		return nil, err
	}
	return &TeamProblem{
		team:    t,
		problem: problem,
	}, nil
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
