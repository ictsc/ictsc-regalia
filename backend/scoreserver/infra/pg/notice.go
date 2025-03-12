package pg

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jmoiron/sqlx"
)

var _ domain.NoticeReader = (*repo)(nil)

// var _ domain.NoticeWriter = (*repo)(nil)

type noticeRow struct {
	Slug          string    `db:"slug"`
	Title         string    `db:"title"`
	EffectiveFrom time.Time `db:"effective_from"`
	Markdown      string    `db:"markdown"`
}

var listNoticesQuery = `
	SELECT slug, title, markdown, effective_from
	FROM notices
	ORDER BY effective_from DESC`

func (r *repo) ListNotices(ctx context.Context) ([]*domain.NoticeData, error) {
	rows, err := r.ext.QueryxContext(ctx, listNoticesQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list notices")
	}
	defer func() { _ = rows.Close() }()

	var notices []*domain.NoticeData
	for rows.Next() {
		var row noticeRow
		if err := rows.StructScan(&row); err != nil {
			return nil, errors.Wrap(err, "failed to scan notice row")
		}
		notices = append(notices, (*domain.NoticeData)(&row))
	}
	return notices, nil
}

var _ domain.NoticeWriter = (*RepositoryTx)(nil)

func (tx *RepositoryTx) SaveNotices(ctx context.Context, notices []*domain.NoticeData) error {
	if _, err := tx.ext.ExecContext(ctx, "DELETE FROM notices"); err != nil {
		return errors.Wrap(err, "failed to delete notices")
	}
	if len(notices) == 0 {
		return nil
	}

	rows := make([]noticeRow, 0, len(notices))
	for _, notice := range notices {
		rows = append(rows, noticeRow(*notice))
	}

	if _, err := sqlx.NamedExecContext(ctx, tx.ext, `
		INSERT INTO notices (slug, title, markdown, effective_from)
		VALUES (:slug, :title, :markdown, :effective_from)
	`, rows); err != nil {
		return errors.Wrap(err, "failed to insert notices")
	}

	return nil
}
