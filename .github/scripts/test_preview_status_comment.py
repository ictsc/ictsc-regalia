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


def run(api):
    return status.update_comment(api, REPO, 3500, SHA, 100, 200)


class PreviewCommentTests(unittest.TestCase):
    def test_create_describes_images_without_claiming_deployment_ready(self):
        api = FakeAPI()
        self.assertEqual(run(api), "created bot comment")
        self.assertEqual(api.writes[0][:2],
                         ("POST", "repos/ictsc/ictsc-regalia/issues/3500/comments"))
        body = api.writes[0][2]["body"]
        self.assertIn("https://pr-3500.k8s-dev.ictsc.net/", body)
        self.assertIn(f"Head: `{SHA}`", body)
        self.assertIn("/actions/runs/100", body)
        self.assertIn("/actions/runs/200", body)
        self.assertIn("Deployment: **not verified", body)

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
        body = status.render(REPO, 3500, SHA, 100, 200)
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


if __name__ == "__main__":
    unittest.main()
