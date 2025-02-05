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

var _ domain.InvitationCodeLister = (*repo)(nil)

func (r *repo) ListInvitationCodes(ctx context.Context, filter domain.InvitationCodeFilter) ([]*domain.InvitationCodeData, error) {
	cond := "TRUE"
	var args []any

	if filter.Code != "" {
		cond += " AND code = ?"
		args = append(args, filter.Code)
	}

	rows, err := r.ext.QueryxContext(ctx, r.ext.Rebind(`
		SELECT
			ic.id,
			ic.code,
			ic.expires_at,
			ic.created_at,
			t.id AS "t.id",
			t.code AS "t.code",
			t.name AS "t.name",
			t.organization AS "t.organization"
		FROM invitation_codes AS ic
		LEFT JOIN teams AS t ON ic.team_id = t.id
		WHERE `+cond,
	), args...)
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
