package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

type invitationCodeRow struct {
	ID        uuid.UUID `db:"id"`
	Code      string    `db:"code"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *invitationCodeRow) asDomain(teamRow *teamRow) (*domain.InvitationCode, error) {
	team, err := teamRow.asDomain()
	if err != nil {
		return nil, err
	}
	return domain.NewInvitationCode(domain.InvitationCodeInput{
		ID:        r.ID,
		Team:      team,
		Code:      r.Code,
		ExpiresAt: r.ExpiresAt,
		CreatedAt: r.CreatedAt,
	})
}

var _ domain.InvitationCodeLister = (*repo)(nil)

func (r *repo) ListInvitationCodes(ctx context.Context, filter domain.InvitationCodeFilter) ([]*domain.InvitationCode, error) {
	cond := "TRUE"
	var args []any

	if filter.Code != "" {
		cond += " AND code = ?"
		args = append(args, filter.Code)
	}

	rows, err := r.ext.QueryxContext(ctx, r.ext.Rebind(`
		SELECT
			ic.id AS "ic.id",
			ic.code AS "ic.code",
			ic.expires_at AS "ic.expires_at",
			ic.created_at AS "ic.created_at",
			t.id AS "t.id",
			t.code AS "t.code",
			t.name AS "t.name",
			t.organization AS "t.organization",
			t.created_at AS "t.created_at",
			t.updated_at AS "t.updated_at"
		FROM invitation_codes AS ic
		LEFT JOIN teams AS t ON ic.team_id = t.id
		WHERE `+cond,
	), args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*domain.InvitationCode{}, nil
		}
		return nil, domain.WrapAsInternal(err, "failed to select invitation_codes")
	}
	defer func() { _ = rows.Close() }()

	var (
		row struct {
			Team           *teamRow           `db:"t"`
			InvitationCode *invitationCodeRow `db:"ic"`
		}
		invitationCodes []*domain.InvitationCode
	)
	for rows.Next() {
		if err := rows.StructScan(&row); err != nil {
			return nil, domain.WrapAsInternal(err, "failed to scan invitation_codes row")
		}

		invitationCode, err := row.InvitationCode.asDomain(row.Team)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid invitation_code; id=%v", row.InvitationCode.ID)
		}

		invitationCodes = append(invitationCodes, invitationCode)
	}
	return invitationCodes, nil
}

var _ domain.InvitationCodeCreator = (*repo)(nil)

func (r *repo) CreateInvitationCode(ctx context.Context, code *domain.InvitationCode) error {
	if _, err := r.ext.ExecContext(ctx, r.ext.Rebind(`
		INSERT INTO invitation_codes (id, team_id, code, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`), code.ID(), code.Team().ID(), code.Code(), code.ExpiresAt(), code.CreatedAt()); err != nil {
		return domain.WrapAsInternal(err, "failed to insert invitation_code")
	}
	return nil
}
