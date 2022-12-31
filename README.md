# ictsc_sachiko_v3

## 構成

- CSS : Tailwind CSS, daisyUI
- 状態管理: SWR

This is a [Next.js](https://nextjs.org/) project bootstrapped
with [`create-next-app`](https://github.com/vercel/next.js/tree/canary/packages/create-next-app).

## development

First, run the development server:

```bash
cp .env .env.local
# 適切にコンフィグを変更して下さい
npm run dev
# or
yarn dev
```

### Docker build

1. Docker をインストール
2. コンテナをビルド: `docker build -t nextjs-docker .`
3. コンテナを起動: `docker run -p 3000:3000 nextjs-docker`

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.
