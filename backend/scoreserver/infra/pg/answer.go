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

var _ domain.AnswerReader = (*repo)(nil)

func (r *repo) ListAnswersForAdmin(ctx context.Context) ([]*domain.AnswerData, error) {
	ctx, span := tracer.Start(ctx, "ListAnswersForAdmin")
	defer span.End()
	return r.listAnswers(ctx, `
SELECT
	`+answerViewColumns.String("answer")+`,
	`+scoreColumns.As("score")+`
FROM answer_view AS answer
LEFT JOIN answer_scores AS answer_score
	ON answer_score.answer_id = answer.id AND answer_score.visibility = 'PRIVATE'
LEFT JOIN scores AS score ON score.marking_result_id = answer_score.marking_result_id
ORDER BY answer.created_at ASC`)
}

func (r *repo) ListAnswersByTeamProblemForAdmin(ctx context.Context, teamCode int64, problemCode string) ([]*domain.AnswerData, error) {
	ctx, span := tracer.Start(ctx, "ListAnswersByTeamProblemForAdmin")
	defer span.End()
	return r.listAnswers(ctx, `
SELECT
	`+answerViewColumns.String("answer")+`,
	`+scoreColumns.As("score")+`
FROM answer_view AS answer
LEFT JOIN answer_scores AS answer_score
	ON answer_score.answer_id = answer.id AND answer_score.visibility = 'PRIVATE'
LEFT JOIN scores AS score ON score.marking_result_id = answer_score.marking_result_id
WHERE answer."team.code" = $1 AND answer."problem.code" = $2
ORDER BY answer.number ASC`,
		teamCode, problemCode)
}

func (r *repo) ListAnswersByTeamProblemForPublic(ctx context.Context, teamCode int64, problemCode string) ([]*domain.AnswerData, error) {
	ctx, span := tracer.Start(ctx, "ListAnswersByTeamProblemForPublic")
	defer span.End()
	return r.listAnswers(ctx, `
SELECT
	`+answerViewColumns.String("answer")+`,
	`+scoreColumns.As("score")+`
FROM answer_view AS answer
LEFT JOIN answer_scores AS answer_score
	ON answer_score.answer_id = answer.id AND answer_score.visibility = 'PUBLIC'
LEFT JOIN scores AS score ON score.marking_result_id = answer_score.marking_result_id
WHERE answer."team.code" = $1 AND answer."problem.code" = $2
ORDER BY answer.number ASC`,
		teamCode, problemCode)
}

func (r *repo) GetLatestAnswerByTeamProblemForPublic(ctx context.Context, teamID, problemID uuid.UUID) (*domain.AnswerData, error) {
	ctx, span := tracer.Start(ctx, "GetLatestAnswerByTeamProblemForPublic")
	defer span.End()
	return r.getAnswer(ctx, `
SELECT
	`+answerViewColumns.String("answer")+`,
	`+scoreColumns.As("score")+`
FROM answer_view AS answer
LEFT JOIN answer_scores AS answer_score
	ON answer_score.answer_id = answer.id AND answer_score.visibility = 'PUBLIC'
LEFT JOIN scores AS score ON score.marking_result_id = answer_score.marking_result_id
WHERE answer."team.id" = $1 AND answer."problem.id" = $2
ORDER BY answer.number DESC
LIMIT 1`, teamID, problemID)
}

func (r *repo) GetAnswerDetailForPublic(
	ctx context.Context, teamCode int64, problemCode string, answerNumber uint32,
) (*domain.AnswerDetailData, error) {
	ctx, span := tracer.Start(ctx, "GetAnswerDetailForPublic")
	defer span.End()
	return r.getAnswerDetail(ctx, `
SELECT
	`+answerViewColumns.String("answer")+`,
	`+scoreColumns.As("score")+`,
	descriptive.body AS "descriptive.body"
FROM answer_view AS answer
LEFT JOIN answer_scores AS answer_score
	ON answer_score.answer_id = answer.id AND answer_score.visibility = 'PUBLIC'
LEFT JOIN scores AS score ON score.marking_result_id = answer_score.marking_result_id
LEFT JOIN descriptive_answers AS descriptive ON descriptive.answer_id = answer.id
WHERE answer."team.code" = $1 AND answer."problem.code" = $2 AND answer.number = $3
LIMIT 1`, teamCode, problemCode, answerNumber)
}

func (r *repo) GetAnswerDetailForAdmin(
	ctx context.Context, teamCode int64, problemCode string, answerNumber uint32,
) (*domain.AnswerDetailData, error) {
	ctx, span := tracer.Start(ctx, "GetAnswerDetailForAdmin")
	defer span.End()
	return r.getAnswerDetail(ctx, `
SELECT
	`+answerViewColumns.String("answer")+`,
	`+scoreColumns.As("score")+`,
	descriptive.body AS "descriptive.body"
FROM answer_view AS answer
LEFT JOIN answer_scores AS answer_score
	ON answer_score.answer_id = answer.id AND answer_score.visibility = 'PRIVATE'
LEFT JOIN scores AS score ON score.marking_result_id = answer_score.marking_result_id
LEFT JOIN descriptive_answers AS descriptive ON descriptive.answer_id = answer.id
WHERE answer."team.code" = $1 AND answer."problem.code" = $2 AND answer.number = $3
LIMIT 1`, teamCode, problemCode, answerNumber)
}

func (r *repo) listAnswers(ctx context.Context, query string, args ...any) ([]*domain.AnswerData, error) {
	rows, err := r.ext.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select answers")
	}
	defer rows.Close() //nolint:errcheck

	var answers []*domain.AnswerData
	for rows.Next() {
		var row answerWithScoreRow
		if err := rows.StructScan(&row); err != nil {
			return nil, errors.Wrap(err, "failed to scan answer")
		}
		answers = append(answers, row.data())
	}

	return answers, nil
}

func (r *repo) getAnswer(ctx context.Context, query string, args ...any) (*domain.AnswerData, error) {
	var row answerWithScoreRow
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
		Score     *domain.ScoreData   `db:"-"`
	}
	answerRow struct {
		answerDataRow
		Team                     teamRow                          `db:"team"`
		Problem                  problemRow                       `db:"problem"`
		ProblemPercentagePenalty redeployPercentagePenaltyNullRow `db:"problem_rpp"`
		User                     userRow                          `db:"author"`
		RateLimitInterval        pgtype.Interval                  `db:"rate_limit_interval"`
	}
	answerWithScoreRow struct {
		answerRow
		MarkingResultID sql.Null[uuid.UUID] `db:"score.marking_result_id"`
		MarkedScore     sql.Null[uint32]    `db:"score.marked_score"`
		Penalty         sql.Null[uint32]    `db:"score.penalty"`
	}
	answerDetailRow struct {
		answerWithScoreRow
		DescriptiveAnswerBody sql.Null[string] `db:"descriptive.body"`
	}
)

var (
	answerViewColumns = columns([]string{
		"id", "number", "created_at", "rate_limit_interval",
		"team.id", "team.code", "team.name", "team.organization", "team.max_members",
		"problem.id", "problem.code", "problem.type", "problem.title", "problem.max_score", "problem.category", "problem.redeploy_rule",
		"problem_rpp.threshold", "problem_rpp.percentage",
		"author.id", "author.name",
	})
	scoreColumns = columns([]string{"marking_result_id", "marked_score", "penalty"})
)

func (r answerRow) data() *domain.AnswerData {
	r.answerDataRow.Team = (*domain.TeamData)(&r.Team)
	r.answerDataRow.Problem = r.Problem.data()
	r.answerDataRow.Problem.PercentagePenalty = r.ProblemPercentagePenalty.data()
	r.answerDataRow.Author = (*domain.UserData)(&r.User)
	r.Interval = time.Microsecond * time.Duration(r.RateLimitInterval.Microseconds)
	return (*domain.AnswerData)(&r.answerDataRow)
}

func (r answerWithScoreRow) data() *domain.AnswerData {
	answer := r.answerRow.data()
	if r.MarkingResultID.Valid && r.MarkedScore.Valid && r.Penalty.Valid {
		answer.Score = &domain.ScoreData{
			MarkingResultID: r.MarkingResultID.V,
			MarkedScore:     r.MarkedScore.V,
			Penalty:         r.Penalty.V,
		}
	}
	return answer
}

func (r answerDetailRow) data() *domain.AnswerDetailData {
	answer := r.answerWithScoreRow.data()
	body := &domain.AnswerBodyData{}
	if r.DescriptiveAnswerBody.Valid {
		body.Descriptive = &domain.DescriptiveAnswerBodyData{Body: r.DescriptiveAnswerBody.V}
	}
	return &domain.AnswerDetailData{Answer: answer, Body: body}
}

var _ domain.AnswerWriter = (*RepositoryTx)(nil)

var (
	createAnswerQuery = `
INSERT INTO answers (id, number, team_id, problem_id, user_id, created_at_range)
	VALUES ($1, $2, $3, $4, $5, tstzrange($6::timestamptz, $6::timestamptz + $7::interval))`
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
