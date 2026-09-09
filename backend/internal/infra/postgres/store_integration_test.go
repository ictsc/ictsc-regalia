package postgres

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const postgresTestImage = "postgres:17-alpine"

type postgresFixture struct {
	pool  *pgxpool.Pool
	store *Store
}

func TestPostgresStoreIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("PostgreSQL integration tests are disabled by -short")
	}

	fixture := newPostgresFixture(t)

	t.Run("applies the canonical migration", func(t *testing.T) {
		tables := []string{
			"teams",
			"contestants",
			"invitations",
			"content_snapshots",
			"answers",
			"marking_results",
			"score_selections",
			"deployments",
			"deployment_events",
			"competition_state",
			"ranking_snapshots",
			"audit_log",
		}
		for _, table := range tables {
			var relation *string
			if err := fixture.pool.QueryRow(t.Context(), `SELECT to_regclass($1)::text`, "public."+table).Scan(&relation); err != nil {
				t.Fatalf("look up table %s: %v", table, err)
			}
			if relation == nil || *relation != table {
				t.Errorf("table %s was not created (got %v)", table, relation)
			}
		}

		var stateRows int
		if err := fixture.pool.QueryRow(t.Context(), `SELECT count(*) FROM competition_state WHERE singleton`).Scan(&stateRows); err != nil {
			t.Fatalf("read initial competition state: %v", err)
		}
		if stateRows != 1 {
			t.Fatalf("competition_state rows = %d, want 1", stateRows)
		}
	})

	t.Run("persists team colors and preserves omitted patches", func(t *testing.T) {
		ctx := t.Context()
		team, err := fixture.store.CreateTeam(ctx, core.Team{Code: 98, Name: "color-team", Organization: "ICTSC", MemberLimit: 4})
		if err != nil {
			t.Fatal(err)
		}
		if team.Color != core.DefaultTeamColor {
			t.Fatalf("default color = %q", team.Color)
		}
		color := "#0083C3"
		if _, err := fixture.store.UpdateTeam(ctx, team.Code, core.TeamPatch{Color: &color}); err != nil {
			t.Fatal(err)
		}
		name := "color-team-renamed"
		if _, err := fixture.store.UpdateTeam(ctx, team.Code, core.TeamPatch{Name: &name}); err != nil {
			t.Fatal(err)
		}
		read, err := fixture.store.GetTeam(ctx, team.Code)
		if err != nil {
			t.Fatal(err)
		}
		if read.Color != color {
			t.Fatalf("persisted color = %q", read.Color)
		}
		if _, err := fixture.pool.Exec(ctx, "UPDATE teams SET color='#FFFFFF' WHERE code=98"); err == nil {
			t.Fatal("invalid color accepted")
		}
	})

	t.Run("serializes concurrent signup and enforces team capacity", func(t *testing.T) {
		ctx := t.Context()
		team := core.Team{Code: 2, Name: "capacity-team", Organization: "ICTSC", MemberLimit: 1}
		if _, err := fixture.store.CreateTeam(ctx, team); err != nil {
			t.Fatalf("create team: %v", err)
		}

		now := time.Now().UTC().Truncate(time.Microsecond)
		invitations := []core.Invitation{
			{Code: "capacity-invite-a", TeamCode: team.Code, CreatedAt: now, ExpiresAt: now.Add(time.Hour)},
			{Code: "capacity-invite-b", TeamCode: team.Code, CreatedAt: now, ExpiresAt: now.Add(time.Hour)},
		}
		for _, invitation := range invitations {
			if _, err := fixture.store.CreateInvitation(ctx, invitation); err != nil {
				t.Fatalf("create invitation %s: %v", invitation.Code, err)
			}
		}

		type signupResult struct {
			contestant core.Contestant
			err        error
		}
		start := make(chan struct{})
		results := make(chan signupResult, len(invitations))
		for index, invitation := range invitations {
			index, invitation := index, invitation
			go func() {
				<-start
				contestant, err := fixture.store.ConsumeInvitation(ctx, invitation.Code, now.Add(time.Minute), core.Contestant{
					Name:        fmt.Sprintf("entrant-%d", index),
					DisplayName: fmt.Sprintf("Entrant %d", index),
					DiscordID:   fmt.Sprintf("10000000000000000%d", index),
				})
				results <- signupResult{contestant: contestant, err: err}
			}()
		}
		close(start)

		successes := 0
		conflicts := 0
		for range invitations {
			result := <-results
			if result.err == nil {
				successes++
				if result.contestant.TeamCode != team.Code {
					t.Errorf("signup team = %d, want %d", result.contestant.TeamCode, team.Code)
				}
				continue
			}
			var domainErr *core.Error
			if !errors.As(result.err, &domainErr) || domainErr.Status != http.StatusConflict {
				t.Errorf("losing signup error = %v, want a 409 domain error", result.err)
				continue
			}
			conflicts++
		}
		if successes != 1 || conflicts != 1 {
			t.Fatalf("concurrent signup results: successes=%d conflicts=%d, want 1 and 1", successes, conflicts)
		}

		var members, consumed int
		if err := fixture.pool.QueryRow(ctx, `SELECT count(*) FROM contestants WHERE team_code=$1`, team.Code).Scan(&members); err != nil {
			t.Fatalf("count team members: %v", err)
		}
		if err := fixture.pool.QueryRow(ctx, `SELECT count(*) FROM invitations WHERE team_code=$1 AND consumed_at IS NOT NULL`, team.Code).Scan(&consumed); err != nil {
			t.Fatalf("count consumed invitations: %v", err)
		}
		if members != 1 || consumed != 1 {
			t.Fatalf("persisted members=%d consumed invitations=%d, want 1 and 1", members, consumed)
		}
	})

	t.Run("allocates answer numbers and enforces the twenty minute interval", func(t *testing.T) {
		ctx := t.Context()
		const (
			teamCode   = int64(3)
			contestant = "answerer"
			problem    = "Q01"
			commit     = "1111111111111111111111111111111111111111"
		)
		seedTeamContestantAndContent(t, fixture, teamCode, contestant, "100000000000000010", commit)

		base := time.Date(2026, time.August, 30, 1, 2, 3, 0, time.UTC)
		answer := core.Answer{
			TeamCode:      teamCode,
			ProblemCode:   problem,
			AuthorName:    contestant,
			Body:          "first answer",
			SubmittedAt:   base,
			ContentCommit: commit,
			MaxScore:      100,
			RedeployRule:  core.RedeployRule{Type: core.RedeployUnredeployable},
		}

		first, retryAfter, err := fixture.store.SubmitAnswer(ctx, answer, core.AnswerInterval)
		if err != nil {
			t.Fatalf("submit first answer: %v", err)
		}
		if first.Number != 1 || retryAfter != 0 {
			t.Fatalf("first answer number=%d retry=%s, want 1 and zero", first.Number, retryAfter)
		}

		tooSoon := answer
		tooSoon.Body = "too soon"
		tooSoon.SubmittedAt = base.Add(10 * time.Minute)
		_, retryAfter, err = fixture.store.SubmitAnswer(ctx, tooSoon, core.AnswerInterval)
		var rateLimit *core.Error
		if !errors.As(err, &rateLimit) || rateLimit.Status != http.StatusTooManyRequests || rateLimit.Code != "answer_rate_limited" {
			t.Fatalf("early submission error = %v, want answer_rate_limited 429", err)
		}
		if retryAfter != 10*time.Minute || rateLimit.RetryAfter != retryAfter {
			t.Fatalf("retry after = %s (error %s), want 10m", retryAfter, rateLimit.RetryAfter)
		}

		second := answer
		second.Body = "second answer"
		second.SubmittedAt = base.Add(core.AnswerInterval)
		second, _, err = fixture.store.SubmitAnswer(ctx, second, core.AnswerInterval)
		if err != nil {
			t.Fatalf("submit answer at interval boundary: %v", err)
		}
		if second.Number != 2 {
			t.Fatalf("second answer number = %d, want 2", second.Number)
		}

		var count, maxNumber int
		if err := fixture.pool.QueryRow(ctx, `SELECT count(*), max(number) FROM answers WHERE team_code=$1 AND problem_code=$2`, teamCode, problem).Scan(&count, &maxNumber); err != nil {
			t.Fatalf("inspect persisted answers: %v", err)
		}
		if count != 2 || maxNumber != 2 {
			t.Fatalf("persisted answers count=%d max number=%d, want 2 and 2", count, maxNumber)
		}
	})

	t.Run("deduplicates deployment events and validates transitions", func(t *testing.T) {
		ctx := t.Context()
		const (
			teamCode = int64(4)
			commit   = "2222222222222222222222222222222222222222"
		)
		seedTeamContestantAndContent(t, fixture, teamCode, "deployer", "100000000000000020", commit)

		requestedAt := time.Date(2026, time.August, 30, 2, 0, 0, 0, time.UTC)
		deployment, err := fixture.store.CreateDeployment(ctx, core.Deployment{
			RequestID:     uuid.NewString(),
			TeamCode:      teamCode,
			ProblemCode:   "Q02",
			ContentCommit: commit,
			RequestedAt:   requestedAt,
			LatestStatus:  core.DeploymentQueued,
		})
		if err != nil {
			t.Fatalf("create deployment: %v", err)
		}
		if deployment.Revision != 1 {
			t.Fatalf("first deployment revision = %d, want 1", deployment.Revision)
		}

		deployingMessage := "provisioning"
		deploying := core.DeploymentEvent{
			EventID:    uuid.NewString(),
			OccurredAt: requestedAt.Add(time.Minute + 789*time.Nanosecond),
			Status:     core.DeploymentDeploying,
			Message:    &deployingMessage,
		}
		updated, duplicate, err := fixture.store.AppendDeploymentEvent(ctx, deployment.RequestID, deploying, false)
		if err != nil || duplicate {
			t.Fatalf("append DEPLOYING event: duplicate=%v err=%v", duplicate, err)
		}
		if updated.LatestStatus != core.DeploymentDeploying {
			t.Fatalf("latest status = %s, want DEPLOYING", updated.LatestStatus)
		}

		updated, duplicate, err = fixture.store.AppendDeploymentEvent(ctx, deployment.RequestID, deploying, false)
		if err != nil || !duplicate {
			t.Fatalf("replay identical event: duplicate=%v err=%v", duplicate, err)
		}
		if len(updated.Events) != 1 {
			t.Fatalf("events after identical replay = %d, want 1", len(updated.Events))
		}

		changedMessage := "different payload"
		mismatch := deploying
		mismatch.Message = &changedMessage
		_, _, err = fixture.store.AppendDeploymentEvent(ctx, deployment.RequestID, mismatch, false)
		assertCoreErrorCode(t, err, http.StatusConflict, "duplicate_event_mismatch")

		invalid := core.DeploymentEvent{
			EventID:    uuid.NewString(),
			OccurredAt: requestedAt.Add(2 * time.Minute),
			Status:     core.DeploymentQueued,
		}
		_, _, err = fixture.store.AppendDeploymentEvent(ctx, deployment.RequestID, invalid, false)
		assertCoreErrorCode(t, err, http.StatusConflict, "invalid_deployment_transition")

		completed := core.DeploymentEvent{
			EventID:    uuid.NewString(),
			OccurredAt: requestedAt.Add(3 * time.Minute),
			Status:     core.DeploymentCompleted,
		}
		updated, duplicate, err = fixture.store.AppendDeploymentEvent(ctx, deployment.RequestID, completed, false)
		if err != nil || duplicate || updated.LatestStatus != core.DeploymentCompleted {
			t.Fatalf("append COMPLETED event: status=%s duplicate=%v err=%v", updated.LatestStatus, duplicate, err)
		}

		second, err := fixture.store.CreateDeployment(ctx, core.Deployment{
			RequestID:     uuid.NewString(),
			TeamCode:      teamCode,
			ProblemCode:   "Q02",
			ContentCommit: commit,
			RequestedAt:   requestedAt.Add(4 * time.Minute),
			LatestStatus:  core.DeploymentQueued,
		})
		if err != nil {
			t.Fatalf("create second deployment: %v", err)
		}
		if second.Revision != 2 {
			t.Fatalf("second deployment revision = %d, want 2", second.Revision)
		}
		secondDeploying := core.DeploymentEvent{EventID: uuid.NewString(), OccurredAt: requestedAt.Add(4*time.Minute + 30*time.Second), Status: core.DeploymentDeploying}
		if _, _, err := fixture.store.AppendDeploymentEvent(ctx, second.RequestID, secondDeploying, false); err != nil {
			t.Fatalf("append DEPLOYING event: %v", err)
		}
		failed := core.DeploymentEvent{EventID: uuid.NewString(), OccurredAt: requestedAt.Add(5 * time.Minute), Status: core.DeploymentFailed}
		if _, _, err := fixture.store.AppendDeploymentEvent(ctx, second.RequestID, failed, false); err != nil {
			t.Fatalf("append FAILED event: %v", err)
		}
		recovering := core.DeploymentEvent{EventID: uuid.NewString(), OccurredAt: requestedAt.Add(6 * time.Minute), Status: core.DeploymentDeploying}
		_, _, err = fixture.store.AppendDeploymentEvent(ctx, second.RequestID, recovering, false)
		assertCoreErrorCode(t, err, http.StatusConflict, "invalid_deployment_transition")
		updated, _, err = fixture.store.AppendDeploymentEvent(ctx, second.RequestID, recovering, true)
		if err != nil || updated.LatestStatus != core.DeploymentDeploying {
			t.Fatalf("apply recovery transition: status=%s err=%v", updated.LatestStatus, err)
		}

		rows, err := fixture.pool.Query(ctx, `SELECT event_id FROM deployment_events WHERE request_id=$1`, deployment.RequestID)
		if err != nil {
			t.Fatalf("list first deployment events: %v", err)
		}
		defer rows.Close()
		var eventIDs []string
		for rows.Next() {
			var eventID string
			if err := rows.Scan(&eventID); err != nil {
				t.Fatalf("scan event ID: %v", err)
			}
			eventIDs = append(eventIDs, eventID)
		}
		if err := rows.Err(); err != nil {
			t.Fatalf("iterate event IDs: %v", err)
		}
		sort.Strings(eventIDs)
		wantEventIDs := []string{completed.EventID, deploying.EventID}
		sort.Strings(wantEventIDs)
		if fmt.Sprint(eventIDs) != fmt.Sprint(wantEventIDs) {
			t.Fatalf("persisted event IDs = %v, want %v", eventIDs, wantEventIDs)
		}
	})
}

func newPostgresFixture(t *testing.T) *postgresFixture {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	databaseURL := os.Getenv("ICTSC_TEST_DATABASE_URL")
	if databaseURL == "" {
		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image: postgresTestImage,
				Env: map[string]string{
					"POSTGRES_DB":       "regalia",
					"POSTGRES_USER":     "regalia",
					"POSTGRES_PASSWORD": "regalia-test-password",
				},
				ExposedPorts: []string{"5432/tcp"},
				WaitingFor: wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(90 * time.Second),
			},
			Started: true,
		})
		testcontainers.CleanupContainer(t, container)
		if err != nil {
			t.Fatalf("start %s (use -short to explicitly disable integration tests): %v", postgresTestImage, err)
		}

		host, err := container.Host(ctx)
		if err != nil {
			t.Fatalf("resolve PostgreSQL container host: %v", err)
		}
		port, err := container.MappedPort(ctx, "5432/tcp")
		if err != nil {
			t.Fatalf("resolve PostgreSQL container port: %v", err)
		}
		databaseURL = fmt.Sprintf("postgres://regalia:regalia-test-password@%s/regalia?sslmode=disable", net.JoinHostPort(host, port.Port()))
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		t.Fatalf("parse PostgreSQL test URL: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("open PostgreSQL test pool: %v", err)
	}
	t.Cleanup(pool.Close)
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping PostgreSQL test container: %v", err)
	}

	paths, err := filepath.Glob(filepath.Join(filepath.Dir(canonicalMigrationPath(t)), "*.sql"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		migration, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := pool.Exec(ctx, string(migration), pgx.QueryExecModeSimpleProtocol); err != nil {
			t.Fatalf("apply %s: %v", path, err)
		}
	}

	return &postgresFixture{pool: pool, store: New(pool)}
}

func canonicalMigrationPath(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve integration test source path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "../../../db/migrations/0001_init.sql"))
}

func seedTeamContestantAndContent(t *testing.T, fixture *postgresFixture, teamCode int64, contestant, discordID, commit string) {
	t.Helper()
	ctx := t.Context()
	team := core.Team{
		Code:         teamCode,
		Name:         fmt.Sprintf("team-%d", teamCode),
		Organization: "ICTSC",
		MemberLimit:  5,
	}
	if _, err := fixture.store.CreateTeam(ctx, team); err != nil {
		t.Fatalf("create seed team: %v", err)
	}
	if _, err := fixture.pool.Exec(ctx, `
		INSERT INTO contestants(name,display_name,self_introduction,discord_id,team_code)
		VALUES($1,$2,'',$3,$4)`, contestant, contestant, discordID, teamCode); err != nil {
		t.Fatalf("create seed contestant: %v", err)
	}
	if _, err := fixture.pool.Exec(ctx, `
		INSERT INTO content_snapshots(commit_sha,repository,git_ref,manifest,fetched_at,activated_at,active)
		VALUES($1,'ictsc/content','refs/heads/main','{"version":1,"sections":[],"problems":[],"announcements":[]}'::jsonb,now(),now(),false)`, commit); err != nil {
		t.Fatalf("create seed content snapshot: %v", err)
	}
}

func assertCoreErrorCode(t *testing.T, err error, status int, code string) {
	t.Helper()
	var domainErr *core.Error
	if !errors.As(err, &domainErr) {
		t.Fatalf("error = %v, want core.Error %s", err, code)
	}
	if domainErr.Status != status || domainErr.Code != code {
		t.Fatalf("error status/code = %d/%s, want %d/%s", domainErr.Status, domainErr.Code, status, code)
	}
}
