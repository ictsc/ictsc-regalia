# 競技者サインインフロー

最長のフローは以下．いわゆる OAuth2 Code Flow になっている
```mermaid
sequenceDiagram
    User->>UA: ログイン
    UA->>Backend: 認証要求
    Backend->>Backend: 認証リクエスト生成
    Backend->>UA: 認証URLへリダイレクト
    UA->>Discord: リダイレクト
    Discord-->>UA: 認証画面
    UA-->>User: 認証画面
    User->>UA: 認証する
    UA->>Discord: 認証
    Discord->>Backend: コールバック
    Backend->>Discord: トークン要求
    Discord-->>Backend: トークン
    Backend->>Discord: ユーザー情報取得
    Discord-->>Backend: ユーザー情報
    Backend->>Backend: セッション生成
    Backend-->>UA: セッション付きリダイレクト
    UA-->>User: ログイン後画面
```

## 認証リクエスト

エンドポイント `GET /api/auth/discord`

フロントエンドからはこのエンドポイントにナビゲーションする形になる

### クエリ
- (optional) `next_url` 認証完了時に遷移するパスを指定する．デフォルトは`/`

### 振舞い

1. 既にセッションが存在するなら`next_url`にリダイレクトして終了
2. セッションが存在しない or セッション切れであるならば認証リクエストを作成する
3. 認証中状態を表すセッション Cookie を付与してリダイレクトする．Discord から戻ってきたときにセッションを見るため，この Cookie は SameSite: Lax になっている必要がある

## 認証コールバック

エンドポイント `GET /api/auth/discord/callback`

### 振舞い

1. セッション Cookie と OAuth2 Code Flow のコールバックとしてのパラメータを検証する
2. トークンや Discord 上ユーザ情報などの情報を取得してユーザとの紐付けをする
3. 紐付けをもとにしてセッションを生成して`next_url`にリダイレクトする．このセッション Cookie は API リクエストからのみ使われるため SameSite: Strict になる
