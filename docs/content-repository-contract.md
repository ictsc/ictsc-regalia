# GitHubコンテンツリポジトリ契約

問題、Section、お知らせの正本は設定済みGitHub repositoryです。Regaliaはbranchの内容を直接配信せず、
GitHub Actionsが指定した完全なcommit object IDから全fileを取得・検証し、成功したsnapshotだけを
有効化します。

## Repository layout

既定のmanifestは`content/manifest.yaml`です。manifest内のMarkdown pathはmanifestが置かれた
directory（既定では`content/`）からの相対pathです。

```text
content/
  manifest.yaml
  rule.md
  problems/A01.md
  explanations/A01.md
  announcements/welcome.md
```

絶対path、`..`、`.` segment、backslash、正規化前後で変わるpathは拒否します。symlinkやdirectoryは
利用せず、GitHub Contents APIが`type: file`として返すUTF-8 fileだけを取得します。

## Manifest version 1

```yaml
version: 1
sections:
  - slug: day1-am
    beginning: 2026-08-30T09:00:00+09:00
    ending: 2026-08-30T12:00:00+09:00
    problems: [A01]
problems:
  - code: A01
    title: Example Problem
    max_score: 100
    category: Network
    section_slug: day1-am
    type: DESCRIPTIVE
    body_path: problems/A01.md
    explanation_path: explanations/A01.md
    redeploy_rule:
      type: PERCENTAGE_PENALTY
      penalty_threshold: 1
      penalty_percentage: 10
announcements:
  - slug: welcome
    title: Welcome
    markdown_path: announcements/welcome.md
    effective_from: 2026-08-30T08:30:00+09:00
rule_path: rule.md
```

- unknown field、複数YAML document、inlineの問題本文・解説・お知らせ・ruleは拒否します。
- Section期間は`[beginning, ending)`で相互に重複できません。slugは全体で一意です。
- problem codeは全体で一意で、1文字から8文字の英数字・`_`・`-`です。各problemはちょうど1つの
  Sectionの`problems`へ記載し、`section_slug`も一致させます。
- `max_score`は1以上、現行のproblem `type`は`DESCRIPTIVE`です。
- `PERCENTAGE_PENALTY`は0以上の`penalty_threshold`と0から99の`penalty_percentage`を必須とします。
  `UNREDEPLOYABLE`と`MANUAL`では両fieldを`null`または省略します。
- announcement slugは一意で、`effective_from`はRFC 3339 date-timeです。
- manifestは1 MiB、各Markdownは2 MiB以下のvalid UTF-8とします。問題本文は空白だけにできません。

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
2. Contents APIの`ref`へ同じcommit IDを指定してmanifestを取得する。
3. manifest metadataと全参照pathを先に検証する。
4. 同じcommit IDから全Markdownを取得し、全て成功した場合だけsnapshotを返す。
5. active commitとの前後関係はCompare APIで確認し、古いまたはdivergeしたworkflowによる巻き戻しを拒否する。

取得途中の値は配信しません。DBへの原子的なactive切替はbackend service/repositoryが担当します。同じ
commitのrefreshは冪等です。GitHub障害や新commitのvalidation失敗時は保存済みの最終正常snapshotを
配信し続け、snapshotが一度もない場合だけcontent APIを利用不可にします。

問題の全fieldは競技開始後も変更できます。回答、採点、再展開は処理時の`content_commit`とmetadataへ
固定され、active catalogだけが現在ランキングの問題集合を決めます。過去commitはAdmin APIから参照します。

branch pollingは行いません。content更新の契機は上記GitHub Actionsからのrefresh HTTP requestだけです。

