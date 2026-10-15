# アーキテクチャ

## API契約

`backend/openapi.json`を唯一の正本とします。Go側はstrict Chi server、
TypeScript側は共有REST clientを生成し、手書きのhandler/usecaseと画面用modelを
生成境界の外側に置きます。APIの公開prefixは同一originの`/api/v1`です。

## 実行コンポーネント

- Contestant Nuxt SPA: 競技者向け画面。Discord認証、回答、再展開、ランキングを提供。
- Admin Nuxt SPA: `/admin/`で配信する運営画面。Discord guild roleで権限を判定。
- Go API: 認証、競技ロジック、コンテンツsnapshot、SState連携を担当。
- PostgreSQL: チーム、回答、採点、得点、再展開、監査履歴とcontent snapshotを保存。
- Redis: opaque sessionとOAuth state、複数API instance間のdeployment Pub/Subに使用。
- Discord: 競技者・運営のAuthorization Code + PKCE認証元。
- GitHub: 問題・Section・お知らせの正本。GitHub Actions OIDCでrefreshを起動。
- SState: 再展開実行サービス。APIがrequestし、SStateがHTTP callbackで状態を通知。

## 主要なデータフロー

1. GitHub Actionsがcommit SHA付きrefresh APIを呼びます。
2. APIはOIDC claimとmanifest/Markdownを検証し、成功時だけactive snapshotを切り替えます。
3. 回答・採点・再展開は処理時のcontent commitへ固定されます。
4. 再展開POSTをSStateが`202 Accepted`で受理し、進捗はBearer付きcallbackでAPIへ返します。
5. APIはcallbackをDBへ冪等記録し、Redis Pub/Sub経由でブラウザのSSEへ配信します。

バックグラウンドの再展開ポーリングは行いません。運営の`sync`操作だけが明示的な
一回限りの状態照会です。

## セキュリティ境界

- 競技者とAdminは別々のHttpOnly Cookie sessionを持ちます。
- Admin session発行時にDiscord guild membershipと許可roleを確認します。
- すべてのmutationでOriginを検証します。
- GitHub refreshはActions OIDC、SState callbackは固定Bearerで認証します。
- 問題本文は開始済みSectionだけ競技者へ返し、未来の予定APIは時刻情報だけを返します。
