CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code BIGINT NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    organization VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE teams IS 'チーム';
COMMENT ON COLUMN teams.id IS 'チーム ID';
COMMENT ON COLUMN teams.code IS 'チーム番号';
COMMENT ON COLUMN teams.name IS 'チーム名';
COMMENT ON COLUMN teams.organization IS 'チームの所属組織名';

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
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	display_name VARCHAR(255) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (user_id)
);
COMMENT ON TABLE user_profiles IS 'ユーザプロフィール';
COMMENT ON COLUMN user_profiles.id IS 'ユーザプロフィール ID';
COMMENT ON COLUMN user_profiles.user_id IS 'ユーザ ID';
COMMENT ON COLUMN user_profiles.display_name IS 'ユーザ表示名';
