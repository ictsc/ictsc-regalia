package domain

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	MarkingResult struct {
		id        uuid.UUID
		answer    *Answer
		score     *Score
		rationale *MarkingRationale
		judge     string
		createdAt time.Time
	}
	Score struct {
		max    uint32
		marked uint32
	}
	MarkingRationale struct {
		problemType ProblemType
		descriptive *DescriptiveMarkingRaionale
	}
	DescriptiveMarkingRaionale struct {
		comment string
	}
)

func (m *MarkingResult) Answer() *Answer {
	return m.answer
}

func (m *MarkingResult) Score() *Score {
	return m.score
}

func (m *MarkingResult) Judge() string {
	return m.judge
}

func (m *MarkingResult) Rationale() *MarkingRationale {
	return m.rationale
}

func (m *MarkingResult) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Score) MarkedScore() uint32 {
	return m.marked
}

func (m *Score) MaxScore() uint32 {
	return m.max
}

func (m *Score) Penalty() uint32 {
	return 0 // 現状は再展開がないので0
}

func (m *MarkingRationale) Type() ProblemType {
	return m.problemType
}

func (m *MarkingRationale) Descriptive() *DescriptiveMarkingRaionale {
	if m.problemType != ProblemTypeDescriptive {
		return &DescriptiveMarkingRaionale{}
	}
	return m.descriptive
}

func (d *DescriptiveMarkingRaionale) Comment() string {
	return d.comment
}

func ListAllMarkingResults(ctx context.Context, eff MarkingResultReader) ([]*MarkingResult, error) {
	markingResultDataList, err := eff.ListMarkingResults(ctx)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to list marking results")
	}

	markingResults := make([]*MarkingResult, 0, len(markingResultDataList))
	for _, data := range markingResultDataList {
		markingResult, err := data.parse()
		if err != nil {
			return nil, err
		}
		markingResults = append(markingResults, markingResult)
	}
	return markingResults, nil
}

type MarkInput struct {
	Score uint32
	Judge string

	// Descriptive
	Comment string
}

func (a *Answer) Mark(ctx context.Context, eff MarkingResultWriter, now time.Time, input *MarkInput) (*MarkingResult, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, WrapAsInternal(err, "failed to generate uuid")
	}

	markingResult, err := (&MarkingResultData{
		ID:     id,
		Judge:  input.Judge,
		Answer: a.Data(),
		Score:  &ScoreData{MarkedScore: input.Score},
		Rationale: &MarkingRationaleData{
			DescriptiveComment: input.Comment,
		},
		CreatedAt: now,
	}).parse()
	if err != nil {
		return nil, err
	}

	if err := eff.CreateMarkingResult(ctx, markingResult.Data()); err != nil {
		return nil, WrapAsInternal(err, "failed to create marking result")
	}

	return markingResult, nil
}

type (
	MarkingResultData struct {
		ID        uuid.UUID
		Judge     string
		Answer    *AnswerData
		Score     *ScoreData
		Rationale *MarkingRationaleData
		CreatedAt time.Time
	}
	ScoreData struct {
		MarkedScore uint32
	}
	MarkingRationaleData struct {
		DescriptiveComment string
	}

	MarkingResultReader interface {
		ListMarkingResults(ctx context.Context) ([]*MarkingResultData, error)
	}
	MarkingResultWriter interface {
		CreateMarkingResult(ctx context.Context, m *MarkingResultData) error
	}
)

func (m *MarkingResult) Data() *MarkingResultData {
	return &MarkingResultData{
		ID:        m.id,
		Judge:     m.judge,
		Answer:    m.answer.Data(),
		Score:     &ScoreData{MarkedScore: m.score.MarkedScore()},
		Rationale: m.rationale.data(),
		CreatedAt: m.createdAt,
	}
}

func (m *MarkingRationale) data() *MarkingRationaleData {
	switch m.problemType {
	case ProblemTypeDescriptive:
		return &MarkingRationaleData{DescriptiveComment: m.descriptive.Comment()}
	case ProblemTypeUnknown:
		fallthrough
	default:
		return &MarkingRationaleData{}
	}
}

func (m *MarkingResultData) parse() (*MarkingResult, error) {
	answer, err := m.Answer.parse()
	if err != nil {
		return nil, err
	}

	score, err := m.Score.parse(answer.Problem())
	if err != nil {
		return nil, err
	}

	rationale := m.Rationale.parse(answer.Problem().Type())

	return &MarkingResult{
		id:        m.ID,
		judge:     m.Judge,
		answer:    answer,
		score:     score,
		rationale: rationale,
		createdAt: m.CreatedAt,
	}, nil
}

func (s *ScoreData) parse(problem *Problem) (*Score, error) {
	if s.MarkedScore > problem.MaxScore() {
		return nil, NewInvalidArgumentError("marked score is over max score", nil)
	}
	return &Score{marked: s.MarkedScore}, nil
}

func (r *MarkingRationaleData) parse(problemType ProblemType) *MarkingRationale {
	switch problemType {
	case ProblemTypeDescriptive:
		return &MarkingRationale{
			problemType: ProblemTypeDescriptive,
			descriptive: &DescriptiveMarkingRaionale{comment: r.DescriptiveComment},
		}
	case ProblemTypeUnknown:
		fallthrough
	default:
		return &MarkingRationale{
			problemType: problemType,
		}
	}
}
