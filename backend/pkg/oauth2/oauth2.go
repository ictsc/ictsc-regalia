// Package oauth2 OAuth2認証プロバイダーインターフェース定義
package oauth2

// Provider OAuth2認証プロバイダーインターフェース
type Provider interface {
	AuthorizationURL() (reqID, url string)
	ExchangeCodeForToken(reqID, code string) (token string, err error)
	GetIdentity(token string) (identity string, err error)
}
