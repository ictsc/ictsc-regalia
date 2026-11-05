import importlib.util
from pathlib import Path
import unittest

SCRIPT = Path(__file__).with_name("preview_status_comment.py")
spec = importlib.util.spec_from_file_location("preview_status_comment", SCRIPT)
status = importlib.util.module_from_spec(spec)
spec.loader.exec_module(status)

REPO = "ictsc/ictsc-regalia"
SHA = "a" * 40


def pr(sha=SHA, state="open", labels=("preview",)):
    return {
        "number": 3500,
        "state": state,
        "base": {"repo": {"full_name": REPO}},
        "head": {"sha": sha},
        "labels": [{"name": name} for name in labels],
    }


def comment(body, user_id=status.BOT_ID, login=status.BOT_LOGIN, kind="Bot"):
    return {"id": 42, "body": body,
            "user": {"id": user_id, "login": login, "type": kind}}


class FakeAPI:
    def __init__(self, prs=None, pages=None):
        self.prs = list(prs or [pr(), pr()])
        self.pages = pages or {1: []}
        self.calls = []

    def __call__(self, method, path, payload=None):
        self.calls.append((method, path, payload))
        if method == "GET" and path.endswith("/pulls/3500"):
            return self.prs.pop(0)
        if method == "GET" and "/comments?" in path:
            page = int(path.split("page=")[-1])
            return self.pages.get(page, [])
        if method in ("POST", "PATCH"):
            return {"id": 42}
        raise AssertionError((method, path))

    @property
    def writes(self):
        return [call for call in self.calls if call[0] in ("POST", "PATCH")]


def run(api, phase="published", build_run_id=100, attempt=1, publish_run_id=200, sha=SHA):
    if phase not in status.PUBLISH_PHASES:
        publish_run_id = None
    return status.update_comment(api, REPO, 3500, sha, phase, build_run_id, attempt,
                                 publish_run_id)


def body_for(phase, build_run_id=100, attempt=1, sha=SHA):
    publish = 200 if phase in status.PUBLISH_PHASES else None
    return status.render(REPO, 3500, sha, phase, build_run_id, attempt, publish)


class PreviewCommentTests(unittest.TestCase):
    def test_create_describes_images_without_claiming_deployment_ready(self):
        api = FakeAPI()
        self.assertEqual(run(api), "created bot comment")
        self.assertEqual(api.writes[0][:2],
                         ("POST", "repos/ictsc/ictsc-regalia/issues/3500/comments"))
        body = api.writes[0][2]["body"]
        self.assertIn("https://pr-3500.k8s-dev.ictsc.net/", body)
        self.assertIn(f"/commit/{SHA}'><code>{SHA[:7]}</code>", body)
        self.assertIn("/actions/runs/100", body)
        self.assertIn("/actions/runs/200", body)
        self.assertIn("✅ Images published", body)
        self.assertIn("not verified", body)
        self.assertNotIn("Ready", body)

    def test_updates_only_own_marked_comment(self):
        api = FakeAPI(pages={1: [comment(status.MARKER + " old")]})
        self.assertEqual(run(api), "updated bot comment")
        self.assertEqual(api.writes[0][:2],
                         ("PATCH", "repos/ictsc/ictsc-regalia/issues/comments/42"))

    def test_does_not_edit_someone_elses_marked_comment(self):
        foreign = comment(status.MARKER + " old", user_id=123, login="operator", kind="User")
        api = FakeAPI(pages={1: [foreign]})
        self.assertEqual(run(api), "created bot comment")
        self.assertEqual(api.writes[0][0], "POST")

    def test_retry_is_idempotent(self):
        body = body_for("published")
        api = FakeAPI(pages={1: [comment(body)]})
        self.assertIn("unchanged", run(api))
        self.assertEqual(api.writes, [])

    def test_finds_bot_comment_on_later_page(self):
        api = FakeAPI(pages={1: [comment("ordinary", user_id=1)] * 100,
                             2: [comment(status.MARKER + " old")]})
        self.assertEqual(run(api), "updated bot comment")
        self.assertIn("page=2", api.calls[-3][1])

    def test_closed_unlabeled_or_stale_pr_is_not_commented(self):
        for stale in (pr(state="closed"), pr(labels=()), pr(sha="b" * 40)):
            with self.subTest(stale=stale):
                api = FakeAPI(prs=[stale])
                self.assertIn("skipped", run(api))
                self.assertEqual(api.writes, [])

    def test_recheck_catches_close_or_new_head_before_write(self):
        for changed in (pr(state="closed"), pr(sha="b" * 40)):
            with self.subTest(changed=changed):
                api = FakeAPI(prs=[pr(), changed])
                self.assertIn("skipped", run(api))
                self.assertEqual(api.writes, [])

    def test_multiple_owned_comments_fails_closed(self):
        api = FakeAPI(pages={1: [comment(status.MARKER), comment(status.MARKER)]})
        with self.assertRaisesRegex(ValueError, "multiple bot-owned"):
            run(api)
        self.assertEqual(api.writes, [])



class RenderTests(unittest.TestCase):
    def test_every_phase_renders_table_and_status(self):
        expected = {
            "building": ("🔨 Building…", True, "Waiting for build", "runs/100)"),
            "publishing": ("📦 Publishing images…", True, "Queued", "runs/100)"),
            "build_failed": ("❌ Build failed", False, "Skipped", "runs/100)"),
            "build_cancelled": ("⏹️ Cancelled", False, "Skipped", "runs/100)"),
            "published": ("✅ Images published (preview will be deployed by Argo CD)",
                          True, "Succeeded</a>", "runs/200)"),
            "publish_failed": ("❌ Publish failed", False, "Failed</a>", "runs/200)"),
            "publish_cancelled": ("⏹️ Cancelled", False, "Cancelled</a>", "runs/200)"),
        }
        self.assertEqual(set(expected), set(status.PHASES))
        for phase, (text, shows_url, publish_cell, logs) in expected.items():
            with self.subTest(phase=phase):
                body = body_for(phase)
                self.assertTrue(body.startswith(status.MARKER + "\n"))
                self.assertIn("## PR preview #3500", body)
                self.assertIn(f"<tr><td><strong>Status:</strong></td><td>{text}</td></tr>", body)
                self.assertIn("<strong>Latest commit:</strong>", body)
                self.assertEqual("Preview URL" in body, shows_url)
                self.assertIn(f"<strong>Image publish:</strong></td><td>", body)
                self.assertIn(publish_cell, body)
                self.assertIn(f"[View logs](https://github.com/{REPO}/actions/{logs}", body)
                self.assertIn("not verified", body)
                self.assertEqual(status.parse_state(body)[0], SHA)

    def test_build_phase_mapping(self):
        cases = {
            ("in_progress", ""): "building",
            ("completed", "success"): "publishing",
            ("completed", "failure"): "build_failed",
            ("completed", "timed_out"): "build_failed",
            ("completed", "cancelled"): "build_cancelled",
            ("completed", "skipped"): None,
            ("completed", "startup_failure"): None,
            ("requested", ""): None,
        }
        for (action, conclusion), phase in cases.items():
            with self.subTest(action=action, conclusion=conclusion):
                self.assertEqual(status.build_phase(action, conclusion), phase)

    def test_rejects_mismatched_publish_run(self):
        with self.assertRaisesRegex(ValueError, "phase"):
            status.update_comment(FakeAPI(), REPO, 3500, SHA, "published", 100, 1, None)
        with self.assertRaisesRegex(ValueError, "phase"):
            status.update_comment(FakeAPI(), REPO, 3500, SHA, "building", 100, 1, 200)
        with self.assertRaisesRegex(ValueError, "phase"):
            status.update_comment(FakeAPI(), REPO, 3500, SHA, "deployed", 100, 1, None)


class OrderingTests(unittest.TestCase):
    def write(self, existing, **kwargs):
        api = FakeAPI(pages={1: [comment(existing)]})
        return run(api, **kwargs), api.writes

    def test_progresses_through_stages_of_one_run(self):
        order = ["building", "publishing", "published"]
        for older, newer in zip(order, order[1:]):
            with self.subTest(older=older, newer=newer):
                result, writes = self.write(body_for(older), phase=newer)
                self.assertEqual(result, "updated bot comment")
                self.assertIn(status.PHASES[newer][1], writes[0][2]["body"])

    def test_late_earlier_stage_does_not_overwrite_later_stage(self):
        for later, earlier in (("published", "building"), ("published", "publishing"),
                               ("publishing", "building"), ("build_failed", "building"),
                               ("publish_failed", "publishing")):
            with self.subTest(later=later, earlier=earlier):
                result, writes = self.write(body_for(later), phase=earlier)
                self.assertIn("skipped", result)
                self.assertEqual(writes, [])

    def test_newer_head_sha_overwrites_older_sha_status(self):
        old = body_for("published", build_run_id=900, sha="b" * 40)
        result, writes = self.write(old, phase="building", build_run_id=100)
        self.assertEqual(result, "updated bot comment")
        self.assertIn(f"<code>{SHA[:7]}</code>", writes[0][2]["body"])

    def test_older_sha_finishing_late_does_not_overwrite_newer_sha(self):
        # The PR head is now SHA; a publish for the older b-SHA is stale.
        api = FakeAPI(pages={1: [comment(body_for("building"))]})
        result = run(api, phase="published", build_run_id=50, sha="b" * 40)
        self.assertIn("skipped", result)
        self.assertEqual(api.writes, [])

    def test_rerun_attempt_replaces_failed_attempt(self):
        result, _ = self.write(body_for("build_failed", attempt=1),
                               phase="building", attempt=2)
        self.assertEqual(result, "updated bot comment")

    def test_older_attempt_or_run_of_same_sha_is_ignored(self):
        for kwargs in ({"attempt": 1}, {"build_run_id": 99, "attempt": 2}):
            with self.subTest(**kwargs):
                result, writes = self.write(body_for("building", attempt=2),
                                            phase="published", **kwargs)
                self.assertIn("skipped", result)
                self.assertEqual(writes, [])

    def test_newer_run_of_same_sha_replaces_older_run(self):
        result, _ = self.write(body_for("published", build_run_id=100),
                               phase="building", build_run_id=101)
        self.assertEqual(result, "updated bot comment")

    def test_legacy_comment_without_state_is_replaced(self):
        result, _ = self.write(status.MARKER + "\n### PR preview #3500", phase="building")
        self.assertEqual(result, "updated bot comment")

    def test_publish_rerun_can_report_failure_after_success(self):
        result, _ = self.write(body_for("published"), phase="publish_failed")
        self.assertEqual(result, "updated bot comment")


if __name__ == "__main__":
    unittest.main()
