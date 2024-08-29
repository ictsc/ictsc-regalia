package oauth2

import (
	"context"
	"net/url"
)

// Provider OAuth2認証プロバイダーインターフェース
type Provider interface {
	// AuthorizationURL フロントエンドをリダイレクトさせる認可URLを発行する
	AuthorizationURL() (reqID string, url *url.URL)

	// GetToken OAuthプロバイダーが発行したcodeをアクセストークンに変換する
	GetToken(ctx context.Context, reqID, code string) (token string, err error)

	// GetIdentity アクセストークンを使用して、プロバイダー上のIdentityを取得する
	GetIdentity(ctx context.Context, token string) (identity string, err error)
}
