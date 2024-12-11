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
