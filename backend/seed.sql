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
INSERT INTO problems (id, code, type, title, max_score, redeploy_rule, created_at, updated_at) VALUES
	('16643c32-c686-44ba-996b-2fbe43b54513', 'ZZA', 'DESCRIPTIVE', '問題A', 100, 'UNREDEPLOYABLE', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00'),
	('24f6aef0-5dcd-4032-825b-d1b19174a6f2', 'ZZB', 'DESCRIPTIVE', '問題B', 200, 'PERCENTAGE_PENALTY', '2025-02-03 00:00:00+00', '2025-02-03 00:00:00+00');

TRUNCATE TABLE redeploy_percentage_penalties CASCADE;
INSERT INTO redeploy_percentage_penalties (problem_id, threshold, percentage) VALUES
	('24f6aef0-5dcd-4032-825b-d1b19174a6f2', 2, 10);

TRUNCATE TABLE problem_contents CASCADE;
INSERT INTO problem_contents (problem_id, page_id, page_path, body, explanation) VALUES
	('16643c32-c686-44ba-996b-2fbe43b54513', 'page1', '/page1', '問題Aの本文', '問題Aの解説'),
	('24f6aef0-5dcd-4032-825b-d1b19174a6f2', 'page2', '/page2', '問題Bの本文', '問題Bの解説');

TRUNCATE TABLE answers CASCADE;
INSERT INTO answers (id, problem_id, team_id, number, user_id,  created_at_range) VALUES
	('7cedf13e-5325-425e-a5d6-fea5fc127e49', '16643c32-c686-44ba-996b-2fbe43b54513', 'a1de8fe6-26c8-42d7-b494-dea48e409091', 1, '3a4ca027-5e02-4ade-8e2d-eddb39adc235', tstzrange('2025-02-03 00:00:00+00', '2025-02-03 00:20:00+00')),
	('4bb7a232-e0de-4b6d-b1a3-8e50737d73b2', '16643c32-c686-44ba-996b-2fbe43b54513', 'a1de8fe6-26c8-42d7-b494-dea48e409091', 2, '3a4ca027-5e02-4ade-8e2d-eddb39adc235', tstzrange('2025-02-03 00:30:00+00', '2025-02-03 00:50:00+00')),
	('abbe9c4e-eef5-40ac-a04e-6d8877b15185', '24f6aef0-5dcd-4032-825b-d1b19174a6f2', 'a1de8fe6-26c8-42d7-b494-dea48e409091', 1, '3a4ca027-5e02-4ade-8e2d-eddb39adc235', tstzrange('2025-02-03 00:00:00+00', '2025-02-03 00:20:00+00'));

TRUNCATE TABLE descriptive_answers CASCADE;
INSERT INTO descriptive_answers (answer_id, body) VALUES
	('7cedf13e-5325-425e-a5d6-fea5fc127e49', '問題Aへのチーム1の回答1'),
	('4bb7a232-e0de-4b6d-b1a3-8e50737d73b2', '問題Aへのチーム1の回答2'),
	('abbe9c4e-eef5-40ac-a04e-6d8877b15185', '問題Bへのチーム1の回答1');
