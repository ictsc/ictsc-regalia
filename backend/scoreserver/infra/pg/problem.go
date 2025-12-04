package pg

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

var _ domain.ProblemReader = (*repo)(nil)

var listProblemsQuery = `
	SELECT
		` + problemCols.String("p") + `,
		` + redeployPercentagePenaltyCols.As("rpp") + `
	FROM problems AS p
	LEFT JOIN redeploy_percentage_penalties AS rpp ON p.id = rpp.problem_id`

func (r *repo) ListProblems(ctx context.Context) ([]*domain.ProblemData, error) {
	rows, err := r.ext.QueryxContext(ctx, listProblemsQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list problems")
	}
	defer func() { _ = rows.Close() }()

	var problems []*domain.ProblemData
	for rows.Next() {
		var row problemDataRow
		if err := rows.StructScan(&row); err != nil {
			return nil, errors.Wrap(err, "failed to scan problem row")
		}
		problems = append(problems, row.data())
	}
	return problems, nil
}

var getProblemByCodeQuery = listProblemsQuery + `
	WHERE code = $1`

func (r *repo) GetProblemByCode(ctx context.Context, code string) (*domain.ProblemData, error) {
	var row problemDataRow
	if err := sqlx.GetContext(ctx, r.ext, &row, getProblemByCodeQuery, code); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("problem", nil)
		}
		return nil, errors.Wrap(err, "failed to get problem by code")
	}
	return row.data(), nil
}

type problemDataRow struct {
	problemRow
	Rpp redeployPercentagePenaltyNullRow `db:"rpp"`
}

func (r *problemDataRow) data() *domain.ProblemData {
	data := r.problemRow.data()
	data.PercentagePenalty = r.Rpp.data()
	return data
}

var getDescriptiveProblemQuery = `
	SELECT
		` + problemCols.String("p") + `,
		` + problemContentCols.As("c") + `,
		` + redeployPercentagePenaltyCols.As("rpp") + `
	FROM problems AS p
	JOIN problem_contents AS c ON p.id = c.problem_id
	LEFT JOIN redeploy_percentage_penalties AS rpp ON p.id = rpp.problem_id
	WHERE p.type = 'DESCRIPTIVE' AND p.id = $1`

func (r *repo) GetDescriptiveProblem(ctx context.Context, id uuid.UUID) (*domain.DescriptiveProblemData, error) {
	var row struct {
		problemDataRow
		Content problemContentRow `db:"c"`
	}
	if err := sqlx.GetContext(ctx, r.ext, &row, getDescriptiveProblemQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("descriptive problem", nil)
		}
		return nil, errors.Wrap(err, "failed to get descriptive problem")
	}
	return &domain.DescriptiveProblemData{
		Problem: row.problemDataRow.data(),
		Content: row.Content.data(),
	}, nil
}

var _ domain.ProblemWriter = (*RepositoryTx)(nil)

var (
	saveProblemQuery = `
	INSERT INTO problems
		(id, code, type, title, max_score, category, redeploy_rule, created_at, updated_at)
	VALUES
		(:id, :code, :type, :title, :max_score, :category, :redeploy_rule, NOW(), NOW())
	ON CONFLICT (id) DO UPDATE SET
		code = EXCLUDED.code,
		type = EXCLUDED.type,
		title = EXCLUDED.title,
		max_score = EXCLUDED.max_score,
		category = EXCLUDED.category,
		redeploy_rule = EXCLUDED.redeploy_rule,
		updated_at = NOW()`

	saveProblemContentQuery = `
	INSERT INTO problem_contents (problem_id, page_id, page_path, body, explanation)
	VALUES (:problem_id, :page_id, :page_path, :body, :explanation)
	ON CONFLICT (problem_id) DO UPDATE SET
		page_id = EXCLUDED.page_id,
		page_path = EXCLUDED.page_path,
		body = EXCLUDED.body,
		explanation = EXCLUDED.explanation`

	saveRedeployPercentagePenaltyQuery = `
	INSERT INTO redeploy_percentage_penalties (problem_id, threshold, percentage)
	VALUES (:problem_id, :threshold, :percentage)
	ON CONFLICT (problem_id) DO UPDATE SET
		threshold = EXCLUDED.threshold,
		percentage = EXCLUDED.percentage`
)

func (r *RepositoryTx) SaveDescriptiveProblem(ctx context.Context, descriptiveProblem *domain.DescriptiveProblemData) error {
	{
		problem := descriptiveProblem.Problem
		if _, err := sqlx.NamedExecContext(ctx, r.ext, saveProblemQuery, problemRow{
			ID:           problem.ID,
			Code:         problem.Code,
			Type:         problemType(problem.ProblemType),
			Title:        problem.Title,
			MaxStore:     problem.MaxScore,
			Category:     problem.Category,
			RedeployRule: redployRule(problem.RedeployRule),
		}); err != nil {
			if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == "problems_code_key" {
					return domain.NewAlreadyExistsError("problem code", nil)
				}
			}
			return errors.Wrap(err, "failed to save problem")
		}
	}
	if penalty := descriptiveProblem.Problem.PercentagePenalty; penalty != nil {
		if _, err := sqlx.NamedExecContext(ctx, r.ext,
			saveRedeployPercentagePenaltyQuery,
			struct {
				ProblemID uuid.UUID `db:"problem_id"`
				redeployPercentagePenaltyRow
			}{
				ProblemID: descriptiveProblem.Problem.ID,
				redeployPercentagePenaltyRow: redeployPercentagePenaltyRow{
					Threshold:  penalty.Threshold,
					Percentage: penalty.Percentage,
				},
			},
		); err != nil {
			return errors.Wrap(err, "failed to save percentage penalty")
		}
	}
	{
		content := descriptiveProblem.Content
		if _, err := sqlx.NamedExecContext(ctx, r.ext, saveProblemContentQuery,
			struct {
				ProblemID uuid.UUID `db:"problem_id"`
				problemContentRow
			}{
				ProblemID: descriptiveProblem.Problem.ID,
				problemContentRow: problemContentRow{
					Body:        content.Body,
					Explanation: content.Explanation,
				},
			},
		); err != nil {
			return errors.Wrap(err, "failed to save problem content")
		}
	}
	return nil
}

func (r *RepositoryTx) DeleteProblem(ctx context.Context, id uuid.UUID) error {
	if _, err := r.ext.ExecContext(ctx, "DELETE FROM problems WHERE id = $1", id); err != nil {
		return errors.Wrap(err, "failed to delete problem")
	}
	return nil
}

type (
	problemRow struct {
		ID           uuid.UUID   `db:"id"`
		Code         string      `db:"code"`
		Type         problemType `db:"type"`
		Title        string      `db:"title"`
		MaxStore     uint32      `db:"max_score"`
		Category     string      `db:"category"`
		RedeployRule redployRule `db:"redeploy_rule"`
	}
	redeployPercentagePenaltyRow struct {
		Threshold  uint32 `db:"threshold"`
		Percentage uint32 `db:"percentage"`
	}
	redeployPercentagePenaltyNullRow struct {
		Threshold  sql.Null[uint32] `db:"threshold"`
		Percentage sql.Null[uint32] `db:"percentage"`
	}
	problemContentRow struct {
		PageID      string `db:"page_id"`
		PagePath    string `db:"page_path"`
		Body        string `db:"body"`
		Explanation string `db:"explanation"`
	}
)

var (
	problemCols                   = columns([]string{"id", "code", "type", "title", "max_score", "category", "redeploy_rule"})
	redeployPercentagePenaltyCols = columns([]string{"threshold", "percentage"})
	problemContentCols            = columns([]string{"body", "explanation"})
)

func (r *problemRow) data() *domain.ProblemData {
	return &domain.ProblemData{
		ID:           r.ID,
		Code:         r.Code,
		ProblemType:  domain.ProblemType(r.Type),
		Title:        r.Title,
		MaxScore:     r.MaxStore,
		Category:     r.Category,
		RedeployRule: domain.RedeployRule(r.RedeployRule),
	}
}

func (r *redeployPercentagePenaltyNullRow) data() *domain.RedeployPenaltyPercentage {
	if !r.Threshold.Valid || !r.Percentage.Valid {
		return nil
	}
	return &domain.RedeployPenaltyPercentage{
		Threshold:  r.Threshold.V,
		Percentage: r.Percentage.V,
	}
}

func (r *problemContentRow) data() *domain.ProblemContentData {
	return &domain.ProblemContentData{
		Body:        r.Body,
		Explanation: r.Explanation,
	}
}

type problemType domain.ProblemType

var (
	_ sql.Scanner   = (*problemType)(nil)
	_ driver.Valuer = problemType(domain.ProblemTypeUnknown)
)

func (t *problemType) Scan(src any) error {
	*t = problemType(domain.ProblemTypeUnknown)
	if src == nil {
		return nil
	}
	v, ok := src.(string)
	if !ok {
		return nil
	}
	if v == "DESCRIPTIVE" {
		*t = problemType(domain.ProblemTypeDescriptive)
	}
	return nil
}

func (t problemType) Value() (driver.Value, error) {
	switch domain.ProblemType(t) {
	case domain.ProblemTypeDescriptive:
		return "DESCRIPTIVE", nil
	case domain.ProblemTypeUnknown:
		fallthrough
	default:
		return nil, errors.New("unknown problem type")
	}
}

type redployRule domain.RedeployRule

var (
	_ sql.Scanner   = (*redployRule)(nil)
	_ driver.Valuer = redployRule(0)
)

func (r *redployRule) Scan(src any) error {
	*r = redployRule(domain.RedeployRuleUnknown)
	if src == nil {
		return nil
	}
	v, ok := src.(string)
	if !ok {
		return nil
	}
	switch v {
	case "UNREDEPLOYABLE":
		*r = redployRule(domain.RedeployRuleUnredeployable)
	case "PERCENTAGE_PENALTY":
		*r = redployRule(domain.RedeployRulePercentagePenalty)
	case "MANUAL":
		*r = redployRule(domain.RedeployRuleManual)
	}
	return nil
}

func (r redployRule) Value() (driver.Value, error) {
	switch domain.RedeployRule(r) {
	case domain.RedeployRuleUnredeployable:
		return "UNREDEPLOYABLE", nil
	case domain.RedeployRulePercentagePenalty:
		return "PERCENTAGE_PENALTY", nil
	case domain.RedeployRuleManual:
		return "MANUAL", nil
	case domain.RedeployRuleUnknown:
		fallthrough
	default:
		return nil, errors.New("unknown redeploy rule")
	}
}
