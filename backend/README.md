# Backend

## アーキテクチャ

```
backend/
├── cmd/              # エントリーポイント
│   └── scoreserver/  # スコアサーバー
├── pkg/              # 共有パッケージ
│   └── proto/        # 生成されたProtocol Buffersコード
├── scoreserver/      # スコアサーバーの実装
│   ├── domain/       # ドメインモデル
│   └── infra/        # インフラストラクチャ層
├── scripts/
│   ├── local-exec    # ローカル実行用スクリプト
│   ├── kube-exec     # リモート実行用スクリプト
│   └── migrate       # マイグレーション
├── schema.sql        # DB スキーマ
└── seed.sql          # 開発用 DB 初期データ
```

## 技術スタック

- **言語**: Go
- **データベース**: PostgreSQL
- **API**: Protocol Buffers / Connect
- **監視**: OpenTelemetry
- **開発ツール**:
  - golangci-lint: リンター

## ローカル開発

### 必要要件

- Go 1.22 以上
- Docker
- Discord App

### 環境変数

- `DISCORD_CLIENT_ID` Discord App の Client ID
- `DISCORD_CLIENT_SECRET` Discord App の Client Secret

### ローカル開発環境の起動

1. 依存サービスの起動:

```console
$ docker compose up -d
```

2. スコアサーバーの起動:

```console
$ ./scripts/local-exec go run ./cmd/scoreserver/ -dev
```

## 開発用スクリプト

### local-exec

ローカル開発環境の PostgreSQL と Redis にアクセスするための環境変数を入れてコマンドを実行します。

#### 例: psql

PostgreSQL のクライアント psql は PG* 系の環境変数を受け付けるためそのまま起動できます。

```console
$ ./scripts/local-exec psql
```

### kube-exec

Kubernetes 環境の依存にアクセスするための環境変数を入れてコマンドを起動します。
kubectl が必要です。

### migrate

スキーマの変更を反映します。内部的には [sqldef](https://github.com/sqldef/sqldef) のラッパーであるため，オプションなどはそちらを参照してください。

PG* 系の環境変数を受け付けるため，local-exec などと組み合わせて利用できます
```console
$ ./scripts/local-exec ./scripts/migrate
```

## テスト

```console
$ go test ./...
```

テストは [Testcontainers](https://testcontainers.com) を使って開発環境とは分離された環境で動作します。

## リント

```console
$ golangci-lint run
```

設定は`.golangci.yaml`で管理されています。

### テレメトリ

ローカル開発の場合，Jaeger UI (http://localhost:16686) でトレースを確認できます。
