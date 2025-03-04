package domain

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
)

type (
	Answer struct {
		id        uuid.UUID
		number    uint32
		team      *Team
		problem   *Problem
		author    *User
		createdAt time.Time
		// 次の回答を受付可能にするまでの時間
		interval time.Duration
	}
	answer = Answer

	AnswerDetail struct {
		*answer
		body *AnswerBody
	}
	AnswerBody struct {
		problemType ProblemType
		descriptive *DescriptiveAnswerBody
	}
	DescriptiveAnswerBody struct {
		body string
	}
)

const (
	AnswerInterval = 20 * time.Minute
)

func ListAnswers(ctx context.Context, eff AnswerReader) ([]*Answer, error) {
	answerDataList, err := eff.ListAnswers(ctx)
	if err != nil {
		return nil, err
	}

	answers := make([]*Answer, 0, len(answerDataList))
	for _, answerData := range answerDataList {
		answer, err := answerData.parse()
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, nil
}

func ListAnswersByTeamProblem(ctx context.Context, eff AnswerReader, teamCode TeamCode, problemCode ProblemCode) ([]*Answer, error) {
	answerDataList, err := eff.ListAnswersByTeamProblem(ctx, int64(teamCode), string(problemCode))
	if err != nil {
		return nil, err
	}

	answers := make([]*Answer, 0, len(answerDataList))
	for _, answerData := range answerDataList {
		answer, err := answerData.parse()
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, nil
}

func (a *Answer) Number() uint32 {
	return a.number
}

func (a *Answer) Team() *Team {
	return a.team
}

func (a *Answer) Problem() *Problem {
	return a.problem
}

func (a *Answer) Author() *TeamMember {
	return &TeamMember{
		user: a.author,
		team: a.team,
	}
}

func (a *Answer) CreatedAt() time.Time {
	return a.createdAt
}

func GetAnswerDetail(ctx context.Context, eff AnswerReader, teamCode TeamCode, problemCode ProblemCode, answerNumber uint32) (*AnswerDetail, error) {
	answerDetailData, err := eff.GetAnswerDetail(ctx, int64(teamCode), string(problemCode), answerNumber)
	if err != nil {
		return nil, err
	}

	return answerDetailData.parse()
}

func (a *AnswerDetail) Answer() *Answer {
	return a.answer
}

func (a *AnswerDetail) Body() *AnswerBody {
	return a.body
}

func (b *AnswerBody) Type() ProblemType {
	return b.problemType
}

func (b *AnswerBody) Descriptive() (*DescriptiveAnswerBody, error) {
	if b.problemType != ProblemTypeDescriptive {
		return nil, newInternalError(errors.New("problem type is not descriptive"))
	}
	return b.descriptive, nil
}

func (b *DescriptiveAnswerBody) Body() string {
	return b.body
}

func (tm *TeamMember) SubmitDescriptiveAnswer(
	ctx context.Context, now time.Time, eff AnswerWriter,
	problem *Problem, body string,
) (*AnswerDetail, error) {
	return tm.submitAnswer(ctx, now, eff, problem, &AnswerBodyData{
		Descriptive: &DescriptiveAnswerBodyData{
			Body: body,
		},
	})
}

func (tm *TeamMember) submitAnswer(
	ctx context.Context, now time.Time, eff AnswerWriter,
	problem *Problem, body *AnswerBodyData,
) (*AnswerDetail, error) {
	prevAnswerData, err := eff.GetLatestAnswerByTeamProblem(
		ctx, uuid.UUID(tm.team.teamID), uuid.UUID(problem.problemID))
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, errors.Wrap(err, "failed to get latest answer")
	}
	var number uint32
	if prevAnswerData != nil {
		prevAnswer, err := prevAnswerData.parse()
		if err != nil {
			return nil, err
		}
		if now.Before(prevAnswer.CreatedAt().Add(prevAnswer.interval)) {
			return nil, errors.WithStack(ErrTooEarlyToSubmitAnswer)
		}
		number = prevAnswer.Number() + 1
	} else {
		number = 1
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, WrapAsInternal(err, "failed to generate answer ID")
	}

	answer, err := (&AnswerDetailData{
		Answer: &AnswerData{
			ID:        id,
			Number:    number,
			Team:      tm.team.Data(),
			Problem:   problem.Data(),
			Author:    tm.user.Data(),
			CreatedAt: now,
			Interval:  AnswerInterval,
		},
		Body: body,
	}).parse()
	if err != nil {
		return nil, err
	}

	if err := eff.CreateAnswer(ctx, answer.Data()); err != nil {
		return nil, WrapAsInternal(err, "failed to create answer")
	}

	return answer, nil
}

var (
	ErrTooEarlyToSubmitAnswer = NewInvalidArgumentError("too early to submit answer", nil)
)

type (
	AnswerData struct {
		ID        uuid.UUID
		Number    uint32
		Team      *TeamData
		Problem   *ProblemData
		Author    *UserData
		CreatedAt time.Time
		Interval  time.Duration
	}
	AnswerDetailData struct {
		Answer *AnswerData
		Body   *AnswerBodyData
	}
	AnswerBodyData struct {
		Descriptive *DescriptiveAnswerBodyData
	}
	DescriptiveAnswerBodyData struct {
		Body string
	}
	AnswerReader interface {
		ListAnswers(ctx context.Context) ([]*AnswerData, error)
		ListAnswersByTeamProblem(ctx context.Context, teamCode int64, problemCode string) ([]*AnswerData, error)
		GetAnswerDetail(ctx context.Context, teamCode int64, problemCode string, answerNumber uint32) (*AnswerDetailData, error)
	}
	AnswerWriter interface {
		GetLatestAnswerByTeamProblem(ctx context.Context, teamID, problemID uuid.UUID) (*AnswerData, error)
		CreateAnswer(ctx context.Context, data *AnswerDetailData) error
	}
)

func (d *AnswerData) parse() (*Answer, error) {
	team, err := d.Team.parse()
	if err != nil {
		return nil, err
	}

	problem, err := d.Problem.parse()
	if err != nil {
		return nil, err
	}

	author, err := d.Author.parse()
	if err != nil {
		return nil, err
	}

	return &Answer{
		id:        d.ID,
		number:    d.Number,
		team:      team,
		problem:   problem,
		author:    author,
		createdAt: d.CreatedAt,
		interval:  d.Interval,
	}, nil
}

func (a *Answer) Data() *AnswerData {
	return &AnswerData{
		ID:        a.id,
		Number:    a.number,
		Team:      a.team.Data(),
		Problem:   a.problem.Data(),
		Author:    a.author.Data(),
		CreatedAt: a.createdAt,
		Interval:  a.interval,
	}
}

func (d *AnswerDetailData) parse() (*AnswerDetail, error) {
	answer, err := d.Answer.parse()
	if err != nil {
		return nil, err
	}

	var body AnswerBody
	switch answer.Problem().Type() {
	case ProblemTypeDescriptive:
		if d.Body.Descriptive == nil {
			return nil, NewInvalidArgumentError("body.descriptive is required", nil)
		}
		descBody := d.Body.Descriptive.Body
		if descBody == "" {
			return nil, NewInvalidArgumentError("body.descriptive.body is required", nil)
		}
		body.descriptive = &DescriptiveAnswerBody{
			body: d.Body.Descriptive.Body,
		}
	case ProblemTypeUnknown:
		break
	}

	return &AnswerDetail{
		answer: answer,
		body:   &body,
	}, nil
}

func (a *AnswerDetail) Data() *AnswerDetailData {
	var bodyData AnswerBodyData
	if a.body.descriptive != nil {
		bodyData.Descriptive = &DescriptiveAnswerBodyData{
			Body: a.body.descriptive.body,
		}
	}
	return &AnswerDetailData{
		Answer: a.answer.Data(),
		Body:   &bodyData,
	}
}
