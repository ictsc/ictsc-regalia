package httpserver

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

func toAPITeam(team core.Team) api.Team {
	color := api.TeamColor(core.TeamColor(team.Color))
	return api.Team{Color: &color, Code: team.Code, Name: team.Name, Organization: team.Organization, MemberLimit: team.MemberLimit}
}

func toAPIProfile(contestant core.Contestant) api.ContestantProfile {
	return api.ContestantProfile{Name: contestant.Name, DisplayName: contestant.DisplayName, SelfIntroduction: contestant.SelfIntroduction}
}

func toAPIInvitation(invitation core.Invitation) api.Invitation {
	return api.Invitation{Code: invitation.Code, TeamCode: invitation.TeamCode, CreatedAt: invitation.CreatedAt, ExpiresAt: invitation.ExpiresAt}
}

func toAPICatalog(problem core.Problem) api.ProblemCatalogEntry {
	return api.ProblemCatalogEntry{Code: problem.Code, Title: problem.Title, Category: problem.Category, MaxScore: problem.MaxScore}
}

func toAPIRedeploy(rule core.RedeployRule) api.RedeployRule {
	return api.RedeployRule{
		Type: api.RedeployRuleType(rule.Type), PenaltyThreshold: rule.Threshold, PenaltyPercentage: rule.Percentage,
	}
}

func toAPISections(snapshot core.ContentSnapshot, now *time.Time) []api.Section {
	problemByCode := make(map[string]core.Problem, len(snapshot.Manifest.Problems))
	for _, problem := range snapshot.Manifest.Problems {
		problemByCode[problem.Code] = problem
	}
	sections := make([]api.Section, 0, len(snapshot.Manifest.Sections))
	for _, section := range snapshot.Manifest.Sections {
		if now != nil && section.Beginning.After(*now) {
			continue
		}
		problems := make([]api.ProblemCatalogEntry, 0, len(section.ProblemIDs))
		for _, code := range section.ProblemIDs {
			if problem, ok := problemByCode[code]; ok {
				problems = append(problems, toAPICatalog(problem))
			}
		}
		sections = append(sections, api.Section{Slug: section.Slug, Beginning: section.Beginning, Ending: section.Ending, Problems: problems})
	}
	sort.SliceStable(sections, func(i, j int) bool {
		if sections[i].Beginning.Equal(sections[j].Beginning) {
			return sections[i].Slug < sections[j].Slug
		}
		return sections[i].Beginning.Before(sections[j].Beginning)
	})
	return sections
}

func toAPIAdminProblem(problem core.Problem, commit string) api.AdminProblem {
	return api.AdminProblem{
		Code: problem.Code, Title: problem.Title, MaxScore: problem.MaxScore, Category: problem.Category,
		SectionSlug: problem.SectionSlug, Type: api.ProblemType(problem.Type), Body: problem.Body,
		Explanation: problem.Explanation, RedeployRule: toAPIRedeploy(problem.Redeploy), ContentCommit: commit,
	}
}

func toAPIAnnouncement(announcement core.Announcement) api.Announcement {
	return api.Announcement{Slug: announcement.Slug, Title: announcement.Title, Markdown: announcement.Markdown, EffectiveFrom: announcement.EffectiveFrom}
}

func toAPIScore(score core.Score) *api.Score {
	return &api.Score{MarkedScore: score.MarkedScore, Penalty: score.Penalty, Score: score.EffectiveScore, MaxScore: score.MaxScore}
}

func toAPIRanking(entries []core.RankingEntry, frozen bool, frozenAt *time.Time) api.RankingResponse {
	ranking := make([]api.Rank, 0, len(entries))
	for _, entry := range entries {
		ranking = append(ranking, api.Rank{
			Rank: entry.Rank, TeamCode: entry.Team.Code, TeamName: entry.Team.Name,
			Organization: entry.Team.Organization, Score: entry.Score,
			LastEffectiveSubmissionAt: entry.LastEffectiveSubmissionAt,
		})
	}
	return api.RankingResponse{Ranking: ranking, Frozen: frozen, FrozenAt: frozenAt}
}

func toAPIAnswer(answer core.Answer, score *core.Score) api.Answer {
	var mapped *api.Score
	if score != nil {
		mapped = toAPIScore(*score)
	}
	return api.Answer{
		Number: answer.Number, SubmittedAt: answer.SubmittedAt, ContentCommit: answer.ContentCommit,
		Body: api.AnswerBody{Type: api.AnswerBodyTypeDESCRIPTIVE, Body: answer.Body}, Score: mapped,
	}
}

func toAPIAnswerSummary(answer core.Answer, score *core.Score) api.AnswerSummary {
	var mapped *api.Score
	if score != nil {
		mapped = toAPIScore(*score)
	}
	return api.AnswerSummary{
		Number: answer.Number, Type: api.ProblemTypeDESCRIPTIVE, SubmittedAt: answer.SubmittedAt,
		ContentCommit: answer.ContentCommit, Score: mapped,
	}
}

func toAPIContestantDeployment(deployment core.Deployment, penalty, allowed int32) api.ContestantDeployment {
	return api.ContestantDeployment{
		Revision: deployment.Revision, Status: api.DeploymentStatus(deployment.LatestStatus),
		RequestedAt: deployment.RequestedAt, ContentCommit: deployment.ContentCommit,
		Penalty: penalty, AllowedRequestCount: allowed,
	}
}

func toAPIAdminDeployment(deployment core.Deployment) api.AdminDeployment {
	orderedEvents := append([]core.DeploymentEvent(nil), deployment.Events...)
	sort.SliceStable(orderedEvents, func(i, j int) bool {
		if orderedEvents[i].OccurredAt.Equal(orderedEvents[j].OccurredAt) {
			return orderedEvents[i].EventID < orderedEvents[j].EventID
		}
		return orderedEvents[i].OccurredAt.Before(orderedEvents[j].OccurredAt)
	})
	events := make([]api.DeploymentEvent, 0, len(orderedEvents))
	for _, event := range orderedEvents {
		events = append(events, toAPIEvent(event))
	}
	return api.AdminDeployment{
		TeamCode: deployment.TeamCode, ProblemCode: deployment.ProblemCode, Revision: deployment.Revision,
		ContentCommit: deployment.ContentCommit, LatestStatus: api.DeploymentStatus(deployment.LatestStatus), Events: events,
	}
}

func toAPIEvent(event core.DeploymentEvent) api.DeploymentEvent {
	id, _ := uuid.Parse(event.EventID)
	return api.DeploymentEvent{EventId: id, OccurredAt: event.OccurredAt, Status: api.DeploymentStatus(event.Status), Message: event.Message}
}

func toAPIMarking(result core.MarkingResult) api.MarkingResult {
	id, _ := uuid.Parse(result.ID)
	return api.MarkingResult{
		Id: id, Answer: api.AdminAnswerReference{TeamCode: result.TeamCode, ProblemCode: result.ProblemCode, AnswerNumber: result.AnswerNumber},
		Judge: api.AdminActor{Name: result.Judge}, Score: result.MarkedScore, Rationale: result.Rationale,
		CreatedAt: result.CreatedAt, Visibility: api.MarkingVisibility(result.Visibility),
	}
}

func (h *Handler) scoreForAnswer(ctx context.Context, answer core.Answer, contestant bool) (*core.Score, error) {
	markings, err := h.service.Store.ListMarkingResults(ctx)
	if err != nil {
		return nil, err
	}
	state, err := h.service.Store.GetCompetitionState(ctx)
	if err != nil {
		return nil, err
	}
	latest := core.LatestMarkings(markings)
	keyed := core.SelectBestScores([]core.Answer{answer}, markings)
	score, ok := keyed[answer.TeamCode][answer.ProblemCode]
	if !ok {
		return nil, nil
	}
	if contestant {
		marking, hasMarking := core.LatestMarkingForAnswer(latest, answer.TeamCode, answer.ProblemCode, answer.Number)
		visibility := core.MarkingVisibilityAt(answer.SubmittedAt, h.service.Now(), state.RankingFreezeAt, state.FinalRevealedAt)
		if !hasMarking || marking.ID == "" || visibility == core.VisibilityPrivate {
			return nil, nil
		}
	}
	return &score, nil
}

func submissionStatus(snapshot core.ContentSnapshot, problem core.Problem, now time.Time) api.SubmissionStatus {
	for _, section := range snapshot.Manifest.Sections {
		if section.Slug == problem.SectionSlug {
			return api.SubmissionStatus{
				IsSubmittable:   !now.Before(section.Beginning) && now.Before(section.Ending),
				SubmittableFrom: &section.Beginning, SubmittableUntil: &section.Ending,
			}
		}
	}
	return api.SubmissionStatus{}
}
