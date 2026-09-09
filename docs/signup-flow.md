# 競技者登録フロー

Discord認証済みで競技者未登録の場合だけ、短期`signup-session`を使って登録する。

1. `GET /api/v1/viewer`が`DISCORD_AUTHENTICATED`を返すことを確認する。
2. 名前、表示名、運営が発行した招待コードを`POST /api/v1/auth/signup`へ送る。
3. PostgreSQL transactionで招待の有効期限・未使用・チーム定員、名前とDiscord IDの一意性を確認する。
4. 成功時は招待を消費して競技者を作成し、`signup-session`を消去して3日間の`user-session`を発行する。

同時登録は招待行とチーム行のlockおよび一意制約で直列化し、競合時はRFC 9457の
`409 application/problem+json`を返す。
