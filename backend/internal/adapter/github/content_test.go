package github

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const testCommit = "0123456789abcdef0123456789abcdef01234567"

func TestFetchExactCommitAndReferencedMarkdown(t *testing.T) {
	files := validContentFiles()
	var requestedFiles []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer github-token" || r.Header.Get("X-GitHub-Api-Version") != "2022-11-28" {
			t.Fatalf("missing GitHub request headers: %#v", r.Header)
		}
		if r.URL.Path == "/repos/ictsc/content/git/commits/"+testCommit {
			writeTestJSON(t, w, map[string]any{"sha": testCommit, "parents": []any{}})
			return
		}
		const prefix = "/repos/ictsc/content/contents/"
		if strings.HasPrefix(r.URL.Path, prefix) {
			if got := r.URL.Query().Get("ref"); got != testCommit {
				t.Fatalf("content ref = %q, want exact commit", got)
			}
			filePath, err := url.PathUnescape(strings.TrimPrefix(r.URL.Path, prefix))
			if err != nil {
				t.Fatal(err)
			}
			requestedFiles = append(requestedFiles, filePath)
			content, exists := files[filePath]
			if !exists {
				http.NotFound(w, r)
				return
			}
			writeTestJSON(t, w, map[string]any{
				"type": "file", "encoding": "base64", "size": len(content),
				"content": base64.StdEncoding.EncodeToString([]byte(content)), "sha": "blob-sha", "download_url": nil,
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	now := time.Date(2026, 8, 30, 2, 0, 0, 0, time.UTC)
	client := newGitHubTestClient(t, server.URL, server.Client(), now)
	client.apiToken = "github-token"
	snapshot, err := client.Fetch(context.Background(), "ictsc/content", "refs/heads/main", testCommit)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.CommitSHA != testCommit || snapshot.Repository != "ictsc/content" || snapshot.Ref != "refs/heads/main" || !snapshot.FetchedAt.Equal(now) {
		t.Fatalf("unexpected snapshot metadata: %#v", snapshot)
	}
	problem, ok := snapshot.Problem("A01")
	if !ok || problem.Body != "# Problem A01\n" || problem.Explanation != "# Explanation A01\n" {
		t.Fatalf("unexpected problem: %#v", problem)
	}
	announcement, ok := snapshot.Announcement("welcome")
	if !ok || announcement.Markdown != "# Welcome\n" || snapshot.Manifest.RuleMarkdown != "# Rule\n" {
		t.Fatalf("referenced markdown was not loaded: %#v %#v", announcement, snapshot.Manifest)
	}
	if len(requestedFiles) != len(files) {
		t.Fatalf("requested files = %#v", requestedFiles)
	}
}

func TestIsAncestorUsesGitHubCompare(t *testing.T) {
	descendant := "abcdef0123456789abcdef0123456789abcdef01"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "/repos/ictsc/content/compare/" + testCommit + "..." + descendant
		if r.URL.Path != want {
			t.Fatalf("path = %q, want %q", r.URL.Path, want)
		}
		writeTestJSON(t, w, map[string]any{"status": "ahead", "ahead_by": 3})
	}))
	defer server.Close()
	client := newGitHubTestClient(t, server.URL, server.Client(), time.Now())

	ancestor, err := client.IsAncestor(context.Background(), "ictsc/content", testCommit, descendant)
	if err != nil {
		t.Fatal(err)
	}
	if !ancestor {
		t.Fatal("expected compare status ahead to mean ancestor")
	}
}

func TestParseManifestRejectsInvalidContracts(t *testing.T) {
	base := validManifest()
	tests := map[string]string{
		"unknown field": strings.Replace(base, "version: 1", "version: 1\nunknown: true", 1),
		"unsafe path":   strings.Replace(base, "body_path: problems/A01.md", "body_path: ../secret.md", 1),
		"overlap": strings.Replace(base, "announcements:", `  - slug: second
    beginning: 2026-08-30T11:00:00Z
    ending: 2026-08-30T13:00:00Z
    problems: []
announcements:`, 1),
		"missing section membership": strings.Replace(base, "problems: [A01]", "problems: []", 1),
		"inline body":                strings.Replace(base, "body_path: problems/A01.md", "body_path: problems/A01.md\n    body: leaked", 1),
	}
	for name, manifest := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := ParseManifest([]byte(manifest))
			if !errors.Is(err, ErrInvalidManifest) {
				t.Fatalf("error = %v, want ErrInvalidManifest", err)
			}
		})
	}
}

func TestFetchRejectsCommitMismatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeTestJSON(t, w, map[string]any{"sha": strings.Repeat("f", 40)})
	}))
	defer server.Close()
	client := newGitHubTestClient(t, server.URL, server.Client(), time.Now())
	_, err := client.Fetch(context.Background(), "ictsc/content", "refs/heads/main", testCommit)
	if err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("error = %v, want exact-commit mismatch", err)
	}
}

func validContentFiles() map[string]string {
	return map[string]string{
		"content/manifest.yaml":            validManifest(),
		"content/problems/A01.md":          "# Problem A01\n",
		"content/explanations/A01.md":      "# Explanation A01\n",
		"content/announcements/welcome.md": "# Welcome\n",
		"content/rule.md":                  "# Rule\n",
	}
}

func validManifest() string {
	return fmt.Sprintf(`version: 1
sections:
  - slug: day1
    beginning: 2026-08-30T09:00:00Z
    ending: 2026-08-30T12:00:00Z
    problems: [A01]
problems:
  - code: A01
    title: Example Problem
    max_score: 100
    category: Network
    section_slug: day1
    type: DESCRIPTIVE
    body_path: problems/A01.md
    explanation_path: explanations/A01.md
    redeploy_rule:
      type: PERCENTAGE_PENALTY
      penalty_threshold: 1
      penalty_percentage: 10
announcements:
  - slug: welcome
    title: Welcome
    markdown_path: announcements/welcome.md
    effective_from: 2026-08-30T08:00:00Z
rule_path: rule.md
`)
}
