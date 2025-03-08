DROP VIEW IF EXISTS latest_marking_results;
CREATE VIEW latest_marking_results AS (
	SELECT DISTINCT ON (answer_id) *
	FROM marking_results
	ORDER BY answer_id, created_at DESC
);

DROP VIEW IF EXISTS answer_view CASCADE;
CREATE VIEW answer_view AS (
	SELECT
		a.id, a.number,
		lower(a.created_at_range) AS "created_at",
		upper(a.created_at_range) - lower(a.created_at_range) AS "rate_limit_interval",

		t.id AS "team.id", t.code AS "team.code", t.name AS "team.name",
		t.organization AS "team.organization", t.max_members AS "team.max_members",

		p.id AS "problem.id", p.code AS "problem.code", p.type AS "problem.type",
		p.title AS "problem.title", p.max_score AS "problem.max_score",
		p.redeploy_rule AS "problem.redeploy_rule",
		rpp.threshold AS "problem_rpp.threshold", rpp.percentage AS "problem_rpp.percentage",

		u.id AS "author.id", u.name AS "author.name"
	FROM answers AS a
	INNER JOIN teams AS t ON a.team_id = t.id
	INNER JOIN problems AS p ON a.problem_id = p.id
	LEFT JOIN redeploy_percentage_penalties AS rpp ON p.id = rpp.problem_id
	INNER JOIN users AS u ON a.user_id = u.id
);
