package core

// RealtimeKind names a change that connected screens should re-read through
// the REST API. Realtime events never carry scores or answer bodies, so the
// REST handlers stay the only place where visibility rules are evaluated.
type RealtimeKind string

const (
	// Published across API instances after a successful write.
	RealtimeAnswerSubmitted       RealtimeKind = "ANSWER_SUBMITTED"
	RealtimeAnswerWorkflowUpdated RealtimeKind = "ANSWER_WORKFLOW_UPDATED"
	RealtimeMarkingUpdated        RealtimeKind = "MARKING_UPDATED"
	RealtimeScoresRecalculated    RealtimeKind = "SCORES_RECALCULATED"
	RealtimeCompetitionUpdated    RealtimeKind = "COMPETITION_UPDATED"
	RealtimeContentUpdated        RealtimeKind = "CONTENT_UPDATED"
	RealtimeTeamsUpdated          RealtimeKind = "TEAMS_UPDATED"

	// Derived by every API instance from the clock and delivered only to its
	// own connections, so multiple instances never duplicate them.
	RealtimeAnswerPublishDelayElapsed RealtimeKind = "ANSWER_PUBLISH_DELAY_ELAPSED"
	RealtimeScheduleReached           RealtimeKind = "SCHEDULE_REACHED"
	RealtimeRankingUpdated            RealtimeKind = "RANKING_UPDATED"
)

type RealtimeEvent struct {
	Kind         RealtimeKind `json:"kind"`
	TeamCode     int64        `json:"team_code,omitempty"`
	ProblemCode  string       `json:"problem_code,omitempty"`
	AnswerNumber int32        `json:"answer_number,omitempty"`
}
