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
	('16643c32-c686-44ba-996b-2fbe43b54513', 'AAA', 'DESCRIPTIVE', 'チーム名をください', 10, 'Server', 'UNREDEPLOYABLE', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('24f6aef0-5dcd-4032-825b-d1b19174a6f2', 'BBB', 'DESCRIPTIVE', '過去のみ問題 (day1-am)', 100, 'Test', 'UNREDEPLOYABLE', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('35f7bf01-6ede-5043-926c-e2c20c285b03', 'CCC', 'DESCRIPTIVE', '未来のみ問題 (day2-am)', 150, 'Test', 'UNREDEPLOYABLE', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('46a8cf12-7fef-6154-a37d-f3d31d396c14', 'DDD', 'DESCRIPTIVE', '全期間問題', 500, 'Test', 'UNREDEPLOYABLE', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00');

TRUNCATE TABLE redeploy_percentage_penalties CASCADE;

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
		E'## 採点基準\n\n何か書いてあれば500点');

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
	('7cedf13e-5325-425e-a5d6-fea5fc127e49', 'PUBLIC', '862b646a-5fdd-4a77-bb2d-7ef5d4f1d069'),
	('abbe9c4e-eef5-40ac-a04e-6d8877b15185', 'PRIVATE', '87cb974d-dedd-4039-a189-b34f3a57e62c'),
	('abbe9c4e-eef5-40ac-a04e-6d8877b15185', 'PUBLIC', '87cb974d-dedd-4039-a189-b34f3a57e62c');

TRUNCATE TABLE problem_scores CASCADE;
INSERT INTO problem_scores (team_id, problem_id, visibility, marking_result_id, updated_at) VALUES
	('a1de8fe6-26c8-42d7-b494-dea48e409091', '16643c32-c686-44ba-996b-2fbe43b54513', 'PRIVATE', '358ebef7-c626-44fa-9a3d-da2d9083fe5e', '2025-02-03 00:00:00+00'),
	('a1de8fe6-26c8-42d7-b494-dea48e409091', '16643c32-c686-44ba-996b-2fbe43b54513', 'PUBLIC', '862b646a-5fdd-4a77-bb2d-7ef5d4f1d069', '2025-02-03 00:00:00+00'),
	('a1de8fe6-26c8-42d7-b494-dea48e409091', '24f6aef0-5dcd-4032-825b-d1b19174a6f2', 'PRIVATE', '87cb974d-dedd-4039-a189-b34f3a57e62c', '2025-02-03 00:10:00+00'),
	('a1de8fe6-26c8-42d7-b494-dea48e409091', '24f6aef0-5dcd-4032-825b-d1b19174a6f2', 'PUBLIC', '87cb974d-dedd-4039-a189-b34f3a57e62c', '2025-02-03 00:10:00+00');

TRUNCATE TABLE notices CASCADE;
INSERT INTO notices (slug, title, markdown, effective_from) VALUES
	('current-notice', 'Current Notice', '現在のお知らせです', '2025-02-03 00:00:00+00'),
	('current-notice2', 'Current Notice2', '現在のお知らせ2です', '2025-02-03 00:00:00+00'),
	('future-notice', 'Future Notice', '未来のお知らせです', '2035-03-03T00:00:00Z');

TRUNCATE TABLE schedules CASCADE;
INSERT INTO schedules (name, start_at, end_at) VALUES
    ('day1-am', NOW() - INTERVAL '4 hours', NOW() - INTERVAL '2 hours'),
    ('day1-pm', NOW() - INTERVAL '1 hour',  NOW() + INTERVAL '2 hours'),
    ('day2-am', NOW() + INTERVAL '18 hours', NOW() + INTERVAL '20 hours'),
    ('day2-pm', NOW() + INTERVAL '21 hours', NOW() + INTERVAL '23 hours');

TRUNCATE TABLE problem_schedules CASCADE;
INSERT INTO problem_schedules (problem_id, schedule_name) VALUES
    -- AAA: day1-pm のみ (現在 → 表示○、提出○)
    ('16643c32-c686-44ba-996b-2fbe43b54513', 'day1-pm'),
    -- BBB: day1-am のみ (過去 → 表示○、提出✗)
    ('24f6aef0-5dcd-4032-825b-d1b19174a6f2', 'day1-am'),
    -- CCC: day2-am のみ (未来 → 表示✗、提出✗)
    ('35f7bf01-6ede-5043-926c-e2c20c285b03', 'day2-am'),
    -- DDD: 全スケジュール (過去+現在+未来 → 表示○、提出○)
    ('46a8cf12-7fef-6154-a37d-f3d31d396c14', 'day1-am'),
    ('46a8cf12-7fef-6154-a37d-f3d31d396c14', 'day1-pm'),
    ('46a8cf12-7fef-6154-a37d-f3d31d396c14', 'day2-am'),
    ('46a8cf12-7fef-6154-a37d-f3d31d396c14', 'day2-pm');

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
