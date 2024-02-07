# Frontend

## 構成

- `octavio`: 競技者用ダッシュボード
- `duardo`: 運営用ダッシュボード

## 開発方法

```bash
git submodule update --init --recursive

cd ictsc-outlands/frontend/octavio
pnpm font

cd ictsc-outlands
pnpm install

pnpm dev
```

- http://localhost:3000 で競技者用ダッシュボードが開きます。
- http://localhost:3000 で運営用ダッシュボードが開きます。