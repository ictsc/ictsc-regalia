import copy
import unittest
from validate_preview_run import validate

REPO = "ictsc/ictsc-regalia"
SHA = "a" * 40
RUN = {"event": "pull_request", "conclusion": "success", "workflow_id": 42,
       "repository": {"full_name": REPO},
       "path": ".github/workflows/pr-preview-build.yaml@refs/pull/7/merge",
       "head_repository": {"full_name": REPO},
       "pull_requests": [{"number": 7, "head": {"sha": SHA}}]}
PR = {"number": 7, "state": "open", "head": {"sha": SHA, "repo": {"full_name": REPO}},
      "labels": [{"name": "preview"}]}
ARTIFACTS = {"total_count": 3, "artifacts": [
    {"id": 11, "name": "preview-backend", "size_in_bytes": 1024,
     "digest": "sha256:" + "b" * 64, "expired": False},
    {"id": 12, "name": "preview-frontend", "size_in_bytes": 1024,
     "digest": "sha256:" + "c" * 64, "expired": False},
    {"id": 13, "name": "preview-migrations", "size_in_bytes": 1024,
     "digest": "sha256:" + "d" * 64, "expired": False}]}


class ValidationTests(unittest.TestCase):
    def setUp(self):
        self.run = copy.deepcopy(RUN)
        self.pr = copy.deepcopy(PR)
        self.artifacts = copy.deepcopy(ARTIFACTS)

    def valid(self):
        return validate(self.run, self.pr, self.artifacts, 42, REPO)

    def test_valid(self):
        self.assertEqual(self.valid()["backend_id"], 11)

    def test_stale_head(self):
        self.pr["head"]["sha"] = "d" * 40
        with self.assertRaises(ValueError): self.valid()

    def test_closed(self):
        self.pr["state"] = "closed"
        with self.assertRaises(ValueError): self.valid()

    def test_unlabeled(self):
        self.pr["labels"] = []
        with self.assertRaises(ValueError): self.valid()

    def test_fork(self):
        self.pr["head"]["repo"]["full_name"] = "someone/fork"
        self.run["head_repository"]["full_name"] = "someone/fork"
        self.assertEqual(self.valid()["number"], 7)

    def test_mismatched_linked_sha(self):
        self.run["pull_requests"][0]["head"]["sha"] = "d" * 40
        with self.assertRaises(ValueError): self.valid()

    def test_wrong_workflow(self):
        self.run["workflow_id"] = 99
        with self.assertRaises(ValueError): self.valid()

    def test_wrong_base_repository(self):
        self.run["repository"]["full_name"] = "someone/other"
        with self.assertRaises(ValueError): self.valid()

    def test_wrong_path(self):
        self.run["path"] = ".github/workflows/other.yaml@refs/pull/7/merge"
        with self.assertRaises(ValueError): self.valid()

    def test_duplicate_artifact(self):
        self.artifacts["artifacts"][1]["name"] = "preview-backend"
        with self.assertRaises(ValueError): self.valid()

    def test_invalid_digest(self):
        self.artifacts["artifacts"][0]["digest"] = "sha256:bad"
        with self.assertRaises(ValueError): self.valid()

    def test_expired(self):
        self.artifacts["artifacts"][0]["expired"] = True
        with self.assertRaises(ValueError): self.valid()


if __name__ == "__main__":
    unittest.main()
