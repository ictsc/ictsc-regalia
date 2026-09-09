# SState HTTP連携契約

この文書はRegalia backendと別リポジトリのSState間のHTTP契約です。旧Basic認証、
`201 Created`、定期的なstatus pollingは廃止します。時刻はRFC 3339、IDとstatusは
大文字・小文字を含め完全一致で扱います。

## 認証とtransport

- RegaliaからSStateへの全requestは`Authorization: Bearer <SSTATE_API_BEARER>`を送ります。
- SStateからRegaliaへのcallbackは`Authorization: Bearer <SSTATE_CALLBACK_BEARER>`を送ります。
  2つのtokenは別々に設定・rotationできる固定共有secretです。
- productionではHTTPSを必須とし、tokenやcallback bodyをlogへ残しません。
- JSON requestには`Content-Type: application/json`、JSON responseには
  `Content-Type: application/json`を設定します。

## 再展開要求

`POST /redeploy`

```json
{
  "request_id": "123e4567-e89b-12d3-a456-426614174000",
  "team_code": 12,
  "problem_code": "A01",
  "revision": 3,
  "content_commit": "0123456789abcdef0123456789abcdef01234567",
  "callback_url": "https://score.example/api/v1/admin/deployments/12/A01/3/events"
}
```

- `request_id`はRegaliaが発行するUUIDで、SState側の冪等性keyです。
- `team_code`は2から99、`problem_code`は`^[A-Za-z0-9_-]{1,8}$`、`revision`は1以上です。
- `content_commit`は要求時の問題metadataを固定する40から64文字の小文字hex object IDです。
- `callback_url`はこのrevision専用です。SStateはhostを組み替えず、そのURLへcallbackします。
- SStateが要求を永続化した場合だけ`202 Accepted`を返します。Regalia clientは他の2xxを含む
  202以外を失敗として扱います。
- 同じ`request_id`と同じbodyの再送は副作用なしで`202`、同じ`request_id`で異なるbodyは
  `409 Conflict`とします。

## 状態callback

SStateは状態が変化するたびに、要求で受け取った`callback_url`へ`POST`します。

```json
{
  "event_id": "018f47d2-7f3d-7c6d-8a2a-4d96cd884a01",
  "occurred_at": "2026-08-30T03:04:05Z",
  "status": "DEPLOYING",
  "message": "creating instances"
}
```

- `event_id`はcallback再送後も変えないUUIDです。
- callbackで送れるstatusは`DEPLOYING`、`COMPLETED`、`FAILED`です。`QUEUED`はRegaliaが要求作成時にだけ生成します。
- 通常遷移は`QUEUED -> DEPLOYING -> COMPLETED|FAILED`です。callbackの再送は同一event IDで
  冪等化されます。障害復旧で`FAILED -> DEPLOYING`を行う場合だけRegaliaの運営syncを
  recoveryとして実行します。その後は通常遷移へ戻り、`COMPLETED`は終端です。
- Regaliaは同じ`event_id`・同じpayloadへ`200 OK`、同じID・異なるpayloadへ`409 Conflict`を返します。
- timeoutまたは5xxでは同じeventを指数backoffで再送します。4xxは設定またはcontract不一致として
  自動再送を止め、運営へ通知します。

callbackの受理後、browserへの反映はRegaliaのRedis Pub/SubとSSEで行います。SStateがbrowserへ
接続する必要はありません。

## 運営による一回限りのsync

`GET /status/{team_code_2digits}/{problem_code}`

例: `GET /status/02/A01`

通常の進捗通知はcallbackを使います。このendpointは障害調査時にAdminが明示的にsyncを押した
1回だけ呼び出され、background job、timer、browser pollingからは呼びません。

```json
{
  "event_id": "018f47d2-7f3d-7c6d-8a2a-4d96cd884a01",
  "occurred_at": "2026-08-30T03:04:05Z",
  "status": "DEPLOYING",
  "message": "creating instances"
}
```

返すのはSStateが永続化した最新eventです。同じ最新状態には同じ`event_id`と`occurred_at`を返し、
sync自体で新しいeventを生成しません。対象がなければ`404 Not Found`です。

## Contract testの受け入れ条件

- `POST /redeploy`のfield、Bearer、202-only動作を相互に検証する。
- callbackの固定Bearer、event ID冪等性、status遷移を検証する。
- `GET /status/...`は1操作につきHTTP request 1回だけで、定期requestが発生しないことを検証する。
- unknown JSON field、無効なUUID・日時・statusは成功扱いにしない。

