package contestant

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

type (
	SignUpRequest struct {
		InvitationCode string `json:"invitation_code"`
		Name           string `json:"name"`
		DisplayName    string `json:"display_name"`
	}
	SignUpResponse struct {
		Message string            `json:"message,omitempty"`
		Codes   []SignUpErrorCode `json:"codes,omitempty"`
	}
	SignUpErrorCode string
)

const (
	SignUpErrorCodeInvalidInvitationCode SignUpErrorCode = "invalid_invitation_code"
	SignUpErrorCodeTeamIsFull            SignUpErrorCode = "team_is_full"
	SignUpErrorCodeInvalidName           SignUpErrorCode = "invalid_name"
	SignUpErrorCodeDuplicateName         SignUpErrorCode = "duplicate_name"
	SignUpErrorCodeInvalidDisplayName    SignUpErrorCode = "invalid_display_name"

	signUpRequestLimit    = 1024
	signUpRateLimitWindow = 10 * time.Minute
	signUpRateLimitCount  = 3
)

func (h *AuthHandler) handleSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	signUpSess, err := session.SignUpSessionStore.Get(r.Context())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(r.Context(), w, SignUpResponse{Message: "session not found"})
			return
		} else {
			slog.ErrorContext(r.Context(), "failed to get signup session", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(r.Context(), w, SignUpResponse{})
			return
		}
	}
	if signUpSess.Discord == nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(r.Context(), w, SignUpResponse{Message: "no discord identity"})
		return
	}

	if isAllowed, err := h.RateLimiter.Check(r.Context(),
		"auth.signup:"+signUpSess.Discord.ID,
		signUpRateLimitWindow, signUpRateLimitCount,
	); err != nil {
		slog.ErrorContext(r.Context(), "failed to check rate limit", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(r.Context(), w, SignUpResponse{})
	} else if !isAllowed {
		w.WriteHeader(http.StatusTooManyRequests)
		writeJSON(r.Context(), w, SignUpResponse{Message: "rate limit exceeded"})
		return
	}

	var req SignUpRequest
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(r.Context(), w, SignUpResponse{Message: "Content-Type must be application/json"})
		return
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, signUpRequestLimit)).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(r.Context(), w, SignUpResponse{})
		return
	}
	_, _ = io.Copy(io.Discard, r.Body)

	userProfile, err := SignUp(r.Context(), h.SignUpEffect, time.Now(), &SignUpInput{
		InvitationCode: req.InvitationCode,
		Name:           req.Name,
		DisplayName:    req.DisplayName,
		DiscordID:      signUpSess.Discord.ID,
	})
	if err != nil {
		handleSignUpError(r.Context(), w, err)
		return
	}

	if err := session.SignUpSessionStore.Write(r, w, nil, h.signUpSessionOption(r)); err != nil {
		slog.ErrorContext(r.Context(), "failed to delete signup session", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(r.Context(), w, SignUpResponse{})
		return
	}

	userSess := &session.UserSession{UserID: uuid.UUID(userProfile.ID())}
	sessOpt := h.userSessionOption(r)
	if err := session.UserSessionStore.Write(r, w, userSess, sessOpt); err != nil {
		slog.ErrorContext(r.Context(), "failed to write user session", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(r.Context(), w, SignUpResponse{})
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(r.Context(), w, SignUpResponse{})
}

func handleSignUpError(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrInternal) {
		slog.ErrorContext(ctx, "failed to sign up", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(ctx, w, SignUpResponse{})
		return
	}
	var codes []SignUpErrorCode
	if errors.Is(err, domain.ErrInvalidUserName) {
		codes = append(codes, SignUpErrorCodeInvalidName)
	}
	if errors.Is(err, domain.ErrDuplicateUserName) {
		codes = append(codes, SignUpErrorCodeDuplicateName)
	}
	if errors.Is(err, domain.ErrInvalidDisplayName) {
		codes = append(codes, SignUpErrorCodeInvalidDisplayName)
	}
	if errors.Is(err, domain.ErrInvitationCodeExpired) || errors.Is(err, domain.ErrInvitationCodeNotFound) {
		codes = append(codes, SignUpErrorCodeInvalidInvitationCode)
	}
	if errors.Is(err, domain.ErrTeamIsFull) {
		codes = append(codes, SignUpErrorCodeTeamIsFull)
	}
	w.WriteHeader(http.StatusBadRequest)
	writeJSON(ctx, w, SignUpResponse{Codes: codes})
}

func writeJSON(ctx context.Context, w http.ResponseWriter, v any) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.ErrorContext(ctx, "failed to write json", "error", errors.WithStack(err))
		http.Error(w, "{}", http.StatusInternalServerError)
	}
}

type (
	// SignUpEffect ユーザー登録処理に必要な依存
	SignUpEffect   domain.Tx[SignUpTxEffect]
	SignUpTxEffect interface {
		domain.InvitationCodeReader
		domain.UserCreator
		domain.DiscordUserLinker
		domain.TeamMemberManager
	}
	SignUpInput struct {
		InvitationCode string
		Name           string
		DisplayName    string
		DiscordID      string
	}
)

func SignUp(ctx context.Context, effect SignUpEffect, now time.Time, input *SignUpInput) (*domain.UserProfile, error) {
	discordUserID, err := domain.NewDiscordID(input.DiscordID)
	if err != nil {
		return nil, err
	}
	codeString := domain.InvitationCodeString(input.InvitationCode)

	return domain.RunTx(ctx, effect, func(effect SignUpTxEffect) (*domain.UserProfile, error) {
		var (
			invitationCode *domain.InvitationCode
			userProfile    *domain.UserProfile
			errs           []error
		)
		{
			ic, err := codeString.Code(ctx, effect)
			if err != nil {
				errs = append(errs, err)
			}
			invitationCode = ic

			up, err := domain.CreateUser(ctx, effect, input.Name, input.DisplayName)
			if err != nil {
				slog.ErrorContext(ctx, "failed to create user", "error", err)
				errs = append(errs, err)
			}
			userProfile = up
		}
		if len(errs) > 0 {
			return nil, errors.Join(errs...)
		}

		if err := userProfile.LinkDiscord(ctx, effect, discordUserID); err != nil {
			return nil, err
		}
		if err := userProfile.JoinTeam(ctx, effect, now, invitationCode); err != nil {
			return nil, err
		}

		return userProfile, nil
	})
}
