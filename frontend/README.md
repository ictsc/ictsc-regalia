# Frontend

競技者用と運営用のNuxt 4 / Vue / TypeScriptアプリをpnpm workspaceで管理します。
両アプリは `ssr: false` のSPAです。Nodeサーバーは本番配信に不要です。

## 構成

- `packages/contestant/nuxt/`: 競技者画面。`pages`、`layouts`、`components`、`composables`、`middleware`、`features`。
- `packages/admin/nuxt/`: 運営画面。同じ構成で、本番base URLは `/admin/`。
- `packages/ui/`: 共通Vue部品、デザインCSS、チームカラーパレット、安全なMarkdown表示、Storybook。
- `packages/api/`: OpenAPI生成型、credentials付きRESTクライアント、RFC 9457エラー、mapper、SSE。
- `packages/config/`: 共通の開発設定。

デザインのコードはfrontend内で完結します。参照ディレクトリへのimportやコピー処理はありません。

問題文、解説、告知、ルール、回答などのMarkdownは共通でGitHub Flavored Markdown（GFM）として表示します。段落内で明示的に改行する場合は、行末に半角スペース2個かバックスラッシュを入れます。GFMに加えて `$...$` と `$$...$$` の数式記法を表示できます。

## 開発

Nodeは `.node-version`、pnpmは `package.json` の指定版を使用します。

```sh
pnpm install --frozen-lockfile
pnpm --filter @ictsc/competition dev
pnpm --filter @ictsc/admin dev
```

競技者用は `http://localhost:3000/`、運営用は `http://localhost:3001/admin/`。
開発時の `/api` はprefixを保持したまま `localhost:8080` へproxyします。

## ライブ開発プレビュー

このマシンでは、リポジトリルートから次のCompose設定で起動します。
Nuxt開発サーバーはDocker内部で動き、3000番が両画面とAPIの共通入口になります。
保存したフロントエンドの変更はブラウザへ自動反映されます。
初回は先に `frontend/` で `pnpm install --frozen-lockfile` を実行してください。

```sh
docker compose --env-file backend/.env.local -f backend/compose.yaml -f backend/compose.live-preview.yaml up -d --wait frontend
```

- 競技者画面: <http://100.96.0.22:3000/>
- 運営画面: <http://100.96.0.22:3000/admin/>
- API: <http://100.96.0.22:3000/api/v1/health>

Discord認証の戻り先も3000番のままです。Nuxtの内部ポート3002/3003は
外部公開しません。このマシン以外で実行するときは、必要に応じて
`ICTSC_DEV_UID` と `ICTSC_DEV_GID` に作業ユーザーのIDを設定してください。

## 静的な統合プレビュー

ビルド済みの画面を確認する場合は、ライブ用の追加Compose設定を外し、
リポジトリルートで次を実行します。

```sh
docker compose --env-file backend/.env.local -f backend/compose.yaml up -d --build --wait frontend
```

両画面のURLとAPIの経路はライブプレビューと同じです。静的プレビューでは
コード変更の反映にフロントイメージの再ビルドが必要です。

Discordを使う場合は、`backend/.env.local` とDiscord Developer Portalの
リダイレクトURLを次の2件に揃えます。

```text
http://100.96.0.22:3000/api/v1/auth/discord/callback
http://100.96.0.22:3000/api/v1/admin/auth/discord/callback
```

## API生成とチームカラー

`backend/openapi.json` が唯一のAPI契約です。変更後はリポジトリrootで `task generate` を実行し、GoとTypeScriptの生成物を更新します。frontend側だけの再生成は `pnpm generate` です。

チームカラーは管理画面のチーム編集で17色から選択できます。未設定時は `#A6E35F`。
既存DBでは、更新したAPIを起動する前に `backend/db/migrations/0002_team_color.sql` を適用してください。
新規開発DBではComposeがマイグレーションを番号順に適用します。

回答下書きは `regalia/draft/v1/<競技者>/<チーム>/<問題>` をキーにlocalStorageへ保存します。
提出後も残り、別の利用者には復元しません。端末間共有はありません。
旧形式は利用者を識別できないため自動移行せず、そのまま残します。

## 検査

```sh
pnpm ci:lint
pnpm ci:test
pnpm build
pnpm exec playwright install chromium
pnpm e2e
pnpm --filter @ictsc/ui story:build
```

PlaywrightはAPIをmockし、認証、回答、下書き、採点、カラー設定、SSE、モバイル表示を確認します。
スクリーンショットは `test-results/` に出力します。

## 静的配信

`pnpm build` で各アプリの `.output/public/` を生成します。DockerfileはそれぞれをNginxへ配置し、
`/admin/` と `/` の直接アクセスを各SPAのindexへフォールバックします。
APIは同一originの `/api/v1`、認証はCookie、再展開更新はSSEです。

## Storybook

```sh
pnpm --filter @ictsc/ui story
```

共通Vueコンポーネントを `http://localhost:6006/` で確認できます。

統合プレビューを検証する場合は `E2E_STATIC_ORIGIN=http://100.96.0.22:3000 pnpm e2e` を使用します。
