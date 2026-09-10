package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	"log"
	"time"
)

// seedDevContent only initializes an empty development store. Real snapshots
// and any later edits are preserved across restarts.
func seedDevContent(ctx context.Context, store core.Store) error {
	if _, err := store.ActiveContent(ctx); err == nil {
		return nil
	} else {
		var e *core.Error
		if !errors.As(err, &e) || e.Code != "content_not_available" {
			return err
		}
	}
	snapshot := devContentSnapshot()
	if err := service.ValidateSnapshot(snapshot); err != nil {
		return err
	}
	if _, err := store.ActivateContent(ctx, snapshot, core.NoActiveContentCommit); err != nil {
		return err
	}
	log.Print("development mock content activated: 3 problems")
	return nil
}
func devContentSnapshot() core.ContentSnapshot {
	beginning := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	ending := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	problems := []core.Problem{
		{Code: "M01", Title: "【モック】デフォルトルートを確認しよう", Category: "network", Body: "# 動作確認用のモック問題\n\n実機環境はありません。以下の情報から原因と復旧方針を記述してください。\n\n- PC: 192.0.2.10/24\n- ゲートウェイ: 192.0.2.1\n- 同一ネットワークには到達できるが、外部には到達できない\n- `ip route` には `192.0.2.0/24 dev eth0` だけが表示される\n\n## 解答してほしいこと\n不足している設定と、修正後に確認する項目を説明してください。", Explanation: "デフォルトルートが不足している。ゲートウェイへの到達性を確認してdefault via 192.0.2.1を設定し、外部への経路と疎通を確認する。"},
		{Code: "M02", Title: "【モック】Webサーバーの設定ミス", Category: "server", Body: "# 動作確認用のモック問題\n\n実機環境はありません。Webサーバーを再起動したところ、次のエラーで起動しなくなりました。\n\n```text\nnginx: [emerg] unexpected end of file, expecting \";\" or \"}\"\n```\n\n## 解答してほしいこと\n疑うべき設定ミス、設定の検査方法、復旧後の確認手順を説明してください。", Explanation: "設定ファイル末尾のセミコロンや閉じ括弧の不足を確認する。nginx -tで構文を検査してから起動し、HTTP応答を確認する。"},
		{Code: "M03", Title: "【モック】名前解決の切り分け", Category: "network", Body: "# 動作確認用のモック問題\n\n実機環境はありません。IPアドレスを指定するとWebサーバーへ接続できますが、ホスト名では接続できません。\n\n```text\n$ dig app.example.test\n;; communications error: timed out\n```\n\n## 解答してほしいこと\nDNSサーバーの設定、DNSサーバーへの通信、問い合わせ結果をどの順に確認するか説明してください。", Explanation: "クライアントのresolver設定を確認し、指定DNSサーバーに対して直接問い合わせる。UDP/TCP 53の到達性とサービス状態を調べ、正しいレコードとHTTP応答を確認する。"},
	}
	for i := range problems {
		p := &problems[i]
		p.MaxScore = 100
		p.Type = "DESCRIPTIVE"
		p.SectionSlug = "mock"
		p.Redeploy = core.RedeployRule{Type: core.RedeployUnredeployable}
	}
	manifest := core.Manifest{Version: 1, Sections: []core.Section{{Slug: "mock", Beginning: beginning, Ending: ending, ProblemIDs: []string{"M01", "M02", "M03"}}}, Problems: problems, Announcements: []core.Announcement{{Slug: "mock-notice", Title: "開発環境のモック問題です", Markdown: "画面表示・回答提出・採点の動作確認用です。本番競技の問題ではなく、VM環境は用意していません。", EffectiveFrom: beginning}}}
	encoded, _ := json.Marshal(manifest)
	return core.ContentSnapshot{CommitSHA: fmt.Sprintf("%x", sha256.Sum256(encoded)), Repository: "development/mock-content", Ref: "refs/heads/mock", Manifest: manifest, FetchedAt: time.Now().UTC()}
}
