# GitHubコンテンツリポジトリ契約

設定・反映・障害対応の手順は [content-operations.md](content-operations.md) を参照してください。

問題、Section、お知らせの正本は設定済みGitHub repositoryです。Regaliaはbranchの内容を直接配信せず、
GitHub Actionsが指定した完全なcommit object IDから全fileを取得・検証し、成功したsnapshotだけを
有効化します。

## Repository layout

ICTSCのマスタリポジトリ（例: `ictsc2025-final-master`）と同じ配置を読みます。パスはリポジトリ直下からです。
`schedule.yaml` の場所だけは `ICTSC_GITHUB_SCHEDULE_PATH`（既定 `schedule.yaml`）で変更できます。

```text
schedule.yaml
rule.md                          （任意）
problems/<name>.problem.md       （frontmatter + 問題文）
problems/<name>.explanation.md   （任意。同じ <name> の問題の解説）
notices/<slug>.md                （frontmatter + 本文。_template.md は読みません）
images/<file>                    （png / jpg / jpeg / gif / webp）
```

上記以外のファイル（`sync/`、`image-downloader/`、README など）は無視します。指定commitのGit treeを
再帰的に列挙し、`problems/`・`notices/`・`images/` は直下のファイルだけを対象にします。

## schedule.yaml

```yaml
ranking_freeze_at: 2026-03-15T12:00:00+09:00   # 任意
schedules:
- name: day1-am
  startAt: 2026-03-14T10:55:00+09:00
  endAt: 2026-03-14T13:00:00+09:00
```

- `name` は一意で、英数字・`_`・`-` です。期間は `[startAt, endAt)` で、互いに重複できません。
- `ranking_freeze_at` が正本です。コンテンツ更新時に順位凍結時刻へ反映し、書かれていなければ凍結を解除します。
- 二次予選モードでは `schedules` はちょうど1つで、`ranking_freeze_at` は書けません。

## rule.md

`rule.md` がある commit を有効化すると、その内容を参加者のルールページに反映します（運営画面の競技設定で
編集したルールは上書きされます）。`rule.md` がない commit では、現在のルールをそのまま残します。

## 問題

```markdown
---
code: 1234
title: 前任者の悪意
max_score: 100
category: Server
redeploy_rule:
  type: PERCENTAGE_PENALTY
  penalty_threshold: 1
submissionable_schedules: [day2]
---

問題文
```

- `code` は大文字に変換し、1〜8文字の英数字・`_`・`-` です。`0000`（例題）と `draft: true` の問題は取り込みません。
- `redeploy_rule.type` は `UNREDEPLOYABLE`・`PERCENTAGE_PENALTY`・`MANUAL` です。
  - `redeploy_rule` を省略すると、`PERCENTAGE_PENALTY`・しきい値999・減点0%（実質無制限）になります。
  - `PERCENTAGE_PENALTY` で `penalty_percentage` を省略（または0）すると10%、`penalty_threshold` を省略すると0です。
  - `UNREDEPLOYABLE` と `MANUAL` では `penalty_*` は無視します。
- `submissionable_schedules` は通常モードでは必須で、現在は1つだけ指定できます。二次予選モードでは書けず、
  全問題がその1つのスケジュールに属します。
- 任意で `default_grader_discord_id`（Discord IDの数字文字列）を指定できます。
- `answer_max_length` は回答の最大文字数です（1〜100000、省略時3000）。改行は数えず、Unicodeのコードポイント数で数えます。
- `allow_attachments: true` の問題だけ、回答にテキストのログを添付できます（省略時は不可）。
- 問題文は空白だけにできません。未知のfrontmatter fieldは拒否します。

## お知らせ

```markdown
---
title: "お知らせタイトル"
from: "2026-03-14T10:00:00+09:00"
---

本文
```

ファイル名（`.md` を除く）がslugです。`from` はRFC 3339のdate-timeで、この時刻から表示します。

## 画像

`images/` 直下の画像（png / jpg / jpeg / gif / webp、ファイル名は英数字・`.`・`_`・`-`）を取り込み、Markdownからは
`/images/<file>` として参照します。サブディレクトリや未対応の形式のファイルが `images/` にあると、取り込みを
エラーにします。1ファイル10 MiB、合計100 MiBまでです。

- 画像の中身はGit blob SHA単位で1回だけ保存し、各commitはファイル名とblobの対応だけを持ちます。
- 競技者画面では `/api/v1/contestant/images/<file>` に置き換えます。返すのは、開始済みの問題文・表示中の
  お知らせ・ルールから参照されている画像だけで、開始前の問題や解説だけで使う画像は404です。
- 運営画面では `/api/v1/admin/images/<file>` に置き換え、過去commitの問題を表示するときはそのcommitの画像を
  `?commit=` で取得します。
- 応答には `ETag`（blob SHA）と `X-Content-Type-Options: nosniff` を付け、`If-None-Match` が一致すれば304を返します。

## 共通

- 各Markdownは2 MiB、`schedule.yaml` は1 MiB以下のvalid UTF-8とします。
- frontmatterはファイル先頭の `---` で始まり `---` で閉じるYAMLで、複数documentは拒否します。

## GitHub Actions OIDC

content refreshはGitHub Actionsが`id-token: write`で取得したOIDC tokenをBearerとして
`POST /api/v1/admin/content/actions/refresh`へ送ります。Regaliaは署名をdiscovery/JWKSからRS256で
検証し、次のclaimを設定値と完全一致させます。

- `iss`: `https://token.actions.githubusercontent.com`
- `aud`: Regalia用の専用audience
- `repository_id`: 設定済みrepository ID
- `repository`: 設定済み`owner/name`
- `ref`: 設定済みpublish ref（例: `refs/heads/main`）
- `job_workflow_ref`: 許可したworkflow fileとref
- `exp`、`nbf`、`iat`、`sub`

workflowではcheckout結果の短縮SHAではなく`${{ github.sha }}`の完全な小文字commit SHAをbodyの
`commit`へ指定します。OIDC tokenとGitHub API read tokenは用途が異なるため、同じsecretとして
扱いません。

## Exact-commit取得と有効化

Regalia adapterは次の順でGitHub APIを呼びます。

1. Git commit endpointで指定object IDが存在し、response SHAが完全一致することを確認する。
2. Git tree APIで同じcommit IDのファイル一覧を取得する。
3. Contents APIの`ref`へ同じcommit IDを指定して `schedule.yaml`・`rule.md`・問題・お知らせを取得し、検証する。
4. 画像はtreeのblob SHAからGit blob APIで取得し、全て成功した場合だけsnapshotを返す。
5. active commitとの前後関係はCompare APIで確認し、古いまたはdivergeしたworkflowによる巻き戻しを拒否する。

取得途中の値は配信しません。DBへの原子的なactive切替はbackend service/repositoryが担当します。同じ
commitのrefreshは冪等です。GitHub障害や新commitのvalidation失敗時は保存済みの最終正常snapshotを
配信し続け、snapshotが一度もない場合だけcontent APIを利用不可にします。

問題の全fieldは競技開始後も変更できます。回答、採点、再展開は処理時の`content_commit`とmetadataへ
固定され、active catalogだけが現在ランキングの問題集合を決めます。過去commitはAdmin APIから参照します。

branch pollingは行いません。content更新の契機は上記GitHub Actionsからのrefresh HTTP requestだけです。

