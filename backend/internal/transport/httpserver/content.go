package httpserver

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

func (h *Handler) GetAdminContentStatus(ctx context.Context, _ api.GetAdminContentStatusRequestObject) (api.GetAdminContentStatusResponseObject, error) {
	return api.GetAdminContentStatus200JSONResponse{Content: toAPIContentStatus(h.service.ContentStatus(ctx))}, nil
}

func (h *Handler) RefreshAdminContent(ctx context.Context, request api.RefreshAdminContentRequestObject) (api.RefreshAdminContentResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	claims := principal(ctx).ActionsClaims
	if claims == nil {
		if principal(ctx).Session.Kind != "admin" {
			return nil, core.NewError(http.StatusUnauthorized, "authentication_required", "Admin or GitHub Actions authentication is required")
		}
		claims = &service.ActionsClaims{Repository: h.service.Config.ContentRepository, Ref: h.service.Config.ContentRef}
	}
	if _, err := h.service.RefreshContent(ctx, *claims, request.Body.Commit); err != nil {
		return nil, err
	}
	return api.RefreshAdminContent200JSONResponse{Content: toAPIContentStatus(h.service.ContentStatus(ctx))}, nil
}

func toAPIContentStatus(status service.ContentStatus) api.ContentStatus {
	return api.ContentStatus{
		State:                api.ContentRefreshState(status.State),
		ActiveCommit:         status.ActiveCommit,
		ActivatedAt:          status.ActivatedAt,
		SourceRef:            status.SourceRef,
		LastAttemptAt:        status.LastAttemptAt,
		LastAttemptCommit:    status.LastAttemptCommit,
		LastError:            status.LastError,
		ServingLastKnownGood: status.ServingLastGood,
	}
}

func (h *Handler) ListAdminSections(ctx context.Context, _ api.ListAdminSectionsRequestObject) (api.ListAdminSectionsResponseObject, error) {
	snapshot, err := h.service.Store.ActiveContent(ctx)
	if err != nil {
		return nil, err
	}
	return api.ListAdminSections200JSONResponse{Sections: toAPISections(snapshot, nil)}, nil
}

func (h *Handler) ListContestantSections(ctx context.Context, _ api.ListContestantSectionsRequestObject) (api.ListContestantSectionsResponseObject, error) {
	snapshot, err := h.service.Store.ActiveContent(ctx)
	if err != nil {
		return nil, err
	}
	now := h.service.Now()
	return api.ListContestantSections200JSONResponse{Sections: toAPISections(snapshot, &now)}, nil
}

func (h *Handler) GetContestantSchedule(ctx context.Context, _ api.GetContestantScheduleRequestObject) (api.GetContestantScheduleResponseObject, error) {
	snapshot, err := h.service.Store.ActiveContent(ctx)
	if err != nil {
		return nil, err
	}
	sections := make([]api.ContestantScheduleSection, 0, len(snapshot.Manifest.Sections))
	for _, section := range snapshot.Manifest.Sections {
		sections = append(sections, api.ContestantScheduleSection{Slug: section.Slug, Beginning: section.Beginning, Ending: section.Ending})
	}
	sort.SliceStable(sections, func(i, j int) bool {
		if sections[i].Beginning.Equal(sections[j].Beginning) {
			return sections[i].Slug < sections[j].Slug
		}
		return sections[i].Beginning.Before(sections[j].Beginning)
	})
	return api.GetContestantSchedule200JSONResponse{Sections: sections}, nil
}

func (h *Handler) ListAdminProblems(ctx context.Context, _ api.ListAdminProblemsRequestObject) (api.ListAdminProblemsResponseObject, error) {
	snapshot, err := h.service.Store.ActiveContent(ctx)
	if err != nil {
		return nil, err
	}
	problems := make([]api.AdminProblem, 0, len(snapshot.Manifest.Problems))
	for _, problem := range orderedContentProblems(snapshot) {
		problems = append(problems, toAPIAdminProblem(problem, snapshot.CommitSHA))
	}
	return api.ListAdminProblems200JSONResponse{Problems: problems}, nil
}

func (h *Handler) GetAdminProblem(ctx context.Context, request api.GetAdminProblemRequestObject) (api.GetAdminProblemResponseObject, error) {
	var snapshot core.ContentSnapshot
	var err error
	if request.Params.Commit != nil {
		snapshot, err = h.service.Store.GetContent(ctx, string(*request.Params.Commit))
	} else {
		snapshot, err = h.service.Store.ActiveContent(ctx)
	}
	if err != nil {
		return nil, err
	}
	problem, ok := snapshot.Problem(request.ProblemCode)
	if !ok {
		return nil, core.NewError(http.StatusNotFound, "content_not_found", "Problem does not exist in the selected content commit")
	}
	return api.GetAdminProblem200JSONResponse{Problem: toAPIAdminProblem(problem, snapshot.CommitSHA)}, nil
}

func (h *Handler) ListContestantProblems(ctx context.Context, _ api.ListContestantProblemsRequestObject) (api.ListContestantProblemsResponseObject, error) {
	contestant, err := h.currentContestant(ctx)
	if err != nil {
		return nil, err
	}
	snapshot, err := h.service.Store.ActiveContent(ctx)
	if err != nil {
		return nil, err
	}
	scores, err := h.service.ContestantScores(ctx, contestant.TeamCode)
	if err != nil {
		return nil, err
	}
	now := h.service.Now()
	started := make(map[string]bool)
	for _, section := range snapshot.Manifest.Sections {
		if !now.Before(section.Beginning) {
			started[section.Slug] = true
		}
	}
	problems := make([]api.ProblemSummary, 0)
	for _, problem := range orderedContentProblems(snapshot) {
		if !started[problem.SectionSlug] {
			continue
		}
		var score *api.Score
		if selected, ok := scores[contestant.TeamCode][problem.Code]; ok {
			score = toAPIScore(selected)
		}
		problems = append(problems, api.ProblemSummary{
			Code: problem.Code, Title: problem.Title, MaxScore: problem.MaxScore, Category: problem.Category,
			SectionSlug: problem.SectionSlug, Score: score, SubmissionStatus: submissionStatus(snapshot, problem, now),
			Deployment: contestantDeploymentState(problem),
		})
	}
	return api.ListContestantProblems200JSONResponse{Problems: problems}, nil
}

func (h *Handler) GetContestantProblem(ctx context.Context, request api.GetContestantProblemRequestObject) (api.GetContestantProblemResponseObject, error) {
	contestant, err := h.currentContestant(ctx)
	if err != nil {
		return nil, err
	}
	snapshot, err := h.service.Store.ActiveContent(ctx)
	if err != nil {
		return nil, err
	}
	problem, ok := snapshot.Problem(request.ProblemCode)
	if !ok {
		return nil, core.NewError(http.StatusNotFound, "resource_not_found", "Problem was not found")
	}
	now := h.service.Now()
	sectionStarted := false
	for _, section := range snapshot.Manifest.Sections {
		if section.Slug == problem.SectionSlug && !now.Before(section.Beginning) {
			sectionStarted = true
			break
		}
	}
	if !sectionStarted {
		return nil, core.NewError(http.StatusNotFound, "resource_not_found", "Problem was not found")
	}
	scores, err := h.service.ContestantScores(ctx, contestant.TeamCode)
	if err != nil {
		return nil, err
	}
	var score *api.Score
	if selected, ok := scores[contestant.TeamCode][problem.Code]; ok {
		score = toAPIScore(selected)
	}
	return api.GetContestantProblem200JSONResponse{Problem: api.ProblemDetail{
		Code: problem.Code, Title: problem.Title, MaxScore: problem.MaxScore, Category: problem.Category,
		SectionSlug: problem.SectionSlug, Type: api.ProblemType(problem.Type), Body: problem.Body,
		Score: score, SubmissionStatus: submissionStatus(snapshot, problem, now), Deployment: contestantDeploymentState(problem),
	}}, nil
}

func contestantDeploymentState(problem core.Problem) api.ContestantProblemDeploymentState {
	return api.ContestantProblemDeploymentState{
		Redeployable:     problem.Redeploy.Type != core.RedeployUnredeployable,
		PenaltyThreshold: problem.Redeploy.Threshold,
	}
}

func (h *Handler) ListAdminAnnouncements(ctx context.Context, _ api.ListAdminAnnouncementsRequestObject) (api.ListAdminAnnouncementsResponseObject, error) {
	snapshot, err := h.service.Store.ActiveContent(ctx)
	if err != nil {
		return nil, err
	}
	return api.ListAdminAnnouncements200JSONResponse{Announcements: mapAnnouncements(snapshot.Manifest.Announcements, nil)}, nil
}

func (h *Handler) GetAdminAnnouncement(ctx context.Context, request api.GetAdminAnnouncementRequestObject) (api.GetAdminAnnouncementResponseObject, error) {
	snapshot, err := h.service.Store.ActiveContent(ctx)
	if err != nil {
		return nil, err
	}
	announcement, ok := snapshot.Announcement(request.AnnouncementSlug)
	if !ok {
		return nil, core.NewError(http.StatusNotFound, "content_not_found", "Announcement was not found")
	}
	return api.GetAdminAnnouncement200JSONResponse{Announcement: toAPIAnnouncement(announcement)}, nil
}

func (h *Handler) ListContestantAnnouncements(ctx context.Context, _ api.ListContestantAnnouncementsRequestObject) (api.ListContestantAnnouncementsResponseObject, error) {
	snapshot, err := h.service.Store.ActiveContent(ctx)
	if err != nil {
		return nil, err
	}
	now := h.service.Now()
	return api.ListContestantAnnouncements200JSONResponse{Announcements: mapAnnouncements(snapshot.Manifest.Announcements, &now)}, nil
}

func (h *Handler) GetContestantAnnouncement(ctx context.Context, request api.GetContestantAnnouncementRequestObject) (api.GetContestantAnnouncementResponseObject, error) {
	snapshot, err := h.service.Store.ActiveContent(ctx)
	if err != nil {
		return nil, err
	}
	announcement, ok := snapshot.Announcement(request.AnnouncementSlug)
	if !ok || announcement.EffectiveFrom.After(h.service.Now()) {
		return nil, core.NewError(http.StatusNotFound, "resource_not_found", "Announcement was not found")
	}
	return api.GetContestantAnnouncement200JSONResponse{Announcement: toAPIAnnouncement(announcement)}, nil
}

func mapAnnouncements(source []core.Announcement, effectiveAt interface{ Before(time.Time) bool }) []api.Announcement {
	mapped := make([]api.Announcement, 0, len(source))
	for _, announcement := range source {
		if effectiveAt != nil && effectiveAt.Before(announcement.EffectiveFrom) {
			continue
		}
		mapped = append(mapped, toAPIAnnouncement(announcement))
	}
	sort.SliceStable(mapped, func(i, j int) bool {
		if mapped[i].EffectiveFrom.Equal(mapped[j].EffectiveFrom) {
			return mapped[i].Slug < mapped[j].Slug
		}
		return mapped[i].EffectiveFrom.After(mapped[j].EffectiveFrom)
	})
	return mapped
}
