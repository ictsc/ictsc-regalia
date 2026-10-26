package github

import (
	"context"
	"crypto/sha1"
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

func TestFetchReadsMasterLayoutAtExactCommit(t *testing.T) {
	files := map[string]string{
		"schedule.yaml":               "schedules:\n- name: day1\n  startAt: 2026-08-30T09:00:00+09:00\n  endAt: 2026-08-30T12:00:00+09:00\n",
		"rule.md":                     "# Rule\n",
		"problems/A01.problem.md":     "---\ncode: a01\ntitle: Example\nmax_score: 100\ncategory: Network\nsubmissionable_schedules: [day1]\n---\n\n# Problem A01\n",
		"problems/A01.explanation.md": "# Explanation A01\n",
		"notices/welcome.md":          "---\ntitle: Welcome\nfrom: \"2026-08-30T08:00:00+09:00\"\n---\n\n# Welcome\n",
		"notices/_template.md":        "---\ntitle: template\nfrom: x\n---\n",
		"images/diagram.png":          "PNGDATA",
		"README.md":                   "ignored",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/repos/ictsc/content/git/commits/"+testCommit:
			writeTestJSON(t, w, map[string]any{"sha": testCommit})
		case r.URL.Path == "/repos/ictsc/content/git/trees/"+testCommit:
			if r.URL.Query().Get("recursive") != "1" {
				t.Fatalf("tree must be listed recursively: %s", r.URL)
			}
			entries := make([]map[string]any, 0, len(files))
			for filePath, content := range files {
				entries = append(entries, map[string]any{"path": filePath, "type": "blob", "size": len(content), "sha": testBlobSHA(filePath)})
			}
			writeTestJSON(t, w, map[string]any{"truncated": false, "tree": entries})
		case strings.HasPrefix(r.URL.Path, "/repos/ictsc/content/git/blobs/"):
			sha, _ := url.PathUnescape(strings.TrimPrefix(r.URL.Path, "/repos/ictsc/content/git/blobs/"))
			content := ""
			for filePath, candidate := range files {
				if testBlobSHA(filePath) == sha {
					content = candidate
				}
			}
			writeTestJSON(t, w, map[string]any{"encoding": "base64", "size": len(content), "content": base64.StdEncoding.EncodeToString([]byte(content))})
		case strings.HasPrefix(r.URL.Path, "/repos/ictsc/content/contents/"):
			if r.URL.Query().Get("ref") != testCommit {
				t.Fatalf("content must be read at the exact commit: %s", r.URL)
			}
			filePath, _ := url.PathUnescape(strings.TrimPrefix(r.URL.Path, "/repos/ictsc/content/contents/"))
			content, ok := files[filePath]
			if !ok {
				http.NotFound(w, r)
				return
			}
			writeTestJSON(t, w, map[string]any{"type": "file", "encoding": "base64", "size": len(content), "content": base64.StdEncoding.EncodeToString([]byte(content))})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client := newGitHubTestClient(t, server.URL, server.Client(), time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC))
	snapshot, err := client.Fetch(context.Background(), "ictsc/content", "refs/heads/main", testCommit)
	if err != nil {
		t.Fatal(err)
	}
	problem, ok := snapshot.Problem("A01")
	if !ok || problem.Body != "# Problem A01\n" || problem.Explanation != "# Explanation A01\n" || problem.SectionSlug != "day1" {
		t.Fatalf("problem = %#v", problem)
	}
	if snapshot.Manifest.RuleMarkdown != "# Rule\n" || len(snapshot.Manifest.Announcements) != 1 || snapshot.Manifest.Announcements[0].Markdown != "# Welcome" {
		t.Fatalf("rule or notices = %#v", snapshot.Manifest)
	}
	if len(snapshot.Assets) != 1 || snapshot.Assets[0].Name != "diagram.png" || snapshot.Assets[0].ContentType != "image/png" || string(snapshot.Assets[0].Data) != "PNGDATA" {
		t.Fatalf("assets = %#v", snapshot.Assets)
	}
}

func TestBuildManifestFollowsMasterSyncRules(t *testing.T) {
	files := masterFiles("schedules:\n- name: day1-am\n  startAt: 2026-03-14T10:55:00+09:00\n  endAt: 2026-03-14T13:00:00+09:00\nranking_freeze_at: 2026-03-14T12:00:00+09:00\n")
	files.Problems["problems/0099.problem.md"] = []byte("---\ncode: 0099\ntitle: Tutorial\nmax_score: 10000\ncategory: Tutorial\nredeploy_rule:\n  type: UNREDEPLOYABLE\n  penalty_threshold: 0\nsubmissionable_schedules: [day1-am]\n---\n\nbody\n")
	files.Problems["problems/1234.problem.md"] = []byte("---\ncode: 1234\ntitle: Penalty\nmax_score: 100\ncategory: Server\nredeploy_rule:\n  type: PERCENTAGE_PENALTY\n  penalty_threshold: 1\nsubmissionable_schedules: [day1-am]\n---\nbody\n")
	files.Problems["problems/3818.problem.md"] = []byte("---\ncode: 3818\ntitle: Default\nmax_score: 200\ncategory: Network\nsubmissionable_schedules: [day1-am]\n---\nbody\n")
	files.Problems["problems/0000.problem.md"] = []byte("---\ncode: \"0000\"\ntitle: Example\nmax_score: 1\nsubmissionable_schedules: [day1-am]\n---\nbody\n")
	files.Problems["problems/draft.problem.md"] = []byte("---\ncode: DRAFT\ntitle: Draft\nmax_score: 1\ndraft: true\n---\nbody\n")
	manifest, err := BuildManifest(files, false)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.RankingFreezeAt == nil || !manifest.RankingFreezeAt.Equal(time.Date(2026, 3, 14, 3, 0, 0, 0, time.UTC)) {
		t.Fatalf("freeze = %v", manifest.RankingFreezeAt)
	}
	byCode := map[string]int{}
	for index, problem := range manifest.Problems {
		byCode[problem.Code] = index
	}
	if len(manifest.Problems) != 3 || len(manifest.Sections[0].ProblemIDs) != 3 {
		t.Fatalf("problems = %#v", manifest.Problems)
	}
	if rule := manifest.Problems[byCode["0099"]].Redeploy; rule.Type != "UNREDEPLOYABLE" || rule.Threshold != nil {
		t.Fatalf("unredeployable = %#v", rule)
	}
	if rule := manifest.Problems[byCode["1234"]].Redeploy; *rule.Threshold != 1 || *rule.Percentage != 10 {
		t.Fatalf("percentage default = %#v", rule)
	}
	if rule := manifest.Problems[byCode["3818"]].Redeploy; rule.Type != "PERCENTAGE_PENALTY" || *rule.Threshold != 999 || *rule.Percentage != 0 {
		t.Fatalf("missing rule default = %#v", rule)
	}
}

func TestBuildManifestSecondaryQualifierRules(t *testing.T) {
	single := "schedules:\n- name: contest\n  startAt: 2026-12-12T13:00:00+09:00\n  endAt: 2026-12-12T18:00:00+09:00\n"
	files := masterFiles(single)
	files.Problems["problems/mid.problem.md"] = []byte("---\ncode: mid\ntitle: Mid\nmax_score: 100\ncategory: Network\nredeploy_rule:\n  type: MANUAL\n---\nbody\n")
	manifest, err := BuildManifest(files, true)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.Problems[0].Code != "MID" || manifest.Problems[0].SectionSlug != "contest" {
		t.Fatalf("problem = %#v", manifest.Problems[0])
	}
	tests := map[string]MasterFiles{
		"named schedule": withProblem(masterFiles(single), "---\ncode: x\ntitle: X\nmax_score: 1\nsubmissionable_schedules: [contest]\n---\nbody\n"),
		"two schedules":  withProblem(masterFiles(single+"- name: extra\n  startAt: 2026-12-12T19:00:00+09:00\n  endAt: 2026-12-12T20:00:00+09:00\n"), "---\ncode: x\ntitle: X\nmax_score: 1\n---\nbody\n"),
		"freeze":         withProblem(masterFiles(single+"ranking_freeze_at: 2026-12-12T17:00:00+09:00\n"), "---\ncode: x\ntitle: X\nmax_score: 1\n---\nbody\n"),
	}
	for name, files := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := BuildManifest(files, true); !errors.Is(err, ErrInvalidManifest) {
				t.Fatalf("error = %v, want ErrInvalidManifest", err)
			}
		})
	}
}

func TestBuildManifestRejectsInvalidMasterFiles(t *testing.T) {
	schedule := "schedules:\n- name: day1\n  startAt: 2026-08-30T09:00:00Z\n  endAt: 2026-08-30T12:00:00Z\n"
	tests := map[string]MasterFiles{
		"unknown frontmatter field": withProblem(masterFiles(schedule), "---\ncode: x\ntitle: X\nmax_score: 1\nsubmissionable_schedules: [day1]\nunknown: 1\n---\nbody\n"),
		"missing frontmatter":       withProblem(masterFiles(schedule), "# no frontmatter\n"),
		"unknown schedule":          withProblem(masterFiles(schedule), "---\ncode: x\ntitle: X\nmax_score: 1\nsubmissionable_schedules: [day9]\n---\nbody\n"),
		"multiple schedules":        withProblem(masterFiles(schedule), "---\ncode: x\ntitle: X\nmax_score: 1\nsubmissionable_schedules: [day1, day2]\n---\nbody\n"),
		"empty body":                withProblem(masterFiles(schedule), "---\ncode: x\ntitle: X\nmax_score: 1\nsubmissionable_schedules: [day1]\n---\n\n"),
		"unknown redeploy type":     withProblem(masterFiles(schedule), "---\ncode: x\ntitle: X\nmax_score: 1\nredeploy_rule:\n  type: 3\nsubmissionable_schedules: [day1]\n---\nbody\n"),
		"overlapping schedules":     masterFiles(schedule + "- name: day2\n  startAt: 2026-08-30T11:00:00Z\n  endAt: 2026-08-30T13:00:00Z\n"),
	}
	invalidNotice := masterFiles(schedule)
	invalidNotice.Notices["notices/bad.md"] = []byte("---\ntitle: Bad\nfrom: tomorrow\n---\nbody\n")
	tests["invalid notice date"] = invalidNotice
	for name, files := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := BuildManifest(files, false); !errors.Is(err, ErrInvalidManifest) {
				t.Fatalf("error = %v, want ErrInvalidManifest", err)
			}
		})
	}
}

func masterFiles(schedule string) MasterFiles {
	return MasterFiles{
		Schedule: []byte(schedule), Problems: map[string][]byte{}, Explanations: map[string][]byte{}, Notices: map[string][]byte{},
	}
}

func withProblem(files MasterFiles, content string) MasterFiles {
	files.Problems["problems/x.problem.md"] = []byte(content)
	return files
}

func testBlobSHA(filePath string) string {
	return fmt.Sprintf("%040x", sha1.Sum([]byte(filePath)))
}

func TestValidateImageEntryRejectsUnservableFiles(t *testing.T) {
	sha := testBlobSHA("x")
	for name, entry := range map[string]treeEntry{
		"nested":      {Path: "images/sub/a.png", SHA: sha},
		"svg":         {Path: "images/a.svg", SHA: sha},
		"text":        {Path: "images/notes.txt", SHA: sha},
		"space":       {Path: "images/a b.png", SHA: sha},
		"invalid sha": {Path: "images/a.png", SHA: "not-a-sha"},
	} {
		t.Run(name, func(t *testing.T) {
			if err := validateImageEntry(entry); !errors.Is(err, ErrInvalidManifest) {
				t.Fatalf("error = %v", err)
			}
		})
	}
	if err := validateImageEntry(treeEntry{Path: "images/problem1658-abc.PNG", SHA: sha}); err != nil {
		t.Fatal(err)
	}
}
