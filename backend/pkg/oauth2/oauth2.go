// Package oauth2 OAuth2認証プロバイダーインターフェース定義
package oauth2

import "context"

// Provider OAuth2認証プロバイダーインターフェース
type Provider interface {
	AuthorizationURL() (reqID, url string)
	ExchangeCodeForToken(ctx context.Context, reqID, code string) (token string, err error)
	GetIdentity(ctx context.Context, token string) (identity string, err error)
}
