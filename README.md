# ictsc-regalia

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
└── Taskfile (全体共通のタスク定義)
```

各フォルダ以下に関しては、それぞれのフォルダにある`README.md`を参照して下さい。

## proto の更新

protobuf のスキーマを更新した後は各言語ごとの定義を再生成する必要があります

```sh
task generate
```
