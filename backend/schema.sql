CREATE EXTENSION pg_stat_statements;
-- CREATE EXTENSION pg_stat_kcache;
-- CREATE EXTENSION set_user;
CREATE EXTENSION btree_gist;

CREATE TABLE rules (
	page_path TEXT,
	markdown TEXT
);
COMMENT ON TABLE rules IS 'ルール';
COMMENT ON COLUMN rules.page_path IS 'Wiki上のページパス';
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
	redeploy_rule redeploy_rule NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE problems IS '問題';
COMMENT ON COLUMN problems.id IS '問題 ID';
COMMENT ON COLUMN problems.type IS '問題の種別';
COMMENT ON COLUMN problems.title IS '問題名';
COMMENT ON COLUMN problems.max_score IS '最大得点';
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
	page_id VARCHAR(255) NOT NULL,
	page_path VARCHAR(255) NOT NULL,
	body TEXT NOT NULL,
	explanation TEXT NOT NULL
);
COMMENT ON TABLE problem_contents IS '問題の内容';
COMMENT ON COLUMN problem_contents.problem_id IS '問題 ID';
COMMENT ON COLUMN problem_contents.page_id IS 'Wiki上のページ ID';
COMMENT ON COLUMN problem_contents.page_path IS 'Wiki上のページパス';
COMMENT ON COLUMN problem_contents.body IS '問題文';
COMMENT ON COLUMN problem_contents.explanation IS '運営向け解説情報';

CREATE TABLE answers (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
	team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
	number INT NOT NULL CHECK (number > 0),
	UNIQUE (problem_id, team_id, number),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	created_at_range TSTZRANGE NOT NULL,
	CONSTRAINT answers_rate_limit EXCLUDE USING GIST
		(problem_id WITH =, team_id WITH =, created_at_range WITH &&)
);
COMMENT ON TABLE answers IS '回答';
COMMENT ON COLUMN answers.id IS '回答 ID';
COMMENT ON COLUMN answers.problem_id IS '回答対象の問題 ID';
COMMENT ON COLUMN answers.team_id IS '回答したチーム ID';
COMMENT ON COLUMN answers.number IS '回答番号';
COMMENT ON COLUMN answers.user_id IS '回答者のユーザ ID';
COMMENT ON COLUMN answers.created_at_range IS '回答日時から次に回答できるまでの期間';

CREATE TABLE descriptive_answers (
	answer_id UUID PRIMARY KEY REFERENCES answers(id) ON DELETE CASCADE,
	body TEXT NOT NULL
);
COMMENT ON TABLE descriptive_answers IS '記述式回答';
COMMENT ON COLUMN descriptive_answers.answer_id IS '回答 ID';
COMMENT ON COLUMN descriptive_answers.body IS '回答内容';

CREATE TABLE marking_results (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	answer_id UUID NOT NULL REFERENCES answers(id) ON DELETE CASCADE,
	score INT NOT NULL CHECK (score >= 0),
	judge_name TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE marking_results IS '採点結果';
COMMENT ON COLUMN marking_results.id IS '採点結果 ID';
COMMENT ON COLUMN marking_results.answer_id IS '回答 ID';
COMMENT ON COLUMN marking_results.score IS '得点';
COMMENT ON COLUMN marking_results.judge_name IS '採点者名';
COMMENT ON COLUMN marking_results.created_at IS '採点日時';

CREATE TABLE descriptive_marking_rationale (
	marking_result_id UUID PRIMARY KEY REFERENCES marking_results(id) ON DELETE CASCADE,
	rationale TEXT NOT NULL
);
COMMENT ON TABLE descriptive_marking_rationale IS '記述式採点の根拠';
COMMENT ON COLUMN descriptive_marking_rationale.marking_result_id IS '採点結果 ID';
COMMENT ON COLUMN descriptive_marking_rationale.rationale IS '採点根拠';

CREATE TYPE deployment_status AS ENUM ('QUEUED', 'DEPLOYING', 'COMPLETED', 'FAILED');

CREATE TABLE redeployment_requests (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
	problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
	status deployment_status NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
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
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	path VARCHAR(255) NOT NULL UNIQUE,
	title VARCHAR(255),
	markdown TEXT NOT NULL,
	effective_from TIMESTAMPTZ NOT NULL,
	effective_until TIMESTAMPTZ NOT NULL
);
COMMENT ON TABLE notices IS 'お知らせ';
COMMENT ON COLUMN notices.id IS 'お知らせ ID';
COMMENT ON COLUMN notices.path IS 'Wiki上のページパス';
COMMENT ON COLUMN notices.title IS 'タイトル';
COMMENT ON COLUMN notices.markdown IS '本文';
COMMENT ON COLUMN notices.effective_from IS '掲示開始時間';
COMMENT ON COLUMN notices.effective_until IS '掲示終了時間';
