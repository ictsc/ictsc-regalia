BEGIN;

-- The job ID a redeploy service without callbacks (rstate) gave the request,
-- so polling reads that job and never an earlier one for the same problem.
ALTER TABLE deployments
    ADD COLUMN gateway_job_id TEXT CHECK (char_length(gateway_job_id) BETWEEN 1 AND 128);

COMMIT;
