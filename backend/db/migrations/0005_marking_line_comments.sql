BEGIN;

ALTER TABLE marking_results
    ADD COLUMN line_comments JSONB NOT NULL DEFAULT '[]'::jsonb
    CHECK (jsonb_typeof(line_comments) = 'array');

COMMIT;
