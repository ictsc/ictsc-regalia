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
		// 次の解答を受付可能にするまでの時間
		interval time.Duration
		score    *Score
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
	AnswerInterval = 1 * time.Millisecond
)

func ListAnswersForAdmin(ctx context.Context, eff AnswerReader) ([]*Answer, error) {
	answerDataList, err := eff.ListAnswersForAdmin(ctx)
	if err != nil {
		return nil, err
	}

	answers, err := parseAnswerDataList(answerDataList)
	if err != nil {
		return nil, err
	}

	return answers, nil
}

func ListAnswersByTeamProblemForPublic(ctx context.Context, eff AnswerReader, teamCode TeamCode, problemCode ProblemCode) ([]*Answer, error) {
	answerDataList, err := eff.ListAnswersByTeamProblemForPublic(ctx, int64(teamCode), string(problemCode))
	if err != nil {
		return nil, err
	}

	answers, err := parseAnswerDataList(answerDataList)
	if err != nil {
		return nil, err
	}

	return answers, nil
}

func (tp *TeamProblem) answersForPublic(ctx context.Context, eff AnswerReader) ([]*Answer, error) {
	return ListAnswersByTeamProblemForPublic(ctx, eff, tp.team.Code(), tp.problem.Code())
}

func (tp *TeamProblem) answersForAdmin(ctx context.Context, eff AnswerReader) ([]*Answer, error) {
	answerDataList, err := eff.ListAnswersByTeamProblemForAdmin(ctx,
		int64(tp.team.Code()), string(tp.problem.Code()))
	if err != nil {
		return nil, err
	}

	answers, err := parseAnswerDataList(answerDataList)
	if err != nil {
		return nil, err
	}

	return answers, nil
}

func parseAnswerDataList(answerDataList []*AnswerData) ([]*Answer, error) {
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

func (a *Answer) TeamProblem() *TeamProblem {
	return &TeamProblem{
		team:    a.team,
		problem: a.problem,
	}
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

func (a *Answer) Score() *Score {
	return a.score
}

func GetAnswerDetailForPublic(
	ctx context.Context, eff AnswerReader,
	teamCode TeamCode, problemCode ProblemCode, answerNumber uint32,
) (*AnswerDetail, error) {
	answerDetailData, err := eff.GetAnswerDetailForPublic(ctx, int64(teamCode), string(problemCode), answerNumber)
	if err != nil {
		return nil, err
	}

	return answerDetailData.parse()
}

func GetAnswerDetailForAdmin(
	ctx context.Context, eff AnswerReader,
	teamCode TeamCode, problemCode ProblemCode, answerNumber uint32,
) (*AnswerDetail, error) {
	answerDetailData, err := eff.GetAnswerDetailForAdmin(ctx, int64(teamCode), string(problemCode), answerNumber)
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
	prevAnswerData, err := eff.GetLatestAnswerByTeamProblemForPublic(
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
		ID        uuid.UUID     `json:"id"`
		Number    uint32        `json:"number"`
		Team      *TeamData     `json:"team"`
		Problem   *ProblemData  `json:"problem"`
		Author    *UserData     `json:"author"`
		CreatedAt time.Time     `json:"created_at"`
		Interval  time.Duration `json:"interval"`
		Score     *ScoreData    `json:"score,omitzero"`
	}
	AnswerDetailData struct {
		Answer *AnswerData     `json:"answer"`
		Body   *AnswerBodyData `json:"body"`
	}
	AnswerBodyData struct {
		Descriptive *DescriptiveAnswerBodyData `json:"descriptive,omitzero"`
	}
	DescriptiveAnswerBodyData struct {
		Body string `json:"body"`
	}
	AnswerReader interface {
		ListAnswersForAdmin(ctx context.Context) ([]*AnswerData, error)
		ListAnswersByTeamProblemForAdmin(ctx context.Context, teamCode int64, problemCode string) ([]*AnswerData, error)
		ListAnswersByTeamProblemForPublic(ctx context.Context, teamCode int64, problemCode string) ([]*AnswerData, error)
		GetAnswerDetailForAdmin(ctx context.Context, teamCode int64, problemCode string, answerNumber uint32) (*AnswerDetailData, error)
		GetAnswerDetailForPublic(ctx context.Context, teamCode int64, problemCode string, answerNumber uint32) (*AnswerDetailData, error)
	}
	AnswerWriter interface {
		GetLatestAnswerByTeamProblemForPublic(ctx context.Context, teamID, problemID uuid.UUID) (*AnswerData, error)
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

	var score *Score
	if d.Score != nil {
		score, err = d.Score.parse(problem)
		if err != nil {
			return nil, err
		}
	}

	return &Answer{
		id:        d.ID,
		number:    d.Number,
		team:      team,
		problem:   problem,
		author:    author,
		createdAt: d.CreatedAt,
		interval:  d.Interval,
		score:     score,
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
	body.problemType = answer.Problem().Type()
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
