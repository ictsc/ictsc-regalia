package httpserver

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

func (h *Handler) ListContestantAnswers(ctx context.Context, request api.ListContestantAnswersRequestObject) (api.ListContestantAnswersResponseObject, error) {
	contestant, err := h.currentContestant(ctx)
	if err != nil {
		return nil, err
	}
	team, problem := contestant.TeamCode, string(request.ProblemCode)
	answers, err := h.service.Store.ListAnswers(ctx, core.AnswerFilter{TeamCode: &team, ProblemCode: &problem})
	if err != nil {
		return nil, err
	}
	sort.SliceStable(answers, func(i, j int) bool { return answers[i].Number > answers[j].Number })
	mapped := make([]api.AnswerSummary, 0, len(answers))
	var last *time.Time
	for _, answer := range answers {
		score, err := h.scoreForAnswer(ctx, answer, true)
		if err != nil {
			return nil, err
		}
		mapped = append(mapped, toAPIAnswerSummary(answer, score))
		if last == nil || answer.SubmittedAt.After(*last) {
			value := answer.SubmittedAt
			last = &value
		}
	}
	return api.ListContestantAnswers200JSONResponse{
		Answers: mapped, SubmitIntervalSeconds: int(core.AnswerInterval / time.Second), LastSubmittedAt: last,
	}, nil
}

func (h *Handler) SubmitContestantAnswer(ctx context.Context, request api.SubmitContestantAnswerRequestObject) (api.SubmitContestantAnswerResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	contestant, err := h.currentContestant(ctx)
	if err != nil {
		return nil, err
	}
	answer, _, err := h.service.SubmitAnswer(ctx, contestant, request.ProblemCode, request.Body.Body)
	if err != nil {
		return nil, err
	}
	return api.SubmitContestantAnswer201JSONResponse{Answer: toAPIAnswer(answer, nil)}, nil
}

func (h *Handler) GetContestantAnswer(ctx context.Context, request api.GetContestantAnswerRequestObject) (api.GetContestantAnswerResponseObject, error) {
	contestant, err := h.currentContestant(ctx)
	if err != nil {
		return nil, err
	}
	answer, err := h.service.Store.GetAnswer(ctx, contestant.TeamCode, request.ProblemCode, request.AnswerNumber)
	if err != nil {
		return nil, err
	}
	score, err := h.scoreForAnswer(ctx, answer, true)
	if err != nil {
		return nil, err
	}
	return api.GetContestantAnswer200JSONResponse{Answer: toAPIAnswer(answer, score)}, nil
}

func (h *Handler) ListAdminAnswers(ctx context.Context, request api.ListAdminAnswersRequestObject) (api.ListAdminAnswersResponseObject, error) {
	answers, err := h.service.Store.ListAnswers(ctx, core.AnswerFilter{TeamCode: request.Params.TeamCode, ProblemCode: request.Params.ProblemCode})
	if err != nil {
		return nil, err
	}
	includeMarked := request.Params.IncludeMarked != nil && *request.Params.IncludeMarked
	mapped := make([]api.AdminAnswer, 0, len(answers))
	for _, answer := range answers {
		score, err := h.scoreForAnswer(ctx, answer, false)
		if err != nil {
			return nil, err
		}
		if score != nil && !includeMarked {
			continue
		}
		mappedAnswer, err := h.toAdminAnswer(ctx, answer, score)
		if err != nil {
			return nil, err
		}
		mapped = append(mapped, mappedAnswer)
	}
	sort.SliceStable(mapped, func(i, j int) bool {
		if !mapped[i].SubmittedAt.Equal(mapped[j].SubmittedAt) {
			return mapped[i].SubmittedAt.After(mapped[j].SubmittedAt)
		}
		if mapped[i].Reference.TeamCode != mapped[j].Reference.TeamCode {
			return mapped[i].Reference.TeamCode < mapped[j].Reference.TeamCode
		}
		if mapped[i].Reference.ProblemCode != mapped[j].Reference.ProblemCode {
			return mapped[i].Reference.ProblemCode < mapped[j].Reference.ProblemCode
		}
		return mapped[i].Reference.AnswerNumber > mapped[j].Reference.AnswerNumber
	})
	return api.ListAdminAnswers200JSONResponse{Answers: mapped}, nil
}

func (h *Handler) GetAdminAnswer(ctx context.Context, request api.GetAdminAnswerRequestObject) (api.GetAdminAnswerResponseObject, error) {
	answer, err := h.service.Store.GetAnswer(ctx, request.TeamCode, request.ProblemCode, request.AnswerNumber)
	if err != nil {
		return nil, err
	}
	score, err := h.scoreForAnswer(ctx, answer, false)
	if err != nil {
		return nil, err
	}
	mapped, err := h.toAdminAnswer(ctx, answer, score)
	if err != nil {
		return nil, err
	}
	return api.GetAdminAnswer200JSONResponse{Answer: mapped}, nil
}

func (h *Handler) toAdminAnswer(ctx context.Context, answer core.Answer, score *core.Score) (api.AdminAnswer, error) {
	team, err := h.service.Store.GetTeam(ctx, answer.TeamCode)
	if err != nil {
		return api.AdminAnswer{}, err
	}
	author, err := h.service.Store.GetContestant(ctx, answer.AuthorName)
	if err != nil {
		return api.AdminAnswer{}, err
	}
	snapshot, err := h.service.Store.GetContent(ctx, answer.ContentCommit)
	if err != nil {
		return api.AdminAnswer{}, err
	}
	problem, ok := snapshot.Problem(answer.ProblemCode)
	if !ok {
		return api.AdminAnswer{}, core.NewError(http.StatusInternalServerError, "internal_error", "Answer references missing problem snapshot")
	}
	var markingScore *api.MarkingScore
	if score != nil {
		markingScore = &api.MarkingScore{Marked: score.MarkedScore, Penalty: score.Penalty, Total: score.EffectiveScore, Max: score.MaxScore}
	}
	workflow, err := h.service.Store.GetAnswerWorkflow(ctx, answer.TeamCode, answer.ProblemCode, answer.Number)
	if err != nil {
		return api.AdminAnswer{}, err
	}
	var assignee *string
	source := api.AdminAnswerWorkflowAssignmentSource("DEFAULT")
	if workflow.ClaimedBy != "" {
		assignee = &workflow.ClaimedBy
		source = api.AdminAnswerWorkflowAssignmentSource("CLAIMED")
	} else if problem.DefaultGraderDiscordID != "" {
		assignee = &problem.DefaultGraderDiscordID
	}
	return api.AdminAnswer{
		Reference: api.AdminAnswerReference{TeamCode: answer.TeamCode, ProblemCode: answer.ProblemCode, AnswerNumber: answer.Number},
		Team:      toAPITeam(team), Author: toAPIProfile(author), Problem: toAPICatalog(problem),
		Body: api.AnswerBody{Type: api.AnswerBodyTypeDESCRIPTIVE, Body: answer.Body}, SubmittedAt: answer.SubmittedAt,
		ContentCommit: answer.ContentCommit, Score: markingScore,
		Workflow: api.AdminAnswerWorkflow{Status: api.AdminAnswerWorkflowStatus(workflow.Status), AssigneeDiscordId: assignee, AssignmentSource: source, Revision: workflow.Revision},
	}, nil
}

func (h *Handler) UpdateAdminAnswerWorkflow(ctx context.Context, request api.UpdateAdminAnswerWorkflowRequestObject) (api.UpdateAdminAnswerWorkflowResponseObject, error) {
	if request.Body == nil || (request.Body.Status == nil && request.Body.Assignment == nil) {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Status or assignment is required")
	}
	if request.Body.Status != nil && !request.Body.Status.Valid() {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Invalid workflow status")
	}
	if request.Body.Assignment != nil && !request.Body.Assignment.Valid() {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Invalid assignment action")
	}
	update := core.WorkflowUpdate{ExpectedRevision: request.Body.ExpectedRevision, ActorDiscordID: principal(ctx).Session.Discord.ID}
	if request.Body.Status != nil {
		status := core.AnswerWorkflowStatus(*request.Body.Status)
		update.Status = &status
	}
	if request.Body.Assignment != nil {
		update.Assignment = string(*request.Body.Assignment)
	}
	if update.Assignment == "CLAIM_SELF" && update.ActorDiscordID == "" {
		return nil, core.NewError(http.StatusForbidden, "admin_identity_required", "Discord identity is required to claim an answer")
	}
	if _, err := h.service.Store.UpdateAnswerWorkflow(ctx, request.TeamCode, request.ProblemCode, request.AnswerNumber, update); err != nil {
		return nil, err
	}
	answer, err := h.service.Store.GetAnswer(ctx, request.TeamCode, request.ProblemCode, request.AnswerNumber)
	if err != nil {
		return nil, err
	}
	score, err := h.scoreForAnswer(ctx, answer, false)
	if err != nil {
		return nil, err
	}
	mapped, err := h.toAdminAnswer(ctx, answer, score)
	if err != nil {
		return nil, err
	}
	return api.UpdateAdminAnswerWorkflow200JSONResponse{Answer: mapped}, nil
}

func (h *Handler) ListAdminMarkingResults(ctx context.Context, request api.ListAdminMarkingResultsRequestObject) (api.ListAdminMarkingResultsResponseObject, error) {
	results, err := h.service.Store.ListMarkingResults(ctx)
	if err != nil {
		return nil, err
	}
	state, err := h.service.Store.GetCompetitionState(ctx)
	if err != nil {
		return nil, err
	}
	mapped := make([]api.MarkingResult, 0, len(results))
	for _, result := range results {
		answer, err := h.service.Store.GetAnswer(ctx, result.TeamCode, result.ProblemCode, result.AnswerNumber)
		if err != nil {
			return nil, err
		}
		result.Visibility = core.MarkingVisibilityAt(answer.SubmittedAt, h.service.Now(), state.RankingFreezeAt, state.FinalRevealedAt)
		if request.Params.Visibility != nil && core.Visibility(*request.Params.Visibility) != result.Visibility {
			continue
		}
		mapped = append(mapped, toAPIMarking(result))
	}
	sort.SliceStable(mapped, func(i, j int) bool {
		if !mapped[i].CreatedAt.Equal(mapped[j].CreatedAt) {
			return mapped[i].CreatedAt.After(mapped[j].CreatedAt)
		}
		return mapped[i].Id.String() < mapped[j].Id.String()
	})
	return api.ListAdminMarkingResults200JSONResponse{MarkingResults: mapped}, nil
}

func (h *Handler) CreateAdminMarkingResult(ctx context.Context, request api.CreateAdminMarkingResultRequestObject) (api.CreateAdminMarkingResultResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	answer, err := h.service.Store.GetAnswer(ctx, request.Body.Answer.TeamCode, request.Body.Answer.ProblemCode, request.Body.Answer.AnswerNumber)
	if err != nil {
		return nil, err
	}
	if request.Body.Score > answer.MaxScore {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Score cannot exceed the answer snapshot maximum")
	}
	lineComments := make([]core.AnswerLineComment, 0)
	if request.Body.LineComments != nil {
		lineCount := strings.Count(strings.ReplaceAll(answer.Body, "\r\n", "\n"), "\n") + 1
		for _, comment := range *request.Body.LineComments {
			endLine := comment.LineNumber
			if comment.EndLineNumber != nil {
				endLine = *comment.EndLineNumber
			}
			if comment.LineNumber < 1 || endLine < comment.LineNumber || int(endLine) > lineCount || strings.TrimSpace(comment.Body) == "" || utf8.RuneCountInString(comment.Body) > 2000 {
				return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Line comment has an invalid line or body")
			}
			lineComments = append(lineComments, core.AnswerLineComment{LineNumber: comment.LineNumber, EndLineNumber: comment.EndLineNumber, Body: comment.Body})
		}
	}
	state, err := h.service.Store.GetCompetitionState(ctx)
	if err != nil {
		return nil, err
	}
	visibility := core.MarkingVisibilityAt(answer.SubmittedAt, h.service.Now(), state.RankingFreezeAt, state.FinalRevealedAt)
	result, err := h.service.Store.CreateMarkingResult(ctx, core.MarkingResult{
		ID: uuid.NewString(), TeamCode: answer.TeamCode, ProblemCode: answer.ProblemCode, AnswerNumber: answer.Number,
		Judge: principal(ctx).Session.AdminName, MarkedScore: request.Body.Score, Rationale: request.Body.Rationale, LineComments: lineComments,
		CreatedAt: h.service.Now(), Visibility: visibility,
	})
	if err != nil {
		return nil, err
	}
	return api.CreateAdminMarkingResult201JSONResponse{MarkingResult: toAPIMarking(result)}, nil
}

func (h *Handler) DeleteAdminMarkingLineComment(ctx context.Context, request api.DeleteAdminMarkingLineCommentRequestObject) (api.DeleteAdminMarkingLineCommentResponseObject, error) {
	result, err := h.service.Store.DeleteMarkingLineComment(ctx, request.MarkingResultId.String(), request.CommentIndex, principal(ctx).Session.AdminName, h.service.Now())
	if err != nil {
		return nil, err
	}
	return api.DeleteAdminMarkingLineComment200JSONResponse{MarkingResult: toAPIMarking(result)}, nil
}

func (h *Handler) ListAdminScores(ctx context.Context, _ api.ListAdminScoresRequestObject) (api.ListAdminScoresResponseObject, error) {
	scores, err := h.service.Scores(ctx, nil)
	if err != nil {
		return nil, err
	}
	snapshot, err := h.service.Store.ActiveContent(ctx)
	if err != nil {
		return nil, err
	}
	mapped := make([]api.AdminScore, 0)
	for teamCode, byProblem := range scores {
		team, err := h.service.Store.GetTeam(ctx, teamCode)
		if err != nil {
			return nil, err
		}
		for problemCode, score := range byProblem {
			problem, ok := snapshot.Problem(problemCode)
			if !ok {
				continue
			}
			mapped = append(mapped, api.AdminScore{
				Team: toAPITeam(team), Problem: toAPICatalog(problem), MarkedScore: score.MarkedScore,
				Penalty: score.Penalty, Score: score.EffectiveScore,
			})
		}
	}
	sort.SliceStable(mapped, func(i, j int) bool {
		if mapped[i].Team.Code != mapped[j].Team.Code {
			return mapped[i].Team.Code < mapped[j].Team.Code
		}
		return mapped[i].Problem.Code < mapped[j].Problem.Code
	})
	return api.ListAdminScores200JSONResponse{Scores: mapped}, nil
}

func (h *Handler) RecalculateAdminScores(ctx context.Context, _ api.RecalculateAdminScoresRequestObject) (api.RecalculateAdminScoresResponseObject, error) {
	scores, err := h.service.Scores(ctx, nil)
	if err != nil {
		return nil, err
	}
	if err := h.service.Store.RecalculateScores(ctx, scores, principal(ctx).Session.AdminName); err != nil {
		return nil, err
	}
	return api.RecalculateAdminScores204Response{}, nil
}

func (h *Handler) RevealAdminFinalScores(ctx context.Context, _ api.RevealAdminFinalScoresRequestObject) (api.RevealAdminFinalScoresResponseObject, error) {
	if _, err := h.service.Store.RevealFinal(ctx, h.service.Now(), principal(ctx).Session.AdminName); err != nil {
		return nil, err
	}
	return api.RevealAdminFinalScores204Response{}, nil
}
