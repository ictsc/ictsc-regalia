package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

var (
	_ domain.AnswerReader = (*repo)(nil)

	answerQueryBase = `
SELECT
	` + answerColumns.String("a") + `,
	` + teamColumns.As("team") + `,
	` + problemCols.As("problem") + `,
	` + redeployPercentagePenaltyCols.As("problem_rpp") + `,
	` + userColumns.As("author") + `
FROM answers AS a
INNER JOIN teams AS team ON a.team_id = team.id
INNER JOIN problems AS problem ON a.problem_id = problem.id
LEFT JOIN redeploy_percentage_penalties AS problem_rpp ON problem.id = problem_rpp.problem_id
INNER JOIN users AS author ON a.user_id = author.id`

	listAnswersQuery = answerQueryBase + `
ORDER BY a.created_at ASC`
	listAnswersByTeamProblemQuery = answerQueryBase + `
WHERE team.code = $1 AND problem.code = $2
ORDER BY a.number ASC`
	getLatestAnswerByTeamProblemQuery = answerQueryBase + `
WHERE a.team_id = $1 AND a.problem_id = $2
ORDER BY a.number DESC
LIMIT 1`

	answerDetailQueryBase = `
SELECT
	` + answerColumns.String("a") + `,
	` + teamColumns.As("team") + `,
	` + problemCols.As("problem") + `,
	` + redeployPercentagePenaltyCols.As("problem_rpp") + `,
	` + userColumns.As("author") + `,
	descriptive.body AS "descriptive.body"
FROM answers AS a
INNER JOIN teams AS team ON a.team_id = team.id
INNER JOIN problems AS problem ON a.problem_id = problem.id
LEFT JOIN redeploy_percentage_penalties AS problem_rpp ON problem.id = problem_rpp.problem_id
INNER JOIN users AS author ON a.user_id = author.id
LEFT JOIN descriptive_answers AS descriptive ON a.id = descriptive.answer_id`

	answerDetailQuery = answerDetailQueryBase + `
WHERE team.code = $1 AND problem.code = $2 AND a.number = $3`
)

func (r *repo) ListAnswers(ctx context.Context) ([]*domain.AnswerData, error) {
	ctx, span := tracer.Start(ctx, "ListAnswers")
	defer span.End()
	return r.listAnswers(ctx, listAnswersQuery)
}

func (r *repo) ListAnswersByTeamProblem(ctx context.Context, teamCode int64, problemCode string) ([]*domain.AnswerData, error) {
	ctx, span := tracer.Start(ctx, "ListAnswersByTeamProblem")
	defer span.End()
	return r.listAnswers(ctx, listAnswersByTeamProblemQuery, teamCode, problemCode)
}

func (r *repo) GetLatestAnswerByTeamProblem(ctx context.Context, teamID, problemID uuid.UUID) (*domain.AnswerData, error) {
	ctx, span := tracer.Start(ctx, "GetLatestAnswerByTeamProblem")
	defer span.End()
	return r.getAnswer(ctx, getLatestAnswerByTeamProblemQuery, teamID, problemID)
}

func (r *repo) GetAnswerDetail(ctx context.Context, teamCode int64, problemCode string, answerNumber uint32) (*domain.AnswerDetailData, error) {
	ctx, span := tracer.Start(ctx, "GetAnswerDetail")
	defer span.End()
	return r.getAnswerDetail(ctx, answerDetailQuery, teamCode, problemCode, answerNumber)
}

func (r *repo) listAnswers(ctx context.Context, query string, args ...any) ([]*domain.AnswerData, error) {
	rows, err := r.ext.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select answers")
	}
	defer rows.Close() //nolint:errcheck

	var answers []*domain.AnswerData
	for rows.Next() {
		var row answerRow
		if err := rows.StructScan(&row); err != nil {
			return nil, errors.Wrap(err, "failed to scan answer")
		}
		answers = append(answers, row.data())
	}

	return answers, nil
}

func (r *repo) getAnswer(ctx context.Context, query string, args ...any) (*domain.AnswerData, error) {
	var row answerRow
	if err := sqlx.GetContext(ctx, r.ext, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("answer", nil)
		}
		return nil, errors.Wrap(err, "failed to select answer")
	}
	return row.data(), nil
}

func (r *repo) getAnswerDetail(ctx context.Context, query string, args ...any) (*domain.AnswerDetailData, error) {
	var row answerDetailRow
	if err := sqlx.GetContext(ctx, r.ext, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("answer", nil)
		}
		return nil, errors.Wrap(err, "failed to select answer")
	}
	return row.data(), nil
}

type (
	answerDataRow struct {
		ID        uuid.UUID           `db:"id"`
		Number    uint32              `db:"number"`
		Team      *domain.TeamData    `db:"-"`
		Problem   *domain.ProblemData `db:"-"`
		Author    *domain.UserData    `db:"-"`
		CreatedAt time.Time           `db:"created_at"`
		Interval  time.Duration       `db:"-"`
	}
	answerRow struct {
		answerDataRow
		Team                     teamRow                          `db:"team"`
		Problem                  problemRow                       `db:"problem"`
		ProblemPercentagePenalty redeployPercentagePenaltyNullRow `db:"problem_rpp"`
		User                     userRow                          `db:"author"`
		RateLimitInterval        pgtype.Interval                  `db:"rate_limit_interval"`
	}
	answerDetailRow struct {
		answerRow
		DescriptiveAnswerBody sql.Null[string] `db:"descriptive.body"`
	}
)

var answerColumns = columns([]string{"id", "number", "created_at", "rate_limit_interval"})

func (r answerRow) data() *domain.AnswerData {
	r.answerDataRow.Team = (*domain.TeamData)(&r.Team)
	r.answerDataRow.Problem = r.Problem.data()
	r.answerDataRow.Problem.PercentagePenalty = r.ProblemPercentagePenalty.data()
	r.answerDataRow.Author = (*domain.UserData)(&r.User)
	r.Interval = time.Microsecond * time.Duration(r.RateLimitInterval.Microseconds)
	return (*domain.AnswerData)(&r.answerDataRow)
}

func (r answerDetailRow) data() *domain.AnswerDetailData {
	answer := r.answerRow.data()
	body := &domain.AnswerBodyData{}
	if r.DescriptiveAnswerBody.Valid {
		body.Descriptive = &domain.DescriptiveAnswerBodyData{Body: r.DescriptiveAnswerBody.V}
	}
	return &domain.AnswerDetailData{Answer: answer, Body: body}
}

var _ domain.AnswerWriter = (*RepositoryTx)(nil)

var (
	createAnswerQuery = `
INSERT INTO answers (id, number, team_id, problem_id, user_id, created_at, rate_limit_interval, created_at_range)
	VALUES ($1, $2, $3, $4, $5, $6, $7, tstzrange($8::timestamptz, $8::timestamptz + $7::interval))`
	createDescriptiveAnswerQuery = `
INSERT INTO descriptive_answers (answer_id, body)
VALUES ($1, $2)`
)

func (r *RepositoryTx) CreateAnswer(ctx context.Context, data *domain.AnswerDetailData) error {
	ctx, span := tracer.Start(ctx, "CreateAnswer")
	defer span.End()

	if _, err := r.ext.ExecContext(ctx, createAnswerQuery,
		data.Answer.ID,
		data.Answer.Number,
		data.Answer.Team.ID,
		data.Answer.Problem.ID,
		data.Answer.Author.ID,
		data.Answer.CreatedAt,
		data.Answer.Interval,
		data.Answer.CreatedAt,
	); err != nil {
		if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ExclusionViolation && pgErr.ConstraintName == "answers_rate_limit" {
				return errors.WithStack(domain.ErrTooEarlyToSubmitAnswer)
			}
		}
		return errors.Wrap(err, "failed to insert answer")
	}

	if data.Body.Descriptive != nil {
		if _, err := r.ext.ExecContext(ctx, createDescriptiveAnswerQuery,
			data.Answer.ID, data.Body.Descriptive.Body,
		); err != nil {
			return errors.Wrap(err, "failed to insert descriptive answer")
		}
	}

	return nil
}
