import copy
import unittest
from validate_preview_run import select_pr_number, select_status_target, validate

REPO = "ictsc/ictsc-regalia"
SHA = "a" * 40
RUN = {"event": "pull_request", "conclusion": "success", "workflow_id": 42,
       "repository": {"full_name": REPO},
       "path": ".github/workflows/pr-preview-build.yaml@refs/pull/7/merge",
       "head_sha": SHA, "head_branch": "preview-branch",
       "head_repository": {"full_name": REPO},
       "pull_requests": [{"number": 7, "head": {"sha": SHA}}]}
PR = {"number": 7, "state": "open", "head": {"sha": SHA, "ref": "preview-branch",
                                              "repo": {"full_name": REPO}},
      "base": {"repo": {"full_name": REPO}},
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
        self.candidates = [self.pr]

    def valid(self):
        return validate(self.run, self.pr, self.artifacts, 42, REPO, self.candidates)

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
        self.run["pull_requests"] = []
        self.assertEqual(self.valid()["number"], 7)

    def test_unlinked_fork_no_candidate(self):
        self.run["pull_requests"] = []
        self.candidates = []
        with self.assertRaises(ValueError): self.valid()

    def test_unlinked_fork_ambiguous_candidate(self):
        self.run["pull_requests"] = []
        other = copy.deepcopy(self.pr)
        other["number"] = 8
        self.candidates.append(other)
        with self.assertRaises(ValueError): self.valid()

    def test_multiple_linked_prs(self):
        self.run["pull_requests"].append({"number": 8, "head": {"sha": SHA}})
        with self.assertRaises(ValueError): self.valid()

    def test_unlinked_fork_wrong_head_repo(self):
        self.run["pull_requests"] = []
        self.run["head_repository"]["full_name"] = "someone/other"
        with self.assertRaises(ValueError): self.valid()

    def test_unlinked_fork_wrong_branch(self):
        self.run["pull_requests"] = []
        self.run["head_branch"] = "other-branch"
        with self.assertRaises(ValueError): self.valid()

    def test_unlinked_fork_wrong_head_sha(self):
        self.run["pull_requests"] = []
        self.run["head_sha"] = "c" * 40
        with self.assertRaises(ValueError): self.valid()

    def test_unlinked_fork_closed(self):
        self.run["pull_requests"] = []
        self.pr["state"] = "closed"
        with self.assertRaises(ValueError): self.valid()

    def test_unlinked_fork_unlabeled(self):
        self.run["pull_requests"] = []
        self.pr["labels"] = []
        with self.assertRaises(ValueError): self.valid()

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


class StatusTargetTests(unittest.TestCase):
    def setUp(self):
        self.run = copy.deepcopy(RUN)
        self.run.update(run_attempt=2, conclusion=None, status="in_progress")
        self.jobs = {"jobs": [{"conclusion": None}, {"conclusion": "skipped"}]}

    def target(self, candidates=None):
        return select_status_target(self.run, candidates or [copy.deepcopy(PR)],
                                    self.jobs, 42, REPO)

    def test_running_build(self):
        self.assertEqual(self.target(), {"number": 7, "sha": SHA, "attempt": 2})

    def test_failed_build(self):
        self.run["conclusion"] = "failure"
        self.assertEqual(self.target()["number"], 7)

    def test_all_jobs_skipped_is_not_a_preview_build(self):
        self.jobs = {"jobs": [{"conclusion": "skipped"}] * 3}
        self.assertIsNone(self.target())

    def test_wrong_workflow_or_event(self):
        for key, value in (("workflow_id", 1), ("event", "push"),
                           ("path", ".github/workflows/other.yaml")):
            with self.subTest(key=key):
                run = copy.deepcopy(self.run)
                run[key] = value
                with self.assertRaises(ValueError):
                    select_status_target(run, [PR], self.jobs, 42, REPO)

    def test_invalid_attempt_or_jobs(self):
        self.run["run_attempt"] = 0
        with self.assertRaises(ValueError): self.target()
        self.run["run_attempt"] = 1
        self.jobs = {}
        with self.assertRaises(ValueError): self.target()

    def test_unlinked_fork_requires_unique_labeled_candidate(self):
        self.run["pull_requests"] = []
        self.assertEqual(self.target()["number"], 7)
        unlabeled = copy.deepcopy(PR)
        unlabeled["labels"] = []
        with self.assertRaises(ValueError): self.target([unlabeled])


if __name__ == "__main__":
    unittest.main()
