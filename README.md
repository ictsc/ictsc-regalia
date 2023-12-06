# ictsc-outlands

ICTSC Score Server (4th Gen)

## ディレクトリ構成

```plain
ictsc-outlands/
│
├── .github/ (Actions関連)
│
├── .vscode/ (VSCodeユーザーのための初期設定テンプレ)
│
├── dev/ (全体共通の開発用ツール・設定)
├── docs/ (ドキュメント)
│
├── frontend/ (フロントエンド)
├── backend/ (バックエンド)
│
├── proto/ (protocol buffersスキーマ定義)
│
├── charts/ (K8デプロイ用Helmチャート)
│
└── Makefile (全体共通のタスク定義)
```

各フォルダ以下に関しては、それぞれのフォルダにある`README.md`を参照して下さい。

## 開発の始め方

`sh init.sh`をすると、`task`のインストールを要求されるのでインストールする。  
その後は、`task`コマンドで各種操作を行う。
