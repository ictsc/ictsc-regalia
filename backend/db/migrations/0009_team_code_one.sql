BEGIN;

-- Team codes start at 1: the secondary qualifier's rstate workspaces run
-- from team01, and team 1 is a contestant team there.
ALTER TABLE teams DROP CONSTRAINT teams_code_check;
ALTER TABLE teams ADD CONSTRAINT teams_code_check CHECK (code BETWEEN 1 AND 99);

COMMIT;
