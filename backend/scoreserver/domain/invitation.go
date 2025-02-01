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
	InvitationCode struct {
		id        uuid.UUID
		team      *Team
		code      string
		expiresAt time.Time
		createdAt time.Time
	}
	InvitationCodeInput struct {
		ID        uuid.UUID
		Team      *Team
		Code      string
		ExpiresAt time.Time
		CreatedAt time.Time
	}
)

func NewInvitationCode(input InvitationCodeInput) (*InvitationCode, error) {
	if input.ID.IsNil() {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("id is required"))
	}
	if input.Code == "" {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("code is required"))
	}
	if input.Team == nil {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("invitation code must belong to a team"))
	}
	if input.ExpiresAt == (time.Time{}) {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("expires_at is required"))
	}
	if input.CreatedAt == (time.Time{}) {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("created_at is required"))
	}
	if input.ExpiresAt.Before(input.CreatedAt) {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("expired before create"))
	}

	return &InvitationCode{
		id:        input.ID,
		team:      input.Team,
		code:      input.Code,
		expiresAt: input.ExpiresAt,
		createdAt: input.CreatedAt,
	}, nil
}

func (c *InvitationCode) ID() uuid.UUID {
	return c.id
}

func (c *InvitationCode) Team() *Team {
	return c.team
}

func (c *InvitationCode) Code() string {
	return c.code
}

func (c *InvitationCode) ExpiresAt() time.Time {
	return c.expiresAt
}

func (c *InvitationCode) CreatedAt() time.Time {
	return c.createdAt
}

const invitationCodeLength = 16

// 誤読しやすい文字を除外した文字セット
const invitationCodeCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func generateInvitationCode() (string, error) {
	charsetLen := big.NewInt(int64(len(invitationCodeCharset)))
	code := make([]byte, invitationCodeLength)
	for i := range code {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", errors.Wrap(err, "failed to generate random number")
		}
		code[i] = invitationCodeCharset[n.Int64()]
	}
	return string(code), nil
}

func createInvitationCode(team *Team, expiresAt time.Time, now time.Time) (*InvitationCode, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, NewError(ErrTypeInternal, errors.Wrap(err, "failed to generate uuid"))
	}

	code, err := generateInvitationCode()
	if err != nil {
		return nil, NewError(ErrTypeInternal, err)
	}

	return NewInvitationCode(InvitationCodeInput{
		ID:        id,
		Team:      team,
		Code:      code,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	})
}

// 招待コードの一覧を取得するワークフロー
type (
	InvitationCodeListWorkflow struct {
		Lister InvitationCodeLister
	}
)

func (w *InvitationCodeListWorkflow) Run(ctx context.Context) ([]*InvitationCode, error) {
	ics, err := w.Lister.ListInvitationCodes(ctx, InvitationCodeFilter{})
	if err != nil {
		return nil, err
	}
	return ics, nil
}

// 招待コードを作成するワークフロー
type (
	InvitationCodeCreateWorkflow struct {
		TeamGetter TeamGetter
		RunTx      TxFunc[InvitationCodeCreator]
		Clock      Clock
	}
	InvitationCodeCreateInput struct {
		TeamCode  int
		ExpiresAt time.Time
	}
)

func (w *InvitationCodeCreateWorkflow) Run(ctx context.Context, input InvitationCodeCreateInput) (*InvitationCode, error) {
	if input.TeamCode == 0 {
		return nil, NewError(ErrTypeInvalidArgument, errors.New("teamCode is required"))
	}
	teamCode, err := NewTeamCode(input.TeamCode)
	if err != nil {
		return nil, err
	}

	team, err := w.TeamGetter.GetTeamByCode(ctx, teamCode)
	if err != nil {
		return nil, err
	}

	invitationCode, err := createInvitationCode(team, input.ExpiresAt, w.Clock.Now())
	if err != nil {
		return nil, err
	}

	if err := w.RunTx(ctx, func(effect InvitationCodeCreator) error {
		return effect.CreateInvitationCode(ctx, invitationCode)
	}); err != nil {
		return nil, err
	}

	return invitationCode, nil
}

type (
	InvitationCodeFilter struct {
		Code string
	}
	InvitationCodeLister interface {
		ListInvitationCodes(ctx context.Context, filter InvitationCodeFilter) ([]*InvitationCode, error)
	}
	InvitationCodeCreator interface {
		CreateInvitationCode(ctx context.Context, invitationCode *InvitationCode) error
	}
)
