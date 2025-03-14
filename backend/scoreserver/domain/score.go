package domain

type Score struct {
	max     uint32
	marked  uint32
	penalty uint32
}

func (s *Score) MarkedScore() uint32 {
	return s.marked
}

func (s *Score) MaxScore() uint32 {
	return s.max
}

func (s *Score) Penalty() uint32 {
	return s.penalty // 現状は再展開がないので0
}

func (s *Score) TotalScore() uint32 {
	return max(0, s.MarkedScore()-s.Penalty())
}

type ScoreData struct {
	MarkedScore uint32 `json:"marked_score"`
	Penalty     uint32 `json:"penalty"`
}

func (s *Score) Data() *ScoreData {
	return &ScoreData{
		MarkedScore: s.marked,
		Penalty:     s.penalty,
	}
}

func (s *ScoreData) parse(problem *Problem) (*Score, error) {
	if s.MarkedScore > problem.MaxScore() {
		return nil, NewInvalidArgumentError("marked score is over max score", nil)
	}
	return &Score{
		max:    problem.maxScore,
		marked: s.MarkedScore,
	}, nil
}
