TRUNCATE TABLE rules CASCADE;
INSERT INTO rules (page_path, markdown) VALUES
	('/rules', '# ルール\nこれはルールです');

TRUNCATE TABLE teams CASCADE;
INSERT INTO teams (id, code, name, organization, max_members, created_at, updated_at) VALUES
    ('a1de8fe6-26c8-42d7-b494-dea48e409091', 1, 'トラブルシューターズ', 'ICTSC Association', 6, '2025-02-02 08:00:00+00', '2025-02-02 08:00:00+00'),
    ('83027d5e-fa32-41d6-b290-fc38ba337f89', 2, 'トラブルメイカーズ', 'ICTSC Association', 2, '2025-02-02 08:00:00+00', '2025-02-02 08:00:00+00');

TRUNCATE TABLE invitation_codes CASCADE;
INSERT INTO invitation_codes (id, team_id, code, created_at, expires_at) VALUES
    ('ad3f83d3-65be-4884-8a03-adb11a8127ef', 'a1de8fe6-26c8-42d7-b494-dea48e409091', 'LHNZXGSF7L59WCG9', '2025-02-02 08:10:00+00', '2038-04-02 15:00:00+00');

TRUNCATE TABLE users CASCADE;
INSERT INTO users (id, name, created_at) VALUES
	('3a4ca027-5e02-4ade-8e2d-eddb39adc235', 'alice', '2025-02-03 00:00:00+00'),
	('c4530ce6-d990-4414-8389-feca26883115', 'bob', '2025-02-03 00:00:00+00');

TRUNCATE TABLE user_profiles CASCADE;
INSERT INTO user_profiles (user_id, display_name, created_at, updated_at) VALUES
	('3a4ca027-5e02-4ade-8e2d-eddb39adc235', 'Alice', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('c4530ce6-d990-4414-8389-feca26883115', 'ボブ', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00');

TRUNCATE TABLE discord_users CASCADE;
INSERT INTO discord_users (user_id, discord_user_id, linked_at) VALUES
	('3a4ca027-5e02-4ade-8e2d-eddb39adc235', 123456789012345678, '2025-02-03 00:00:00+00'),
	('c4530ce6-d990-4414-8389-feca26883115', 234567890123456789, '2025-02-03 00:00:00+00');

TRUNCATE TABLE team_members CASCADE;
INSERT INTO team_members (user_id, team_id, invitation_code_id) VALUES
    ('3a4ca027-5e02-4ade-8e2d-eddb39adc235', 'a1de8fe6-26c8-42d7-b494-dea48e409091', 'ad3f83d3-65be-4884-8a03-adb11a8127ef');

TRUNCATE TABLE problems CASCADE;
INSERT INTO problems (id, code, type, title, max_score, category, redeploy_rule, created_at, updated_at) VALUES
	('16643c32-c686-44ba-996b-2fbe43b54513', '0001', 'DESCRIPTIVE', 'チーム名をください', 10, 'Server', 'UNREDEPLOYABLE', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('24f6aef0-5dcd-4032-825b-d1b19174a6f2', '0002', 'DESCRIPTIVE', '過去のみ問題 (day1-am)', 100, 'Test', 'PERCENTAGE_PENALTY', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('35f7bf01-6ede-5043-926c-e2c20c285b03', '0003', 'DESCRIPTIVE', '未来のみ問題 (day2-am)', 150, 'Test', 'MANUAL', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('46a8cf12-7fef-6154-a37d-f3d31d396c14', '0004', 'DESCRIPTIVE', '全期間問題', 500, 'Test', 'UNREDEPLOYABLE', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	-- 追加: 現在提出可能 (day1-pm)
	('57b9d023-8100-4265-b48e-e4e42e507d25', '0005', 'DESCRIPTIVE', 'OSPF経路制御', 200, 'Network', 'UNREDEPLOYABLE', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('68cae134-9211-4376-a59f-f5f53f618e36', '0006', 'DESCRIPTIVE', 'Webサーバー構築', 150, 'Server', 'MANUAL', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('79dbe245-a322-4487-b6a0-060640729f47', '0007', 'DESCRIPTIVE', 'ファイアウォール設定', 300, 'Security', 'PERCENTAGE_PENALTY', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('8aecf356-b433-4598-87b1-171751830058', '0008', 'DESCRIPTIVE', 'コンテナ運用管理', 250, 'Server', 'UNREDEPLOYABLE', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	-- 追加: 過去 (day1-am のみ → 提出終了)
	('9bfd0467-c544-46a9-98c2-282862941169', '0009', 'DESCRIPTIVE', 'BGPピアリング (day1-am)', 200, 'Network', 'UNREDEPLOYABLE', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('ac0e1578-d655-47ba-a9d3-39397305227a', '0010', 'DESCRIPTIVE', 'DNS権威サーバー (day1-am)', 100, 'Network', 'MANUAL', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('bd1f2689-e766-48cb-bae4-4a4a8416338b', '0011', 'DESCRIPTIVE', 'RADIUS認証 (day1-am)', 150, 'Security', 'UNREDEPLOYABLE', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	-- 追加: 未来 (day2-am → 提出予定)
	('ce30379a-f877-49dc-8bf5-5b5b9527449c', '0012', 'DESCRIPTIVE', 'IDS/IPSチューニング (day2-am)', 200, 'Security', 'MANUAL', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('df41489b-0988-4aed-9c06-6c6ca63855ad', '0013', 'DESCRIPTIVE', 'IPv6移行 (day2-am)', 250, 'Network', 'UNREDEPLOYABLE', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	-- 追加: 未来 (day2-pm → 別の提出予定グループ)
	('e05259ac-1a99-4bfe-ad17-7d7db74966be', '0014', 'DESCRIPTIVE', 'ロードバランサ冗長化 (day2-pm)', 300, 'Server', 'PERCENTAGE_PENALTY', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('f1636abd-2baa-4c0f-be28-8e8ec85a77cf', '0015', 'DESCRIPTIVE', 'TLS証明書管理 (day2-pm)', 150, 'Security', 'UNREDEPLOYABLE', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	-- 追加: 複数スケジュール (day1-pm + day2-am → 現在提出可能)
	('02747bce-3cbb-4d10-8f39-9f9fd96b88d0', '0016', 'DESCRIPTIVE', '総合トラブルシューティング', 500, 'Network', 'MANUAL', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00');

TRUNCATE TABLE redeploy_percentage_penalties CASCADE;
INSERT INTO redeploy_percentage_penalties (problem_id, threshold, percentage) VALUES
	('24f6aef0-5dcd-4032-825b-d1b19174a6f2', 1, 10),
	('79dbe245-a322-4487-b6a0-060640729f47', 2, 15),
	('e05259ac-1a99-4bfe-ad17-7d7db74966be', 1, 20);

TRUNCATE TABLE problem_contents CASCADE;
INSERT INTO problem_contents (problem_id, page_id, page_path, body, explanation) VALUES
	('16643c32-c686-44ba-996b-2fbe43b54513', 'page1', '/page1',
		E'## 概要\n\nトラコン君はチーム名を求めています。チーム名を教えてあげましょう。\n\n## 制約\n\n- チーム名の読みがひらがなで書かれていること',
		E'## 採点基準\n\nチーム名の読みが書かれていれば10点'),
	('24f6aef0-5dcd-4032-825b-d1b19174a6f2', 'page2', '/page2',
		E'## 概要\n\nこの問題はday1-amのみで提出可能でした。\n\n現在は**表示されるが提出できない**はずです。\n\n## 回答\n\n何か書いてください。',
		E'## 採点基準\n\n何か書いてあれば100点'),
	('35f7bf01-6ede-5043-926c-e2c20c285b03', 'page3', '/page3',
		E'## 概要\n\nこの問題はday2-amで提出可能になります。\n\n現在は**表示されない**はずです。\n\n## 回答\n\n何か書いてください。',
		E'## 採点基準\n\n何か書いてあれば150点'),
	('46a8cf12-7fef-6154-a37d-f3d31d396c14', 'page4', '/page4',
		E'## 概要\n\nこの問題は全てのスケジュールで提出可能です。\n\n常に**表示も提出も可能**なはずです。\n\n## 回答\n\n何か書いてください。',
		E'## 採点基準\n\n何か書いてあれば500点'),
	('57b9d023-8100-4265-b48e-e4e42e507d25', 'page5', '/page5',
		E'## 概要\n\nOSPF (Open Shortest Path First) の経路制御に関する問題です。\n\n指定されたネットワークトポロジでOSPFを正しく設定してください。\n\n## 回答\n\n設定内容を記述してください。',
		E'## 採点基準\n\n正しい経路が確立されていれば200点'),
	('68cae134-9211-4376-a59f-f5f53f618e36', 'page6', '/page6',
		E'## 概要\n\nApache/Nginx を用いたWebサーバーの構築問題です。\n\nバーチャルホストの設定とSSLの有効化を行ってください。\n\n## 回答\n\n設定ファイルの内容を記述してください。',
		E'## 採点基準\n\nHTTPSでアクセスできれば150点'),
	('79dbe245-a322-4487-b6a0-060640729f47', 'page7', '/page7',
		E'## 概要\n\niptables/nftables を用いたファイアウォールの設定問題です。\n\n指定された通信のみを許可するルールを作成してください。\n\n## 回答\n\nルールの設定内容を記述してください。',
		E'## 採点基準\n\n正しいフィルタリングができていれば300点'),
	('8aecf356-b433-4598-87b1-171751830058', 'page8', '/page8',
		E'## 概要\n\nDockerコンテナの運用管理に関する問題です。\n\n停止したコンテナの原因を特定し、復旧してください。\n\n## 回答\n\n原因と対処方法を記述してください。',
		E'## 採点基準\n\nコンテナが正常稼働していれば250点'),
	('9bfd0467-c544-46a9-98c2-282862941169', 'page9', '/page9',
		E'## 概要\n\nBGPピアリングの設定問題です（day1-amで終了）。\n\nAS間のピアリングを確立してください。\n\n## 回答\n\n設定内容を記述してください。',
		E'## 採点基準\n\nピアリングが確立されていれば200点'),
	('ac0e1578-d655-47ba-a9d3-39397305227a', 'page10', '/page10',
		E'## 概要\n\nDNS権威サーバーの構築問題です（day1-amで終了）。\n\nBINDまたはNSDを用いてゾーンを設定してください。\n\n## 回答\n\nゾーンファイルの内容を記述してください。',
		E'## 採点基準\n\n名前解決ができれば100点'),
	('bd1f2689-e766-48cb-bae4-4a4a8416338b', 'page11', '/page11',
		E'## 概要\n\nRADIUS認証サーバーの設定問題です（day1-amで終了）。\n\nFreeRADIUSを用いて802.1X認証を設定してください。\n\n## 回答\n\n設定内容を記述してください。',
		E'## 採点基準\n\n認証が通れば150点'),
	('ce30379a-f877-49dc-8bf5-5b5b9527449c', 'page12', '/page12',
		E'## 概要\n\nIDS/IPSのチューニング問題です（day2-amから提出可能）。\n\nSnort/Suricataのルールを最適化してください。\n\n## 回答\n\nルールの内容を記述してください。',
		E'## 採点基準\n\n検知精度が基準を満たせば200点'),
	('df41489b-0988-4aed-9c06-6c6ca63855ad', 'page13', '/page13',
		E'## 概要\n\nIPv4からIPv6への移行問題です（day2-amから提出可能）。\n\nデュアルスタック環境を構築してください。\n\n## 回答\n\n設定内容を記述してください。',
		E'## 採点基準\n\nIPv6で通信できれば250点'),
	('e05259ac-1a99-4bfe-ad17-7d7db74966be', 'page14', '/page14',
		E'## 概要\n\nロードバランサの冗長化問題です（day2-pmから提出可能）。\n\nHAProxyまたはKeepalivedを用いて冗長構成を構築してください。\n\n## 回答\n\n設定内容を記述してください。',
		E'## 採点基準\n\nフェイルオーバーが動作すれば300点'),
	('f1636abd-2baa-4c0f-be28-8e8ec85a77cf', 'page15', '/page15',
		E'## 概要\n\nTLS証明書の管理問題です（day2-pmから提出可能）。\n\nLet''s Encryptを用いた証明書の自動更新を設定してください。\n\n## 回答\n\n手順を記述してください。',
		E'## 採点基準\n\n証明書が有効であれば150点'),
	('02747bce-3cbb-4d10-8f39-9f9fd96b88d0', 'page16', '/page16',
		E'## 概要\n\n総合トラブルシューティング問題です（day1-pm〜day2-amの長期間）。\n\n複数のサービスが停止しています。全て復旧してください。\n\n## 回答\n\n原因と対処方法を記述してください。',
		E'## 採点基準\n\n全サービスが復旧していれば500点');

TRUNCATE TABLE answers CASCADE;
INSERT INTO answers (id, problem_id, team_id, number, user_id,  created_at_range) VALUES
	('7cedf13e-5325-425e-a5d6-fea5fc127e49', '16643c32-c686-44ba-996b-2fbe43b54513', 'a1de8fe6-26c8-42d7-b494-dea48e409091', 1, '3a4ca027-5e02-4ade-8e2d-eddb39adc235', tstzrange('2025-02-03 00:00:00+00', '2025-02-03 00:20:00+00')),
	('4bb7a232-e0de-4b6d-b1a3-8e50737d73b2', '16643c32-c686-44ba-996b-2fbe43b54513', 'a1de8fe6-26c8-42d7-b494-dea48e409091', 2, '3a4ca027-5e02-4ade-8e2d-eddb39adc235', tstzrange('2025-02-03 00:30:00+00', '2025-02-03 00:50:00+00')),
	('abbe9c4e-eef5-40ac-a04e-6d8877b15185', '24f6aef0-5dcd-4032-825b-d1b19174a6f2', 'a1de8fe6-26c8-42d7-b494-dea48e409091', 1, '3a4ca027-5e02-4ade-8e2d-eddb39adc235', tstzrange('2025-02-03 00:10:00+00', '2025-02-03 00:30:00+00'));

TRUNCATE TABLE descriptive_answers CASCADE;
INSERT INTO descriptive_answers (answer_id, body) VALUES
	('7cedf13e-5325-425e-a5d6-fea5fc127e49', '問題Aへのチーム1の解答1'),
	('4bb7a232-e0de-4b6d-b1a3-8e50737d73b2', '問題Aへのチーム1の解答2'),
	('abbe9c4e-eef5-40ac-a04e-6d8877b15185', '問題Bへのチーム1の解答1');

TRUNCATE TABLE marking_results CASCADE;
INSERT INTO marking_results (id, answer_id, judge_name, visibility, created_at) VALUES
	('862b646a-5fdd-4a77-bb2d-7ef5d4f1d069', '7cedf13e-5325-425e-a5d6-fea5fc127e49', 'judge', 'PUBLIC',  '2025-02-03 01:00:00+00'),
	('358ebef7-c626-44fa-9a3d-da2d9083fe5e', '7cedf13e-5325-425e-a5d6-fea5fc127e49', 'judge', 'PRIVATE', '2025-02-03 01:30:00+00'),
	('87cb974d-dedd-4039-a189-b34f3a57e62c', 'abbe9c4e-eef5-40ac-a04e-6d8877b15185', 'judge', 'PUBLIC',  '2025-02-03 01:20:00+00');

TRUNCATE TABLE descriptive_marking_rationales CASCADE;
INSERT INTO descriptive_marking_rationales (marking_result_id, rationale) VALUES
	('862b646a-5fdd-4a77-bb2d-7ef5d4f1d069', 'comment'),
	('358ebef7-c626-44fa-9a3d-da2d9083fe5e', 'comment2'),
	('87cb974d-dedd-4039-a189-b34f3a57e62c', 'comment3');

TRUNCATE TABLE scores CASCADE;
INSERT INTO scores (marking_result_id, marked_score) VALUES
	('862b646a-5fdd-4a77-bb2d-7ef5d4f1d069', 8),
	('358ebef7-c626-44fa-9a3d-da2d9083fe5e', 7),
	('87cb974d-dedd-4039-a189-b34f3a57e62c', 80);

TRUNCATE TABLE answer_scores CASCADE;
INSERT INTO answer_scores (answer_id, visibility, marking_result_id) VALUES
	('7cedf13e-5325-425e-a5d6-fea5fc127e49', 'PRIVATE', '358ebef7-c626-44fa-9a3d-da2d9083fe5e'),
	('7cedf13e-5325-425e-a5d6-fea5fc127e49', 'TEAM', '862b646a-5fdd-4a77-bb2d-7ef5d4f1d069'),
	('7cedf13e-5325-425e-a5d6-fea5fc127e49', 'PUBLIC', '862b646a-5fdd-4a77-bb2d-7ef5d4f1d069'),
	('abbe9c4e-eef5-40ac-a04e-6d8877b15185', 'PRIVATE', '87cb974d-dedd-4039-a189-b34f3a57e62c'),
	('abbe9c4e-eef5-40ac-a04e-6d8877b15185', 'TEAM', '87cb974d-dedd-4039-a189-b34f3a57e62c'),
	('abbe9c4e-eef5-40ac-a04e-6d8877b15185', 'PUBLIC', '87cb974d-dedd-4039-a189-b34f3a57e62c');

TRUNCATE TABLE problem_scores CASCADE;
INSERT INTO problem_scores (team_id, problem_id, visibility, marking_result_id, updated_at) VALUES
	('a1de8fe6-26c8-42d7-b494-dea48e409091', '16643c32-c686-44ba-996b-2fbe43b54513', 'PRIVATE', '358ebef7-c626-44fa-9a3d-da2d9083fe5e', '2025-02-03 00:00:00+00'),
	('a1de8fe6-26c8-42d7-b494-dea48e409091', '16643c32-c686-44ba-996b-2fbe43b54513', 'TEAM', '862b646a-5fdd-4a77-bb2d-7ef5d4f1d069', '2025-02-03 00:00:00+00'),
	('a1de8fe6-26c8-42d7-b494-dea48e409091', '16643c32-c686-44ba-996b-2fbe43b54513', 'PUBLIC', '862b646a-5fdd-4a77-bb2d-7ef5d4f1d069', '2025-02-03 00:00:00+00'),
	('a1de8fe6-26c8-42d7-b494-dea48e409091', '24f6aef0-5dcd-4032-825b-d1b19174a6f2', 'PRIVATE', '87cb974d-dedd-4039-a189-b34f3a57e62c', '2025-02-03 00:10:00+00'),
	('a1de8fe6-26c8-42d7-b494-dea48e409091', '24f6aef0-5dcd-4032-825b-d1b19174a6f2', 'TEAM', '87cb974d-dedd-4039-a189-b34f3a57e62c', '2025-02-03 00:10:00+00'),
	('a1de8fe6-26c8-42d7-b494-dea48e409091', '24f6aef0-5dcd-4032-825b-d1b19174a6f2', 'PUBLIC', '87cb974d-dedd-4039-a189-b34f3a57e62c', '2025-02-03 00:10:00+00');

TRUNCATE TABLE score_visibility_settings CASCADE;
INSERT INTO score_visibility_settings (id, ranking_freeze_at) VALUES
	(TRUE, NULL);

TRUNCATE TABLE notices CASCADE;
INSERT INTO notices (slug, title, markdown, effective_from) VALUES
	('current-notice', 'Current Notice', '現在のお知らせです', '2025-02-03 00:00:00+00'),
	('current-notice2', 'Current Notice2', '現在のお知らせ2です', '2025-02-03 00:00:00+00'),
	('future-notice', 'Future Notice', '未来のお知らせです', '2035-03-03T00:00:00Z');

TRUNCATE TABLE schedules CASCADE;
INSERT INTO schedules (name, start_at, end_at) VALUES
    ('day1-am', '2026-01-01 01:00:00+00', '2026-01-01 03:00:00+00'),
    ('day1-pm', '2026-01-01 04:00:00+00', '2099-12-31 23:59:59+00'),
    ('day2-am', '2100-01-01 00:00:00+00', '2100-01-01 03:00:00+00'),
    ('day2-pm', '2100-01-01 04:00:00+00', '2100-01-01 07:00:00+00');

TRUNCATE TABLE problem_schedules CASCADE;
INSERT INTO problem_schedules (problem_id, schedule_name) VALUES
    -- 0001: day1-pm のみ (現在 → 提出○)
    ('16643c32-c686-44ba-996b-2fbe43b54513', 'day1-pm'),
    -- 0002: day1-am のみ (過去 → 提出✗)
    ('24f6aef0-5dcd-4032-825b-d1b19174a6f2', 'day1-am'),
    -- 0003: day2-am のみ (未来 → 提出予定)
    ('35f7bf01-6ede-5043-926c-e2c20c285b03', 'day2-am'),
    -- 0004: 全スケジュール (過去+現在+未来 → 提出○)
    ('46a8cf12-7fef-6154-a37d-f3d31d396c14', 'day1-am'),
    ('46a8cf12-7fef-6154-a37d-f3d31d396c14', 'day1-pm'),
    ('46a8cf12-7fef-6154-a37d-f3d31d396c14', 'day2-am'),
    ('46a8cf12-7fef-6154-a37d-f3d31d396c14', 'day2-pm'),
    -- 0005〜0008: day1-pm のみ (現在 → 提出○)
    ('57b9d023-8100-4265-b48e-e4e42e507d25', 'day1-pm'),
    ('68cae134-9211-4376-a59f-f5f53f618e36', 'day1-pm'),
    ('79dbe245-a322-4487-b6a0-060640729f47', 'day1-pm'),
    ('8aecf356-b433-4598-87b1-171751830058', 'day1-pm'),
    -- 0009〜0011: day1-am のみ (過去 → 提出✗)
    ('9bfd0467-c544-46a9-98c2-282862941169', 'day1-am'),
    ('ac0e1578-d655-47ba-a9d3-39397305227a', 'day1-am'),
    ('bd1f2689-e766-48cb-bae4-4a4a8416338b', 'day1-am'),
    -- 0012〜0013: day2-am のみ (未来 → 提出予定)
    ('ce30379a-f877-49dc-8bf5-5b5b9527449c', 'day2-am'),
    ('df41489b-0988-4aed-9c06-6c6ca63855ad', 'day2-am'),
    -- 0014〜0015: day2-pm のみ (未来・別グループ → 提出予定)
    ('e05259ac-1a99-4bfe-ad17-7d7db74966be', 'day2-pm'),
    ('f1636abd-2baa-4c0f-be28-8e8ec85a77cf', 'day2-pm'),
    -- 0016: day1-pm + day2-am (現在+未来 → 提出○)
    ('02747bce-3cbb-4d10-8f39-9f9fd96b88d0', 'day1-pm'),
    ('02747bce-3cbb-4d10-8f39-9f9fd96b88d0', 'day2-am');

TRUNCATE TABLE redeployment_requests CASCADE;
INSERT INTO redeployment_requests (id, team_id, problem_id, revision) VALUES
	('27fe747d-12c3-4dca-9b51-d3d13e8d815a', 'a1de8fe6-26c8-42d7-b494-dea48e409091', '24f6aef0-5dcd-4032-825b-d1b19174a6f2', 1),
	('427e945c-5bf0-46aa-9965-3b0fc9e6b869', '83027d5e-fa32-41d6-b290-fc38ba337f89', '24f6aef0-5dcd-4032-825b-d1b19174a6f2', 1);

TRUNCATE TABLE redeployment_events CASCADE;
INSERT INTO redeployment_events (id, request_id, status, created_at) VALUES
	('ebdd6be5-c1f1-4af6-97f8-968c5cb6a871', '27fe747d-12c3-4dca-9b51-d3d13e8d815a', 'QUEUED', '2025-02-03 00:00:00+00'),
	('ccbcaf7d-d20b-40f1-a626-e7db4744d99b', '27fe747d-12c3-4dca-9b51-d3d13e8d815a', 'DEPLOYING', '2025-02-03 00:01:00+00'),
	('ab7d5da4-8ee7-4017-865a-4dbac044d26c', '27fe747d-12c3-4dca-9b51-d3d13e8d815a', 'COMPLETED', '2025-02-03 00:06:00+00'),
	('19c3ea6d-02d6-4eee-8799-6a07d84634a7', '427e945c-5bf0-46aa-9965-3b0fc9e6b869', 'QUEUED', '2025-02-03 00:10:00+00'),
	('7b8f0998-f487-43f6-a2b3-92e5000615e3', '427e945c-5bf0-46aa-9965-3b0fc9e6b869', 'DEPLOYING', '2025-02-03 00:11:00+00');
