package domain

import "context"

type (
	TeamProblem struct {
		team    *Team
		problem *Problem
	}
	TeamProblemDetail struct {
		team          *Team
		problemDetail *DescriptiveProblem
	}
)

func (t *Team) Problems(ctx context.Context, eff ProblemReader) ([]*TeamProblem, error) {
	ps, err := ListProblems(ctx, eff)
	if err != nil {
		return nil, err
	}
	tps := make([]*TeamProblem, 0, len(ps))
	for _, p := range ps {
		tps = append(tps, &TeamProblem{team: t, problem: p})
	}
	return tps, nil
}

func (tp *TeamProblem) Team() *Team {
	return tp.team
}

func (tp *TeamProblem) Problem() *Problem {
	return tp.problem
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
