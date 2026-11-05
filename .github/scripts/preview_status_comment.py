#!/usr/bin/env python3
"""Post or update the PR preview status comment on a validated preview PR.

The comment mirrors the Cloudflare Pages deploy comment and is edited in place
as the preview progresses (building, publishing, published or failed). It never
claims that the deployment is Ready: no workflow here can see Argo CD health.
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
STATE_PREFIX = "<!-- ictsc-regalia-pr-preview-state"
STATE_RE = re.compile(re.escape(STATE_PREFIX)
                      + r" sha=([0-9a-f]{40}) run=([0-9]+) attempt=([0-9]+) rank=([0-9]+) -->")

# phase: (rank, status text, build cell, publish cell, show preview URL)
# Within one build run attempt a higher rank is a later stage, so a late
# delivery of an earlier stage never overwrites a later one.
PHASES = {
    "building": (0, "🔨 Building…", "In progress", "Waiting for build", True),
    "publishing": (1, "📦 Publishing images…", "Succeeded", "Queued", True),
    "build_failed": (2, "❌ Build failed", "Failed", "Skipped", False),
    "build_cancelled": (2, "⏹️ Cancelled", "Cancelled", "Skipped", False),
    "published": (3, "✅ Images published (preview will be deployed by Argo CD)",
                  "Succeeded", "Succeeded", True),
    "publish_failed": (3, "❌ Publish failed", "Succeeded", "Failed", False),
    "publish_cancelled": (3, "⏹️ Cancelled", "Succeeded", "Cancelled", False),
}
PUBLISH_PHASES = {"published", "publish_failed", "publish_cancelled"}


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


def build_phase(action, conclusion):
    """Map a PR preview build workflow_run event to a phase, or None to ignore."""
    if action == "in_progress":
        return "building"
    if action != "completed":
        return None
    return {"success": "publishing", "failure": "build_failed",
            "timed_out": "build_failed", "cancelled": "build_cancelled"}.get(conclusion)


def parse_state(body):
    match = STATE_RE.search(body)
    if not match:
        return None
    sha, run, attempt, rank = match.groups()
    return sha, (int(run), int(attempt), int(rank))


def supersedes(body, sha, key):
    """Return whether a status for (sha, key) may replace the comment body.

    The caller has already checked that sha is the PR's current head, so a
    comment about any other SHA is older and always replaced. For the same SHA,
    a newer run/attempt or a later stage of the same attempt wins.
    """
    state = parse_state(body)
    if state is None or state[0] != sha:
        return True
    return key >= state[1]


def render(repo, number, sha, phase, build_run_id, build_run_attempt=1,
           publish_run_id=None):
    rank, status, build_cell, publish_cell, show_url = PHASES[phase]
    base = f"https://github.com/{repo}"
    build_url = f"{base}/actions/runs/{build_run_id}"
    publish_url = publish_run_id and f"{base}/actions/runs/{publish_run_id}"
    url = f"https://pr-{number}.k8s-dev.ictsc.net/"
    rows = [
        ("Latest commit", f"<a href='{base}/commit/{sha}'><code>{sha[:7]}</code></a>"),
        ("Status", status),
    ]
    if show_url:
        rows.append(("Preview URL", f"<a href='{url}'>{url}</a>"))
    rows.append(("Build", f"<a href='{build_url}'>{build_cell}</a>"))
    rows.append(("Image publish", f"<a href='{publish_url}'>{publish_cell}</a>"
                 if publish_url else publish_cell))
    table = "\n".join(f"<tr><td><strong>{name}:</strong></td><td>{value}</td></tr>"
                      for name, value in rows)
    return "\n".join([
        MARKER,
        f"{STATE_PREFIX} sha={sha} run={build_run_id} "
        f"attempt={build_run_attempt} rank={rank} -->",
        f"## PR preview #{number}",
        "",
        f"<table>{table}</table>",
        "",
        f"[View logs]({publish_url or build_url})",
        "",
        "<sub>Argo CD deployment and application health are not verified by this "
        "comment; the URL may be unavailable and requires Dex sign-in. Closing the "
        "PR or removing the <code>preview</code> label tears the preview down.</sub>",
    ]) + "\n"


def update_comment(api, repo, number, sha, phase, build_run_id,
                   build_run_attempt=1, publish_run_id=None):
    if not REPO_RE.fullmatch(repo) or not SHA_RE.fullmatch(sha):
        raise ValueError("invalid repository or SHA")
    if phase not in PHASES or (phase in PUBLISH_PHASES) != (publish_run_id is not None):
        raise ValueError("invalid phase or missing publish run ID")
    ids = [number, build_run_id, build_run_attempt]
    if publish_run_id is not None:
        ids.append(publish_run_id)
    if any(not isinstance(n, int) or isinstance(n, bool) or n < 1 for n in ids):
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
    key = (build_run_id, build_run_attempt, PHASES[phase][0])
    if owned and not supersedes(owned[0]["body"], sha, key):
        return "skipped: comment already shows a newer run or later stage"
    # Recheck immediately before a write; a close or new commit must not get a
    # new status based on the previous workflow run.
    if not eligible(api("GET", pr_path), repo, number, sha):
        return "skipped: PR changed before comment write"

    body = render(repo, number, sha, phase, build_run_id, build_run_attempt,
                  publish_run_id)
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
    publish_run = os.environ.get("PUBLISH_RUN_ID")
    print(update_comment(
        github_api,
        os.environ["GITHUB_REPOSITORY"],
        int(os.environ["PR_NUMBER"]),
        os.environ["HEAD_SHA"],
        os.environ["PREVIEW_PHASE"],
        int(os.environ["BUILD_RUN_ID"]),
        int(os.environ["BUILD_RUN_ATTEMPT"]),
        int(publish_run) if publish_run else None,
    ))
