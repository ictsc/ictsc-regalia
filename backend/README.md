# Backend

バックエンドは Go 言語で実装されたスコアサーバーです。競技参加者のスコアを管理し、リアルタイムな採点機能を提供します。フロントエンドとは Protocol Buffers/Connect を使用してコミュニケーションを行い、高性能で信頼性の高い API を提供します。

## プロジェクトの目的

- 競技参加者のスコアをリアルタイムに管理・表示
- 安全で高速な API 提供
- 運営スタッフによる効率的な採点作業の実現
- 詳細なモニタリングとロギングによる安定運用

## アーキテクチャ

```
backend/
├── cmd/                    # エントリーポイント
│   └── scoreserver/       # スコアサーバーのメイン実装
├── pkg/                   # 共有パッケージ
│   └── proto/            # 生成されたProtocol Buffersコード
├── scoreserver/          # スコアサーバーのコアロジック
│   ├── admin/           # 管理者向けAPI実装
│   ├── domain/         # ドメインモデル（ビジネスロジック）
│   └── infra/         # インフラストラクチャ層
│       └── pg/        # PostgreSQL実装
└── scripts/           # 開発用スクリプト
    ├── local-exec    # ローカル実行スクリプト
    ├── migrate       # DBマイグレーション
    └── psqlw         # PostgreSQL操作
```

## 技術スタック

- **言語**: Go
- **データベース**: PostgreSQL
- **API**: Protocol Buffers / Connect
- **監視**: OpenTelemetry
- **開発ツール**:
  - golangci-lint: リンター
  - buf: Protocol Buffers ツール

## セットアップ

### 必要要件

- Go 1.22 以上
- Docker

### ローカル開発環境の起動

1. 依存サービスの起動:

```console
$ docker compose up -d
```

これにより以下のサービスが起動します:

- PostgreSQL (任意のエフェメラル ポートにランダムに割り当てられます。./scripts/local-exec 経由でそのポートを自動的に参照しています。)
- Adminer (データベース管理 UI、8080 ポート)
- Jaeger (分散トレーシング、16686 ポート)

2. スコアサーバーの起動:

```console
$ ./scripts/local-exec go run ./cmd/scoreserver/ -dev
```

### データベースマイグレーション

スキーマの変更を適用するには:

```console
$ ./scripts/local-exec ./scripts/migrate
```

### 開発用コマンド

- `./scripts/psqlw`: PostgreSQL への接続
- `./scripts/local-exec`: ローカル環境での実行
- `./scripts/migrate`: データベースマイグレーション

## 主要コンポーネント

### スコアサーバー (`scoreserver/`)

- `domain/`: ドメインモデルとビジネスロジック

  - `team.go`: チーム関連のドメインモデル
  - `error.go`: カスタムエラー定義
  - `tx.go`: トランザクション管理

- `admin/`: 管理者向け API

  - `team.go`: チーム管理 API
  - `server_test.go`: サーバーテスト

- `infra/`: インフラストラクチャ層
  - `pg/`: PostgreSQL 実装
    - `repository.go`: リポジトリベース実装
    - `team.go`: チーム関連のデータアクセス

### 監視とトレーシング

OpenTelemetry を使用して以下の機能を提供:

- メトリクス収集
- 分散トレーシング
- ログ集約

Jaeger UI (http://localhost:16686) でトレースを確認できます。

### データベース

- PostgreSQL を使用
- スキーマは`schema.sql`で管理
- Adminer (http://localhost:8080) でデータベースを管理可能

## テスト

```console
$ go test ./...
```

## リンター

```console
$ golangci-lint run
```

設定は`.golangci.yaml`で管理されています。
