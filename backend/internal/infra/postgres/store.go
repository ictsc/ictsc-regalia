package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

func Open(ctx context.Context, databaseURL string) (*Store, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse postgres configuration: %w", err)
	}
	config.ConnConfig.RuntimeParams["application_name"] = "ictsc-regalia-api"
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	store := &Store{pool: pool, now: func() time.Time { return time.Now().UTC() }}
	if err := store.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return store, nil
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool, now: func() time.Time { return time.Now().UTC() }}
}

func (s *Store) Close() { s.pool.Close() }

func (s *Store) Ping(ctx context.Context) error {
	if err := s.pool.Ping(ctx); err != nil {
		return core.WrapError(http.StatusServiceUnavailable, "database_unavailable", "PostgreSQL is unavailable", err)
	}
	return nil
}

func (s *Store) ListTeams(ctx context.Context) ([]core.Team, error) {
	rows, err := s.pool.Query(ctx, `SELECT code, name, organization, member_limit, color FROM teams ORDER BY code`)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()
	teams := make([]core.Team, 0)
	for rows.Next() {
		var team core.Team
		if err := rows.Scan(&team.Code, &team.Name, &team.Organization, &team.MemberLimit, &team.Color); err != nil {
			return nil, dbError(err)
		}
		teams = append(teams, team)
	}
	return teams, dbError(rows.Err())
}

func (s *Store) GetTeam(ctx context.Context, code int64) (core.Team, error) {
	var team core.Team
	err := s.pool.QueryRow(ctx, `SELECT code, name, organization, member_limit, color FROM teams WHERE code=$1`, code).
		Scan(&team.Code, &team.Name, &team.Organization, &team.MemberLimit, &team.Color)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.Team{}, core.NewError(http.StatusNotFound, "team_not_found", "Team was not found")
	}
	return team, dbError(err)
}

func (s *Store) CreateTeam(ctx context.Context, team core.Team) (core.Team, error) {
	err := s.pool.QueryRow(ctx, `
		INSERT INTO teams (code, name, organization, member_limit, color)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING code, name, organization, member_limit, color`, team.Code, team.Name, team.Organization, team.MemberLimit, core.TeamColor(team.Color)).
		Scan(&team.Code, &team.Name, &team.Organization, &team.MemberLimit, &team.Color)
	if err == nil {
		return team, nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		if pgErr.ConstraintName == "teams_pkey" {
			return core.Team{}, core.WrapError(http.StatusConflict, "team_code_conflict", "Team code already exists", err)
		}
		if pgErr.ConstraintName == "teams_name_key" {
			return core.Team{}, core.WrapError(http.StatusConflict, "conflict", "Team name already exists", err)
		}
	}
	return core.Team{}, dbError(err)
}

func (s *Store) UpdateTeam(ctx context.Context, code int64, patch core.TeamPatch) (core.Team, error) {
	return s.updateTeamTransactional(ctx, code, patch)
}

func (s *Store) DeleteTeam(ctx context.Context, code int64) error {
	result, err := s.pool.Exec(ctx, `DELETE FROM teams WHERE code=$1`, code)
	if err != nil {
		return conflictError(err, "team_not_empty", "Team has related contestants or competition data")
	}
	if result.RowsAffected() == 0 {
		return core.NewError(http.StatusNotFound, "team_not_found", "Team was not found")
	}
	return nil
}

func (s *Store) ListInvitations(ctx context.Context) ([]core.Invitation, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT code, team_code, created_at, expires_at FROM invitations
		WHERE consumed_at IS NULL ORDER BY created_at DESC, code`)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()
	invitations := make([]core.Invitation, 0)
	for rows.Next() {
		var invitation core.Invitation
		if err := rows.Scan(&invitation.Code, &invitation.TeamCode, &invitation.CreatedAt, &invitation.ExpiresAt); err != nil {
			return nil, dbError(err)
		}
		invitations = append(invitations, invitation)
	}
	return invitations, dbError(rows.Err())
}

func (s *Store) CreateInvitation(ctx context.Context, invitation core.Invitation) (core.Invitation, error) {
	err := s.pool.QueryRow(ctx, `
		INSERT INTO invitations(code,team_code,created_at,expires_at) VALUES($1,$2,$3,$4)
		RETURNING code,team_code,created_at,expires_at`, invitation.Code, invitation.TeamCode, invitation.CreatedAt, invitation.ExpiresAt).
		Scan(&invitation.Code, &invitation.TeamCode, &invitation.CreatedAt, &invitation.ExpiresAt)
	return invitation, conflictError(err, "invitation_conflict", "Invitation code already exists")
}

func (s *Store) ConsumeInvitation(ctx context.Context, code string, now time.Time, contestant core.Contestant) (core.Contestant, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return core.Contestant{}, dbError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var teamCode int64
	var expiresAt time.Time
	var consumedAt *time.Time
	err = tx.QueryRow(ctx, `
		SELECT team_code,expires_at,consumed_at FROM invitations
		WHERE code=$1 FOR UPDATE`, code).Scan(&teamCode, &expiresAt, &consumedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.Contestant{}, core.NewError(http.StatusConflict, "conflict", "Invitation code does not exist")
	}
	if err != nil {
		return core.Contestant{}, dbError(err)
	}
	if consumedAt != nil {
		return core.Contestant{}, core.NewError(http.StatusConflict, "invitation_already_used", "Invitation was already used")
	}
	if !now.Before(expiresAt) {
		return core.Contestant{}, core.NewError(http.StatusConflict, "invitation_expired", "Invitation has expired")
	}
	var memberLimit int32
	err = tx.QueryRow(ctx, `
		SELECT member_limit FROM teams WHERE code=$1 FOR UPDATE`, teamCode).Scan(&memberLimit)
	if err != nil {
		return core.Contestant{}, dbError(err)
	}
	var memberCount int32
	err = tx.QueryRow(ctx, `SELECT count(*)::int FROM contestants WHERE team_code=$1`, teamCode).Scan(&memberCount)
	if err != nil {
		return core.Contestant{}, dbError(err)
	}
	if memberCount >= memberLimit {
		return core.Contestant{}, core.NewError(http.StatusConflict, "team_full", "Team member limit has been reached")
	}
	contestant.TeamCode = teamCode
	_, err = tx.Exec(ctx, `
		INSERT INTO contestants(name,display_name,self_introduction,discord_id,team_code)
		VALUES($1,$2,$3,$4,$5)`, contestant.Name, contestant.DisplayName, contestant.SelfIntroduction, contestant.DiscordID, contestant.TeamCode)
	if err != nil {
		return core.Contestant{}, conflictError(err, "contestant_already_registered", "Contestant name or Discord account is already registered")
	}
	_, err = tx.Exec(ctx, `UPDATE invitations SET consumed_at=$2, consumed_by=$3 WHERE code=$1`, code, now, contestant.Name)
	if err != nil {
		return core.Contestant{}, dbError(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return core.Contestant{}, conflictError(err, "signup_conflict", "Concurrent signup could not be completed")
	}
	return contestant, nil
}

func (s *Store) ListContestants(ctx context.Context) ([]core.Contestant, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT name,display_name,self_introduction,discord_id,team_code
		FROM contestants ORDER BY team_code,name`)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()
	contestants := make([]core.Contestant, 0)
	for rows.Next() {
		var contestant core.Contestant
		if err := rows.Scan(&contestant.Name, &contestant.DisplayName, &contestant.SelfIntroduction, &contestant.DiscordID, &contestant.TeamCode); err != nil {
			return nil, dbError(err)
		}
		contestants = append(contestants, contestant)
	}
	return contestants, dbError(rows.Err())
}

func (s *Store) GetContestant(ctx context.Context, name string) (core.Contestant, error) {
	return s.getContestant(ctx, `WHERE name=$1`, name)
}

func (s *Store) GetContestantByDiscord(ctx context.Context, discordID string) (core.Contestant, error) {
	return s.getContestant(ctx, `WHERE discord_id=$1`, discordID)
}

func (s *Store) getContestant(ctx context.Context, predicate string, value any) (core.Contestant, error) {
	var contestant core.Contestant
	err := s.pool.QueryRow(ctx, `SELECT name,display_name,self_introduction,discord_id,team_code FROM contestants `+predicate, value).
		Scan(&contestant.Name, &contestant.DisplayName, &contestant.SelfIntroduction, &contestant.DiscordID, &contestant.TeamCode)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.Contestant{}, core.NewError(http.StatusNotFound, "contestant_not_found", "Contestant was not found")
	}
	return contestant, dbError(err)
}

func (s *Store) UpdateContestant(ctx context.Context, name string, patch core.ContestantPatch) (core.Contestant, error) {
	var contestant core.Contestant
	err := s.pool.QueryRow(ctx, `
		UPDATE contestants SET display_name=COALESCE($2,display_name),
			self_introduction=COALESCE($3,self_introduction),updated_at=now()
		WHERE name=$1 RETURNING name,display_name,self_introduction,discord_id,team_code`, name, patch.DisplayName, patch.SelfIntroduction).
		Scan(&contestant.Name, &contestant.DisplayName, &contestant.SelfIntroduction, &contestant.DiscordID, &contestant.TeamCode)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.Contestant{}, core.NewError(http.StatusNotFound, "contestant_not_found", "Contestant was not found")
	}
	return contestant, dbError(err)
}

func (s *Store) ActiveContent(ctx context.Context) (core.ContentSnapshot, error) {
	return s.readContent(ctx, http.StatusBadGateway, "content_not_available", "No valid content snapshot is available", `WHERE active=true`)
}

func (s *Store) GetContent(ctx context.Context, commit string) (core.ContentSnapshot, error) {
	return s.readContent(ctx, http.StatusNotFound, "content_not_found", "Content snapshot was not found", `WHERE commit_sha=$1`, commit)
}

func (s *Store) readContent(ctx context.Context, notFoundStatus int, notFoundCode, notFoundMessage, predicate string, args ...any) (core.ContentSnapshot, error) {
	var snapshot core.ContentSnapshot
	var manifest []byte
	err := s.pool.QueryRow(ctx, `
		SELECT commit_sha,repository,git_ref,manifest,fetched_at,activated_at
		FROM content_snapshots `+predicate, args...).Scan(
		&snapshot.CommitSHA, &snapshot.Repository, &snapshot.Ref, &manifest, &snapshot.FetchedAt, &snapshot.ActivatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.ContentSnapshot{}, core.NewError(notFoundStatus, notFoundCode, notFoundMessage)
	}
	if err != nil {
		return core.ContentSnapshot{}, dbError(err)
	}
	if err := json.Unmarshal(manifest, &snapshot.Manifest); err != nil {
		return core.ContentSnapshot{}, core.WrapError(http.StatusInternalServerError, "internal_error", "Stored content snapshot is invalid", err)
	}
	return snapshot, nil
}

func (s *Store) ActivateContent(ctx context.Context, snapshot core.ContentSnapshot, expectedCurrentCommit string) (core.ContentSnapshot, error) {
	manifest, err := json.Marshal(snapshot.Manifest)
	if err != nil {
		return core.ContentSnapshot{}, core.WrapError(http.StatusUnprocessableEntity, "content_invalid", "Manifest cannot be encoded", err)
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return core.ContentSnapshot{}, dbError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var current string
	hasCurrent := true
	err = tx.QueryRow(ctx, `SELECT commit_sha FROM content_snapshots WHERE active=true FOR UPDATE`).Scan(&current)
	if errors.Is(err, pgx.ErrNoRows) {
		hasCurrent = false
	} else if err != nil {
		return core.ContentSnapshot{}, dbError(err)
	}
	if hasCurrent && current == snapshot.CommitSHA {
		if err := tx.Rollback(ctx); err != nil {
			return core.ContentSnapshot{}, dbError(err)
		}
		return s.GetContent(ctx, snapshot.CommitSHA)
	}
	if expectedCurrentCommit == core.NoActiveContentCommit && hasCurrent {
		return core.ContentSnapshot{}, core.NewError(http.StatusConflict, "content_refresh_in_progress", "Active content appeared while refresh was in progress")
	}
	if expectedCurrentCommit != "" && expectedCurrentCommit != core.NoActiveContentCommit && (!hasCurrent || current != expectedCurrentCommit) {
		return core.ContentSnapshot{}, core.NewError(http.StatusConflict, "content_refresh_in_progress", "Active content changed while refresh was in progress")
	}
	activatedAt := s.now()
	_, err = tx.Exec(ctx, `UPDATE content_snapshots SET active=false WHERE active=true`)
	if err != nil {
		return core.ContentSnapshot{}, dbError(err)
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO content_snapshots(commit_sha,repository,git_ref,manifest,fetched_at,activated_at,active)
		VALUES($1,$2,$3,$4,$5,$6,true)
		ON CONFLICT(commit_sha) DO UPDATE SET repository=excluded.repository,git_ref=excluded.git_ref,
			manifest=excluded.manifest,fetched_at=excluded.fetched_at,activated_at=excluded.activated_at,active=true`,
		snapshot.CommitSHA, snapshot.Repository, snapshot.Ref, manifest, snapshot.FetchedAt, activatedAt)
	if err != nil {
		return core.ContentSnapshot{}, dbError(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return core.ContentSnapshot{}, conflictError(err, "content_refresh_in_progress", "Content activation conflicted with another refresh")
	}
	snapshot.ActivatedAt = activatedAt
	return snapshot, nil
}

func (s *Store) ListAnswers(ctx context.Context, filter core.AnswerFilter) ([]core.Answer, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT team_code,problem_code,number,author_name,body,submitted_at,content_commit,max_score,redeploy_rule,deployments_before
		FROM answers WHERE ($1::bigint IS NULL OR team_code=$1) AND ($2::text IS NULL OR problem_code=$2)
			AND ($3::timestamptz IS NULL OR submitted_at<=$3)
		ORDER BY submitted_at DESC,team_code,problem_code,number DESC`, filter.TeamCode, filter.ProblemCode, filter.Before)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()
	answers := make([]core.Answer, 0)
	for rows.Next() {
		answer, err := scanAnswer(rows)
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}
	return answers, dbError(rows.Err())
}

func (s *Store) GetAnswer(ctx context.Context, teamCode int64, problemCode string, number int32) (core.Answer, error) {
	answer, err := scanAnswer(s.pool.QueryRow(ctx, `
		SELECT team_code,problem_code,number,author_name,body,submitted_at,content_commit,max_score,redeploy_rule,deployments_before
		FROM answers WHERE team_code=$1 AND problem_code=$2 AND number=$3`, teamCode, problemCode, number))
	if errors.Is(err, pgx.ErrNoRows) {
		return core.Answer{}, core.NewError(http.StatusNotFound, "answer_not_found", "Answer was not found")
	}
	return answer, dbError(err)
}

type scanner interface{ Scan(...any) error }

func scanAnswer(row scanner) (core.Answer, error) {
	var answer core.Answer
	var rule []byte
	err := row.Scan(&answer.TeamCode, &answer.ProblemCode, &answer.Number, &answer.AuthorName, &answer.Body,
		&answer.SubmittedAt, &answer.ContentCommit, &answer.MaxScore, &rule, &answer.DeploymentsBefore)
	if err != nil {
		return core.Answer{}, err
	}
	if err := json.Unmarshal(rule, &answer.RedeployRule); err != nil {
		return core.Answer{}, core.WrapError(http.StatusInternalServerError, "answer_snapshot_corrupt", "Stored answer metadata is invalid", err)
	}
	return answer, nil
}

func (s *Store) SubmitAnswer(ctx context.Context, answer core.Answer, interval time.Duration) (core.Answer, time.Duration, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return core.Answer{}, 0, dbError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = tx.Exec(ctx, `
		INSERT INTO answer_counters(team_code,problem_code,next_number,last_submitted_at)
		VALUES($1,$2,1,NULL) ON CONFLICT DO NOTHING`, answer.TeamCode, answer.ProblemCode)
	if err != nil {
		return core.Answer{}, 0, dbError(err)
	}
	var next int32
	var last *time.Time
	err = tx.QueryRow(ctx, `
		SELECT next_number,last_submitted_at FROM answer_counters
		WHERE team_code=$1 AND problem_code=$2 FOR UPDATE`, answer.TeamCode, answer.ProblemCode).Scan(&next, &last)
	if err != nil {
		return core.Answer{}, 0, dbError(err)
	}
	if last != nil {
		nextAllowed := last.Add(interval)
		if answer.SubmittedAt.Before(nextAllowed) {
			retry := nextAllowed.Sub(answer.SubmittedAt)
			return core.Answer{}, retry, &core.Error{Status: http.StatusTooManyRequests, Code: "answer_rate_limited", Message: "Answer submission interval has not elapsed", RetryAfter: retry}
		}
	}
	rule, err := json.Marshal(answer.RedeployRule)
	if err != nil {
		return core.Answer{}, 0, core.WrapError(http.StatusUnprocessableEntity, "answer_metadata_invalid", "Answer snapshot metadata is invalid", err)
	}
	answer.Number = next
	_, err = tx.Exec(ctx, `
		INSERT INTO answers(team_code,problem_code,number,author_name,body,submitted_at,content_commit,max_score,redeploy_rule,deployments_before)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, answer.TeamCode, answer.ProblemCode, answer.Number,
		answer.AuthorName, answer.Body, answer.SubmittedAt, answer.ContentCommit, answer.MaxScore, rule, answer.DeploymentsBefore)
	if err != nil {
		return core.Answer{}, 0, conflictError(err, "answer_conflict", "Answer could not be submitted concurrently")
	}
	_, err = tx.Exec(ctx, `
		UPDATE answer_counters SET next_number=$3,last_submitted_at=$4
		WHERE team_code=$1 AND problem_code=$2`, answer.TeamCode, answer.ProblemCode, answer.Number+1, answer.SubmittedAt)
	if err != nil {
		return core.Answer{}, 0, dbError(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return core.Answer{}, 0, conflictError(err, "answer_conflict", "Answer could not be submitted concurrently")
	}
	return answer, 0, nil
}

func (s *Store) ListMarkingResults(ctx context.Context) ([]core.MarkingResult, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id::text,team_code,problem_code,answer_number,judge,marked_score,rationale,created_at,visibility
		FROM marking_results ORDER BY created_at DESC,id ASC`)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()
	results := make([]core.MarkingResult, 0)
	for rows.Next() {
		var result core.MarkingResult
		if err := rows.Scan(&result.ID, &result.TeamCode, &result.ProblemCode, &result.AnswerNumber, &result.Judge,
			&result.MarkedScore, &result.Rationale, &result.CreatedAt, &result.Visibility); err != nil {
			return nil, dbError(err)
		}
		results = append(results, result)
	}
	return results, dbError(rows.Err())
}

func (s *Store) CreateMarkingResult(ctx context.Context, result core.MarkingResult) (core.MarkingResult, error) {
	err := s.pool.QueryRow(ctx, `
		INSERT INTO marking_results(id,team_code,problem_code,answer_number,judge,marked_score,rationale,created_at,visibility)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id::text,team_code,problem_code,answer_number,judge,marked_score,rationale,created_at,visibility`,
		result.ID, result.TeamCode, result.ProblemCode, result.AnswerNumber, result.Judge, result.MarkedScore,
		result.Rationale, result.CreatedAt, result.Visibility).Scan(&result.ID, &result.TeamCode, &result.ProblemCode,
		&result.AnswerNumber, &result.Judge, &result.MarkedScore, &result.Rationale, &result.CreatedAt, &result.Visibility)
	return result, conflictError(err, "marking_conflict", "Marking result could not be saved")
}

func (s *Store) ListDeployments(ctx context.Context, filter core.DeploymentFilter) ([]core.Deployment, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT request_id::text,team_code,problem_code,revision,content_commit,requested_at,latest_status
		FROM deployments WHERE ($1::bigint IS NULL OR team_code=$1) AND ($2::text IS NULL OR problem_code=$2)
		ORDER BY requested_at DESC,team_code,problem_code,revision DESC`, filter.TeamCode, filter.ProblemCode)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()
	deployments := make([]core.Deployment, 0)
	for rows.Next() {
		var deployment core.Deployment
		if err := rows.Scan(&deployment.RequestID, &deployment.TeamCode, &deployment.ProblemCode, &deployment.Revision,
			&deployment.ContentCommit, &deployment.RequestedAt, &deployment.LatestStatus); err != nil {
			return nil, dbError(err)
		}
		deployments = append(deployments, deployment)
	}
	if err := rows.Err(); err != nil {
		return nil, dbError(err)
	}
	for index := range deployments {
		events, err := s.deploymentEvents(ctx, deployments[index].RequestID)
		if err != nil {
			return nil, err
		}
		deployments[index].Events = events
	}
	return deployments, nil
}

func (s *Store) CreateDeployment(ctx context.Context, deployment core.Deployment) (core.Deployment, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return core.Deployment{}, dbError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = tx.Exec(ctx, `
		INSERT INTO deployment_counters(team_code,problem_code,next_revision)
		VALUES($1,$2,1) ON CONFLICT DO NOTHING`, deployment.TeamCode, deployment.ProblemCode)
	if err != nil {
		return core.Deployment{}, dbError(err)
	}
	err = tx.QueryRow(ctx, `SELECT next_revision FROM deployment_counters
		WHERE team_code=$1 AND problem_code=$2 FOR UPDATE`, deployment.TeamCode, deployment.ProblemCode).Scan(&deployment.Revision)
	if err != nil {
		return core.Deployment{}, dbError(err)
	}
	var active bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(
		SELECT 1 FROM deployments WHERE team_code=$1 AND problem_code=$2
		AND latest_status IN ('QUEUED','DEPLOYING'))`, deployment.TeamCode, deployment.ProblemCode).Scan(&active)
	if err != nil {
		return core.Deployment{}, dbError(err)
	}
	if active {
		return core.Deployment{}, core.NewError(http.StatusConflict, "deployment_in_progress", "A deployment is already in progress")
	}
	_, err = tx.Exec(ctx, `UPDATE deployment_counters SET next_revision=$3 WHERE team_code=$1 AND problem_code=$2`, deployment.TeamCode, deployment.ProblemCode, deployment.Revision+1)
	if err != nil {
		return core.Deployment{}, dbError(err)
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO deployments(request_id,team_code,problem_code,revision,content_commit,requested_at,latest_status)
		VALUES($1,$2,$3,$4,$5,$6,$7)`, deployment.RequestID, deployment.TeamCode, deployment.ProblemCode,
		deployment.Revision, deployment.ContentCommit, deployment.RequestedAt, deployment.LatestStatus)
	if err != nil {
		return core.Deployment{}, conflictError(err, "deployment_in_progress", "A deployment is already in progress")
	}
	for index := range deployment.Events {
		event := &deployment.Events[index]
		event.OccurredAt = event.OccurredAt.UTC().Truncate(time.Microsecond)
		if event.Status != core.DeploymentQueued {
			return core.Deployment{}, core.NewError(http.StatusInternalServerError, "internal_error", "Initial deployment event must be QUEUED")
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO deployment_events(event_id,request_id,occurred_at,status,message)
			VALUES($1,$2,$3,$4,$5)`, event.EventID, deployment.RequestID, event.OccurredAt, event.Status, event.Message)
		if err != nil {
			return core.Deployment{}, conflictError(err, "deployment_event_conflict", "Initial deployment event could not be saved")
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return core.Deployment{}, conflictError(err, "deployment_in_progress", "Deployment could not be queued concurrently")
	}
	return deployment, nil
}

func (s *Store) AppendDeploymentEvent(ctx context.Context, requestID string, event core.DeploymentEvent, recovery bool) (core.Deployment, bool, error) {
	event.OccurredAt = event.OccurredAt.UTC().Truncate(time.Microsecond)
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return core.Deployment{}, false, dbError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var deployment core.Deployment
	err = tx.QueryRow(ctx, `
		SELECT request_id::text,team_code,problem_code,revision,content_commit,requested_at,latest_status
		FROM deployments WHERE request_id=$1 FOR UPDATE`, requestID).Scan(&deployment.RequestID, &deployment.TeamCode,
		&deployment.ProblemCode, &deployment.Revision, &deployment.ContentCommit, &deployment.RequestedAt, &deployment.LatestStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.Deployment{}, false, core.NewError(http.StatusNotFound, "deployment_not_found", "Deployment was not found")
	}
	if err != nil {
		return core.Deployment{}, false, dbError(err)
	}
	var existingRequestID string
	existing := core.DeploymentEvent{EventID: event.EventID}
	err = tx.QueryRow(ctx, `
		SELECT request_id::text,occurred_at,status,message
		FROM deployment_events WHERE event_id=$1`, event.EventID).Scan(
		&existingRequestID, &existing.OccurredAt, &existing.Status, &existing.Message)
	if err == nil {
		if existingRequestID != requestID || !sameDeploymentEvent(existing, event) {
			return core.Deployment{}, false, core.NewError(http.StatusConflict, "duplicate_event_mismatch", "Event ID was already used with a different payload")
		}
		events, listErr := deploymentEventsFrom(ctx, tx, requestID)
		if listErr != nil {
			return core.Deployment{}, false, listErr
		}
		deployment.Events = events
		return deployment, true, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return core.Deployment{}, false, dbError(err)
	}
	var latestOccurredAt time.Time
	var latestEventID string
	err = tx.QueryRow(ctx, `SELECT occurred_at,event_id FROM deployment_events
		WHERE request_id=$1 ORDER BY occurred_at DESC,event_id DESC LIMIT 1`, requestID).Scan(&latestOccurredAt, &latestEventID)
	if err == nil && (event.OccurredAt.Before(latestOccurredAt) || event.OccurredAt.Equal(latestOccurredAt) && event.EventID < latestEventID) {
		return core.Deployment{}, false, core.NewError(http.StatusConflict, "invalid_deployment_transition", "Deployment event occurred before the latest persisted event")
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return core.Deployment{}, false, dbError(err)
	}
	if !core.CanTransitionDeployment(deployment.LatestStatus, event.Status, recovery) {
		return core.Deployment{}, false, core.NewError(http.StatusConflict, "invalid_deployment_transition", "Deployment status transition is not allowed")
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO deployment_events(event_id,request_id,occurred_at,status,message)
		VALUES($1,$2,$3,$4,$5)`, event.EventID, requestID, event.OccurredAt, event.Status, event.Message)
	if err != nil {
		return core.Deployment{}, false, conflictError(err, "deployment_event_conflict", "Deployment event already exists")
	}
	_, err = tx.Exec(ctx, `UPDATE deployments SET latest_status=$2,updated_at=now() WHERE request_id=$1`, requestID, event.Status)
	if err != nil {
		return core.Deployment{}, false, dbError(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return core.Deployment{}, false, conflictError(err, "deployment_event_conflict", "Deployment event could not be applied concurrently")
	}
	deployment.LatestStatus = event.Status
	deployment.Events, err = s.deploymentEvents(ctx, requestID)
	return deployment, false, err
}

func (s *Store) deploymentEvents(ctx context.Context, requestID string) ([]core.DeploymentEvent, error) {
	return deploymentEventsFrom(ctx, s.pool, requestID)
}

func (s *Store) GetCompetitionState(ctx context.Context) (core.CompetitionState, error) {
	var state core.CompetitionState
	err := s.pool.QueryRow(ctx, `
		SELECT rule_markdown,ranking_freeze_at,final_revealed_at,updated_at,updated_by
		FROM competition_state WHERE singleton=true`).Scan(&state.RuleMarkdown, &state.RankingFreezeAt,
		&state.FinalRevealedAt, &state.UpdatedAt, &state.UpdatedBy)
	return state, dbError(err)
}

func (s *Store) ReplaceRule(ctx context.Context, markdown, actor string) (core.CompetitionState, error) {
	return s.mutateCompetitionState(ctx, actor, "replace_rule", &markdown, true, nil, false, nil)
}

func (s *Store) ReplaceFreezeAt(ctx context.Context, freezeAt *time.Time, actor string) (core.CompetitionState, error) {
	if freezeAt != nil {
		value := freezeAt.UTC().Truncate(time.Microsecond)
		freezeAt = &value
	}
	return s.mutateCompetitionState(ctx, actor, "replace_freeze", nil, false, freezeAt, true, nil)
}

func (s *Store) RevealFinal(ctx context.Context, at time.Time, actor string) (core.CompetitionState, error) {
	return s.mutateCompetitionState(ctx, actor, "reveal_final", nil, false, nil, false, &at)
}

func dbError(err error) error {
	if err == nil {
		return nil
	}
	return core.WrapError(http.StatusInternalServerError, "database_error", "Database operation failed", err)
}

func conflictError(err error, code, message string) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505" || pgErr.Code == "23503" || pgErr.Code == "23514" || pgErr.Code == "40001") {
		return core.WrapError(http.StatusConflict, code, message, err)
	}
	return dbError(err)
}

var _ core.Store = (*Store)(nil)

// RegisterContestant serializes membership capacity checks on the team row.
func (s *Store) RegisterContestant(ctx context.Context, contestant core.Contestant) (core.Contestant, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return core.Contestant{}, dbError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	teamCode := contestant.TeamCode
	var memberLimit int32
	err = tx.QueryRow(ctx, `
		SELECT member_limit FROM teams WHERE code=$1 FOR UPDATE`, teamCode).Scan(&memberLimit)
	if err != nil {
		return core.Contestant{}, dbError(err)
	}
	var memberCount int32
	err = tx.QueryRow(ctx, `SELECT count(*)::int FROM contestants WHERE team_code=$1`, teamCode).Scan(&memberCount)
	if err != nil {
		return core.Contestant{}, dbError(err)
	}
	if memberCount >= memberLimit {
		return core.Contestant{}, core.NewError(http.StatusConflict, "team_full", "Team member limit has been reached")
	}
	contestant.TeamCode = teamCode
	_, err = tx.Exec(ctx, `
		INSERT INTO contestants(name,display_name,self_introduction,discord_id,team_code)
		VALUES($1,$2,$3,$4,$5)`, contestant.Name, contestant.DisplayName, contestant.SelfIntroduction, contestant.DiscordID, contestant.TeamCode)
	if err != nil {
		return core.Contestant{}, conflictError(err, "contestant_already_registered", "Contestant name or Discord account is already registered")
	}

	if err := tx.Commit(ctx); err != nil {
		return core.Contestant{}, dbError(err)
	}
	return contestant, nil
}
