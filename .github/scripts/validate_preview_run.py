#!/usr/bin/env python3
"""Validate trusted GitHub API snapshots for a preview publish run."""
import json
import re
import sys

MAX_ARTIFACT_BYTES = 2_000_000_000
EXPECTED_ARTIFACTS = {"preview-backend", "preview-frontend", "preview-migrations"}
SHA = re.compile(r"[a-f0-9]{40}\Z")
DIGEST = re.compile(r"sha256:[a-f0-9]{64}\Z")


def validate(run, pr, artifacts, workflow_id, repository):
    number = pr["number"]
    sha = pr["head"]["sha"]
    if not isinstance(number, int) or number < 1 or not SHA.fullmatch(sha):
        raise ValueError("invalid PR number or head SHA")
    if (run.get("event") != "pull_request" or run.get("conclusion") != "success"
            or run.get("workflow_id") != workflow_id
            or run.get("path", "").split("@", 1)[0] != ".github/workflows/pr-preview-build.yaml"):
        raise ValueError("unexpected workflow run")
    linked = run.get("pull_requests", [])
    if (len(linked) != 1 or linked[0].get("number") != number
            or linked[0].get("head", {}).get("sha") != sha):
        raise ValueError("run does not belong to current PR head")
    if (pr.get("state") != "open"
            or not any(label.get("name") == "preview" for label in pr.get("labels", []))):
        raise ValueError("PR closed, unlabeled, or forked")
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
    run, pr, artifacts = (json.load(open(p, encoding="utf-8")) for p in sys.argv[1:4])
    try:
        result = validate(run, pr, artifacts, int(sys.argv[4]), sys.argv[5])
    except (KeyError, TypeError, ValueError) as exc:
        raise SystemExit(f"preview validation failed: {exc}") from exc
    for key, value in result.items():
        print(f"{key}={value}")
