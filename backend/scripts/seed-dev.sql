-- Development-only sample data. Run with psql -v ON_ERROR_STOP=1.
-- An existing contest is never overwritten.
BEGIN;

DO $$
DECLARE
    sample_commit constant text := 'dddddddddddddddddddddddddddddddddddddddd';
    sample_manifest jsonb;
BEGIN
    IF EXISTS (SELECT 1 FROM teams)
       OR EXISTS (SELECT 1 FROM content_snapshots)
       OR EXISTS (SELECT 1 FROM answers) THEN
        RAISE EXCEPTION 'seed-dev requires an empty contest database; existing data was left untouched';
    END IF;

    sample_manifest := jsonb_build_object(
        'version', 1,
        'sections', jsonb_build_array(
            jsonb_build_object('slug', 'practice', 'beginning', now() - interval '1 day',
                'ending', now() + interval '30 days', 'problem_codes', jsonb_build_array('NET1', 'WEB1', 'DB1', 'OPS1'))
        ),
        'problems', jsonb_build_array(
            jsonb_build_object('code', 'NET1', 'title', 'DNS の名前解決', 'max_score', 100,
                'category', 'Network', 'section_slug', 'practice', 'type', 'DESCRIPTIVE',
                'body', '社内 DNS で example.test が引けません。原因と修正方法を説明してください。',
                'explanation', 'ゾーン定義とリゾルバの参照先を確認します。',
                'redeploy_rule', jsonb_build_object('type', 'PERCENTAGE_PENALTY', 'penalty_threshold', 1, 'penalty_percentage', 10)),
            jsonb_build_object('code', 'WEB1', 'title', 'HTTP 502 の調査', 'max_score', 150,
                'category', 'Web', 'section_slug', 'practice', 'type', 'DESCRIPTIVE',
                'body', 'リバースプロキシが 502 を返しています。調査手順と復旧方法を示してください。',
                'explanation', 'upstream の接続先とサービスの稼働状態を調べます。',
                'redeploy_rule', jsonb_build_object('type', 'MANUAL')),
            jsonb_build_object('code', 'DB1', 'title', '遅い SQL クエリ', 'max_score', 200,
                'category', 'Database', 'section_slug', 'practice', 'type', 'DESCRIPTIVE',
                'body', '一覧取得が遅くなりました。実行計画を使った調査と改善案を記述してください。',
                'explanation', '実行計画と索引の利用状況を確認します。',
                'redeploy_rule', jsonb_build_object('type', 'UNREDEPLOYABLE')),
            jsonb_build_object('code', 'OPS1', 'title', '監視アラートの設計', 'max_score', 100,
                'category', 'Operations', 'section_slug', 'practice', 'type', 'DESCRIPTIVE',
                'body', 'API の障害を検知するためのメトリクスと通知条件を提案してください。',
                'explanation', 'エラー率、レイテンシ、稼働率を組み合わせます。',
                'redeploy_rule', jsonb_build_object('type', 'MANUAL'))
        ),
        'announcements', jsonb_build_array(
            jsonb_build_object('slug', 'welcome', 'title', 'デバッグ用データへようこそ',
                'markdown', 'これはローカル開発用のサンプル大会です。', 'effective_from', now() - interval '1 day')
        ),
        'rule_markdown', '## サンプル大会\n回答・採点・ランキング表示の確認に利用してください。'
    );

    INSERT INTO teams (code, name, organization, member_limit, color) VALUES
        (2, 'チーム・アルファ', 'サンプル大学', 4, '#0083C3'),
        (3, 'チーム・ベータ', 'デモ工業高専', 4, '#EE536B'),
        (4, 'チーム・ガンマ', 'テスト研究会', 4, '#00A95C');

    INSERT INTO contestants (name, display_name, self_introduction, discord_id, team_code) VALUES
        ('preview-contestant', 'プレビュー参加者', '開発用の仮ログインで利用できます。', '900000000000000001', 2),
        ('alice', 'アリス', 'ネットワークを担当しています。', '900000000000000011', 2),
        ('bob', 'ボブ', 'Web を担当しています。', '900000000000000012', 3),
        ('carol', 'キャロル', 'DB を担当しています。', '900000000000000013', 4);

    INSERT INTO content_snapshots (commit_sha, repository, git_ref, manifest, fetched_at, activated_at, active)
    VALUES (sample_commit, 'local/dev-seed', 'refs/heads/dev', sample_manifest, now(), now(), true);

    INSERT INTO answers (team_code, problem_code, number, author_name, body, submitted_at,
                         content_commit, max_score, redeploy_rule, deployments_before) VALUES
        (2, 'NET1', 1, 'preview-contestant', 'リゾルバの参照先とゾーン設定を確認します。', now() - interval '4 hours', sample_commit, 100,
            '{"type":"PERCENTAGE_PENALTY","penalty_threshold":1,"penalty_percentage":10}', 0),
        (2, 'WEB1', 1, 'alice', 'upstream のプロセスとポートを確認し、設定を修正します。', now() - interval '3 hours', sample_commit, 150,
            '{"type":"MANUAL"}', 0),
        (3, 'NET1', 1, 'bob', 'DNS サーバーの応答とレコードを調査します。', now() - interval '5 hours', sample_commit, 100,
            '{"type":"PERCENTAGE_PENALTY","penalty_threshold":1,"penalty_percentage":10}', 0),
        (3, 'DB1', 1, 'bob', 'EXPLAIN ANALYZE で実行計画を調べます。', now() - interval '2 hours', sample_commit, 200,
            '{"type":"UNREDEPLOYABLE"}', 0),
        (4, 'OPS1', 1, 'carol', '5xx エラー率と p95 レイテンシを監視します。', now() - interval '90 minutes', sample_commit, 100,
            '{"type":"MANUAL"}', 0);

    INSERT INTO answer_counters (team_code, problem_code, next_number, last_submitted_at)
    SELECT team_code, problem_code, max(number) + 1, max(submitted_at)
    FROM answers GROUP BY team_code, problem_code;

    INSERT INTO marking_results (team_code, problem_code, answer_number, judge, marked_score, rationale, created_at, visibility) VALUES
        (2, 'NET1', 1, 'preview-admin', 80, '原因の特定と修正方針が明確です。', now() - interval '3 hours', 'PUBLIC'),
        (3, 'NET1', 1, 'preview-admin', 60, '確認手順は適切ですが修正方針を補足してください。', now() - interval '4 hours', 'TEAM'),
        (3, 'DB1', 1, 'preview-admin', 120, '実行計画を使った調査ができています。', now() - interval '1 hour', 'PUBLIC');

    UPDATE competition_state
    SET rule_markdown = '## サンプル大会\n回答・採点・ランキング表示の確認に利用してください。',
        updated_at = now(), updated_by = 'seed-dev'
    WHERE singleton;
END $$;

COMMIT;
