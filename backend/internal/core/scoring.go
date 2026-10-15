package core

import (
	"sort"
	"strconv"
	"time"
)

func DeploymentPenalty(maxScore int32, rule RedeployRule, deploymentsBefore int32) int32 {
	if rule.Type != RedeployPercentage || rule.Threshold == nil || rule.Percentage == nil {
		return 0
	}
	excess := deploymentsBefore - *rule.Threshold
	if excess <= 0 {
		return 0
	}
	perDeployment := int32((int64(maxScore) * int64(*rule.Percentage)) / 100)
	return perDeployment * excess
}

func EffectiveScore(answer Answer, marking MarkingResult) Score {
	marked := marking.MarkedScore
	if marked > answer.MaxScore {
		marked = answer.MaxScore
	}
	penalty := DeploymentPenalty(answer.MaxScore, answer.RedeployRule, answer.DeploymentsBefore)
	effective := marked - penalty
	if effective < 0 {
		effective = 0
	}
	return Score{
		TeamCode:        answer.TeamCode,
		ProblemCode:     answer.ProblemCode,
		AnswerNumber:    answer.Number,
		MarkingResultID: marking.ID,
		MarkedScore:     marked,
		Penalty:         penalty,
		EffectiveScore:  effective,
		MaxScore:        answer.MaxScore,
		SubmittedAt:     answer.SubmittedAt,
		ContentCommit:   answer.ContentCommit,
	}
}

func LatestMarkings(markings []MarkingResult) map[string]MarkingResult {
	result := make(map[string]MarkingResult)
	for _, marking := range markings {
		key := answerKey(marking.TeamCode, marking.ProblemCode, marking.AnswerNumber)
		current, ok := result[key]
		if !ok || marking.CreatedAt.After(current.CreatedAt) ||
			(marking.CreatedAt.Equal(current.CreatedAt) && marking.ID > current.ID) {
			result[key] = marking
		}
	}
	return result
}

func SelectBestScores(answers []Answer, markings []MarkingResult) map[int64]map[string]Score {
	latest := LatestMarkings(markings)
	selected := make(map[int64]map[string]Score)
	for _, answer := range answers {
		marking, ok := latest[answerKey(answer.TeamCode, answer.ProblemCode, answer.Number)]
		if !ok {
			continue
		}
		score := EffectiveScore(answer, marking)
		if _, ok := selected[answer.TeamCode]; !ok {
			selected[answer.TeamCode] = make(map[string]Score)
		}
		current, ok := selected[answer.TeamCode][answer.ProblemCode]
		if !ok || score.EffectiveScore > current.EffectiveScore ||
			(score.EffectiveScore == current.EffectiveScore && score.SubmittedAt.Before(current.SubmittedAt)) ||
			(score.EffectiveScore == current.EffectiveScore && score.SubmittedAt.Equal(current.SubmittedAt) && score.AnswerNumber < current.AnswerNumber) {
			selected[answer.TeamCode][answer.ProblemCode] = score
		}
	}
	return selected
}

func BuildRanking(teams []Team, activeProblems map[string]struct{}, selected map[int64]map[string]Score) []RankingEntry {
	entries := make([]RankingEntry, 0, len(teams))
	for _, team := range teams {
		entry := RankingEntry{Team: team}
		for problemCode, score := range selected[team.Code] {
			if _, active := activeProblems[problemCode]; !active {
				continue
			}
			entry.Score += int64(score.EffectiveScore)
			if entry.LastEffectiveSubmissionAt == nil || score.SubmittedAt.After(*entry.LastEffectiveSubmissionAt) {
				t := score.SubmittedAt
				entry.LastEffectiveSubmissionAt = &t
			}
		}
		entries = append(entries, entry)
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		left, right := entries[i].LastEffectiveSubmissionAt, entries[j].LastEffectiveSubmissionAt
		if left == nil && right != nil {
			return false
		}
		if left != nil && right == nil {
			return true
		}
		if left != nil && right != nil && !left.Equal(*right) {
			return left.Before(*right)
		}
		return entries[i].Team.Code < entries[j].Team.Code
	})
	var previous *RankingEntry
	var denseRank int32
	for index := range entries {
		if previous == nil || !sameRank(*previous, entries[index]) {
			denseRank++
		}
		entries[index].Rank = denseRank
		previous = &entries[index]
	}
	return entries
}

func sameRank(left, right RankingEntry) bool {
	if left.Score != right.Score {
		return false
	}
	if left.LastEffectiveSubmissionAt == nil || right.LastEffectiveSubmissionAt == nil {
		return left.LastEffectiveSubmissionAt == nil && right.LastEffectiveSubmissionAt == nil
	}
	return left.LastEffectiveSubmissionAt.Equal(*right.LastEffectiveSubmissionAt)
}

func MarkingVisibilityAt(answerSubmittedAt, now time.Time, freezeAt, finalRevealedAt *time.Time) Visibility {
	if finalRevealedAt != nil {
		return VisibilityPublic
	}
	if now.Before(answerSubmittedAt.Add(PublishDelay)) {
		return VisibilityPrivate
	}
	if freezeAt != nil && !now.Before(*freezeAt) {
		return VisibilityTeam
	}
	return VisibilityPublic
}

func answerKey(teamCode int64, problemCode string, answerNumber int32) string {
	return strconv.FormatInt(teamCode, 10) + ":" + problemCode + ":" + strconv.FormatInt(int64(answerNumber), 10)
}
