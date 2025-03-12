package pg

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jmoiron/sqlx"
)

type ruleRow struct {
	Markdown string `db:"markdown"`
}

var _ domain.RuleReader = (*repo)(nil)

func (r *repo) GetRule(ctx context.Context) (*domain.RuleData, error) {
	var row ruleRow
	if err := sqlx.GetContext(ctx, r.ext, &row, `
		SELECT markdown FROM rules LIMIT 1
	`); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("rule", nil)
		}
		return nil, errors.Wrap(err, "failed to get rule")
	}
	return (*domain.RuleData)(&row), nil
}

func (r *RepositoryTx) SaveRule(ctx context.Context, data *domain.RuleData) error {
	if _, err := r.ext.ExecContext(ctx, `DELETE FROM rules`); err != nil {
		return errors.Wrap(err, "failed to delete rule")
	}
	if _, err := sqlx.NamedExecContext(ctx, r.ext, `
		INSERT INTO rules (markdown) VALUES (:markdown)
	`, (*ruleRow)(data)); err != nil {
		return errors.Wrap(err, "failed to insert rule")
	}
	return nil
}
