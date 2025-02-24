DROP VIEW IF EXISTS latest_marking_results;
CREATE VIEW latest_marking_results AS (
	SELECT DISTINCT ON (answer_id) *
	FROM marking_results
	ORDER BY answer_id, created_at DESC
);
