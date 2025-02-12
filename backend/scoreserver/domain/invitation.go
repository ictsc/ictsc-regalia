package domain

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
)

type (
	InvitationCodeString = invitationCodeString
	InvitationCode       = invitationCode
)

func (i InvitationCodeString) Code(ctx context.Context, eff InvitationCodeReader) (*InvitationCode, error) {
	data, err := eff.GetInvitationCode(ctx, string(i))
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get invitation code")
	}
	return data.parse()
}

func (i *InvitationCode) Team() *Team {
	return i.team
}

func (i *InvitationCode) Code() string {
	return string(i.code)
}

func (i *InvitationCode) ExpiresAt() time.Time {
	return i.expiresAt
}

func (i *InvitationCode) Expired(now time.Time) bool {
	return i.expiresAt.Before(now)
}

// ListInvitationCodes - 招待コードの一覧を取得する
func ListInvitationCodes(ctx context.Context, eff InvitationCodeLister) ([]*InvitationCode, error) {
	return listInvitationCodes(ctx, eff)
}

// CreateInvitationCode - 招待コードの作成
func (t *Team) CreateInvitationCode(
	ctx context.Context, eff InvitationCodeCreator, now time.Time, expiresAt time.Time,
) (*InvitationCode, error) {
	return t.createInvitationCode(ctx, eff, now, expiresAt)
}

type InvitationCodeCreateEffect interface {
	InvitationCodeCreator
}

type (
	InvitationCodeData struct {
		ID        uuid.UUID
		Team      *TeamData
		Code      string
		ExpiresAt time.Time
		CreatedAt time.Time
	}
	InvitationCodeFilter struct {
		Code string
	}
	InvitationCodeLister interface {
		ListInvitationCodes(ctx context.Context, filter InvitationCodeFilter) ([]*InvitationCodeData, error)
	}
	InvitationCodeReader interface {
		InvitationCodeLister
		GetInvitationCode(ctx context.Context, codeString string) (*InvitationCodeData, error)
	}
	InvitationCodeCreator interface {
		CreateInvitationCode(ctx context.Context, code *InvitationCodeData) error
	}
)

type (
	invitationCodeString string
	invitationCode       struct {
		id        uuid.UUID
		team      *Team
		code      invitationCodeString
		expiresAt time.Time
		createdAt time.Time
	}
)

func (i *invitationCode) Data() *InvitationCodeData {
	return &InvitationCodeData{
		ID:        i.id,
		Team:      i.team.Data(),
		Code:      string(i.code),
		ExpiresAt: i.expiresAt,
		CreatedAt: i.createdAt,
	}
}

func (d *InvitationCodeData) parse() (*invitationCode, error) {
	team, err := d.Team.parse()
	if err != nil {
		return nil, err
	}

	if d.ExpiresAt.Before(d.CreatedAt) {
		return nil, NewInvalidArgumentError("expired before create", nil)
	}
	return &invitationCode{
		id:        d.ID,
		team:      team,
		code:      invitationCodeString(d.Code),
		expiresAt: d.ExpiresAt,
		createdAt: d.CreatedAt,
	}, nil
}

func listInvitationCodes(
	ctx context.Context, eff InvitationCodeLister,
) ([]*invitationCode, error) {
	list, err := eff.ListInvitationCodes(ctx, InvitationCodeFilter{})
	if err != nil {
		return nil, WrapAsInternal(err, "failed to list invitation codes")
	}
	ics := make([]*InvitationCode, 0, len(list))
	for _, d := range list {
		ic, err := d.parse()
		if err != nil {
			return nil, err
		}
		ics = append(ics, ic)
	}
	return ics, nil
}

func (t *team) createInvitationCode(
	ctx context.Context,
	eff InvitationCodeCreator, now time.Time,
	expiresAt time.Time,
) (*invitationCode, error) {
	if expiresAt.Before(now) {
		return nil, NewInvalidArgumentError("already expired", nil)
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, WrapAsInternal(err, "failed to generate uuid")
	}

	code, err := generateInvitationCode()
	if err != nil {
		return nil, WrapAsInternal(err, "failed to generate invitation code")
	}

	invitationCode := &invitationCode{
		id:        id,
		team:      t,
		code:      code,
		expiresAt: expiresAt,
		createdAt: now,
	}
	if err := eff.CreateInvitationCode(ctx, invitationCode.Data()); err != nil {
		return nil, WrapAsInternal(err, "failed to create invitation code")
	}

	return invitationCode, nil
}

const (
	invitationCodeLength = 16
	// 誤読しやすい文字を除外した文字セット
	invitationCodeCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
)

func generateInvitationCode() (invitationCodeString, error) {
	charsetLen := big.NewInt(int64(len(invitationCodeCharset)))
	code := make([]byte, invitationCodeLength)
	for i := range code {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", errors.Wrap(err, "failed to generate random number")
		}
		code[i] = invitationCodeCharset[n.Int64()]
	}
	return invitationCodeString(code), nil
}
