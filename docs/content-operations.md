# コンテンツ反映の運用手順

問題・ルール・スケジュール・お知らせ・画像を、マスタリポジトリ（例: `ictsc/ictsc2026-pr2-master`）から
Regalia に反映するための手順です。リポジトリの書き方は [content-repository-contract.md](content-repository-contract.md)、
各マスタリポジトリの README を参照してください。

## 全体の流れ

```text
運営が main に push
  └─ GitHub Actions（regalia-refresh.yaml）
       └─ POST /api/v1/admin/content/actions/refresh  {"commit": "<40桁のSHA>"}
            Authorization: Bearer <GitHub Actions の OIDC token>
              └─ Regalia が token を検証
                   └─ GitHub API から同じ commit の中身を取得・検証
                        └─ 問題なければ公開中の内容を切り替え
```

Actions を使わずに、運営画面「コンテンツ」で commit SHA を指定して同じ処理を呼ぶこともできます。

## 誰が何を持っているか

| もの | 持ち主 | 用途 |
| ---- | ------ | ---- |
| マスタリポジトリの書き込み権限 | 作問・運営 | 内容を push する。push がそのまま反映になる |
| 運営画面のログイン（Discord） | 運営 | 状態の確認と手動反映 |
| GitHub Actions の OIDC token | workflow が実行のたびに発行 | Regalia への反映依頼（保存しない・有効期限は数分） |
| `ICTSC_GITHUB_API_TOKEN` | Regalia サーバー | private なマスタリポジトリを読む。**読み取り専用**で発行する |

OIDC token と API token は別物です。前者は「どの workflow からの依頼か」を証明し、後者は Regalia が中身を
取りに行くためだけに使います。依頼する側が中身を送ることはないので、依頼の経路から内容を差し替えることは
できません。

## 認証と認可

反映の依頼は1つのエンドポイント `POST /api/v1/admin/content/actions/refresh` で受けます。
次のどちらかで認証できます。

### GitHub Actions（OIDC）

`Authorization: Bearer <token>` を付け、cookie は送りません。Regalia は次をすべて確認します。

- 署名: issuer の JWKS（`ICTSC_GITHUB_OIDC_DISCOVERY_URL`、既定は issuer の `/.well-known/openid-configuration`）で RS256 を検証
- `iss` = `ICTSC_GITHUB_OIDC_ISSUER`（既定 `https://token.actions.githubusercontent.com`）
- `aud` に `ICTSC_GITHUB_OIDC_AUDIENCE` を含む
- `repository_id` = `ICTSC_GITHUB_REPOSITORY_ID`、`repository` = `ICTSC_GITHUB_REPOSITORY`
- `ref` = `ICTSC_GITHUB_REF`（例: `refs/heads/main`）。他のブランチからの依頼は拒否
- `job_workflow_ref` = `ICTSC_GITHUB_WORKFLOW_REF`。指定した workflow ファイル以外からの依頼は拒否
- `exp`・`nbf`・`iat`・`sub`

`repository_id` も照合するので、同名のリポジトリを作り直されても通りません。
この経路はブラウザの Origin 検証の対象外です（cookie を持たない machine request のため）。

### 運営画面（Discord ログイン）

`admin-session` cookie で認証します。ログインは Discord OAuth で、`ICTSC_DISCORD_ADMIN_GUILD_ID` の
サーバーに所属し、`ICTSC_DISCORD_ADMIN_ROLE_IDS` のいずれかのロールを持つ人だけが通ります。
ブラウザからの変更なので、`Origin` が `ICTSC_ALLOWED_ORIGINS` に含まれている必要があります。

運営画面から呼んだ場合、対象のリポジトリと ref は Regalia の設定値（`ICTSC_GITHUB_REPOSITORY`・
`ICTSC_GITHUB_REF`）に固定され、指定できるのは commit SHA だけです。

## 反映のときに Regalia が行うこと

どちらの経路でも同じです。

1. 反映中の依頼が他にあれば `409 content_refresh_in_progress` で断る。
2. 公開中の commit と同じなら、取得はせずにルールだけ `rule.md` に合わせて終わる。
3. 公開中の commit の**子孫でない** commit（古い commit、force push で履歴が分かれた commit）は
   `409 content_rollback_rejected` で断る。
4. GitHub API で commit・tree・各ファイル・画像を取得し、[契約](content-repository-contract.md)に沿って検証する。
   ひとつでも不正なら `422 content_invalid` で断り、公開中の内容はそのまま。
5. DB で公開中の内容を原子的に切り替える。
6. `schedule.yaml` の `ranking_freeze_at` を順位凍結時刻に反映する（書かれていなければ凍結を解除）。
7. `rule.md` があればルールページに反映する（運営画面で編集したルールは上書きされる）。
   `rule.md` がなければ今のルールを残す。

回答・採点・再展開は、その時点の commit に固定されます。問題を直して反映しても、過去の回答が
書き換わることはありません。

## 初期設定

### 1. GitHub

1. マスタリポジトリを private で作る。
2. Regalia が読むための token を作る（fine-grained token で、そのリポジトリの Contents: Read-only だけ）。
3. `.github/workflows/regalia-refresh.yaml` を置く（マスタリポジトリの雛形に含まれています）。
4. Settings → Secrets and variables → Actions → **Variables** に次を設定する（Secret ではなく Variable）。
   - `REGALIA_API_URL`: Regalia の URL（例: `https://contest.example.com`）
   - `REGALIA_OIDC_AUDIENCE`: 下の `ICTSC_GITHUB_OIDC_AUDIENCE` と同じ値
5. リポジトリ ID を控える: `gh api repos/<owner>/<repo> --jq .id`

### 2. Regalia

| 環境変数 | 値の例 |
| -------- | ------ |
| `ICTSC_GITHUB_REPOSITORY` | `ictsc/ictsc2026-pr2-master` |
| `ICTSC_GITHUB_REPOSITORY_ID` | `1418858779` |
| `ICTSC_GITHUB_REF` | `refs/heads/main` |
| `ICTSC_GITHUB_WORKFLOW_REF` | `ictsc/ictsc2026-pr2-master/.github/workflows/regalia-refresh.yaml@refs/heads/main` |
| `ICTSC_GITHUB_OIDC_AUDIENCE` | 大会ごとに決めた値（例: `regalia-2026-pr2`） |
| `ICTSC_GITHUB_API_TOKEN` | 1-2 で作った読み取り専用 token（Secret で管理） |
| `ICTSC_GITHUB_SCHEDULE_PATH` | 既定 `schedule.yaml` のままでよい |
| `ICTSC_SECONDARY_QUALIFIER_MODE` | 二次予選なら `true` |

設定は起動時に検証されます。値を変えたら API を再起動してください。

### 3. 疎通確認

1. マスタリポジトリの main に何か push する（README の修正など）。
2. Actions の「Regalia refresh」が成功することを確認する。
3. 運営画面「コンテンツ」で、配信状態が「配信中」、有効な commit が push した SHA になっていることを確認する。

## 普段の運用

- **内容を変える**: main に push するだけ。Actions の実行が終われば反映されています。
- **書きかけの問題**: frontmatter に `draft: true` を付けて push する（取り込まれない）。
- **公開前に確認したい**: 別ブランチで作業し、PR で確認してから main にマージする。main 以外の push は反映されません。
- **反映の状況を見る**: 運営画面「コンテンツ」、または Actions の実行ログ。

## 困ったとき

### 反映されない

運営画面「コンテンツ」の配信状態と、Actions の実行ログを見ます。

| 症状 | 原因 | 対処 |
| ---- | ---- | ---- |
| Actions の job が skip される | `REGALIA_API_URL` か `REGALIA_OIDC_AUDIENCE` が未設定 | Variables を設定して workflow を再実行 |
| `401 invalid_machine_token` | token の claim が設定と合わない（workflow のファイル名・ブランチ・audience・リポジトリ ID の違いが多い） | Regalia の `ICTSC_GITHUB_*` と workflow を照合する |
| `409 content_refresh_in_progress` | 別の反映が実行中 | 少し待って再実行 |
| `409 content_rollback_rejected` | 古い commit か、force push で履歴が分かれた | 下の「戻したい」を参照 |
| `422 content_invalid` | 内容が契約に合わない（未知の frontmatter 項目、code の重複、画像の置き場所など） | 運営画面「コンテンツ」の配信状態に出るエラー文を見て直し、push し直す |
| `502 upstream_unavailable` | GitHub か OIDC の発行元に繋がらない、または API token に読み取り権限がない | GitHub の状態と token を確認して再実行 |

失敗しても、公開中の内容（最終正常版）は出たままです。参加者側が壊れることはありません。
このとき配信状態には「最終正常版を配信中」とエラー文が出ます。エラー文は API のメモリに持っているため、API を再起動すると消えます。

### 間違えて反映した・戻したい

古い commit への巻き戻しは拒否されるので、**前に進めて戻します**。

```sh
git revert <間違えた commit>
git push origin main
```

force push で履歴を書き換えると、以後の反映がすべて `content_rollback_rejected` になります。
force push はしないでください。

### Actions が使えないとき（手動反映）

1. 運営画面「コンテンツ」を開く。
2. 「Git commit SHA」に反映したい commit の **40桁の SHA** を入れる（`git rev-parse HEAD`）。
3. 「指定commitを取得」を押す。

手動でも、巻き戻しの禁止と内容の検証は同じように働きます。

### ルールページを直したい

運営画面「競技設定」のルール編集でも直せますが、次に `rule.md` を含む commit を反映すると上書きされます。
恒久的な修正はマスタリポジトリの `rule.md` で行ってください。急ぐときは、`rule.md` を直して push するか、
今の commit を運営画面からもう一度反映すると、`rule.md` の内容に戻ります。

## 開発環境で試す

開発 fake mode（`ICTSC_DEV_FAKE_MODE=true`）では、`ICTSC_GITHUB_REPOSITORY` を設定したときだけ
GitHub から取り込みます。OIDC の項目も起動時に必須なので、手動反映だけを試す場合も仮の値を入れてください。
`ICTSC_ALLOW_INSECURE_UPSTREAMS=true` にすると、`ICTSC_GITHUB_API_BASE_URL` に http の URL を使えます。
