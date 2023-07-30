# Octavio

## 構成

- CSS : Tailwind CSS, daisyUI
- 状態管理: SWR

This is a [Next.js](https://nextjs.org/) project bootstrapped
with [`create-next-app`](https://github.com/vercel/next.js/tree/canary/packages/create-next-app).

## development

First, run the development server:

```bash
# フォントのインストール
npm run font
#or
pnpm font

cp .env .env.local
# 適切にコンフィグを変更して下さい
npm run dev
# or
pnpm dev
```

#### テスト実行

```bash
pnpm test
```

### Storybook の起動

#### 1. Storybook の起動

```bash
pnpm storybook
```

#### 2. アクセス

[http://localhost:6006](http://localhost:6006)

### Docker build

1. Docker をインストール
2. コンテナをビルド:

```
docker build \
  --build-arg=next_public_api_url=http://localhost:8080/api \
  -t octavio .
```

3. コンテナを起動: `docker run -p 3000:3000 octavio`

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.
