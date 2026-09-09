package core

import "time"

const (
	AnswerInterval = 20 * time.Minute
	PublishDelay   = 20 * time.Minute
)

type Team struct {
	Code         int64
	Name         string
	Organization string
	MemberLimit  int32
	Color        string
}

type DiscordIdentity struct {
	ID          string
	Username    string
	DisplayName string
}

type Contestant struct {
	Name             string
	DisplayName      string
	SelfIntroduction string
	DiscordID        string
	TeamCode         int64
}

type Invitation struct {
	Code      string
	TeamCode  int64
	CreatedAt time.Time
	ExpiresAt time.Time
}

type RedeployRuleType string

const (
	RedeployUnredeployable RedeployRuleType = "UNREDEPLOYABLE"
	RedeployPercentage     RedeployRuleType = "PERCENTAGE_PENALTY"
	RedeployManual         RedeployRuleType = "MANUAL"
)

type RedeployRule struct {
	Type       RedeployRuleType `json:"type" yaml:"type"`
	Threshold  *int32           `json:"penalty_threshold" yaml:"penalty_threshold"`
	Percentage *int32           `json:"penalty_percentage" yaml:"penalty_percentage"`
}

type Problem struct {
	Code            string       `json:"code" yaml:"code"`
	Title           string       `json:"title" yaml:"title"`
	MaxScore        int32        `json:"max_score" yaml:"max_score"`
	Category        string       `json:"category" yaml:"category"`
	SectionSlug     string       `json:"section_slug" yaml:"section_slug"`
	Type            string       `json:"type" yaml:"type"`
	Body            string       `json:"body" yaml:"-"`
	Explanation     string       `json:"explanation" yaml:"-"`
	BodyPath        string       `json:"body_path,omitempty" yaml:"body_path"`
	ExplanationPath string       `json:"explanation_path,omitempty" yaml:"explanation_path"`
	Redeploy        RedeployRule `json:"redeploy_rule" yaml:"redeploy_rule"`
}

type Section struct {
	Slug       string    `json:"slug" yaml:"slug"`
	Beginning  time.Time `json:"beginning" yaml:"beginning"`
	Ending     time.Time `json:"ending" yaml:"ending"`
	ProblemIDs []string  `json:"problem_codes,omitempty" yaml:"problems"`
}

type Announcement struct {
	Slug          string    `json:"slug" yaml:"slug"`
	Title         string    `json:"title" yaml:"title"`
	Markdown      string    `json:"markdown" yaml:"-"`
	MarkdownPath  string    `json:"markdown_path,omitempty" yaml:"markdown_path"`
	EffectiveFrom time.Time `json:"effective_from" yaml:"effective_from"`
}

type Manifest struct {
	Version       int            `json:"version" yaml:"version"`
	Sections      []Section      `json:"sections" yaml:"sections"`
	Problems      []Problem      `json:"problems" yaml:"problems"`
	Announcements []Announcement `json:"announcements" yaml:"announcements"`
	RuleMarkdown  string         `json:"rule_markdown,omitempty" yaml:"rule_markdown,omitempty"`
	RulePath      string         `json:"rule_path,omitempty" yaml:"rule_path,omitempty"`
}

type ContentSnapshot struct {
	CommitSHA   string
	Repository  string
	Ref         string
	Manifest    Manifest
	FetchedAt   time.Time
	ActivatedAt time.Time
}

func (s ContentSnapshot) Problem(code string) (Problem, bool) {
	for _, problem := range s.Manifest.Problems {
		if problem.Code == code {
			return problem, true
		}
	}
	return Problem{}, false
}

func (s ContentSnapshot) Announcement(slug string) (Announcement, bool) {
	for _, announcement := range s.Manifest.Announcements {
		if announcement.Slug == slug {
			return announcement, true
		}
	}
	return Announcement{}, false
}

type Answer struct {
	TeamCode          int64
	ProblemCode       string
	Number            int32
	AuthorName        string
	Body              string
	SubmittedAt       time.Time
	ContentCommit     string
	MaxScore          int32
	RedeployRule      RedeployRule
	DeploymentsBefore int32
}

type Visibility string

const (
	VisibilityPrivate Visibility = "PRIVATE"
	VisibilityTeam    Visibility = "TEAM"
	VisibilityPublic  Visibility = "PUBLIC"
)

type MarkingResult struct {
	ID           string
	TeamCode     int64
	ProblemCode  string
	AnswerNumber int32
	Judge        string
	MarkedScore  int32
	Rationale    string
	CreatedAt    time.Time
	Visibility   Visibility
}

type Score struct {
	TeamCode        int64
	ProblemCode     string
	AnswerNumber    int32
	MarkingResultID string
	MarkedScore     int32
	Penalty         int32
	EffectiveScore  int32
	MaxScore        int32
	SubmittedAt     time.Time
	ContentCommit   string
}

type DeploymentStatus string

const (
	DeploymentQueued    DeploymentStatus = "QUEUED"
	DeploymentDeploying DeploymentStatus = "DEPLOYING"
	DeploymentCompleted DeploymentStatus = "COMPLETED"
	DeploymentFailed    DeploymentStatus = "FAILED"
)

type DeploymentEvent struct {
	EventID    string
	OccurredAt time.Time
	Status     DeploymentStatus
	Message    *string
}

type Deployment struct {
	RequestID     string
	TeamCode      int64
	ProblemCode   string
	Revision      int32
	ContentCommit string
	RequestedAt   time.Time
	LatestStatus  DeploymentStatus
	Events        []DeploymentEvent
}

type CompetitionState struct {
	RuleMarkdown    string
	RankingFreezeAt *time.Time
	FinalRevealedAt *time.Time
	UpdatedAt       time.Time
	UpdatedBy       string
}

type RankingEntry struct {
	Rank                      int32
	Team                      Team
	Score                     int64
	LastEffectiveSubmissionAt *time.Time
}

type RankingSnapshot struct {
	FrozenAt      time.Time
	ContentCommit string
	Entries       []RankingEntry
}
