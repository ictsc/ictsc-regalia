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
| 再展開の依頼 | `POST /admin/postJob`（Basic認証、form: `team_id`・`prob_id`） | `200`で受理。`400`は確定した拒否として`FAILED`にする |
| 状態の取得 | `GET /admin/list`（Basic認証） | 該当チーム・問題で`created_time`が最新のjobを使う |

2024年にDiscord bot（[ictsc/ictsc-discord-bot](https://github.com/ictsc/ictsc-discord-bot)）が
再展開に使っていたのと同じ入口です。

- チームは`team_code`を2桁にした`team02`のような名前で送ります。rstateはこれをTerraform
  workspace名として使うので、rstate側のworkspaceも同じ名前にしてください。
- 問題コードは小文字にして送ります（`A01`なら`a01`）。rstateはこれをTerraformの
  `module.<問題コード>`として使います。長さはrstateの制約に合わせて3〜8文字です。
- `/admin/postJob`の`400`は、本文が`BadRequest!`ならチームか問題の指定が不正、それ以外
  （本文はチームのstatusページのURL）なら同じ問題の再展開が待機中か実行中、またはrstateが
  準備期間と判定したことを表します。URLにはtokenが含まれるので、Regaliaは本文を記録せず
  決まった文言で`FAILED`にします。

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

rstateはcallbackを送らないため、APIは`ICTSC_RSTATE_POLL_INTERVAL`ごとに、各チーム・問題の
最新revisionのうち`QUEUED`か`DEPLOYING`のものだけrstateへ問い合わせて状態を進めます。

- 依頼時刻より2分以上前に作られたjobは前回の再展開のものとして無視します。
- `QUEUED`から直接`終了`や`エラー`を見たときは、`DEPLOYING`を挟んで記録します。
- eventの時刻はAPI側の時計、event IDは再展開と状態から決まる値です。複数のAPI instanceが
  同時に記録しても重複しません。

運営画面のsyncも同じ`GET /admin/list`を1回呼んで状態を反映します。
