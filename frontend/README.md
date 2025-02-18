# Frontend

## 構成

```
frontend/
└── packages/
    ├── config/
    ├── proto/
    └── contestant/
        ├── app/
        │   ├── components/ # 複数のページで横断的に使われるロジックを持たないコンポーネント
        │   ├── features/   # 外界とのやりとりなどのロジック
        │   └── routes/     # 各ページの定義
        └── assets/
```

## 必要なもの

- [Node.js](https://nodejs.org)
- [pnpm](https://pnpm.io)

## 開発方法

### 競技者用ダッシュボード

#### 立ち上げ

`http:localhost:8080`でバックエンドサーバーが起動していることを前提とします。
以下のコマンドで`localhost:3000`で起動します。

```sh
cd packages/contestant
pnpm dev
```

#### Storybook

Storybook を起動するには
```sh
cd packages/contestant
pnpm story
```
