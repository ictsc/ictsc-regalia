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

```sh
docker compose -f backend/compose.yaml up --build --wait
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
