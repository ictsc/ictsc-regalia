package admin

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

type ProblemServiceHandler struct {
	Enforcer     *auth.Enforcer
	ListEffect   ProblemListEffect
	GetEffect    ProblemGetEffect
	CreateEffect ProblemCreateEffect
	UpdateEffect ProblemUpdateEffect
	DeleteEffect ProblemDeleteEffect

	adminv1connect.UnimplementedProblemServiceHandler
}

var _ adminv1connect.ProblemServiceHandler = (*ProblemServiceHandler)(nil)

func NewProblemServiceHandler(
	enforcer *auth.Enforcer,
	repo *pg.Repository,
) *ProblemServiceHandler {
	createEffect := pg.Tx(repo, func(rt *pg.RepositoryTx) domain.ProblemWriter { return rt })
	return &ProblemServiceHandler{
		Enforcer:     enforcer,
		ListEffect:   repo,
		GetEffect:    repo,
		CreateEffect: createEffect,
		UpdateEffect: createEffect,
		DeleteEffect: createEffect,
	}
}

type ProblemListEffect = domain.ProblemReader

func (h *ProblemServiceHandler) ListProblems(
	ctx context.Context,
	req *connect.Request[adminv1.ListProblemsRequest],
) (*connect.Response[adminv1.ListProblemsResponse], error) {
	if err := enforce(ctx, h.Enforcer, "problems", "list"); err != nil {
		return nil, err
	}

	problems, err := domain.ListProblems(ctx, h.ListEffect)
	if err != nil {
		return nil, err
	}

	protoProblems := make([]*adminv1.Problem, 0, len(problems))
	for _, problem := range problems {
		protoProblems = append(protoProblems, convertProblem(problem))
	}
	return connect.NewResponse(&adminv1.ListProblemsResponse{
		Problems: protoProblems,
	}), nil
}

type ProblemGetEffect = domain.ProblemReader

func (h *ProblemServiceHandler) GetProblem(
	ctx context.Context,
	req *connect.Request[adminv1.GetProblemRequest],
) (*connect.Response[adminv1.GetProblemResponse], error) {
	if err := enforce(ctx, h.Enforcer, "problems", "get"); err != nil {
		return nil, err
	}

	protoCode := req.Msg.GetCode()
	if protoCode == "" {
		return nil, domain.NewInvalidArgumentError("code is required", nil)
	}

	code, err := domain.NewProblemCode(protoCode)
	if err != nil {
		return nil, err
	}

	problem, err := code.Problem(ctx, h.GetEffect)
	if err != nil {
		return nil, err
	}

	switch problem.Type() {
	case domain.ProblemTypeDescriptive:
		descriptiveProblem, err := problem.DescriptiveProblem(ctx, h.GetEffect)
		if err != nil {
			return nil, err
		}
		return connect.NewResponse(&adminv1.GetProblemResponse{
			Problem: convertDescriptiveProblem(descriptiveProblem),
		}), nil
	case domain.ProblemTypeUnknown:
		fallthrough
	default:
		return connect.NewResponse(&adminv1.GetProblemResponse{
			Problem: convertProblem(problem),
		}), nil
	}
}

type ProblemCreateEffect interface {
	domain.Tx[domain.ProblemWriter]
}

func (h *ProblemServiceHandler) CreateProblem(
	ctx context.Context,
	req *connect.Request[adminv1.CreateProblemRequest],
) (*connect.Response[adminv1.CreateProblemResponse], error) {
	if err := enforce(ctx, h.Enforcer, "problems", "create"); err != nil {
		return nil, err
	}

	code, err := domain.NewProblemCode(req.Msg.GetProblem().GetCode())
	if err != nil {
		return nil, err
	}

	if req.Msg.GetProblem().GetBody().GetType() != adminv1.ProblemType_PROBLEM_TYPE_DESCRIPTIVE {
		return nil, domain.NewInvalidArgumentError("only descriptive problem is supported", nil)
	}

	content, err := domain.NewProblemContent(
		req.Msg.GetProblem().GetBody().GetDescriptive().GetProblemMarkdown(),
		req.Msg.GetProblem().GetBody().GetDescriptive().GetExplanationMarkdown(),
	)
	if err != nil {
		return nil, err
	}

	rule, penalty := parseRedeployRule(req.Msg.GetProblem().GetRedeployRule())
	descriptiveProblem, err := domain.CreateDescriptiveProblem(domain.CreateDescriptiveProblemInput{
		Code:              code,
		Title:             req.Msg.GetProblem().GetTitle(),
		MaxScore:          req.Msg.GetProblem().GetMaxScore(),
		Category:          req.Msg.GetProblem().GetCategory(),
		RedeployRule:      rule,
		PercentagePenalty: penalty,
		Content:           content,
	})
	if err != nil {
		return nil, err
	}

	if err := h.CreateEffect.RunInTx(ctx, func(tx domain.ProblemWriter) error {
		return descriptiveProblem.Save(ctx, tx)
	}); err != nil {
		return nil, err
	}

	return connect.NewResponse(&adminv1.CreateProblemResponse{
		Problem: convertDescriptiveProblem(descriptiveProblem),
	}), nil
}

type ProblemUpdateEffect = ProblemCreateEffect

func (h *ProblemServiceHandler) UpdateProblem(
	ctx context.Context,
	req *connect.Request[adminv1.UpdateProblemRequest],
) (*connect.Response[adminv1.UpdateProblemResponse], error) {
	if err := enforce(ctx, h.Enforcer, "problems", "update"); err != nil {
		return nil, err
	}

	code, err := domain.NewProblemCode(req.Msg.GetProblem().GetCode())
	if err != nil {
		return nil, err
	}

	rule, penalty := parseRedeployRule(req.Msg.GetProblem().GetRedeployRule())

	input := domain.UpdateDescriptiveProblemInput{
		Title:             req.Msg.GetProblem().GetTitle(),
		MaxScore:          req.Msg.GetProblem().GetMaxScore(),
		Category:          req.Msg.GetProblem().GetCategory(),
		RedeployRule:      rule,
		PercentagePenalty: penalty,
		Body:              req.Msg.GetProblem().GetBody().GetDescriptive().GetProblemMarkdown(),
		Explanation:       req.Msg.GetProblem().GetBody().GetDescriptive().GetExplanationMarkdown(),
	}

	descriptiveProblem, err := domain.RunTx(ctx, h.UpdateEffect, func(tx domain.ProblemWriter) (*domain.DescriptiveProblem, error) {
		problem, err := code.Problem(ctx, tx)
		if err != nil {
			return nil, err
		}
		descriptiveProblem, err := problem.DescriptiveProblem(ctx, tx)
		if err != nil {
			return nil, err
		}

		descriptiveProblem, err = descriptiveProblem.Update(input)
		if err != nil {
			return nil, err
		}

		if err := descriptiveProblem.Save(ctx, tx); err != nil {
			return nil, err
		}

		return descriptiveProblem, nil
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&adminv1.UpdateProblemResponse{
		Problem: convertDescriptiveProblem(descriptiveProblem),
	}), nil
}

type ProblemDeleteEffect = domain.Tx[domain.ProblemWriter]

func (h *ProblemServiceHandler) DeleteProblem(
	ctx context.Context,
	req *connect.Request[adminv1.DeleteProblemRequest],
) (*connect.Response[adminv1.DeleteProblemResponse], error) {
	if err := enforce(ctx, h.Enforcer, "problems", "delete"); err != nil {
		return nil, err
	}

	code, err := domain.NewProblemCode(req.Msg.GetCode())
	if err != nil {
		return nil, err
	}

	if err := h.DeleteEffect.RunInTx(ctx, func(tx domain.ProblemWriter) error {
		problem, err := code.Problem(ctx, tx)
		if err != nil {
			return err
		}

		return problem.Delete(ctx, tx)
	}); err != nil {
		return nil, err
	}

	return connect.NewResponse(&adminv1.DeleteProblemResponse{}), nil
}

func convertProblem(problem *domain.Problem) *adminv1.Problem {
	var typ adminv1.ProblemType
	switch problem.Type() {
	case domain.ProblemTypeDescriptive:
		typ = adminv1.ProblemType_PROBLEM_TYPE_DESCRIPTIVE
	case domain.ProblemTypeUnknown:
		fallthrough
	default:
		typ = adminv1.ProblemType_PROBLEM_TYPE_UNSPECIFIED
	}

	return &adminv1.Problem{
		Code:         string(problem.Code()),
		Title:        problem.Title(),
		MaxScore:     problem.MaxScore(),
		Category:     problem.Category(),
		RedeployRule: convertRedeployRule(problem),
		Body: &adminv1.ProblemBody{
			Type: typ,
		},
	}
}

func convertDescriptiveProblem(problem *domain.DescriptiveProblem) *adminv1.Problem {
	proto := convertProblem(problem.Problem())
	proto.Body.Body = &adminv1.ProblemBody_Descriptive{
		Descriptive: &adminv1.DescriptiveProblem{
			ProblemMarkdown:     problem.Body(),
			ExplanationMarkdown: problem.Explanation(),
		},
	}
	return proto
}

func convertRedeployRule(problem *domain.Problem) *adminv1.RedeployRule {
	switch problem.RedeployRule() {
	case domain.RedeployRuleUnredeployable:
		return &adminv1.RedeployRule{
			Type: adminv1.RedeployRuleType_REDEPLOY_RULE_TYPE_UNREDEPLOYABLE,
		}
	case domain.RedeployRulePercentagePenalty:
		return &adminv1.RedeployRule{
			Type:              adminv1.RedeployRuleType_REDEPLOY_RULE_TYPE_PERCENTAGE_PENALTY,
			PenaltyThreshold:  problem.PercentagePenalty().Threshold,
			PenaltyPercentage: problem.PercentagePenalty().Percentage,
		}
	case domain.RedeployRuleManual:
		return &adminv1.RedeployRule{
			Type: adminv1.RedeployRuleType_REDEPLOY_RULE_TYPE_MANUAL,
		}
	case domain.RedeployRuleUnknown:
		fallthrough
	default:
		return &adminv1.RedeployRule{
			Type: adminv1.RedeployRuleType_REDEPLOY_RULE_TYPE_UNSPECIFIED,
		}
	}
}

func parseRedeployRule(proto *adminv1.RedeployRule) (domain.RedeployRule, *domain.RedeployPenaltyPercentage) {
	switch proto.GetType() {
	case adminv1.RedeployRuleType_REDEPLOY_RULE_TYPE_UNREDEPLOYABLE:
		return domain.RedeployRuleUnredeployable, nil
	case adminv1.RedeployRuleType_REDEPLOY_RULE_TYPE_PERCENTAGE_PENALTY:
		return domain.RedeployRulePercentagePenalty, &domain.RedeployPenaltyPercentage{
			Threshold:  proto.GetPenaltyThreshold(),
			Percentage: proto.GetPenaltyPercentage(),
		}
	case adminv1.RedeployRuleType_REDEPLOY_RULE_TYPE_MANUAL:
		return domain.RedeployRuleManual, nil
	case adminv1.RedeployRuleType_REDEPLOY_RULE_TYPE_UNSPECIFIED:
		fallthrough
	default:
		return domain.RedeployRuleUnknown, nil
	}
}
