package domain

import (
	"context"
	"log/slog"
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

func (m *MarkingResult) IsPublic(now time.Time) bool {
	return now.After(m.answer.createdAt.Add(AnswerInterval))
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

type MarkingResultUpdatePenaltyEffect interface {
	MarkingResultPenaltyUpdator
	DeploymentReader
}

func (m *MarkingResult) UpdatePenalty(ctx context.Context, eff MarkingResultUpdatePenaltyEffect) error {
	deploymentRevision, err := m.Answer().TeamProblem().
		DeploymentCountAt(ctx, eff, m.Answer().CreatedAt())
	if err != nil {
		return err
	}
	penalty := m.Answer().Problem().Penalty(uint32(deploymentRevision)) //nolint:gosec
	if m.Score().Penalty() != penalty {
		slog.InfoContext(ctx, "Update penalty", "answer", m.Answer().id, "prev", m.Score().Penalty(), "next", penalty)
		if err := eff.UpdatePenalty(ctx, m.id, penalty); err != nil {
			return WrapAsInternal(err, "failed to update penalty")
		}
		m.Score().penalty = penalty
	}
	return nil
}

func (a *Answer) latestMarkingResult(ctx context.Context, eff MarkingResultReader) (*MarkingResult, error) {
	results, err := ListAllMarkingResults(ctx, eff)
	if err != nil {
		return nil, err
	}
	var latest *MarkingResult
	for _, result := range results {
		if result.Answer().id != a.id {
			continue
		}
		if latest == nil || result.CreatedAt().After(latest.CreatedAt()) {
			latest = result
		}
	}
	if latest == nil {
		return nil, NewNotFoundError("marking result for answer", nil)
	}
	return latest, nil
}

func (a *Answer) latestPublicMarkingResult(ctx context.Context, eff MarkingResultReader, now time.Time) (*MarkingResult, error) {
	results, err := ListAllMarkingResults(ctx, eff)
	if err != nil {
		return nil, err
	}
	var latest *MarkingResult
	for _, result := range results {
		if result.Answer().id != a.id {
			continue
		}
		if !result.IsPublic(now) {
			continue
		}
		if latest == nil || result.CreatedAt().After(latest.CreatedAt()) {
			latest = result
		}
	}
	if latest == nil {
		return nil, NewNotFoundError("marking result for answer", nil)
	}
	return latest, nil
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
		ID        uuid.UUID             `json:"id"`
		Judge     string                `json:"judge"`
		Answer    *AnswerData           `json:"answer"`
		Score     *ScoreData            `json:"score"`
		Rationale *MarkingRationaleData `json:"rationale"`
		CreatedAt time.Time             `json:"created_at"`
	}
	MarkingRationaleData struct {
		DescriptiveComment string `json:"descriptive_comment,omitempty"`
	}

	MarkingResultReader interface {
		ListMarkingResults(ctx context.Context) ([]*MarkingResultData, error)
	}
	MarkingResultPenaltyUpdator interface {
		UpdatePenalty(ctx context.Context, id uuid.UUID, penalty uint32) error
	}
	MarkingResultWriter interface {
		CreateMarkingResult(ctx context.Context, m *MarkingResultData) error
		MarkingResultPenaltyUpdator
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
