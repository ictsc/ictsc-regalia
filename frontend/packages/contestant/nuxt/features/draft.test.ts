import { describe, expect, it } from "vitest";
import { draftKey, loadDraft, saveDraft } from "./draft";
describe("private browser drafts", () => {
  it("round trips Unicode and keeps identities, teams, and problems separate", () => {
    localStorage.clear();
    const key = draftKey("alice", 12, "A01");
    saveDraft(localStorage, key, "原因と復旧確認\n日本語");
    expect(loadDraft(localStorage, key)?.body).toBe("原因と復旧確認\n日本語");
    for (const other of [
      draftKey("bob", 12, "A01"),
      draftKey("alice", 13, "A01"),
      draftKey("alice", 12, "A02"),
    ])
      expect(loadDraft(localStorage, other)).toBeNull();
  });
  it("reports storage errors without claiming a save succeeded", () => {
    expect(() =>
      saveDraft(
        {
          setItem() {
            throw new DOMException("quota");
          },
        },
        "draft",
        "body",
      ),
    ).toThrow("quota");
    expect(() =>
      loadDraft(
        {
          getItem() {
            throw new Error("disabled");
          },
        },
        "draft",
      ),
    ).toThrow("disabled");
  });
  it("rejects malformed drafts and does not import unowned legacy data", () => {
    localStorage.clear();
    localStorage.setItem("/problems/A01/answer", "legacy");
    expect(loadDraft(localStorage, draftKey("alice", 12, "A01"))).toBeNull();
    localStorage.setItem("bad", '{"body":3}');
    expect(() => loadDraft(localStorage, "bad")).toThrow();
  });
});
