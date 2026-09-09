package core

import (
	"context"
	"time"
)

type Store interface {
	Ping(ctx context.Context) error
	ListTeams(ctx context.Context) ([]Team, error)
	GetTeam(ctx context.Context, code int64) (Team, error)
	CreateTeam(ctx context.Context, team Team) (Team, error)
	UpdateTeam(ctx context.Context, code int64, patch TeamPatch) (Team, error)
	DeleteTeam(ctx context.Context, code int64) error
	ListInvitations(ctx context.Context) ([]Invitation, error)
	CreateInvitation(ctx context.Context, invitation Invitation) (Invitation, error)
	ConsumeInvitation(ctx context.Context, code string, now time.Time, contestant Contestant) (Contestant, error)
	ListContestants(ctx context.Context) ([]Contestant, error)
	GetContestant(ctx context.Context, name string) (Contestant, error)
	GetContestantByDiscord(ctx context.Context, discordID string) (Contestant, error)
	UpdateContestant(ctx context.Context, name string, patch ContestantPatch) (Contestant, error)
	ActiveContent(ctx context.Context) (ContentSnapshot, error)
	GetContent(ctx context.Context, commit string) (ContentSnapshot, error)
	ContentAt(ctx context.Context, at time.Time) (ContentSnapshot, error)
	ActivateContent(ctx context.Context, snapshot ContentSnapshot, expectedCurrentCommit string) (ContentSnapshot, error)
	ListAnswers(ctx context.Context, filter AnswerFilter) ([]Answer, error)
	GetAnswer(ctx context.Context, teamCode int64, problemCode string, number int32) (Answer, error)
	SubmitAnswer(ctx context.Context, answer Answer, interval time.Duration) (Answer, time.Duration, error)
	ListMarkingResults(ctx context.Context) ([]MarkingResult, error)
	CreateMarkingResult(ctx context.Context, result MarkingResult) (MarkingResult, error)
	RecalculateScores(ctx context.Context, selections map[int64]map[string]Score, actor string) error
	GetRankingSnapshot(ctx context.Context, frozenAt time.Time) (RankingSnapshot, bool, error)
	PutRankingSnapshot(ctx context.Context, snapshot RankingSnapshot, actor string) (RankingSnapshot, error)
	ListDeployments(ctx context.Context, filter DeploymentFilter) ([]Deployment, error)
	CreateDeployment(ctx context.Context, deployment Deployment) (Deployment, error)
	AppendDeploymentEvent(ctx context.Context, requestID string, event DeploymentEvent, recovery bool) (Deployment, bool, error)
	GetCompetitionState(ctx context.Context) (CompetitionState, error)
	ReplaceRule(ctx context.Context, markdown, actor string) (CompetitionState, error)
	ReplaceFreezeAt(ctx context.Context, freezeAt *time.Time, actor string) (CompetitionState, error)
	RevealFinal(ctx context.Context, at time.Time, actor string) (CompetitionState, error)
}

type TeamPatch struct {
	Name         *string
	Organization *string
	MemberLimit  *int32
	Color        *string
}

type ContestantPatch struct {
	DisplayName      *string
	SelfIntroduction *string
}

type AnswerFilter struct {
	TeamCode    *int64
	ProblemCode *string
	Before      *time.Time
}

type DeploymentFilter struct {
	TeamCode    *int64
	ProblemCode *string
}
