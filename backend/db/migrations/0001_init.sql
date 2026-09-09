BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE teams (
    code BIGINT PRIMARY KEY CHECK (code BETWEEN 2 AND 99),
    name VARCHAR(255) NOT NULL UNIQUE CHECK (btrim(name) <> ''),
    organization VARCHAR(255) NOT NULL CHECK (btrim(organization) <> ''),
    member_limit INTEGER NOT NULL CHECK (member_limit > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE contestants (
    name VARCHAR(32) PRIMARY KEY CHECK (btrim(name) <> ''),
    display_name VARCHAR(255) NOT NULL CHECK (btrim(display_name) <> ''),
    self_introduction VARCHAR(2000) NOT NULL DEFAULT '',
    discord_id TEXT NOT NULL UNIQUE CHECK (discord_id ~ '^[0-9]+$'),
    team_code BIGINT NOT NULL REFERENCES teams(code) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX contestants_team_idx ON contestants(team_code, name);

CREATE TABLE invitations (
    code VARCHAR(255) PRIMARY KEY CHECK (btrim(code) <> ''),
    team_code BIGINT NOT NULL REFERENCES teams(code) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL CHECK (expires_at > created_at),
    consumed_at TIMESTAMPTZ,
    consumed_by VARCHAR(32) REFERENCES contestants(name) ON DELETE SET NULL,
    CHECK ((consumed_at IS NULL) = (consumed_by IS NULL))
);
CREATE INDEX invitations_team_idx ON invitations(team_code, created_at DESC);

CREATE TABLE content_snapshots (
    commit_sha VARCHAR(64) PRIMARY KEY CHECK (commit_sha ~ '^[0-9a-f]{40,64}$'),
    repository TEXT NOT NULL,
    git_ref TEXT NOT NULL,
    manifest JSONB NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL,
    activated_at TIMESTAMPTZ NOT NULL,
    active BOOLEAN NOT NULL DEFAULT false
);
CREATE UNIQUE INDEX content_one_active_idx ON content_snapshots(active) WHERE active;
CREATE INDEX content_activated_idx ON content_snapshots(activated_at DESC);

CREATE TABLE answer_counters (
    team_code BIGINT NOT NULL REFERENCES teams(code) ON DELETE RESTRICT,
    problem_code VARCHAR(8) NOT NULL,
    next_number INTEGER NOT NULL CHECK (next_number > 0),
    last_submitted_at TIMESTAMPTZ,
    PRIMARY KEY (team_code, problem_code)
);

CREATE TABLE answers (
    team_code BIGINT NOT NULL REFERENCES teams(code) ON DELETE RESTRICT,
    problem_code VARCHAR(8) NOT NULL,
    number INTEGER NOT NULL CHECK (number > 0),
    author_name VARCHAR(32) NOT NULL REFERENCES contestants(name) ON DELETE RESTRICT,
    body TEXT NOT NULL CHECK (length(body) BETWEEN 1 AND 100000),
    submitted_at TIMESTAMPTZ NOT NULL,
    content_commit VARCHAR(64) NOT NULL REFERENCES content_snapshots(commit_sha) ON DELETE RESTRICT,
    max_score INTEGER NOT NULL CHECK (max_score > 0),
    redeploy_rule JSONB NOT NULL,
    deployments_before INTEGER NOT NULL CHECK (deployments_before >= 0),
    PRIMARY KEY (team_code, problem_code, number)
);
CREATE INDEX answers_submitted_idx ON answers(submitted_at DESC);
CREATE INDEX answers_problem_idx ON answers(problem_code, submitted_at DESC);

CREATE TYPE marking_visibility AS ENUM ('PRIVATE', 'TEAM', 'PUBLIC');
CREATE TABLE marking_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_code BIGINT NOT NULL,
    problem_code VARCHAR(8) NOT NULL,
    answer_number INTEGER NOT NULL,
    judge TEXT NOT NULL CHECK (btrim(judge) <> ''),
    marked_score INTEGER NOT NULL CHECK (marked_score >= 0),
    rationale TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    visibility marking_visibility NOT NULL,
    FOREIGN KEY (team_code, problem_code, answer_number)
        REFERENCES answers(team_code, problem_code, number) ON DELETE RESTRICT
);
CREATE INDEX marking_answer_latest_idx
    ON marking_results(team_code, problem_code, answer_number, created_at DESC, id DESC);

CREATE TABLE score_selections (
    team_code BIGINT NOT NULL REFERENCES teams(code) ON DELETE RESTRICT,
    problem_code VARCHAR(8) NOT NULL,
    answer_number INTEGER NOT NULL,
    marking_result_id UUID NOT NULL REFERENCES marking_results(id) ON DELETE RESTRICT,
    marked_score INTEGER NOT NULL CHECK (marked_score >= 0),
    penalty INTEGER NOT NULL CHECK (penalty >= 0),
    effective_score INTEGER NOT NULL CHECK (effective_score >= 0),
    max_score INTEGER NOT NULL CHECK (max_score > 0),
    content_commit VARCHAR(64) NOT NULL REFERENCES content_snapshots(commit_sha) ON DELETE RESTRICT,
    calculated_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (team_code, problem_code),
    FOREIGN KEY (team_code, problem_code, answer_number)
        REFERENCES answers(team_code, problem_code, number) ON DELETE RESTRICT
);

CREATE TYPE deployment_status AS ENUM ('QUEUED', 'DEPLOYING', 'COMPLETED', 'FAILED');
CREATE TABLE deployment_counters (
    team_code BIGINT NOT NULL REFERENCES teams(code) ON DELETE RESTRICT,
    problem_code VARCHAR(8) NOT NULL,
    next_revision INTEGER NOT NULL CHECK (next_revision > 0),
    PRIMARY KEY (team_code, problem_code)
);

CREATE TABLE deployments (
    request_id UUID PRIMARY KEY,
    team_code BIGINT NOT NULL REFERENCES teams(code) ON DELETE RESTRICT,
    problem_code VARCHAR(8) NOT NULL,
    revision INTEGER NOT NULL CHECK (revision > 0),
    content_commit VARCHAR(64) NOT NULL REFERENCES content_snapshots(commit_sha) ON DELETE RESTRICT,
    requested_at TIMESTAMPTZ NOT NULL,
    latest_status deployment_status NOT NULL DEFAULT 'QUEUED',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(team_code, problem_code, revision)
);
CREATE INDEX deployments_latest_idx ON deployments(team_code, problem_code, revision DESC);
CREATE UNIQUE INDEX deployments_one_active_idx
    ON deployments(team_code, problem_code)
    WHERE latest_status IN ('QUEUED', 'DEPLOYING');

CREATE TABLE deployment_events (
    event_id UUID PRIMARY KEY,
    request_id UUID NOT NULL REFERENCES deployments(request_id) ON DELETE CASCADE,
    occurred_at TIMESTAMPTZ NOT NULL,
    status deployment_status NOT NULL,
    message VARCHAR(2000)
);
CREATE INDEX deployment_events_request_idx ON deployment_events(request_id, occurred_at, event_id);

CREATE TABLE competition_state (
    singleton BOOLEAN PRIMARY KEY DEFAULT true CHECK (singleton),
    rule_markdown TEXT NOT NULL DEFAULT '',
    ranking_freeze_at TIMESTAMPTZ,
    final_revealed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by TEXT NOT NULL DEFAULT 'system'
);
INSERT INTO competition_state(singleton) VALUES(true);

CREATE TABLE ranking_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    frozen_at TIMESTAMPTZ NOT NULL,
    content_commit VARCHAR(64) NOT NULL REFERENCES content_snapshots(commit_sha) ON DELETE RESTRICT,
    ranking JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by TEXT NOT NULL,
    UNIQUE(frozen_at, content_commit)
);

CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor TEXT NOT NULL,
    action TEXT NOT NULL,
    resource TEXT NOT NULL,
    before_json JSONB,
    after_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX audit_log_created_idx ON audit_log(created_at DESC);

COMMIT;
