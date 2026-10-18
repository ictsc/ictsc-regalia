import { describe, expect, it } from "vitest";
import {
  activeCommentCount,
  clearLineDrafts,
  draftLineComments,
  emptyLineDraftState,
  lineHistory,
  rangeLabel,
  removeDraft,
  selectLine,
  type LineDraftState,
} from "./line-drafts";

const write = (state: LineDraftState, body: string): LineDraftState => ({
  ...state,
  drafts: state.drafts.map((d) =>
    d.id === state.activeId ? { ...d, body } : d,
  ),
});

describe("selectLine", () => {
  it("opens a draft for a line and extends it with shift", () => {
    let state = selectLine(emptyLineDraftState(), 3, false);
    expect(state.drafts).toEqual([
      { id: 1, startLine: 3, endLine: 3, body: "" },
    ]);
    expect(state.activeId).toBe(1);
    state = selectLine(state, 1, true);
    expect(state.drafts[0]).toMatchObject({ startLine: 1, endLine: 3 });
  });

  it("closes the open draft when clicking inside it", () => {
    const empty = selectLine(emptyLineDraftState(), 3, false);
    expect(selectLine(empty, 3, false).drafts).toEqual([]);
    const written = write(empty, "comment");
    const closed = selectLine(written, 3, false);
    expect(closed.drafts).toHaveLength(1);
    expect(closed.activeId).toBeNull();
    expect(selectLine(closed, 3, false).activeId).toBe(1);
  });

  it("drops empty drafts when another line is selected", () => {
    let state = selectLine(emptyLineDraftState(), 3, false);
    state = selectLine(state, 7, false);
    expect(state.drafts.map((d) => d.id)).toEqual([2]);
    state = selectLine(write(state, "kept"), 9, false);
    expect(state.drafts.map((d) => [d.id, d.startLine])).toEqual([
      [2, 7],
      [3, 9],
    ]);
  });

  it("keeps draft ids increasing after removal and clearing", () => {
    let state = selectLine(emptyLineDraftState(), 1, false);
    state = removeDraft(state, 1);
    expect(state.activeId).toBeNull();
    state = clearLineDrafts(selectLine(state, 2, false));
    expect(state.drafts).toEqual([]);
    expect(selectLine(state, 4, false).drafts[0]!.id).toBe(3);
  });
});

describe("draftLineComments", () => {
  it("returns trimmed written drafts sorted by line", () => {
    expect(
      draftLineComments([
        { id: 1, startLine: 5, endLine: 7, body: " range " },
        { id: 2, startLine: 2, endLine: 2, body: "single" },
        { id: 3, startLine: 1, endLine: 1, body: "  " },
      ]),
    ).toEqual([
      { line_number: 2, body: "single" },
      { line_number: 5, end_line_number: 7, body: "range" },
    ]);
  });
});

describe("saved line comments", () => {
  const marking = {
    id: "m1",
    answer: { team_code: 2, problem_code: "A01", answer_number: 1 },
    judge: { name: "staff" },
    score: 80,
    rationale: "",
    created_at: "2026-09-02T00:00:00Z",
    visibility: "PRIVATE",
    line_comments: [
      { line_number: 2, body: "a" },
      { line_number: 3, end_line_number: 5, body: "b" },
      { line_number: 5, body: "gone", deleted_at: "2026-09-02T01:00:00Z" },
    ],
  } as unknown as Parameters<typeof activeCommentCount>[0];
  it("groups active comments by their last line", () => {
    expect(lineHistory([marking])).toEqual({
      2: [{ id: "m1", judge: "staff", body: "a", startLine: 2, endLine: 2 }],
      5: [{ id: "m1", judge: "staff", body: "b", startLine: 3, endLine: 5 }],
    });
    expect(activeCommentCount(marking)).toBe(2);
  });
  it("labels single lines and ranges", () => {
    expect(rangeLabel(4, 4)).toBe("4 行目");
    expect(rangeLabel(3, 5)).toBe("3〜5 行目");
  });
});
