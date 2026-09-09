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
