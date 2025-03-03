package pg

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

var _ domain.NoticeReader = (*repo)(nil)
var _ domain.NoticeWriter = (*repo)(nil)

type (
	noticeRow struct {
		ID             uuid.UUID  `db:"id"`
		Path           string     `db:"path"`
		Title          string     `db:"title"`
		Markdown       string     `db:"markdown"`
		EffectiveFrom  *time.Time `db:"effective_from"`
		EffectiveUntil *time.Time `db:"effective_until"`
	}
)

// TODO: データを昇順降順などするか聞く
var listNoticesQuery = `
	SELECT 
		id, path, title, markdown, effective_from, effective_until
	FROM 
		notices
	WHERE 
		effective_from <= NOW() AND effective_until >= NOW();
	`

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
		notices = append(notices, row.toDomain())
	}
	return notices, nil
}

var saveNoticeQuery = `
	INSERT INTO notices
		(id, path, title, markdown, effective_from, effective_until)
	VALUES
		(:id, :path, :title, :markdown, :effective_from, :effective_until)
	ON CONFLICT (id) DO UPDATE SET
		path = EXCLUDED.path,
		title = EXCLUDED.title,
		markdown = EXCLUDED.markdown,
		effective_from = EXCLUDED.effective_from,
		effective_until = EXCLUDED.effective_until;
`

func (r *repo) SaveNotice(ctx context.Context, notice *domain.NoticeData) error {
	nRow := fromDomain(notice) // domain.NoticeData -> noticeRow に変換
	if _, err := sqlx.NamedExecContext(ctx, r.ext, saveNoticeQuery, nRow); err != nil {
		if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.NewAlreadyExistsError("notice", nil)
			}
		}
		return errors.Wrap(err, "failed to save notice")
	}
	return nil
}

func (n *noticeRow) toDomain() *domain.NoticeData {
	return &domain.NoticeData{
		ID:             n.ID,
		Path:           n.Path,
		Title:          n.Title,
		Markdown:       n.Markdown,
		EffectiveFrom:  n.EffectiveFrom,
		EffectiveUntil: n.EffectiveUntil,
	}
}

func fromDomain(notice *domain.NoticeData) *noticeRow {
	return &noticeRow{
		Path:           notice.Path,
		Title:          notice.Title,
		ID:             notice.ID,
		Markdown:       notice.Markdown,
		EffectiveFrom:  notice.EffectiveFrom,
		EffectiveUntil: notice.EffectiveUntil,
	}
}
