# ictsc-regalia

ICTSCの競技者ダッシュボード、運営コンソール、スコアサーバーをまとめた
モノレポです。APIはOpenAPI契約ファーストで管理します。

## ディレクトリ構成

```text
.
├── backend/                 Go API、DB schema、外部サービスadapter
│   └── openapi.json         APIの唯一の正本
├── frontend/packages/
│   ├── api/                 OpenAPI生成型と共有REST client
│   ├── ui/                  共通Vue UI・デザイン・Storybook
│   ├── contestant/          競技者Nuxt SPA
│   ├── admin/               運営Nuxt SPA（/admin/）
│   └── config/              共通lint/format設定
├── docs/                    設計・外部連携契約
└── Taskfile.yaml            開発タスク
```

## 生成

OpenAPIを変更した後は、リポジトリルートで両言語の生成物を更新します。

```sh
task generate
```

生成後に差分が残る場合は、その差分もAPI契約の変更としてコミットします。

## 開発

依存サービスとAPIをDockerで起動します。

```sh
task dev-backend
```

詳細は [backend/README.md](backend/README.md)、
[frontend/README.md](frontend/README.md)、
[docs/architecture.md](docs/architecture.md) を参照してください。

## Discordロールによるチーム登録

`ICTSC_DISCORD_CONTESTANT_GUILD_ID` と `ICTSC_DISCORD_ROLE_TEAMS`
（例: `{"1547264645960695871":2}`）を設定すると、招待コードの代わりに
Discordロールから登録チームを決定します。管理APIでチームを先に作成してください。
登録画面にはチーム名を表示し、競技者名・表示名のみを入力します。
対応するロールがない場合、複数チームに対応する場合、登録済みチームと違う場合は拒否します。
上限人数と競技者名・Discordアカウントの重複チェックは招待登録と共通です。
マッピング未設定の場合は従来の招待コード登録を使用します。

参加者の入口ではチームロールを優先し、Staffとチームロールの両方がある場合も参加者として進みます。
Staffのみの場合は管理者セッションを発行して `/admin/` へ移動します。管理者専用の入口ではStaffで認証します。
管理者は `ICTSC_DISCORD_ADMIN_GUILD_ID` / `ICTSC_DISCORD_ADMIN_ROLE_IDS` で別途判定します。
OAuthのリダイレクトURIは `/api/v1/auth/discord/callback` と
`/api/v1/admin/auth/discord/callback` をDiscord側にも登録してください。
開発環境では `ICTSC_DEV_FAKE_MODE=true` のままDiscord資格情報を設定すると、
Discord認証のみ実接続になり、GitHubコンテンツ/SStateは無効のままです。
仮認証から切り替える場合は既存の仮セッションを無効にしてください。
