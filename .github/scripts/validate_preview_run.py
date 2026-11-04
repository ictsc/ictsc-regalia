#!/usr/bin/env python3
"""Validate trusted GitHub API snapshots for a preview publish run."""
import json
import re
import sys

MAX_ARTIFACT_BYTES = 2_000_000_000
EXPECTED_ARTIFACTS = {"preview-backend", "preview-frontend", "preview-migrations"}
SHA = re.compile(r"[a-f0-9]{40}\Z")
DIGEST = re.compile(r"sha256:[a-f0-9]{64}\Z")


def validate_run(run, workflow_id, repository):
    if (run.get("event") != "pull_request" or run.get("conclusion") != "success"
            or run.get("workflow_id") != workflow_id
            or run.get("repository", {}).get("full_name") != repository
            or run.get("path", "").split("@", 1)[0] != ".github/workflows/pr-preview-build.yaml"):
        raise ValueError("unexpected workflow run")


def matches_unlinked_run(run, pr):
    head = pr.get("head", {})
    return (pr.get("state") == "open"
            and any(label.get("name") == "preview" for label in pr.get("labels", []))
            and head.get("sha") == run.get("head_sha")
            and head.get("repo", {}).get("full_name")
            == run.get("head_repository", {}).get("full_name")
            and head.get("ref") == run.get("head_branch"))


def select_pr_number(run, candidates, workflow_id, repository):
    validate_run(run, workflow_id, repository)
    linked = run.get("pull_requests")
    if not isinstance(linked, list):
        raise ValueError("missing run PR association")
    if len(linked) == 1:
        number = linked[0].get("number")
        if not isinstance(number, int) or isinstance(number, bool) or number < 1:
            raise ValueError("invalid linked PR number")
        return number
    if linked:
        raise ValueError("ambiguous linked PR association")
    if (not SHA.fullmatch(run.get("head_sha", ""))
            or not run.get("head_repository", {}).get("full_name")
            or not run.get("head_branch") or not isinstance(candidates, list)):
        raise ValueError("missing unlinked run identity or PR candidates")
    matched = [pr["number"] for pr in candidates if matches_unlinked_run(run, pr)]
    if (len(matched) != 1 or not isinstance(matched[0], int)
            or isinstance(matched[0], bool) or matched[0] < 1):
        raise ValueError("unlinked run has no unique current PR")
    return matched[0]


def validate(run, pr, artifacts, workflow_id, repository, candidates):
    number = pr["number"]
    sha = pr["head"]["sha"]
    if (not isinstance(number, int) or isinstance(number, bool) or number < 1
            or not SHA.fullmatch(sha)
            or pr.get("base", {}).get("repo", {}).get("full_name") != repository):
        raise ValueError("invalid PR number, base repository, or head SHA")
    if select_pr_number(run, candidates, workflow_id, repository) != number:
        raise ValueError("selected PR differs from current PR")
    linked = run["pull_requests"]
    if linked and linked[0].get("head", {}).get("sha") != sha:
        raise ValueError("run does not belong to current PR head")
    if not linked and not matches_unlinked_run(run, pr):
        raise ValueError("unlinked run does not match current PR head")
    if (pr.get("state") != "open"
            or not any(label.get("name") == "preview" for label in pr.get("labels", []))):
        raise ValueError("PR closed or unlabeled")
    found = artifacts.get("artifacts", [])
    if artifacts.get("total_count") != 3 or len(found) != 3:
        raise ValueError("unexpected artifact count")
    by_name = {item.get("name"): item for item in found}
    if set(by_name) != EXPECTED_ARTIFACTS:
        raise ValueError("missing or duplicate artifact")
    for item in found:
        if (item.get("expired") is not False
                or not isinstance(item.get("id"), int)
                or not 0 < item.get("size_in_bytes", 0) <= MAX_ARTIFACT_BYTES
                or not DIGEST.fullmatch(item.get("digest", ""))):
            raise ValueError("invalid artifact metadata")
    return {"number": number, "sha": sha,
            "backend_id": by_name["preview-backend"]["id"],
            "frontend_id": by_name["preview-frontend"]["id"],
            "migrations_id": by_name["preview-migrations"]["id"]}


if __name__ == "__main__":
    try:
        mode = sys.argv[1]
        if mode == "select":
            run = json.load(open(sys.argv[2], encoding="utf-8"))
            candidates = json.load(open(sys.argv[3], encoding="utf-8"))
            print(select_pr_number(run, candidates, int(sys.argv[4]), sys.argv[5]))
            raise SystemExit(0)
        if mode != "validate":
            raise ValueError("unknown validation mode")
        run, pr, artifacts, candidates = (
            json.load(open(p, encoding="utf-8")) for p in sys.argv[2:6])
        result = validate(run, pr, artifacts, int(sys.argv[6]), sys.argv[7], candidates)
    except (KeyError, TypeError, ValueError) as exc:
        raise SystemExit(f"preview validation failed: {exc}") from exc
    for key, value in result.items():
        print(f"{key}={value}")
