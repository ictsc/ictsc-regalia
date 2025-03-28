package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jmoiron/sqlx"
)

type (
	invitationCodeRow struct {
		ID        uuid.UUID        `db:"id"`
		Team      *domain.TeamData `db:"-"`
		Code      string           `db:"code"`
		ExpiresAt time.Time        `db:"expires_at"`
		CreatedAt time.Time        `db:"created_at"`
	}
	invitationCodeWithTeamRow struct {
		invitationCodeRow
		Team teamRow `db:"t"`
	}
)

var invitationCodeColumns = columns([]string{"id", "code", "expires_at", "created_at"})

var selectInvitationCode = `
SELECT
	` + invitationCodeColumns.String("ic") + `,
	` + teamColumns.As("t") + `
FROM invitation_codes AS ic
INNER JOIN teams AS t ON ic.team_id = t.id`

var _ domain.InvitationCodeLister = (*repo)(nil)

func (r *repo) ListInvitationCodes(ctx context.Context, filter domain.InvitationCodeFilter) ([]*domain.InvitationCodeData, error) {
	cond := "TRUE"
	var args []any

	if filter.Code != "" {
		cond += " AND ic.code = ?"
		args = append(args, filter.Code)
	}

	rows, err := r.ext.QueryxContext(ctx, r.ext.Rebind(selectInvitationCode+" WHERE "+cond), args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*domain.InvitationCodeData{}, nil
		}
		return nil, domain.WrapAsInternal(err, "failed to select invitation_codes")
	}
	defer func() { _ = rows.Close() }()

	var ics []*domain.InvitationCodeData
	for rows.Next() {
		var row invitationCodeWithTeamRow
		if err := rows.StructScan(&row); err != nil {
			return nil, domain.WrapAsInternal(err, "failed to scan invitation_codes row")
		}

		ic := (*domain.InvitationCodeData)(&row.invitationCodeRow)
		ic.Team = (*domain.TeamData)(&row.Team)
		ics = append(ics, ic)
	}
	return ics, nil
}

var _ domain.InvitationCodeReader = (*repo)(nil)

func (r *repo) GetInvitationCode(ctx context.Context, codeString string) (*domain.InvitationCodeData, error) {
	var row invitationCodeWithTeamRow
	if err := sqlx.GetContext(ctx, r.ext, &row,
		r.ext.Rebind(selectInvitationCode+" WHERE ic.code = ? LIMIT 1"), codeString,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.WithStack(domain.ErrInvitationCodeNotFound)
		}
		return nil, domain.WrapAsInternal(err, "failed to select invitation_code")
	}

	ic := (*domain.InvitationCodeData)(&row.invitationCodeRow)
	ic.Team = (*domain.TeamData)(&row.Team)
	return ic, nil
}

var _ domain.InvitationCodeCreator = (*repo)(nil)

func (r *repo) CreateInvitationCode(ctx context.Context, code *domain.InvitationCodeData) error {
	if _, err := sqlx.NamedExecContext(ctx, r.ext, `
		INSERT INTO invitation_codes (id, team_id, code, expires_at, created_at)
		VALUES (:id, :t.id, :code, :expires_at, :created_at)
	`, invitationCodeWithTeamRow{
		invitationCodeRow: invitationCodeRow(*code),
		Team:              teamRow(*code.Team),
	}); err != nil {
		return domain.WrapAsInternal(err, "failed to insert invitation_code")
	}
	return nil
}
