BEGIN;

-- Text logs submitted with an answer. They are written once in the answer's
-- transaction and never changed; names never become file system paths.
CREATE TABLE answer_attachments (
    team_code BIGINT NOT NULL,
    problem_code VARCHAR(8) NOT NULL,
    answer_number INTEGER NOT NULL,
    position SMALLINT NOT NULL CHECK (position BETWEEN 0 AND 2),
    name TEXT NOT NULL CHECK (name ~ '^[A-Za-z0-9_-][A-Za-z0-9._-]{0,95}\.(txt|log)$'),
    size_bytes INTEGER NOT NULL CHECK (size_bytes BETWEEN 1 AND 1048576),
    content TEXT NOT NULL,
    PRIMARY KEY (team_code, problem_code, answer_number, position),
    UNIQUE (team_code, problem_code, answer_number, name),
    FOREIGN KEY (team_code, problem_code, answer_number)
        REFERENCES answers(team_code, problem_code, number) ON DELETE CASCADE
);

COMMIT;
