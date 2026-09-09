package core

// LatestMarkingForAnswer returns the newest marking selected for an answer.
// Callers do not need to duplicate the private composite-key representation.
func LatestMarkingForAnswer(latest map[string]MarkingResult, teamCode int64, problemCode string, answerNumber int32) (MarkingResult, bool) {
	marking, ok := latest[answerKey(teamCode, problemCode, answerNumber)]
	return marking, ok
}
