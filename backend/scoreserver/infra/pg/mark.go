package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	_ domain.MarkingResultReader = (*repo)(nil)
)

var (
	listMarkingResultsQuery = `
SELECT
	mr.id, mr.judge_name, mr.created_at,
	` + answerViewColumns.As("answer") + `,
	score.marked_score AS "score.marked_score",
	score.penalty AS "score.penalty",
	dr.rationale AS "rationale.descriptive_comment"
FROM marking_results AS mr
INNER JOIN answer_view AS answer ON mr.answer_id = answer.id
INNER JOIN scores AS score ON mr.id = score.marking_result_id
LEFT JOIN descriptive_marking_rationales AS dr ON mr.id = dr.marking_result_id
ORDER BY mr.created_at DESC`
)

type (
	markingResultDataRow struct {
		ID        uuid.UUID                    `db:"id"`
		Judge     string                       `db:"judge_name"`
		Answer    *domain.AnswerData           `db:"-"`
		Score     *domain.ScoreData            `db:"-"`
		Rationale *domain.MarkingRationaleData `db:"-"`
		CreatedAt time.Time                    `db:"created_at"`
	}
	scoreRow struct {
		MarkedScore uint32 `db:"marked_score"`
		Penalty     uint32 `db:"penalty"`
	}
	rationaleRow struct {
		DescriptiveComment sql.NullString `db:"descriptive_comment"`
	}
	markingResultRow struct {
		markingResultDataRow
		Answer    answerRow    `db:"answer"`
		Score     scoreRow     `db:"score"`
		Rationale rationaleRow `db:"rationale"`
	}
)

func (r *markingResultRow) data() *domain.MarkingResultData {
	r.markingResultDataRow.Answer = r.Answer.data()
	r.markingResultDataRow.Score = (*domain.ScoreData)(&r.Score)
	r.markingResultDataRow.Rationale = r.Rationale.data()
	return (*domain.MarkingResultData)(&r.markingResultDataRow)
}

func (dr *rationaleRow) data() *domain.MarkingRationaleData {
	return &domain.MarkingRationaleData{DescriptiveComment: dr.DescriptiveComment.String}
}

func (r *repo) ListMarkingResults(ctx context.Context) ([]*domain.MarkingResultData, error) {
	ctx, span := tracer.Start(ctx, "ListMarkingResults")
	defer span.End()

	rows, err := r.ext.QueryxContext(ctx, listMarkingResultsQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select marking results")
	}
	defer rows.Close() //nolint:errcheck

	var markingResults []*domain.MarkingResultData
	for rows.Next() {
		var row markingResultRow
		if err := rows.StructScan(&row); err != nil {
			return nil, errors.Wrap(err, "failed to scan marking result")
		}
		markingResults = append(markingResults, row.data())
	}

	return markingResults, nil
}

var _ domain.MarkingResultWriter = (*RepositoryTx)(nil)

func (r *RepositoryTx) CreateMarkingResult(ctx context.Context, data *domain.MarkingResultData) error {
	ctx, span := tracer.Start(ctx, "CreateMarkingResult")
	defer span.End()

	if _, err := r.ext.ExecContext(
		ctx,
		`INSERT INTO marking_results (id, judge_name, answer_id, created_at)
		VALUES ($1, $2, $3, $4)`,
		data.ID, data.Judge, data.Answer.ID, data.CreatedAt,
	); err != nil {
		return errors.Wrap(err, "failed to insert into marking_results")
	}

	if _, err := r.ext.ExecContext(
		ctx,
		`INSERT INTO scores (marking_result_id, marked_score)
		VALUES ($1, $2)`,
		data.ID, data.Score.MarkedScore,
	); err != nil {
		return errors.Wrap(err, "failed to insert into scores")
	}

	if data.Rationale.DescriptiveComment != "" {
		if _, err := r.ext.ExecContext(
			ctx,
			`INSERT INTO descriptive_marking_rationales (marking_result_id, rationale)
			VALUES ($1, $2)`,
			data.ID, data.Rationale.DescriptiveComment,
		); err != nil {
			return errors.Wrap(err, "failed to insert into descriptive_marking_rationales")
		}
	}

	return nil
}

func (r *repo) UpdatePenalty(ctx context.Context, id uuid.UUID, penalty uint32) error {
	ctx, span := tracer.Start(ctx, "UpdatePenalty", trace.WithAttributes(
		attribute.String("marking_result_id", id.String()),
		attribute.Int64("penalty", int64(penalty)),
	))
	defer span.End()

	if _, err := r.ext.ExecContext(
		ctx,
		"UPDATE scores SET penalty = $1 WHERE marking_result_id = $2",
		penalty, id,
	); err != nil {
		return errors.Wrap(err, "failed to update penalty")
	}
	return nil
}
