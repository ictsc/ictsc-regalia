# rstate連携（二次予選モード）

`ICTSC_SECONDARY_QUALIFIER_MODE=true`のとき、再展開はSStateではなく
[ictsc/rstate](https://github.com/ictsc/rstate)（Terraformによる再展開ツール）へ依頼します。
公開API（`/api/v1`）と画面の動きは本選と同じで、変わるのはbackendと再展開サービスの間だけです。

## 設定

| 環境変数 | 内容 |
| -------- | ---- |
| `ICTSC_RSTATE_BASE_URL` | rstateのURL（例: `https://rstate.example`）。二次予選モードで必須 |
| `ICTSC_RSTATE_ADMIN_USER` | Basic認証のユーザー。既定`ictsc`（rstateの固定アカウント）。前段のproxyで別の認証をかける場合だけ変える |
| `ICTSC_RSTATE_ADMIN_PASSWORD` | rstateの`adminpass`。二次予選モードで必須 |
| `ICTSC_RSTATE_POLL_INTERVAL` | 状態を取りに行く間隔。既定`10s` |
| `ICTSC_RSTATE_ALLOW_HTTP` | `true`でrstateへのHTTP接続を許可（rstateは既定でHTTPのみで待ち受ける）。既定`false` |

二次予選モードでは`ICTSC_SSTATE_*`と`ICTSC_CALLBACK_BASE_URL`は不要です。

## 使うrstateのAPI

| 用途 | request | Regaliaの扱い |
| ---- | ------- | ------------- |
| 再展開の依頼 | `POST /admin/postJob`（Basic認証、form: `team_id`・`prob_id`） | `200`で受理し、応答のjob IDを再展開ごとに保存する |
| 状態の取得 | `GET /admin/list`（Basic認証） | 保存したjob IDのjobの状態を使う |

2024年にDiscord bot（[ictsc/ictsc-discord-bot](https://github.com/ictsc/ictsc-discord-bot)）が
再展開に使っていたのと同じ入口です。

- チームは`team_code`を2桁にした`team02`のような名前で送ります。rstateはこれをTerraform
  workspace名として使うので、rstate側のworkspaceも同じ名前にしてください。
- team_code は1〜99を送ります（SStateの2〜99とは違い、rstateに下限はありません）。
- 問題コードは小文字にして送ります（`A01`なら`a01`）。rstateはこれをTerraformの
  `module.<問題コード>`として使います。長さはrstateの制約に合わせて3〜8文字で、1〜2文字の
  問題は送る前に`FAILED`にします。
- `/admin/postJob`の`400`は、本文が`BadRequest!`ならチームか問題の指定が不正です。それ以外
  （本文はチームのstatusページのURL）は、同じ問題のjobが待機中か実行中、またはrstateが準備期間と
  判定したことを表します。Regaliaはこのとき`GET /admin/list`を確認し、そのチーム・問題のjobが
  待機中か実行中なら受理されたものとしてそのjob IDを保存します（応答が失われて再送した場合も
  ここに来ます）。無ければ準備期間として`FAILED`にします。URLにはtokenが含まれるので、本文は
  記録しません。

### rstate側の準備

rstateの準備期間判定（`pkg/utils/util.go`）は2022年の固定日時と比較しているため、そのままでは
運営チーム以外の`/admin/postJob`がすべて`400`になります。大会ごとに日時を書き換えてから
デプロイしてください（設定から読めるようにする件は
[ictsc/rstate#43](https://github.com/ictsc/rstate/issues/43)）。

`POST /backend/:teamid/:probcode`はBasic認証の外にあり、誰でも再展開を依頼できます。
Regaliaは使わないので、rstateは外部から届かない場所に置いてください。

## 状態の対応

| rstate | Regalia |
| ------ | ------- |
| `開始待ち` | `QUEUED`のまま |
| `実行中` | `DEPLOYING` |
| `終了` | `COMPLETED` |
| `エラー`、`Limit` | `FAILED` |

## ポーリング

rstateはcallbackを送らないため、APIは`ICTSC_RSTATE_POLL_INTERVAL`ごとに、`QUEUED`か
`DEPLOYING`の再展開だけをDBから読み、`GET /admin/list`を1回呼んで状態を進めます。処理中の
再展開が無いときはrstateを呼びません。

- 再展開とjobは、依頼時に保存したjob ID（`deployments.gateway_job_id`）で対応づけます。
  作成時刻では判定しないので、直前に終わった前回のjobを取り違えず、サーバー間の時計のずれにも
  左右されません。
- job IDがまだ無い再展開は、そのチーム・問題で待機中か実行中のjobだけを自分のjobとみなし、
  そのIDを保存します。rstateは同じチーム・問題のjobを同時に1つしか動かしません。
- `QUEUED`から直接`終了`や`エラー`を見たときは、`DEPLOYING`を挟んで記録します。
- eventの時刻はAPI側の時計、event IDは再展開と状態から決まる値です。
- DBに記録したあとRedisへの配信に失敗した更新は、次の周期で配信し直します。

運営画面のsyncも同じ処理で、`GET /admin/list`を1回呼んで状態を反映します。同じ状態のまま
もう一度syncしても何も変わりません。`QUEUED`のままrstateにjobが無い場合は`FAILED`にします。

## 既存DBの更新

job IDを保存する列は`backend/db/migrations/0008_deployment_job_id.sql`で追加します。既存DBには
API更新前に適用してください。
