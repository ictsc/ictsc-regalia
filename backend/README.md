# Backend

`openapi.json` を唯一のAPI契約として実装したGo HTTP APIです。transportは
`oapi-codegen` のstrict Chi server、永続化はPostgreSQL、sessionとdeployment
SSEのPub/SubはRedisを使用します。

## 生成

OpenAPIを変更したときはリポジトリルートで実行します。

```sh
task generate
```

Go側の生成物は `internal/transport/api/api.gen.go`、TypeScript側は
`frontend/packages/api/src/schema.d.ts` です。生成後に同じコマンドを再実行し、
生成差分が増えないことを確認してください。

## Docker Composeで起動

`.env.example`をコピーして、ローカル用の設定を用意します。

```sh
cp backend/.env.example backend/.env.local
```

Discord認証を使用する場合は、`.env.local`にclient ID、client secret、
contestant/adminのredirect URIを設定します。`.env.local`はGitの管理対象外です。

```sh
docker compose --env-file backend/.env.local -f backend/compose.yaml up --build --wait
```

競技時間と再回答可能までの時間を固定するデモモードは、起動時に環境変数で
切り替えます。

```sh
ICTSC_DEMO_MODE=true docker compose --env-file backend/.env.local -f backend/compose.yaml up -d --build --wait
```

無効に戻す場合は`ICTSC_DEMO_MODE=false`で同じコマンドを実行します。この設定は
frontendコンテナの起動時に反映されるため、値を変えた場合はコンテナを再作成して
ください。URLのクエリでは切り替わりません。

frontendをローカルで起動する場合も同じ環境変数を使用できます。

```sh
cd frontend
ICTSC_DEMO_MODE=true pnpm --filter @ictsc/competition dev
```

- 競技者SPA: <http://localhost:3000/>
- Admin SPA: <http://localhost:3000/admin/>
- API: <http://localhost:8080/api/v1/health>

ComposeのAPIは開発fake modeで外部Discord/GitHub/SState adapterを無効化しますが、
PostgreSQLとRedisは実サービスを使用します。旧schema用volumeは削除せず、
`postgres-openapi-data` と `redis-openapi-data` を新規使用します。旧データの移行は
行いません。

依存サービスだけを起動し、APIをホスト上で実行する場合はリポジトリルートで
次を実行します。

```sh
task dev-backend-local
```

`scripts/dev-exec` がComposeの動的ポートから `ICTSC_DATABASE_URL` と
`ICTSC_REDIS_URL` を設定します。

## 設定

設定値は起動時に一括検証されます。雛形は `.env.example` にあります。本番では
`ICTSC_DEV_FAKE_MODE=false` とし、少なくとも次を設定します。

- PostgreSQL / Redis URL
- 許可するsame-originのOriginとSecure Cookie
- Discord client、callback URI、Admin guild/role ID
- GitHub Actions OIDCのissuer/audience/repository/ref/workflow refとAPI token
- SState request/callback用Bearer token、base URL、callback base URL

本番のURLはHTTPSのみ許可されます。mutationはブラウザsession利用時にOriginを
検証し、SState callbackだけは固定Bearerで認証します。

## テスト

```sh
cd backend
go test -race ./...
go vet ./...
```

PostgreSQL/Redis integration testはTestcontainersを使うため、Docker engineが必要です。
Dockerなしでunit/contract testだけを実行する場合は `go test -short ./...` を使います。

## チームカラーのマイグレーション

既存DBにはAPI更新前に `db/migrations/0002_team_color.sql` を適用します。新規Compose DBは `db/migrations/` のSQLを番号順に実行します。`Team.color` は管理画面から17色のパレットで設定し、未設定時は `#A6E35F` です。

PostgreSQL統合テストは通常Dockerを起動します。Dockerを使用できない場合は、空の使い捨てDBのURLを `ICTSC_TEST_DATABASE_URL` に設定して `go test ./internal/infra/postgres -count=1` を実行できます。指定DBに全マイグレーションとテストデータを作成します。

## Web Push

競技者向けのお知らせはWeb Pushで配信できます。VAPID鍵は一度だけ生成し、公開鍵と秘密鍵を同じ組み合わせのまま使います。

```sh
cd backend
go run ./cmd/vapid-keygen
```

出力された値を`ICTSC_WEB_PUSH_VAPID_PUBLIC_KEY`と`ICTSC_WEB_PUSH_VAPID_PRIVATE_KEY`へ設定し、連絡先を`ICTSC_WEB_PUSH_VAPID_SUBJECT`へ`mailto:admin@example.com`の形式で設定します。秘密鍵はSecretで管理し、Gitへ追加しません。

APIは公開時刻を迎えたお知らせを15秒ごとに確認し、購読後に公開されたものを各ブラウザへ一度だけ配信します。

## 保護された開発プレビュー

fake mode で Discord client ID/secret を空にし、admin guild/role ID と
contestant/admin redirect URI を設定すると、開発専用の仮ログインを利用できる。
外部公開時は必ず運営だけを許可する認証 Gateway で全経路を保護すること。
仮ユーザーは `preview-admin` / `preview-contestant` の共有IDで、実際のDiscord認証、
GitHubコンテンツ同期、SState操作は行わない。本番では利用しない。

## 参加者のDiscordサーバー所属制限

`ICTSC_DISCORD_CONTESTANT_GUILD_ID` にサーバーIDを設定すると、参加者ログインでも
`guilds.members.read` を要求し、そのサーバーへの所属を確認する。非所属やDiscord APIの
取得失敗ではログインさせない。空の場合は従来どおり参加者の所属を制限しない。
運営ログインは引き続き admin guild と運営ロールの両方を要求する。
これはログインの制限で、チーム紐付けや登録時の招待コード要件は変更しない。
所属はログイン時に確認するため、脱退しても発行済みセッションは期限まで有効。
