# Protobuf

フロントエンドとバックエンド間で共有される API 定義

## 構成

- `admin/v1`: 運営用 API
- `contestant/v1`: 競技者用 API

## 開発方法

[Buf](https://buf.build/) を使っています。[Buf CLI のインストール手順](https://buf.build/docs/installation/)に従ってインストールしてください。

### フォーマット

```console
$ buf format -w
```

### リント

```console
$ buf lint
```
