CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code BIGINT NOT NULL UNIQUE CHECK (code BETWEEN 1 AND 99),
    name VARCHAR(255) NOT NULL UNIQUE CHECK (name <> ''),
    organization VARCHAR(255) NOT NULL CHECK (organization <> ''),
    max_members INTEGER NOT NULL CHECK (max_members > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE teams IS 'チーム';
COMMENT ON COLUMN teams.id IS '内部参照用のチームID';
COMMENT ON COLUMN teams.code IS 'APIで使用するチーム番号';
COMMENT ON COLUMN teams.name IS 'チーム名';
COMMENT ON COLUMN teams.organization IS '所属組織名';
COMMENT ON COLUMN teams.max_members IS '参加できる最大人数';
