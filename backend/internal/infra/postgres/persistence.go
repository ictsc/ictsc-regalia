package postgres

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/jackc/pgx/v5"
)

func (s *Store) ContentAt(ctx context.Context, at time.Time) (core.ContentSnapshot, error) {
	return s.readContent(ctx, http.StatusBadGateway, "content_not_available", "No valid content snapshot existed at the requested time",
		`WHERE activated_at <= $1 ORDER BY activated_at DESC LIMIT 1`, at)
}

func (s *Store) RecalculateScores(ctx context.Context, selections map[int64]map[string]core.Score, actor string) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return dbError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err = tx.Exec(ctx, `LOCK TABLE score_selections IN EXCLUSIVE MODE`); err != nil {
		return dbError(err)
	}
	if _, err = tx.Exec(ctx, `DELETE FROM score_selections`); err != nil {
		return dbError(err)
	}
	teams := make([]int64, 0, len(selections))
	for team := range selections {
		teams = append(teams, team)
	}
	sort.Slice(teams, func(i, j int) bool { return teams[i] < teams[j] })
	calculatedAt := s.now()
	for _, team := range teams {
		problems := make([]string, 0, len(selections[team]))
		for problem := range selections[team] {
			problems = append(problems, problem)
		}
		sort.Strings(problems)
		for _, problem := range problems {
			score := selections[team][problem]
			if score.MarkingResultID == "" {
				return core.NewError(http.StatusInternalServerError, "internal_error", "Selected score has no marking result")
			}
			_, err = tx.Exec(ctx, `
				INSERT INTO score_selections(team_code,problem_code,answer_number,marking_result_id,marked_score,penalty,effective_score,max_score,content_commit,calculated_at)
				VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, team, problem, score.AnswerNumber, score.MarkingResultID,
				score.MarkedScore, score.Penalty, score.EffectiveScore, score.MaxScore, score.ContentCommit, calculatedAt)
			if err != nil {
				return dbError(err)
			}
		}
	}
	after, err := json.Marshal(selections)
	if err != nil {
		return core.WrapError(http.StatusInternalServerError, "internal_error", "Could not encode score audit", err)
	}
	if _, err = tx.Exec(ctx, `INSERT INTO audit_log(actor,action,resource,after_json) VALUES($1,'recalculate_scores','scores',$2)`, actor, after); err != nil {
		return dbError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return conflictError(err, "conflict", "Score recalculation conflicted with another operation")
	}
	return nil
}

func (s *Store) GetRankingSnapshot(ctx context.Context, frozenAt time.Time) (core.RankingSnapshot, bool, error) {
	var snapshot core.RankingSnapshot
	var encoded []byte
	err := s.pool.QueryRow(ctx, `
		SELECT frozen_at,content_commit,ranking FROM ranking_snapshots
		WHERE frozen_at=$1 ORDER BY created_at DESC LIMIT 1`, frozenAt).Scan(&snapshot.FrozenAt, &snapshot.ContentCommit, &encoded)
	if err == pgx.ErrNoRows {
		return core.RankingSnapshot{}, false, nil
	}
	if err != nil {
		return core.RankingSnapshot{}, false, dbError(err)
	}
	if err := json.Unmarshal(encoded, &snapshot.Entries); err != nil {
		return core.RankingSnapshot{}, false, core.WrapError(http.StatusInternalServerError, "internal_error", "Stored ranking snapshot is invalid", err)
	}
	return snapshot, true, nil
}

func (s *Store) PutRankingSnapshot(ctx context.Context, snapshot core.RankingSnapshot, actor string) (core.RankingSnapshot, error) {
	encoded, err := json.Marshal(snapshot.Entries)
	if err != nil {
		return core.RankingSnapshot{}, core.WrapError(http.StatusInternalServerError, "internal_error", "Could not encode ranking snapshot", err)
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO ranking_snapshots(frozen_at,content_commit,ranking,created_by)
		VALUES($1,$2,$3,$4) ON CONFLICT(frozen_at,content_commit) DO NOTHING`, snapshot.FrozenAt, snapshot.ContentCommit, encoded, actor)
	if err != nil {
		return core.RankingSnapshot{}, dbError(err)
	}
	stored, ok, err := s.GetRankingSnapshot(ctx, snapshot.FrozenAt)
	if err != nil {
		return core.RankingSnapshot{}, err
	}
	if !ok {
		return core.RankingSnapshot{}, core.NewError(http.StatusInternalServerError, "internal_error", "Ranking snapshot was not persisted")
	}
	return stored, nil
}

func (s *Store) mutateCompetitionState(ctx context.Context, actor, action string, rule *string, setRule bool, freezeAt *time.Time, setFreeze bool, revealAt *time.Time) (core.CompetitionState, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return core.CompetitionState{}, dbError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var before core.CompetitionState
	err = tx.QueryRow(ctx, `
		SELECT rule_markdown,ranking_freeze_at,final_revealed_at,updated_at,updated_by
		FROM competition_state WHERE singleton=true FOR UPDATE`).Scan(&before.RuleMarkdown, &before.RankingFreezeAt,
		&before.FinalRevealedAt, &before.UpdatedAt, &before.UpdatedBy)
	if err != nil {
		return core.CompetitionState{}, dbError(err)
	}
	after := before
	changed := false
	if setRule && rule != nil && after.RuleMarkdown != *rule {
		after.RuleMarkdown, changed = *rule, true
	}
	if setFreeze && !sameOptionalTime(after.RankingFreezeAt, freezeAt) {
		after.RankingFreezeAt, changed = freezeAt, true
	}
	if revealAt != nil && after.FinalRevealedAt == nil {
		value := revealAt.UTC()
		after.FinalRevealedAt, changed = &value, true
	}
	if !changed {
		return before, nil
	}
	after.UpdatedAt, after.UpdatedBy = s.now(), actor
	_, err = tx.Exec(ctx, `
		UPDATE competition_state SET rule_markdown=$1,ranking_freeze_at=$2,final_revealed_at=$3,updated_at=$4,updated_by=$5
		WHERE singleton=true`, after.RuleMarkdown, after.RankingFreezeAt, after.FinalRevealedAt, after.UpdatedAt, after.UpdatedBy)
	if err != nil {
		return core.CompetitionState{}, dbError(err)
	}
	if setFreeze {
		if _, err = tx.Exec(ctx, `DELETE FROM ranking_snapshots`); err != nil {
			return core.CompetitionState{}, dbError(err)
		}
	}
	beforeJSON, _ := json.Marshal(before)
	afterJSON, _ := json.Marshal(after)
	if _, err = tx.Exec(ctx, `
		INSERT INTO audit_log(actor,action,resource,before_json,after_json)
		VALUES($1,$2,'competition_state',$3,$4)`, actor, action, beforeJSON, afterJSON); err != nil {
		return core.CompetitionState{}, dbError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return core.CompetitionState{}, conflictError(err, "conflict", "Competition state update conflicted with another operation")
	}
	return after, nil
}

func sameOptionalTime(left, right *time.Time) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Equal(*right)
}
