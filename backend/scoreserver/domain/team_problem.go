package domain

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type (
	TeamProblem struct {
		*problem
		team  *Team
		score *TeamProblemScore
	}
	TeamProblemDetail struct {
		team          *Team
		problemDetail *DescriptiveProblem
	}
	TeamProblemScore struct {
		problemID   uuid.UUID
		markedScore uint32
		penalty     uint32
		totalScore  uint32
		maxScore    uint32
	}
	TeamProblemScoreData struct {
		ProblemID   uuid.UUID `db:"problem_id"`
		MarkedScore uint32    `db:"marked_score"`
		Penalty     uint32    `db:"penalty"`
		TotalScore  uint32    `db:"total_score"`
		MaxScore    uint32    `db:"max_score"`
	}
)

func (t *Team) Problems(ctx context.Context, eff ProblemReader) ([]*TeamProblem, error) {
	ps, err := ListProblems(ctx, eff)
	if err != nil {
		return nil, err
	}
	scoreData, err := t.listProblemsScore(ctx, eff)
	if err != nil {
		return nil, err
	}
	scoreMap := make(map[uuid.UUID]*TeamProblemScore, len(scoreData))
	for _, s := range scoreData {
		scoreMap[s.ProblemID] = s.parse()
	}
	tps := make([]*TeamProblem, 0, len(ps))
	for _, p := range ps {
		tps = append(tps, &TeamProblem{team: t, problem: p, score: scoreMap[uuid.UUID(p.problemID)]})
	}
	return tps, nil
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

func (tp *TeamProblem) Score() *TeamProblemScore {
	return tp.score
}

func (tp *TeamProblemScore) MarkedScore() uint32 {
	return tp.markedScore
}

func (tp *TeamProblemScore) Penalty() uint32 {
	return tp.penalty
}

func (tp *TeamProblemScore) TotalScore() uint32 {
	return tp.totalScore
}

func (tp *TeamProblemScore) MaxScore() uint32 {
	return tp.maxScore
}

func (tp *TeamProblemScoreData) parse() *TeamProblemScore {
	return &TeamProblemScore{
		problemID:   tp.ProblemID,
		markedScore: tp.MarkedScore,
		penalty:     tp.Penalty,
		totalScore:  tp.TotalScore,
		maxScore:    tp.MaxScore,
	}
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

func (tp *TeamProblemDetail) Team() *Team {
	return tp.team
}

func (tp *TeamProblemDetail) ProblemDetail() *DescriptiveProblem {
	return tp.problemDetail
}
