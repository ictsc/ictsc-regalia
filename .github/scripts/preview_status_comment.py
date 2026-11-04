#!/usr/bin/env python3
"""Post an image-published status on a validated preview PR.

This runs only after the trusted publisher succeeds. It never claims that the
deployment is Ready: the publisher has no access to Argo CD health.
"""

import json
import os
import re
from urllib import request

MARKER = "<!-- ictsc-regalia-pr-preview-images -->"
BOT_ID = 41898282  # github-actions[bot]
BOT_LOGIN = "github-actions[bot]"
SHA_RE = re.compile(r"[0-9a-f]{40}\Z")
REPO_RE = re.compile(r"[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+\Z")


def eligible(pr, repo, number, sha):
    return (pr.get("number") == number
            and pr.get("state") == "open"
            and pr.get("base", {}).get("repo", {}).get("full_name") == repo
            and pr.get("head", {}).get("sha") == sha
            and any(label.get("name") == "preview" for label in pr.get("labels", [])))


def owned_comment(comment):
    user = comment.get("user") or {}
    body = comment.get("body")
    return (isinstance(body, str) and MARKER in body
            and user.get("id") == BOT_ID
            and user.get("login") == BOT_LOGIN
            and user.get("type") == "Bot")


def render(repo, number, sha, build_run_id, publish_run_id):
    host = f"pr-{number}.k8s-dev.ictsc.net"
    return "\n".join([
        MARKER,
        f"### PR preview #{number}",
        f"Expected URL: https://{host}/ (Dex sign-in required once deployed)",
        f"Head: `{sha}`",
        f"Build: [succeeded](https://github.com/{repo}/actions/runs/{build_run_id})",
        f"Image publish: [succeeded](https://github.com/{repo}/actions/runs/{publish_run_id})",
        "Deployment: **not verified by this workflow**. The URL may be unavailable, "
        "and this comment does not indicate Argo CD or application health. "
        "Closing the PR or removing its preview label starts teardown; do not "
        "treat this URL as available after that point.",
    ]) + "\n"


def update_comment(api, repo, number, sha, build_run_id, publish_run_id):
    if not REPO_RE.fullmatch(repo) or not SHA_RE.fullmatch(sha):
        raise ValueError("invalid repository or SHA")
    if any(not isinstance(n, int) or isinstance(n, bool) or n < 1
           for n in (number, build_run_id, publish_run_id)):
        raise ValueError("invalid PR or workflow run ID")

    pr_path = f"repos/{repo}/pulls/{number}"
    if not eligible(api("GET", pr_path), repo, number, sha):
        return "skipped: PR is closed, unlabeled, or has a different head"

    comments = []
    for page in range(1, 101):
        batch = api("GET", f"repos/{repo}/issues/{number}/comments?per_page=100&page={page}")
        if not isinstance(batch, list):
            raise ValueError("invalid comment listing")
        comments.extend(batch)
        if len(batch) < 100:
            break
    else:
        raise ValueError("comment listing exceeded page limit")

    owned = [comment for comment in comments if owned_comment(comment)]
    if len(owned) > 1:
        raise ValueError("multiple bot-owned preview comments")
    # Recheck immediately before a write; a close or new commit must not get a
    # new image-published status based on the previous workflow run.
    if not eligible(api("GET", pr_path), repo, number, sha):
        return "skipped: PR changed before comment write"

    body = render(repo, number, sha, build_run_id, publish_run_id)
    if owned:
        if owned[0]["body"] == body:
            return "unchanged: bot comment already current"
        api("PATCH", f"repos/{repo}/issues/comments/{owned[0]['id']}", {"body": body})
        return "updated bot comment"
    api("POST", f"repos/{repo}/issues/{number}/comments", {"body": body})
    return "created bot comment"


def github_api(method, path, payload=None):
    data = None if payload is None else json.dumps(payload).encode("utf-8")
    req = request.Request(
        f"https://api.github.com/{path}",
        data=data,
        method=method,
        headers={
            "Accept": "application/vnd.github+json",
            "Authorization": f"Bearer {os.environ['GITHUB_TOKEN']}",
            "X-GitHub-Api-Version": "2022-11-28",
            "Content-Type": "application/json",
        },
    )
    with request.urlopen(req, timeout=20) as response:
        return json.load(response)


if __name__ == "__main__":
    print(update_comment(
        github_api,
        os.environ["GITHUB_REPOSITORY"],
        int(os.environ["PR_NUMBER"]),
        os.environ["HEAD_SHA"],
        int(os.environ["BUILD_RUN_ID"]),
        int(os.environ["GITHUB_RUN_ID"]),
    ))
