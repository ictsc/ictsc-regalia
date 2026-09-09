package memory

import (
	"context"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

type CompetitionStore struct {
	mu               sync.RWMutex
	teams            map[int64]core.Team
	invitations      map[string]core.Invitation
	usedInvites      map[string]bool
	contestants      map[string]core.Contestant
	contents         map[string]core.ContentSnapshot
	active           string
	answers          []core.Answer
	markings         []core.MarkingResult
	scoreSelections  map[int64]map[string]core.Score
	rankingSnapshots map[string]core.RankingSnapshot
	deployments      map[string]core.Deployment
	state            core.CompetitionState
}

func NewCompetitionStore() *CompetitionStore {
	return &CompetitionStore{
		teams:            make(map[int64]core.Team),
		invitations:      make(map[string]core.Invitation),
		usedInvites:      make(map[string]bool),
		contestants:      make(map[string]core.Contestant),
		contents:         make(map[string]core.ContentSnapshot),
		scoreSelections:  make(map[int64]map[string]core.Score),
		rankingSnapshots: make(map[string]core.RankingSnapshot),
		deployments:      make(map[string]core.Deployment),
		state:            core.CompetitionState{UpdatedAt: time.Now().UTC(), UpdatedBy: "system"},
	}
}

func (s *CompetitionStore) Ping(context.Context) error { return nil }

func (s *CompetitionStore) ListTeams(context.Context) ([]core.Team, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	teams := make([]core.Team, 0, len(s.teams))
	for _, team := range s.teams {
		teams = append(teams, team)
	}
	sort.Slice(teams, func(i, j int) bool { return teams[i].Code < teams[j].Code })
	return teams, nil
}

func (s *CompetitionStore) GetTeam(_ context.Context, code int64) (core.Team, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	team, ok := s.teams[code]
	if !ok {
		return core.Team{}, core.NewError(http.StatusNotFound, "team_not_found", "Team was not found")
	}
	return team, nil
}

func (s *CompetitionStore) CreateTeam(_ context.Context, team core.Team) (core.Team, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.teams[team.Code]; ok {
		return core.Team{}, core.NewError(http.StatusConflict, "team_code_conflict", "Team code already exists")
	}
	for _, existing := range s.teams {
		if existing.Name == team.Name {
			return core.Team{}, core.NewError(http.StatusConflict, "conflict", "Team name already exists")
		}
	}
	s.teams[team.Code] = team
	return team, nil
}

func (s *CompetitionStore) UpdateTeam(_ context.Context, code int64, patch core.TeamPatch) (core.Team, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	team, ok := s.teams[code]
	if !ok {
		return core.Team{}, core.NewError(http.StatusNotFound, "team_not_found", "Team was not found")
	}
	if patch.Name != nil {
		for existingCode, existing := range s.teams {
			if existingCode != code && existing.Name == *patch.Name {
				return core.Team{}, core.NewError(http.StatusConflict, "conflict", "Team name is already in use")
			}
		}
		team.Name = *patch.Name
	}
	if patch.Color != nil {
		team.Color = *patch.Color
	}
	if patch.Organization != nil {
		team.Organization = *patch.Organization
	}
	if patch.MemberLimit != nil {
		members := int32(0)
		for _, contestant := range s.contestants {
			if contestant.TeamCode == code {
				members++
			}
		}
		if *patch.MemberLimit < members {
			return core.Team{}, core.NewError(http.StatusConflict, "team_full", "Member limit cannot be lower than the current member count")
		}
		team.MemberLimit = *patch.MemberLimit
	}
	s.teams[code] = team
	return team, nil
}

func (s *CompetitionStore) DeleteTeam(_ context.Context, code int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, contestant := range s.contestants {
		if contestant.TeamCode == code {
			return core.NewError(http.StatusConflict, "team_not_empty", "Team still has contestants")
		}
	}
	if _, ok := s.teams[code]; !ok {
		return core.NewError(http.StatusNotFound, "team_not_found", "Team was not found")
	}
	delete(s.teams, code)
	return nil
}

func (s *CompetitionStore) ListInvitations(context.Context) ([]core.Invitation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]core.Invitation, 0)
	for code, invitation := range s.invitations {
		if !s.usedInvites[code] {
			result = append(result, invitation)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.After(result[j].CreatedAt) })
	return result, nil
}

func (s *CompetitionStore) CreateInvitation(_ context.Context, invitation core.Invitation) (core.Invitation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.teams[invitation.TeamCode]; !ok {
		return core.Invitation{}, core.NewError(http.StatusNotFound, "team_not_found", "Team was not found")
	}
	if _, ok := s.invitations[invitation.Code]; ok {
		return core.Invitation{}, core.NewError(http.StatusConflict, "invitation_conflict", "Invitation code already exists")
	}
	s.invitations[invitation.Code] = invitation
	return invitation, nil
}

func (s *CompetitionStore) ConsumeInvitation(_ context.Context, code string, now time.Time, contestant core.Contestant) (core.Contestant, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	invitation, ok := s.invitations[code]
	if !ok {
		return core.Contestant{}, core.NewError(http.StatusConflict, "conflict", "Invitation code does not exist")
	}
	if s.usedInvites[code] {
		return core.Contestant{}, core.NewError(http.StatusConflict, "invitation_already_used", "Invitation was already used")
	}
	if !now.Before(invitation.ExpiresAt) {
		return core.Contestant{}, core.NewError(http.StatusConflict, "invitation_expired", "Invitation has expired")
	}
	team := s.teams[invitation.TeamCode]
	count := int32(0)
	for _, existing := range s.contestants {
		if existing.Name == contestant.Name || existing.DiscordID == contestant.DiscordID {
			return core.Contestant{}, core.NewError(http.StatusConflict, "contestant_already_registered", "Contestant already exists")
		}
		if existing.TeamCode == invitation.TeamCode {
			count++
		}
	}
	if count >= team.MemberLimit {
		return core.Contestant{}, core.NewError(http.StatusConflict, "team_full", "Team member limit has been reached")
	}
	contestant.TeamCode = invitation.TeamCode
	s.contestants[contestant.Name] = contestant
	s.usedInvites[code] = true
	return contestant, nil
}

func (s *CompetitionStore) ListContestants(context.Context) ([]core.Contestant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]core.Contestant, 0, len(s.contestants))
	for _, contestant := range s.contestants {
		result = append(result, contestant)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].TeamCode != result[j].TeamCode {
			return result[i].TeamCode < result[j].TeamCode
		}
		return result[i].Name < result[j].Name
	})
	return result, nil
}

func (s *CompetitionStore) GetContestant(_ context.Context, name string) (core.Contestant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	contestant, ok := s.contestants[name]
	if !ok {
		return core.Contestant{}, core.NewError(http.StatusNotFound, "contestant_not_found", "Contestant was not found")
	}
	return contestant, nil
}

func (s *CompetitionStore) GetContestantByDiscord(_ context.Context, discordID string) (core.Contestant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, contestant := range s.contestants {
		if contestant.DiscordID == discordID {
			return contestant, nil
		}
	}
	return core.Contestant{}, core.NewError(http.StatusNotFound, "contestant_not_found", "Contestant was not found")
}

func (s *CompetitionStore) UpdateContestant(_ context.Context, name string, patch core.ContestantPatch) (core.Contestant, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	contestant, ok := s.contestants[name]
	if !ok {
		return core.Contestant{}, core.NewError(http.StatusNotFound, "contestant_not_found", "Contestant was not found")
	}
	if patch.DisplayName != nil {
		contestant.DisplayName = *patch.DisplayName
	}
	if patch.SelfIntroduction != nil {
		contestant.SelfIntroduction = *patch.SelfIntroduction
	}
	s.contestants[name] = contestant
	return contestant, nil
}

func (s *CompetitionStore) ActiveContent(context.Context) (core.ContentSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snapshot, ok := s.contents[s.active]
	if !ok {
		return core.ContentSnapshot{}, core.NewError(http.StatusBadGateway, "content_not_available", "No valid content snapshot is available")
	}
	return snapshot, nil
}

func (s *CompetitionStore) GetContent(_ context.Context, commit string) (core.ContentSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snapshot, ok := s.contents[commit]
	if !ok {
		return core.ContentSnapshot{}, core.NewError(http.StatusNotFound, "content_not_found", "Content snapshot was not found")
	}
	return snapshot, nil
}

func (s *CompetitionStore) ActivateContent(_ context.Context, snapshot core.ContentSnapshot, expected string) (core.ContentSnapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if snapshot.CommitSHA == s.active {
		return s.contents[s.active], nil
	}
	if expected == core.NoActiveContentCommit && s.active != "" {
		return core.ContentSnapshot{}, core.NewError(http.StatusConflict, "content_refresh_in_progress", "Active content appeared while refresh was in progress")
	}
	if expected != "" && expected != core.NoActiveContentCommit && expected != s.active {
		return core.ContentSnapshot{}, core.NewError(http.StatusConflict, "content_refresh_in_progress", "Active content changed while refresh was in progress")
	}
	snapshot.ActivatedAt = time.Now().UTC()
	s.contents[snapshot.CommitSHA] = snapshot
	s.active = snapshot.CommitSHA
	return snapshot, nil
}

func (s *CompetitionStore) ListAnswers(_ context.Context, filter core.AnswerFilter) ([]core.Answer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]core.Answer, 0)
	for _, answer := range s.answers {
		if filter.TeamCode != nil && answer.TeamCode != *filter.TeamCode ||
			filter.ProblemCode != nil && answer.ProblemCode != *filter.ProblemCode ||
			filter.Before != nil && answer.SubmittedAt.After(*filter.Before) {
			continue
		}
		result = append(result, answer)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].SubmittedAt.After(result[j].SubmittedAt) })
	return result, nil
}

func (s *CompetitionStore) GetAnswer(_ context.Context, team int64, problem string, number int32) (core.Answer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, answer := range s.answers {
		if answer.TeamCode == team && answer.ProblemCode == problem && answer.Number == number {
			return answer, nil
		}
	}
	return core.Answer{}, core.NewError(http.StatusNotFound, "answer_not_found", "Answer was not found")
}

func (s *CompetitionStore) SubmitAnswer(_ context.Context, answer core.Answer, interval time.Duration) (core.Answer, time.Duration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var last *core.Answer
	var number int32 = 1
	for index := range s.answers {
		existing := &s.answers[index]
		if existing.TeamCode == answer.TeamCode && existing.ProblemCode == answer.ProblemCode {
			if existing.Number >= number {
				number = existing.Number + 1
			}
			if last == nil || existing.SubmittedAt.After(last.SubmittedAt) {
				last = existing
			}
		}
	}
	if last != nil && answer.SubmittedAt.Before(last.SubmittedAt.Add(interval)) {
		retry := last.SubmittedAt.Add(interval).Sub(answer.SubmittedAt)
		return core.Answer{}, retry, &core.Error{Status: http.StatusTooManyRequests, Code: "answer_rate_limited", Message: "Answer submission interval has not elapsed", RetryAfter: retry}
	}
	answer.Number = number
	s.answers = append(s.answers, answer)
	return answer, 0, nil
}

func (s *CompetitionStore) ListMarkingResults(context.Context) ([]core.MarkingResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := append([]core.MarkingResult(nil), s.markings...)
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.After(result[j].CreatedAt) })
	return result, nil
}

func (s *CompetitionStore) CreateMarkingResult(_ context.Context, result core.MarkingResult) (core.MarkingResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.markings = append(s.markings, result)
	return result, nil
}

func (s *CompetitionStore) ListDeployments(_ context.Context, filter core.DeploymentFilter) ([]core.Deployment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]core.Deployment, 0)
	for _, deployment := range s.deployments {
		if filter.TeamCode != nil && deployment.TeamCode != *filter.TeamCode ||
			filter.ProblemCode != nil && deployment.ProblemCode != *filter.ProblemCode {
			continue
		}
		result = append(result, deployment)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].RequestedAt.After(result[j].RequestedAt) })
	return result, nil
}

func (s *CompetitionStore) CreateDeployment(_ context.Context, deployment core.Deployment) (core.Deployment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.deployments[deployment.RequestID]; ok {
		return core.Deployment{}, core.NewError(http.StatusConflict, "deployment_conflict", "Deployment already exists")
	}
	for _, existing := range s.deployments {
		if existing.TeamCode == deployment.TeamCode && existing.ProblemCode == deployment.ProblemCode && (existing.LatestStatus == core.DeploymentQueued || existing.LatestStatus == core.DeploymentDeploying) {
			return core.Deployment{}, core.NewError(http.StatusConflict, "deployment_in_progress", "A deployment is already in progress")
		}
		if existing.TeamCode == deployment.TeamCode && existing.ProblemCode == deployment.ProblemCode && existing.Revision >= deployment.Revision {
			deployment.Revision = existing.Revision + 1
		}
	}
	if deployment.Revision < 1 {
		deployment.Revision = 1
	}
	s.deployments[deployment.RequestID] = deployment
	return deployment, nil
}

func (s *CompetitionStore) AppendDeploymentEvent(_ context.Context, requestID string, event core.DeploymentEvent, recovery bool) (core.Deployment, bool, error) {
	event.OccurredAt = event.OccurredAt.UTC().Truncate(time.Microsecond)
	s.mu.Lock()
	defer s.mu.Unlock()
	deployment, ok := s.deployments[requestID]
	if !ok {
		return core.Deployment{}, false, core.NewError(http.StatusNotFound, "deployment_not_found", "Deployment was not found")
	}
	for _, existing := range deployment.Events {
		if existing.EventID == event.EventID {
			if !sameDeploymentEvent(existing, event) {
				return core.Deployment{}, false, core.NewError(http.StatusConflict, "duplicate_event_mismatch", "Event ID was already used with a different payload")
			}
			return deployment, true, nil
		}
	}
	for _, existing := range deployment.Events {
		if event.OccurredAt.Before(existing.OccurredAt) || event.OccurredAt.Equal(existing.OccurredAt) && event.EventID < existing.EventID {
			return core.Deployment{}, false, core.NewError(http.StatusConflict, "invalid_deployment_transition", "Deployment event occurred before the latest persisted event")
		}
	}
	if !core.CanTransitionDeployment(deployment.LatestStatus, event.Status, recovery) {
		return core.Deployment{}, false, core.NewError(http.StatusConflict, "invalid_deployment_transition", "Deployment status transition is not allowed")
	}
	deployment.LatestStatus = event.Status
	deployment.Events = append(deployment.Events, event)
	s.deployments[requestID] = deployment
	return deployment, false, nil
}

func (s *CompetitionStore) GetCompetitionState(context.Context) (core.CompetitionState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state, nil
}

func (s *CompetitionStore) ReplaceRule(_ context.Context, markdown, actor string) (core.CompetitionState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.RuleMarkdown, s.state.UpdatedBy, s.state.UpdatedAt = markdown, actor, time.Now().UTC()
	return s.state, nil
}

func (s *CompetitionStore) ReplaceFreezeAt(_ context.Context, freezeAt *time.Time, actor string) (core.CompetitionState, error) {
	if freezeAt != nil {
		value := freezeAt.UTC().Truncate(time.Microsecond)
		freezeAt = &value
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.RankingFreezeAt = freezeAt
	s.rankingSnapshots = make(map[string]core.RankingSnapshot)
	s.state.UpdatedBy, s.state.UpdatedAt = actor, time.Now().UTC()
	return s.state, nil
}

func (s *CompetitionStore) RevealFinal(_ context.Context, at time.Time, actor string) (core.CompetitionState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state.FinalRevealedAt == nil {
		s.state.FinalRevealedAt = &at
	}
	s.state.UpdatedBy, s.state.UpdatedAt = actor, time.Now().UTC()
	return s.state, nil
}

var _ core.Store = (*CompetitionStore)(nil)
