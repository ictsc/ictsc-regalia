package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
)

const (
	OAuthTTL      = 10 * time.Minute
	SignupTTL     = 10 * time.Minute
	ContestantTTL = 3 * 24 * time.Hour
	AdminTTL      = 8 * time.Hour
)

type Config struct {
	ContestantRedirectURI string
	AdminRedirectURI      string
	ContentRepository     string
	ContentRef            string
	CallbackBaseURL       string
	AdminGuildID          string
	ContestantGuildID     string
	DiscordRoleTeams      map[string]int64
	AdminRoleIDs          map[string]struct{}
}

type Service struct {
	Store       core.Store
	Sessions    session.Store
	Discord     Discord
	Content     ContentSource
	Deployments DeploymentGateway
	Events      EventBus
	Config      Config
	Now         Clock

	statusMu      sync.RWMutex
	contentStatus ContentStatus
}

type ContentStatus struct {
	State             string
	ActiveCommit      *string
	ActivatedAt       *time.Time
	SourceRef         *string
	LastAttemptAt     *time.Time
	LastAttemptCommit *string
	LastError         *string
	ServingLastGood   bool
}

type AuthStart struct {
	URL          string
	SessionToken string
}

type AuthComplete struct {
	Next         string
	SessionToken string
	SessionKind  session.Kind
	Result       DiscordResult
}

func New(store core.Store, sessions session.Store, config Config) *Service {
	return &Service{
		Store: store, Sessions: sessions, Config: config,
		Now:           func() time.Time { return time.Now().UTC() },
		contentStatus: ContentStatus{State: "NEVER"},
	}
}

func (s *Service) BeginDiscord(ctx context.Context, admin bool, next string) (AuthStart, error) {
	if next == "" {
		if admin {
			next = "/admin/"
		} else {
			next = "/"
		}
	}
	if !validOAuthNext(next) {
		return AuthStart{}, core.NewError(http.StatusUnprocessableEntity, "validation_error", "next must be an origin-relative path without a query")
	}
	state, err := randomURLToken(32)
	if err != nil {
		return AuthStart{}, core.WrapError(http.StatusInternalServerError, "internal_error", "Could not create OAuth state", err)
	}
	verifier, err := randomURLToken(48)
	if err != nil {
		return AuthStart{}, core.WrapError(http.StatusInternalServerError, "internal_error", "Could not create PKCE verifier", err)
	}
	kind := session.KindOAuth
	redirectURI := s.Config.ContestantRedirectURI
	if admin {
		kind = session.KindAdminOAuth
		redirectURI = s.Config.AdminRedirectURI
	}
	token, err := s.Sessions.Create(ctx, session.Data{Kind: kind, OAuthState: state, PKCEVerifier: verifier, Next: next}, OAuthTTL)
	if err != nil {
		return AuthStart{}, core.WrapError(http.StatusInternalServerError, "internal_error", "Could not persist OAuth state", err)
	}
	digest := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(digest[:])
	return AuthStart{URL: s.Discord.AuthorizationURL(state, challenge, redirectURI, admin), SessionToken: token}, nil
}

func (s *Service) CompleteDiscord(ctx context.Context, admin bool, oauthToken, code, state string) (AuthComplete, error) {
	kind := session.KindOAuth
	redirectURI := s.Config.ContestantRedirectURI
	if admin {
		kind = session.KindAdminOAuth
		redirectURI = s.Config.AdminRedirectURI
	}
	oauthData, err := s.Sessions.Consume(ctx, oauthToken, kind)
	if err != nil {
		if errors.Is(err, session.ErrNotFound) {
			return AuthComplete{}, core.NewError(http.StatusUnauthorized, "oauth_state_invalid", "OAuth state is invalid or expired")
		}
		return AuthComplete{}, core.WrapError(http.StatusInternalServerError, "internal_error", "Could not read OAuth state", err)
	}
	if state == "" || state != oauthData.OAuthState || code == "" {
		return AuthComplete{}, core.NewError(http.StatusUnauthorized, "oauth_state_invalid", "OAuth state is invalid or expired")
	}
	result, err := s.Discord.Exchange(ctx, code, oauthData.PKCEVerifier, redirectURI, admin)
	if err != nil {
		switch discordOAuthErrorKind(err) {
		case "guild_membership":
			return AuthComplete{}, core.WrapError(http.StatusForbidden, "guild_membership_required", "Configured Discord guild membership is required", err)
		case "authentication_rejected":
			return AuthComplete{}, core.WrapError(http.StatusUnauthorized, "oauth_state_invalid", "Discord rejected authentication", err)
		}
		return AuthComplete{}, core.WrapError(http.StatusBadGateway, "upstream_unavailable", "Discord authentication failed", err)
	}
	staff := false
	if len(s.Config.DiscordRoleTeams) > 0 && result.GuildID == s.Config.AdminGuildID {
		for _, role := range result.RoleIDs {
			if _, ok := s.Config.AdminRoleIDs[role]; ok {
				staff = true
				break
			}
		}
	}
	if admin || staff {
		if !admin {
			oauthData.Next = "/admin/"
		}
		if result.GuildID != s.Config.AdminGuildID {
			return AuthComplete{}, core.NewError(http.StatusForbidden, "guild_membership_required", "Configured Discord guild membership is required")
		}
		allowed := false
		for _, role := range result.RoleIDs {
			if _, ok := s.Config.AdminRoleIDs[role]; ok {
				allowed = true
				break
			}
		}
		if !allowed {
			return AuthComplete{}, core.NewError(http.StatusForbidden, "admin_role_required", "Configured Discord admin role is required")
		}
		token, createErr := s.Sessions.Create(ctx, session.Data{
			Kind: session.KindAdmin, AdminName: result.Identity.DisplayName, Discord: result.Identity,
			GuildID: result.GuildID, RoleIDs: append([]string(nil), result.RoleIDs...),
		}, AdminTTL)
		return AuthComplete{Next: oauthData.Next, SessionToken: token, SessionKind: session.KindAdmin, Result: result}, createErr
	}
	if s.Config.ContestantGuildID != "" && result.GuildID != s.Config.ContestantGuildID {
		return AuthComplete{}, core.NewError(http.StatusForbidden, "guild_membership_required", "Configured Discord guild membership is required")
	}
	var registrationTeam int64
	if len(s.Config.DiscordRoleTeams) > 0 {
		registrationTeam, err = s.TeamFromDiscordRoles(result.RoleIDs)
		if err != nil {
			return AuthComplete{}, err
		}
	}
	contestant, contestantErr := s.Store.GetContestantByDiscord(ctx, result.Identity.ID)
	if contestantErr == nil {
		if registrationTeam != 0 && contestant.TeamCode != registrationTeam {
			return AuthComplete{}, core.NewError(http.StatusForbidden, "team_role_mismatch", "Discord team role does not match the registered team")
		}
		token, createErr := s.Sessions.Create(ctx, session.Data{
			Kind: session.KindContestant, ContestantName: contestant.Name, Discord: result.Identity,
		}, ContestantTTL)
		return AuthComplete{Next: oauthData.Next, SessionToken: token, SessionKind: session.KindContestant, Result: result}, createErr
	}
	var notFound *core.Error
	if !errors.As(contestantErr, &notFound) || notFound.Status != http.StatusNotFound {
		return AuthComplete{}, contestantErr
	}
	token, createErr := s.Sessions.Create(ctx, session.Data{Kind: session.KindSignup, Discord: result.Identity, GuildID: result.GuildID, RoleIDs: result.RoleIDs}, SignupTTL)
	return AuthComplete{Next: oauthData.Next, SessionToken: token, SessionKind: session.KindSignup, Result: result}, createErr
}

func (s *Service) SignUp(ctx context.Context, signupToken, name, displayName, invitationCode string) (core.Contestant, string, error) {
	signup, err := s.Sessions.Get(ctx, signupToken, session.KindSignup)
	if err != nil {
		if errors.Is(err, session.ErrNotFound) {
			return core.Contestant{}, "", core.NewError(http.StatusUnauthorized, "invalid_session", "Signup session is invalid or expired")
		}
		return core.Contestant{}, "", core.WrapError(http.StatusInternalServerError, "internal_error", "Could not consume signup session", err)
	}

	input := core.Contestant{Name: name, DisplayName: displayName, DiscordID: signup.Discord.ID}
	var contestant core.Contestant
	if len(s.Config.DiscordRoleTeams) > 0 {
		if s.Config.ContestantGuildID == "" || signup.GuildID != s.Config.ContestantGuildID {
			return core.Contestant{}, "", core.NewError(http.StatusForbidden, "guild_membership_required", "Discord guild membership is required")
		}
		input.TeamCode, err = s.TeamFromDiscordRoles(signup.RoleIDs)
		if err != nil {
			return core.Contestant{}, "", err
		}
		contestant, err = s.Store.RegisterContestant(ctx, input)
	} else {
		if invitationCode == "" {
			return core.Contestant{}, "", core.NewError(http.StatusUnprocessableEntity, "validation_error", "Invitation code is required")
		}
		contestant, err = s.Store.ConsumeInvitation(ctx, invitationCode, s.Now(), input)
	}
	if err != nil {
		return core.Contestant{}, "", err
	}
	consumed, err := s.Sessions.Consume(ctx, signupToken, session.KindSignup)
	if err != nil {
		if errors.Is(err, session.ErrNotFound) {
			return core.Contestant{}, "", core.NewError(http.StatusUnauthorized, "invalid_session", "Signup session was already used")
		}
		return core.Contestant{}, "", core.WrapError(http.StatusInternalServerError, "internal_error", "Could not consume signup session", err)
	}
	if consumed.Discord.ID != signup.Discord.ID {
		return core.Contestant{}, "", core.NewError(http.StatusInternalServerError, "internal_error", "Consumed signup identity did not match")
	}
	token, err := s.Sessions.Create(ctx, session.Data{
		Kind: session.KindContestant, ContestantName: contestant.Name, Discord: signup.Discord,
	}, ContestantTTL)
	return contestant, token, err
}

func (s *Service) Impersonate(ctx context.Context, admin session.Data, contestantName string) (string, error) {
	contestant, err := s.Store.GetContestant(ctx, contestantName)
	if err != nil {
		return "", err
	}
	return s.Sessions.Create(ctx, session.Data{
		Kind: session.KindContestant, ContestantName: contestant.Name, ImpersonatedBy: admin.AdminName,
	}, AdminTTL)
}

func (s *Service) RefreshContent(ctx context.Context, claims ActionsClaims, commit string) (core.ContentSnapshot, error) {
	now := s.Now()
	s.statusMu.Lock()
	if s.contentStatus.State == "REFRESHING" {
		s.statusMu.Unlock()
		return core.ContentSnapshot{}, core.NewError(http.StatusConflict, "content_refresh_in_progress", "Another content refresh is in progress")
	}
	s.contentStatus.State = "REFRESHING"
	s.contentStatus.LastAttemptAt = &now
	s.contentStatus.LastAttemptCommit = &commit
	s.contentStatus.LastError = nil
	s.statusMu.Unlock()

	active, activeErr := s.Store.ActiveContent(ctx)
	expectedCurrent := core.NoActiveContentCommit
	if activeErr != nil {
		var domainErr *core.Error
		if !errors.As(activeErr, &domainErr) || domainErr.Code != "content_not_available" {
			s.setContentFailure(commit, activeErr, false)
			return core.ContentSnapshot{}, activeErr
		}
	}
	if activeErr == nil {
		expectedCurrent = active.CommitSHA
		if active.CommitSHA == commit {
			s.setContentSuccess(active)
			return active, nil
		}
		isDescendant, err := s.Content.IsAncestor(ctx, claims.Repository, active.CommitSHA, commit)
		if err != nil {
			s.setContentFailure(commit, err, true)
			return core.ContentSnapshot{}, core.WrapError(http.StatusBadGateway, "upstream_unavailable", "Could not compare Git commits", err)
		}
		if !isDescendant {
			err := core.NewError(http.StatusConflict, "content_rollback_rejected", "Content refresh cannot roll back the active commit")
			s.setContentFailure(commit, err, true)
			return core.ContentSnapshot{}, err
		}
	}
	snapshot, err := s.Content.Fetch(ctx, claims.Repository, claims.Ref, commit)
	if err != nil {
		hasLastGood := activeErr == nil
		s.setContentFailure(commit, err, hasLastGood)
		if IsInvalidContentError(err) {
			return core.ContentSnapshot{}, core.WrapError(http.StatusUnprocessableEntity, "content_invalid", "Content snapshot is invalid", err)
		}
		return core.ContentSnapshot{}, core.WrapError(http.StatusBadGateway, "upstream_unavailable", "Could not fetch GitHub content", err)
	}
	if err := ValidateSnapshot(snapshot); err != nil {
		s.setContentFailure(commit, err, activeErr == nil)
		return core.ContentSnapshot{}, core.WrapError(http.StatusUnprocessableEntity, "content_invalid", "Content snapshot is invalid", err)
	}
	snapshot, err = s.Store.ActivateContent(ctx, snapshot, expectedCurrent)
	if err != nil {
		s.setContentFailure(commit, err, activeErr == nil)
		return core.ContentSnapshot{}, err
	}
	s.setContentSuccess(snapshot)
	return snapshot, nil
}

func (s *Service) ContentStatus(ctx context.Context) ContentStatus {
	s.statusMu.RLock()
	status := s.contentStatus
	s.statusMu.RUnlock()
	if status.ActiveCommit == nil {
		if snapshot, err := s.Store.ActiveContent(ctx); err == nil {
			status.ActiveCommit = &snapshot.CommitSHA
			status.ActivatedAt = &snapshot.ActivatedAt
			status.SourceRef = &snapshot.Ref
			if status.State == "NEVER" {
				status.State = "SUCCEEDED"
			}
		}
	}
	return status
}

func (s *Service) setContentSuccess(snapshot core.ContentSnapshot) {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	s.contentStatus.State = "SUCCEEDED"
	s.contentStatus.ActiveCommit = &snapshot.CommitSHA
	s.contentStatus.ActivatedAt = &snapshot.ActivatedAt
	s.contentStatus.SourceRef = &snapshot.Ref
	s.contentStatus.ServingLastGood = false
	s.contentStatus.LastError = nil
}

func (s *Service) setContentFailure(commit string, err error, lastGood bool) {
	message := err.Error()
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	s.contentStatus.State = "FAILED"
	s.contentStatus.LastAttemptCommit = &commit
	s.contentStatus.LastError = &message
	s.contentStatus.ServingLastGood = lastGood
}

func ValidateSnapshot(snapshot core.ContentSnapshot) error {
	if matched, _ := regexp.MatchString(`^[0-9a-f]{40,64}$`, snapshot.CommitSHA); !matched {
		return fmt.Errorf("commit must be 40 to 64 lowercase hexadecimal characters")
	}
	sections := append([]core.Section(nil), snapshot.Manifest.Sections...)
	sort.Slice(sections, func(i, j int) bool { return sections[i].Beginning.Before(sections[j].Beginning) })
	problemCodes := make(map[string]struct{}, len(snapshot.Manifest.Problems))
	for _, problem := range snapshot.Manifest.Problems {
		if _, exists := problemCodes[problem.Code]; exists {
			return fmt.Errorf("duplicate problem code %q", problem.Code)
		}
		problemCodes[problem.Code] = struct{}{}
		if problem.MaxScore <= 0 || strings.TrimSpace(problem.Title) == "" || strings.TrimSpace(problem.Body) == "" {
			return fmt.Errorf("problem %q is incomplete", problem.Code)
		}
		for _, file := range []string{problem.BodyPath, problem.ExplanationPath} {
			if file != "" && !safeRelativePath(file) {
				return fmt.Errorf("problem %q contains unsafe path %q", problem.Code, file)
			}
		}
		if err := validateRedeployRule(problem.Redeploy); err != nil {
			return fmt.Errorf("problem %q: %w", problem.Code, err)
		}
	}
	sectionSlugs := make(map[string]struct{}, len(sections))
	problemReferences := make(map[string]int, len(problemCodes))
	for index, section := range sections {
		if section.Slug == "" || !section.Beginning.Before(section.Ending) {
			return fmt.Errorf("section %q has invalid interval", section.Slug)
		}
		if _, exists := sectionSlugs[section.Slug]; exists {
			return fmt.Errorf("duplicate section slug %q", section.Slug)
		}
		sectionSlugs[section.Slug] = struct{}{}
		if index > 0 && section.Beginning.Before(sections[index-1].Ending) {
			return fmt.Errorf("sections %q and %q overlap", sections[index-1].Slug, section.Slug)
		}
		for _, code := range section.ProblemIDs {
			if _, ok := problemCodes[code]; !ok {
				return fmt.Errorf("section %q references unknown problem %q", section.Slug, code)
			}
			problem, _ := snapshot.Problem(code)
			if problem.SectionSlug != section.Slug {
				return fmt.Errorf("problem %q declares section %q but is referenced by %q", code, problem.SectionSlug, section.Slug)
			}
			problemReferences[code]++
		}
	}
	for _, problem := range snapshot.Manifest.Problems {
		if _, exists := sectionSlugs[problem.SectionSlug]; !exists || problemReferences[problem.Code] != 1 {
			return fmt.Errorf("problem %q must be referenced exactly once by declared section %q", problem.Code, problem.SectionSlug)
		}
	}
	announcements := make(map[string]struct{}, len(snapshot.Manifest.Announcements))
	for _, announcement := range snapshot.Manifest.Announcements {
		if _, exists := announcements[announcement.Slug]; exists {
			return fmt.Errorf("duplicate announcement slug %q", announcement.Slug)
		}
		announcements[announcement.Slug] = struct{}{}
		if announcement.MarkdownPath != "" && !safeRelativePath(announcement.MarkdownPath) {
			return fmt.Errorf("announcement %q contains unsafe path", announcement.Slug)
		}
	}
	if snapshot.Manifest.RulePath != "" && !safeRelativePath(snapshot.Manifest.RulePath) {
		return fmt.Errorf("rule path is unsafe")
	}
	return nil
}

func validateRedeployRule(rule core.RedeployRule) error {
	if rule.Type == core.RedeployPercentage {
		if rule.Threshold == nil || rule.Percentage == nil || *rule.Threshold < 0 || *rule.Percentage < 0 || *rule.Percentage > 99 {
			return fmt.Errorf("percentage rule requires valid threshold and percentage")
		}
		return nil
	}
	if rule.Type != core.RedeployManual && rule.Type != core.RedeployUnredeployable {
		return fmt.Errorf("unknown redeploy rule")
	}
	if rule.Threshold != nil || rule.Percentage != nil {
		return fmt.Errorf("non-percentage rule must have null threshold and percentage")
	}
	return nil
}

func safeRelativePath(value string) bool {
	clean := path.Clean(value)
	return value != "" && clean == value && clean != "." && !strings.HasPrefix(clean, "../") && !strings.HasPrefix(clean, "/")
}

func (s *Service) SubmitAnswer(ctx context.Context, contestant core.Contestant, problemCode, body string) (core.Answer, time.Duration, error) {
	snapshot, err := s.Store.ActiveContent(ctx)
	if err != nil {
		return core.Answer{}, 0, err
	}
	problem, ok := snapshot.Problem(problemCode)
	if !ok {
		return core.Answer{}, 0, core.NewError(http.StatusNotFound, "resource_not_found", "Problem was not found")
	}
	now := s.Now()
	section, ok := findSection(snapshot, problem.SectionSlug)
	if !ok || now.Before(section.Beginning) || !now.Before(section.Ending) {
		return core.Answer{}, 0, core.NewError(http.StatusConflict, "submission_closed", "Problem is not currently submittable")
	}
	team := contestant.TeamCode
	deployments, err := s.Store.ListDeployments(ctx, core.DeploymentFilter{TeamCode: &team, ProblemCode: &problemCode})
	if err != nil {
		return core.Answer{}, 0, err
	}
	answer := core.Answer{
		TeamCode: contestant.TeamCode, ProblemCode: problemCode, AuthorName: contestant.Name,
		Body: body, SubmittedAt: now, ContentCommit: snapshot.CommitSHA, MaxScore: problem.MaxScore,
		RedeployRule: problem.Redeploy, DeploymentsBefore: int32(len(deployments)),
	}
	return s.Store.SubmitAnswer(ctx, answer, core.AnswerInterval)
}

func (s *Service) Scores(ctx context.Context, before *time.Time) (map[int64]map[string]core.Score, error) {
	answers, err := s.Store.ListAnswers(ctx, core.AnswerFilter{Before: before})
	if err != nil {
		return nil, err
	}
	markings, err := s.Store.ListMarkingResults(ctx)
	if err != nil {
		return nil, err
	}
	if before != nil {
		filtered := markings[:0]
		for _, marking := range markings {
			if !marking.CreatedAt.After(*before) {
				filtered = append(filtered, marking)
			}
		}
		markings = filtered
	}
	return core.SelectBestScores(answers, markings), nil
}

func (s *Service) Ranking(ctx context.Context, admin bool) ([]core.RankingEntry, bool, *time.Time, error) {
	state, err := s.Store.GetCompetitionState(ctx)
	if err != nil {
		return nil, false, nil, err
	}
	frozen := false
	var before *time.Time
	if !admin && state.FinalRevealedAt == nil && state.RankingFreezeAt != nil && !s.Now().Before(*state.RankingFreezeAt) {
		frozen, before = true, state.RankingFreezeAt
	}
	if frozen {
		entries, freezeErr := s.frozenRanking(ctx, *before)
		return entries, true, before, freezeErr
	}
	snapshot, err := s.Store.ActiveContent(ctx)
	if err != nil {
		return nil, false, nil, err
	}
	var scores map[int64]map[string]core.Score
	if admin {
		scores, err = s.Scores(ctx, nil)
	} else {
		scores, err = s.publicScores(ctx, before, state.FinalRevealedAt != nil)
	}
	if err != nil {
		return nil, false, nil, err
	}
	teams, err := s.Store.ListTeams(ctx)
	if err != nil {
		return nil, false, nil, err
	}
	active := make(map[string]struct{}, len(snapshot.Manifest.Problems))
	for _, problem := range snapshot.Manifest.Problems {
		active[problem.Code] = struct{}{}
	}
	return core.BuildRanking(teams, active, scores), frozen, before, nil
}

func (s *Service) QueueDeployment(ctx context.Context, teamCode int64, problemCode string, enforceContestWindow bool) (core.Deployment, error) {
	snapshot, err := s.Store.ActiveContent(ctx)
	if err != nil {
		return core.Deployment{}, err
	}
	problem, ok := snapshot.Problem(problemCode)
	if !ok {
		return core.Deployment{}, core.NewError(http.StatusNotFound, "resource_not_found", "Problem was not found")
	}
	if enforceContestWindow {
		now := s.Now()
		open := false
		for _, section := range snapshot.Manifest.Sections {
			if section.Slug == problem.SectionSlug {
				open = !now.Before(section.Beginning) && now.Before(section.Ending)
				break
			}
		}
		if !open {
			return core.Deployment{}, core.NewError(http.StatusConflict, "deployment_not_allowed", "Problem is outside its deployment window")
		}
	}
	if problem.Redeploy.Type == core.RedeployUnredeployable {
		return core.Deployment{}, core.NewError(http.StatusConflict, "deployment_not_allowed", "Problem cannot be redeployed")
	}
	requestID, err := randomUUID()
	if err != nil {
		return core.Deployment{}, core.WrapError(http.StatusInternalServerError, "internal_error", "Could not create deployment ID", err)
	}
	queuedEventID, err := randomUUID()
	if err != nil {
		return core.Deployment{}, core.WrapError(http.StatusInternalServerError, "internal_error", "Could not create deployment event ID", err)
	}
	requestedAt := s.Now().UTC().Truncate(time.Microsecond)
	deployment, err := s.Store.CreateDeployment(ctx, core.Deployment{
		RequestID: requestID, TeamCode: teamCode, ProblemCode: problemCode, ContentCommit: snapshot.CommitSHA,
		RequestedAt: requestedAt, LatestStatus: core.DeploymentQueued,
		Events: []core.DeploymentEvent{{EventID: queuedEventID, OccurredAt: requestedAt, Status: core.DeploymentQueued}},
	})
	if err != nil {
		return core.Deployment{}, err
	}
	_ = s.publishDeployment(ctx, deployment)
	gatewayRequest := DeploymentRequest{
		RequestID: requestID, TeamCode: teamCode, ProblemCode: problemCode, Revision: deployment.Revision,
		ContentCommit: snapshot.CommitSHA, CallbackURL: fmt.Sprintf("%s/api/v1/admin/deployments/%d/%s/%d/events", strings.TrimRight(s.Config.CallbackBaseURL, "/"), teamCode, problemCode, deployment.Revision),
	}
	queueErr := s.Deployments.Queue(ctx, gatewayRequest)
	if queueErr != nil && !definitiveDeploymentRejection(queueErr) {
		queueErr = s.Deployments.Queue(ctx, gatewayRequest)
	}
	if queueErr != nil {
		if definitiveDeploymentRejection(queueErr) {
			if failErr := s.FailQueuedDeployment(ctx, deployment, queueErr.Error()); failErr != nil {
				return core.Deployment{}, failErr
			}
		}
		return core.Deployment{}, core.WrapError(http.StatusBadGateway, "upstream_unavailable", "SState rejected deployment request", queueErr)
	}
	filterTeam, filterProblem := teamCode, problemCode
	if current, listErr := s.Store.ListDeployments(ctx, core.DeploymentFilter{TeamCode: &filterTeam, ProblemCode: &filterProblem}); listErr == nil {
		for _, candidate := range current {
			if candidate.RequestID == requestID {
				deployment = candidate
				break
			}
		}
	}
	return deployment, nil
}

func (s *Service) ApplyDeploymentEvent(ctx context.Context, teamCode int64, problemCode string, revision int32, event core.DeploymentEvent, recovery bool) (core.Deployment, bool, error) {
	team, problem := teamCode, problemCode
	deployments, err := s.Store.ListDeployments(ctx, core.DeploymentFilter{TeamCode: &team, ProblemCode: &problem})
	if err != nil {
		return core.Deployment{}, false, err
	}
	for _, deployment := range deployments {
		if deployment.Revision != revision {
			continue
		}
		updated, duplicate, err := s.Store.AppendDeploymentEvent(ctx, deployment.RequestID, event, recovery)
		if err != nil {
			return updated, duplicate, err
		}
		// Republish duplicate callbacks too: the first attempt may have committed
		// to PostgreSQL while Redis was temporarily unavailable. Reusing the
		// idempotent DB result restores live SSE without adding another event.
		if err := s.publishDeployment(ctx, updated); err != nil {
			return updated, duplicate, core.WrapError(http.StatusInternalServerError, "internal_error", "Could not publish deployment update", err)
		}
		return updated, duplicate, nil
	}
	return core.Deployment{}, false, core.NewError(http.StatusNotFound, "resource_not_found", "Deployment was not found")
}

func (s *Service) publishDeployment(ctx context.Context, deployment core.Deployment) error {
	if s.Events == nil {
		return nil
	}
	encoded, err := json.Marshal(deployment)
	if err != nil {
		return err
	}
	return s.Events.Publish(ctx, encoded)
}

func findSection(snapshot core.ContentSnapshot, slug string) (core.Section, bool) {
	for _, section := range snapshot.Manifest.Sections {
		if section.Slug == slug {
			return section, true
		}
	}
	return core.Section{}, false
}

func randomURLToken(size int) (string, error) {
	value := make([]byte, size)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func randomUUID() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80
	hexValue := hex.EncodeToString(value)
	return hexValue[0:8] + "-" + hexValue[8:12] + "-" + hexValue[12:16] + "-" + hexValue[16:20] + "-" + hexValue[20:], nil
}

func (s *Service) TeamFromDiscordRoles(roles []string) (int64, error) {
	var team int64
	for _, role := range roles {
		code, ok := s.Config.DiscordRoleTeams[role]
		if !ok {
			continue
		}
		if team != 0 && team != code {
			return 0, core.NewError(http.StatusForbidden, "ambiguous_team_roles", "Multiple team roles are assigned; contact staff")
		}
		team = code
	}
	if team == 0 {
		return 0, core.NewError(http.StatusForbidden, "team_role_required", "A registered team role is required")
	}
	return team, nil
}
