CREATE EXTENSION btree_gist;

CREATE TABLE rules (
	page_path TEXT, -- Deprecated: いずれ消す
	markdown TEXT NOT NULL
);
COMMENT ON TABLE rules IS 'ルール';
COMMENT ON COLUMN rules.markdown IS 'Markdown形式のルール';

CREATE TABLE teams (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	code BIGINT NOT NULL UNIQUE,
	name VARCHAR(255) NOT NULL UNIQUE,
	organization VARCHAR(255) NOT NULL,
	max_members INT NOT NULL DEFAULT 1 CHECK (max_members > 0),
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE teams IS 'チーム';
COMMENT ON COLUMN teams.id IS 'チーム ID';
COMMENT ON COLUMN teams.code IS 'チーム番号';
COMMENT ON COLUMN teams.name IS 'チーム名';
COMMENT ON COLUMN teams.organization IS '所属組織名';

CREATE TABLE invitation_codes (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
	code VARCHAR(255) NOT NULL UNIQUE,
	expires_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE invitation_codes IS '招待コード';
COMMENT ON COLUMN invitation_codes.id IS '招待コード ID';
COMMENT ON COLUMN invitation_codes.team_id IS 'チーム ID';
COMMENT ON COLUMN invitation_codes.code IS '招待コード';
COMMENT ON COLUMN invitation_codes.expires_at IS '有効期限';

CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(32) NOT NULL UNIQUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE users IS 'ユーザ';
COMMENT ON COLUMN users.id IS 'ユーザ ID';
COMMENT ON COLUMN users.name IS 'ユーザ名';

CREATE TABLE user_profiles (
	user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
	display_name VARCHAR(255) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE user_profiles IS 'ユーザプロフィール';
COMMENT ON COLUMN user_profiles.user_id IS 'ユーザ ID';
COMMENT ON COLUMN user_profiles.display_name IS 'ユーザ表示名';

CREATE TABLE discord_users (
	user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
	discord_user_id BIGINT NOT NULL UNIQUE,
	linked_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE discord_users IS 'Discord 上のユーザ情報';
COMMENT ON COLUMN discord_users.discord_user_id IS 'Discord ユーザ ID';
COMMENT ON COLUMN discord_users.user_id IS 'ユーザ ID';

CREATE TABLE team_members (
	user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
	team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
	invitation_code_id UUID NOT NULL REFERENCES invitation_codes(id) ON DELETE CASCADE,
	invited_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE team_members IS 'チームメンバ';
COMMENT ON COLUMN team_members.user_id IS 'ユーザ ID';
COMMENT ON COLUMN team_members.team_id IS 'チーム ID';
COMMENT ON COLUMN team_members.invitation_code_id IS '招待コード ID';
COMMENT ON COLUMN team_members.invited_at IS '招待日時';

CREATE TYPE problem_type AS ENUM ('DESCRIPTIVE');
CREATE TYPE redeploy_rule AS ENUM ('UNREDEPLOYABLE', 'PERCENTAGE_PENALTY');
CREATE TABLE problems (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	code VARCHAR(8) NOT NULL UNIQUE,
	type problem_type NOT NULL,
	title VARCHAR(255) NOT NULL,
	max_score INT NOT NULL CHECK (max_score > 0),
	category TEXT NOT NULL DEFAULT '',
	redeploy_rule redeploy_rule NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE problems IS '問題';
COMMENT ON COLUMN problems.id IS '問題 ID';
COMMENT ON COLUMN problems.type IS '問題の種別';
COMMENT ON COLUMN problems.title IS '問題名';
COMMENT ON COLUMN problems.max_score IS '最大得点';
COMMENT ON COLUMN problems.category IS '問題のカテゴリー';
COMMENT ON COLUMN problems.redeploy_rule IS '再展開時のルール';

CREATE TABLE redeploy_percentage_penalties (
	problem_id UUID PRIMARY KEY REFERENCES problems(id) ON DELETE CASCADE,
	threshold INT NOT NULL CHECK (threshold >= 0),
	percentage INT NOT NULL CHECK (percentage >= 0 AND percentage < 100)
);
COMMENT ON TABLE redeploy_percentage_penalties IS '再展開時の割合ペナルティ';
COMMENT ON COLUMN redeploy_percentage_penalties.problem_id IS 'ペナルティ適用対象の問題 ID';
COMMENT ON COLUMN redeploy_percentage_penalties.threshold IS 'ペナルティが適用される再展開回数の閾値';
COMMENT ON COLUMN redeploy_percentage_penalties.percentage IS '再展開一回あたりの最大得点に対する減点率(%)';

CREATE TABLE problem_contents (
	problem_id UUID PRIMARY KEY REFERENCES problems(id) ON DELETE CASCADE,
	page_id VARCHAR(255) NULL, -- Deprecated
	page_path VARCHAR(255) NULL, -- Deprecated
	body TEXT NOT NULL,
	explanation TEXT NOT NULL
);
COMMENT ON TABLE problem_contents IS '問題の内容';
COMMENT ON COLUMN problem_contents.problem_id IS '問題 ID';
COMMENT ON COLUMN problem_contents.body IS '問題文';
COMMENT ON COLUMN problem_contents.explanation IS '運営向け解説情報';

CREATE TABLE answers (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE RESTRICT,
	team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
	number INT NOT NULL CHECK (number > 0),
	UNIQUE (problem_id, team_id, number),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	created_at_range TSTZRANGE NOT NULL,
	CONSTRAINT answers_rate_limit EXCLUDE USING GIST
		(problem_id WITH =, team_id WITH =, created_at_range WITH &&)
);
COMMENT ON TABLE answers IS '解答';
COMMENT ON COLUMN answers.id IS '解答 ID';
COMMENT ON COLUMN answers.problem_id IS '解答対象の問題 ID';
COMMENT ON COLUMN answers.team_id IS '解答したチーム ID';
COMMENT ON COLUMN answers.number IS '解答番号';
COMMENT ON COLUMN answers.user_id IS '解答者のユーザ ID';
COMMENT ON COLUMN answers.created_at_range IS '回答日時から次に解答できるまでの期間';

CREATE TABLE descriptive_answers (
	answer_id UUID PRIMARY KEY REFERENCES answers(id) ON DELETE CASCADE,
	body TEXT NOT NULL
);
COMMENT ON TABLE descriptive_answers IS '記述式解答';
COMMENT ON COLUMN descriptive_answers.answer_id IS '解答 ID';
COMMENT ON COLUMN descriptive_answers.body IS '解答内容';

CREATE TYPE marking_result_visilibity AS ENUM ('PRIVATE', 'PUBLIC');
CREATE TABLE marking_results (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	answer_id UUID NOT NULL REFERENCES answers(id) ON DELETE CASCADE,
	judge_name TEXT NOT NULL,
	visibility marking_result_visilibity NOT NULL DEFAULT 'PRIVATE',
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE marking_results IS '採点結果';
COMMENT ON COLUMN marking_results.id IS '採点結果 ID';
COMMENT ON COLUMN marking_results.answer_id IS '解答 ID';
COMMENT ON COLUMN marking_results.judge_name IS '採点者名';
COMMENT ON COLUMN marking_results.created_at IS '採点日時';

CREATE TABLE scores (
	marking_result_id UUID PRIMARY KEY REFERENCES marking_results(id) ON DELETE CASCADE,
	marked_score INT NOT NULL CHECK (marked_score >= 0),
	penalty INT NOT NULL CHECK (penalty >= 0) DEFAULT 0,
	redeploy_count INT NOT NULL CHECK (redeploy_count >= 0) DEFAULT 0,
	total_score INT GENERATED ALWAYS AS (GREATEST(0, marked_score - penalty)) STORED
);
COMMENT ON TABLE scores IS '採点結果の得点';
COMMENT ON COLUMN scores.marking_result_id IS '採点結果 ID';
COMMENT ON COLUMN scores.marked_score IS '採点による得点';
COMMENT ON COLUMN scores.penalty IS 'ペナルティ';
COMMENT ON COLUMN scores.redeploy_count IS '再展開回数(ペナルティの計算根拠)';

CREATE TABLE answer_scores (
	answer_id UUID NOT NULL REFERENCES answers(id) ON DELETE CASCADE,
	visibility marking_result_visilibity NOT NULL DEFAULT 'PRIVATE',
	PRIMARY KEY (answer_id, visibility),
	marking_result_id UUID NOT NULL REFERENCES marking_results(id) ON DELETE CASCADE
);

CREATE TABLE problem_scores (
	team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
	problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
	visibility marking_result_visilibity NOT NULL DEFAULT 'PRIVATE',
	PRIMARY KEY (team_id, problem_id, visibility),
	marking_result_id UUID NOT NULL REFERENCES marking_results(id) ON DELETE CASCADE,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE team_scores (
	team_id UUID PRIMARY KEY REFERENCES teams(id) ON DELETE CASCADE,
	total_score INT NOT NULL CHECK (total_score >= 0) DEFAULT 0,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE descriptive_marking_rationales (
	marking_result_id UUID PRIMARY KEY REFERENCES marking_results(id) ON DELETE CASCADE,
	rationale TEXT NOT NULL
);
COMMENT ON TABLE descriptive_marking_rationales IS '記述式採点の根拠';
COMMENT ON COLUMN descriptive_marking_rationales.marking_result_id IS '採点結果 ID';
COMMENT ON COLUMN descriptive_marking_rationales.rationale IS '採点根拠';

CREATE TYPE deployment_status AS ENUM ('QUEUED', 'DEPLOYING', 'COMPLETED', 'FAILED');

CREATE TABLE redeployment_requests (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
	problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
	revision INT NOT NULL CHECK (revision > 0) DEFAULT 1,
	UNIQUE (team_id, problem_id, revision)
);

CREATE TABLE redeployment_events (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	request_id UUID NOT NULL REFERENCES redeployment_requests(id) ON DELETE CASCADE,
	status deployment_status NOT NULL,
	UNIQUE (request_id, status),
	message TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notices (
	slug VARCHAR(255) PRIMARY KEY,
	path VARCHAR(255), -- Deprecated
	title VARCHAR(255) NOT NULL,
	markdown TEXT NOT NULL,
	effective_from TIMESTAMPTZ NOT NULL,
	effective_until TIMESTAMPTZ -- Deprecated
);
COMMENT ON TABLE notices IS 'お知らせ';
COMMENT ON COLUMN notices.slug IS '名前';
COMMENT ON COLUMN notices.title IS 'タイトル';
COMMENT ON COLUMN notices.markdown IS '本文';
COMMENT ON COLUMN notices.effective_from IS '掲示開始時間';

CREATE TYPE contest_phase AS ENUM ('UNSPECIFIED', 'OUT_OF_CONTEST', 'IN_CONTEST', 'BREAK', 'AFTER_CONTEST');

CREATE TABLE schedules (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phase contest_phase NOT NULL,
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT schedules_start_end CHECK (start_at < end_at)
);

COMMENT ON TABLE schedules IS 'コンテストスケジュール';
COMMENT ON COLUMN schedules.id IS 'スケジュール ID';
COMMENT ON COLUMN schedules.phase IS 'フェーズ';
COMMENT ON COLUMN schedules.start_at IS '開始時刻';
COMMENT ON COLUMN schedules.end_at IS '終了時刻';
