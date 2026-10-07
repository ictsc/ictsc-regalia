BEGIN;

CREATE TABLE answer_workflows (
    team_code BIGINT NOT NULL,
    problem_code VARCHAR(8) NOT NULL,
    answer_number INTEGER NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('WAITING', 'IN_PROGRESS', 'COMPLETED')),
    claimed_by_discord_id TEXT CHECK (claimed_by_discord_id ~ '^[0-9]+$'),
    revision BIGINT NOT NULL DEFAULT 1 CHECK (revision > 0),
    PRIMARY KEY (team_code, problem_code, answer_number),
    FOREIGN KEY (team_code, problem_code, answer_number)
        REFERENCES answers(team_code, problem_code, number) ON DELETE RESTRICT
);

INSERT INTO answer_workflows (team_code, problem_code, answer_number, status)
SELECT DISTINCT team_code, problem_code, answer_number, 'COMPLETED'
FROM marking_results;

COMMIT;
