package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/gofrs/uuid/v5"
)

type problemID uuid.UUID

func (id problemID) Delete(ctx context.Context, eff ProblemWriter) error {
	if err := eff.DeleteProblem(ctx, uuid.UUID(id)); err != nil {
		return WrapAsInternal(err, "failed to delete problem")
	}
	return nil
}

type ProblemCode string

var (
	problemCodeRegexp = regexp.MustCompile(`^[a-zA-Z]{3}$`)
)

func NewProblemCode(code string) (ProblemCode, error) {
	if !problemCodeRegexp.MatchString(code) {
		return "", NewInvalidArgumentError("invalid problem code", nil)
	}
	code = strings.ToUpper(code)
	return ProblemCode(code), nil
}

func (pc ProblemCode) Problem(ctx context.Context, eff ProblemReader) (*Problem, error) {
	data, err := eff.GetProblemByCode(ctx, string(pc))
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get descriptive problem")
	}
	return data.parse()
}

type ProblemType int

const (
	ProblemTypeUnknown     ProblemType = iota
	ProblemTypeDescriptive             // 記述問題
)

var (
	_ fmt.Stringer     = ProblemType(0)
	_ json.Marshaler   = ProblemType(0)
	_ json.Unmarshaler = (*ProblemType)(nil)
)

func (pt ProblemType) String() string {
	switch pt {
	case ProblemTypeDescriptive:
		return "Descriptive"
	case ProblemTypeUnknown:
		fallthrough
	default:
		return "Unknown"
	}
}

func (pt ProblemType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + pt.String() + `"`), nil
}

func (pt *ProblemType) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	switch s {
	case "Descriptive":
		*pt = ProblemTypeDescriptive
	default:
		*pt = ProblemTypeUnknown
	}
	return nil
}

type RedeployRule int

const (
	RedeployRuleUnknown           RedeployRule = iota
	RedeployRuleUnredeployable                 // 再展開不可
	RedeployRulePercentagePenalty              // 再展開回数に応じて最大得点の一定割合を減点
	RedeployRuleManual                         // 再展開のペナルティを採点者が手動で計算
)

func (r RedeployRule) String() string {
	switch r {
	case RedeployRuleUnredeployable:
		return "Unredeployable"
	case RedeployRulePercentagePenalty:
		return "PercentagePenalty"
	case RedeployRuleManual:
		return "Manual"
	case RedeployRuleUnknown:
		fallthrough
	default:
		return "Unknown"
	}
}

func (r RedeployRule) MarshalJSON() ([]byte, error) {
	return []byte(`"` + r.String() + `"`), nil
}

func (r *RedeployRule) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	switch s {
	case "Unredeployable":
		*r = RedeployRuleUnredeployable
	case "PercentagePenalty":
		*r = RedeployRulePercentagePenalty
	case "Manual":
		*r = RedeployRuleManual
	default:
		*r = RedeployRuleUnknown
	}
	return nil
}

type (
	Problem struct {
		problemID
		code        ProblemCode
		problemType ProblemType
		title       string
		maxScore    uint32
		category    string

		redeployRule      RedeployRule
		percentagePenalty *RedeployPenaltyPercentage
	}
	RedeployPenaltyPercentage struct {
		Threshold  uint32 `json:"threshold"`
		Percentage uint32 `json:"percentage"`
	}
	problem = Problem
)

func ListProblems(ctx context.Context, eff ProblemReader) ([]*Problem, error) {
	data, err := eff.ListProblems(ctx)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to list problems")
	}
	problems := make([]*Problem, 0, len(data))
	for _, d := range data {
		problem, err := d.parse()
		if err != nil {
			return nil, err
		}
		problems = append(problems, problem)
	}
	return problems, nil
}

func (p *Problem) Code() ProblemCode {
	return p.code
}

func (p *Problem) Title() string {
	return p.title
}

func (p *Problem) MaxScore() uint32 {
	return p.maxScore
}

func (p *Problem) Category() string {
	return p.category
}

func (p *Problem) Type() ProblemType {
	return p.problemType
}

func (p *Problem) RedeployRule() RedeployRule {
	return p.redeployRule
}

func (p *Problem) Redeployable() bool {
	return p.redeployRule != RedeployRuleUnredeployable
}

func (p *Problem) PercentagePenalty() *RedeployPenaltyPercentage {
	if p.redeployRule != RedeployRulePercentagePenalty {
		return nil
	}
	penalty := *p.percentagePenalty
	return &penalty
}

func (p *Problem) RemainingDeployments(revision uint32) int32 {
	switch p.redeployRule {
	case RedeployRuleUnredeployable:
		return 0
	case RedeployRulePercentagePenalty:
		threshold := p.percentagePenalty.Threshold
		return int32(int64(threshold) - int64(revision)) //nolint:gosec // uint32 - uint32 なので int32 に収まる
	case RedeployRuleManual:
		return 0 // マニュアルモードでは常に0(再展開ペナルティは常に発生する可能性がある)
	default:
		panic(fmt.Sprintf("unknown redeploy rule: %s", p.redeployRule.String()))
	}
}

func (p *Problem) Penalty(deploymentCount uint32) uint32 {
	switch p.redeployRule {
	case RedeployRulePercentagePenalty:
		penalty := (p.maxScore * p.percentagePenalty.Percentage) / 100            //nolint:mnd
		return penalty * uint32(max(0, -p.RemainingDeployments(deploymentCount))) //nolint:gosec
	case RedeployRuleManual, RedeployRuleUnredeployable:
		return 0
	default:
		panic(fmt.Sprintf("unknown redeploy rule: %s", p.redeployRule.String()))
	}
}

type (
	ProblemContent struct {
		body        string
		explanation string
	}
	problemContent = ProblemContent
)

func NewProblemContent(body, explanation string) (*ProblemContent, error) {
	if body == "" {
		return nil, NewInvalidArgumentError("body is required", nil)
	}
	if explanation == "" {
		return nil, NewInvalidArgumentError("explanation is required", nil)
	}
	return &ProblemContent{
		body:        body,
		explanation: explanation,
	}, nil
}

func (pc *ProblemContent) Body() string {
	return pc.body
}

func (pc *ProblemContent) Explanation() string {
	return pc.explanation
}

type DescriptiveProblem struct {
	*problem
	*problemContent
}

func (dp *DescriptiveProblem) Problem() *Problem {
	return dp.problem
}

func (dp *DescriptiveProblem) Content() *ProblemContent {
	return dp.problemContent
}

type CreateDescriptiveProblemInput struct {
	Code              ProblemCode
	Title             string
	MaxScore          uint32
	Category          string
	RedeployRule      RedeployRule
	PercentagePenalty *RedeployPenaltyPercentage
	Content           *ProblemContent
}

func CreateDescriptiveProblem(input CreateDescriptiveProblemInput) (*DescriptiveProblem, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, WrapAsInternal(err, "failed to generate ID")
	}

	problem, err := (&ProblemData{
		ID:                id,
		Code:              string(input.Code),
		ProblemType:       ProblemTypeDescriptive,
		Title:             input.Title,
		MaxScore:          input.MaxScore,
		Category:          input.Category,
		RedeployRule:      input.RedeployRule,
		PercentagePenalty: input.PercentagePenalty,
	}).parse()
	if err != nil {
		return nil, err
	}

	if input.Content == nil {
		return nil, NewInvalidArgumentError("content is required", nil)
	}

	return &DescriptiveProblem{problem: problem, problemContent: input.Content}, nil
}

type UpdateDescriptiveProblemInput struct {
	Code              *ProblemCode
	Title             string
	MaxScore          uint32
	Category          string
	RedeployRule      RedeployRule
	PercentagePenalty *RedeployPenaltyPercentage
	Body              string
	Explanation       string
}

func (dp *DescriptiveProblem) Update(input UpdateDescriptiveProblemInput) (*DescriptiveProblem, error) {
	data := dp.Data()

	if input.Code != nil {
		data.Problem.Code = string(*input.Code)
	}
	if input.Title != "" {
		data.Problem.Title = input.Title
	}
	if input.MaxScore != 0 && input.MaxScore != data.Problem.MaxScore {
		// TODO: 最大得点の更新
		// 最大得点が更新される場合，既に採点が行われているならば，それらの採点結果を適切に再計算するか，更新を拒否するかを決定する必要がある
		// この操作の必要性が認識され，どのような挙動を取るべきかが明確になるまで，最大得点の更新は許可しない
		return nil, NewInvalidArgumentError("max score cannot be updated", nil)
	}
	if input.Category != "" {
		data.Problem.Category = input.Category
	}
	if input.RedeployRule != RedeployRuleUnknown && input.RedeployRule != data.Problem.RedeployRule {
		// TODO: 再展開ルールを変更する場合，Undeployable になるならば再展開が既に行われていないことを確認する必要がある
		// これは複雑であるため，操作の必要性が認識されるまで再展開ルールの変更は許可しない
		return nil, NewInvalidArgumentError("redeploy rule cannot be updated", nil)
	}
	if input.PercentagePenalty != nil && *input.PercentagePenalty != *data.Problem.PercentagePenalty {
		// TODO(#1298): 減点率の変更
		// 減点率が変更される場合，既に行われている採点結果を再計算する必要がある
		// この操作は問題難易度による微調整として必要性が高いため，優先度が高いが，現状は採点ロジックが未実装であるため，許可しない
		return nil, NewInvalidArgumentError("percentage penalty cannot be updated", nil)
	}
	if input.Body != "" {
		data.Content.Body = input.Body
	}
	if input.Explanation != "" {
		data.Content.Explanation = input.Explanation
	}

	return data.parse()
}

func (p *Problem) DescriptiveProblem(ctx context.Context, eff ProblemReader) (*DescriptiveProblem, error) {
	if p.problemType != ProblemTypeDescriptive {
		return nil, NewInvalidArgumentError("problem type must be descriptive", nil)
	}

	data, err := eff.GetDescriptiveProblem(ctx, uuid.UUID(p.problemID))
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get descriptive problem")
	}
	return data.parse()
}

func (dp *DescriptiveProblem) Save(ctx context.Context, eff ProblemWriter) error {
	if err := eff.SaveDescriptiveProblem(ctx, dp.Data()); err != nil {
		return WrapAsInternal(err, "failed to save descriptive problem")
	}
	return nil
}

type (
	ProblemReader interface {
		ListProblems(ctx context.Context) ([]*ProblemData, error)
		GetProblemByCode(ctx context.Context, code string) (*ProblemData, error)
		GetDescriptiveProblem(ctx context.Context, id uuid.UUID) (*DescriptiveProblemData, error)
	}
	ProblemWriter interface {
		ProblemReader
		SaveDescriptiveProblem(ctx context.Context, data *DescriptiveProblemData) error
		DeleteProblem(ctx context.Context, id uuid.UUID) error
	}
)

type ProblemData struct {
	ID          uuid.UUID   `json:"id"`
	Code        string      `json:"code"`
	ProblemType ProblemType `json:"problem_type"`
	Title       string      `json:"title"`
	MaxScore    uint32      `json:"max_score"`
	Category    string      `json:"category"`

	RedeployRule      RedeployRule               `json:"redeploy_rule"`
	PercentagePenalty *RedeployPenaltyPercentage `json:"percentage_penalty,omitempty"`
}

func (d *ProblemData) parse() (*problem, error) {
	code, err := NewProblemCode(d.Code)
	if err != nil {
		return nil, err
	}

	if d.ProblemType == ProblemTypeUnknown {
		return nil, NewInvalidArgumentError("problem type is required", nil)
	}

	if d.Title == "" {
		return nil, NewInvalidArgumentError("title is required", nil)
	}

	if d.MaxScore == 0 {
		return nil, NewInvalidArgumentError("max score is required", nil)
	}

	if d.Category == "" {
		return nil, NewInvalidArgumentError("category is required", nil)
	}

	if d.RedeployRule == RedeployRuleUnknown {
		return nil, NewInvalidArgumentError("redeploy rule is required", nil)
	}

	var percentagePenalty *RedeployPenaltyPercentage
	if d.RedeployRule == RedeployRulePercentagePenalty {
		if d.PercentagePenalty == nil {
			return nil, NewInvalidArgumentError("percentage penalty is required", nil)
		}
		p := d.PercentagePenalty.Percentage
		if p > 100 { //nolint:mnd
			return nil, NewInvalidArgumentError("invalid percentage penalty", nil)
		}
		if (d.MaxScore*p)%100 != 0 {
			// 減点は整数でなければならない
			return nil, NewInvalidArgumentError("percentage penalty must be a value that makes the deduction points an integer", nil)
		}
		percentagePenalty = d.PercentagePenalty
	}

	return &problem{
		problemID:   problemID(d.ID),
		code:        code,
		problemType: d.ProblemType,
		title:       d.Title,
		maxScore:    d.MaxScore,
		category:    d.Category,

		redeployRule:      d.RedeployRule,
		percentagePenalty: percentagePenalty,
	}, nil
}

func (p *Problem) Data() *ProblemData {
	return &ProblemData{
		ID:          uuid.UUID(p.problemID),
		Code:        string(p.code),
		ProblemType: p.problemType,
		Title:       p.title,
		MaxScore:    p.maxScore,
		Category:    p.category,

		RedeployRule:      p.redeployRule,
		PercentagePenalty: p.PercentagePenalty(),
	}
}

type ProblemContentData struct {
	Body        string `json:"body"`
	Explanation string `json:"explanation"`
}

func (d *ProblemContentData) parse() (*ProblemContent, error) {
	if d.Body == "" {
		return nil, NewInvalidArgumentError("body is required", nil)
	}
	if d.Explanation == "" {
		return nil, NewInvalidArgumentError("explanation is required", nil)
	}
	return &ProblemContent{
		body:        d.Body,
		explanation: d.Explanation,
	}, nil
}

func (pc *ProblemContent) Data() *ProblemContentData {
	return &ProblemContentData{
		Body:        pc.body,
		Explanation: pc.explanation,
	}
}

type DescriptiveProblemData struct {
	Problem *ProblemData        `json:"problem"`
	Content *ProblemContentData `json:"content"`
}

func (d *DescriptiveProblemData) parse() (*DescriptiveProblem, error) {
	if d.Problem.ProblemType != ProblemTypeDescriptive {
		return nil, NewInvalidArgumentError("problem type must be descriptive", nil)
	}

	problem, err := d.Problem.parse()
	if err != nil {
		return nil, err
	}

	content, err := d.Content.parse()
	if err != nil {
		return nil, err
	}

	return &DescriptiveProblem{
		problem:        problem,
		problemContent: content,
	}, nil
}

func (dp *DescriptiveProblem) Data() *DescriptiveProblemData {
	return &DescriptiveProblemData{
		Problem: dp.problem.Data(),
		Content: dp.problemContent.Data(),
	}
}
