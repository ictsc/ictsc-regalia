import type { components } from "@ictsc/api";

type Schemas = components["schemas"];

export type LineDraft = {
  id: number;
  startLine: number;
  endLine: number;
  body: string;
};

/** Line comment drafts on the answer pane, kept until the marking is sent. */
export type LineDraftState = {
  drafts: LineDraft[];
  /** Draft whose editor is open. */
  activeId: number | null;
  /** Line where the current range selection started. */
  anchor: number | null;
  nextId: number;
};

export type LineHistoryNote = {
  id: string;
  judge: string;
  body: string;
  startLine: number;
  endLine: number;
};

export const emptyLineDraftState = (): LineDraftState => ({
  drafts: [],
  activeId: null,
  anchor: null,
  nextId: 1,
});

/** Clear drafts after a marking is saved; ids keep increasing. */
export function clearLineDrafts(state: LineDraftState): LineDraftState {
  return { ...state, drafts: [], activeId: null, anchor: null };
}

export function removeDraft(state: LineDraftState, id: number): LineDraftState {
  const drafts = state.drafts.filter((draft) => draft.id !== id);
  return state.activeId === id
    ? { ...state, drafts, activeId: null, anchor: null }
    : { ...state, drafts };
}

/**
 * Click on a line number. A plain click opens a draft for the line, closes
 * the open draft when clicking inside it, or reopens a written draft; a
 * shift-click extends the open draft to a range.
 */
export function selectLine(
  state: LineDraftState,
  line: number,
  extend: boolean,
): LineDraftState {
  const active = state.drafts.find((draft) => draft.id === state.activeId);
  if (extend && active && state.anchor !== null) {
    const resized = {
      ...active,
      startLine: Math.min(state.anchor, line),
      endLine: Math.max(state.anchor, line),
    };
    return {
      ...state,
      drafts: state.drafts.map((draft) =>
        draft.id === active.id ? resized : draft,
      ),
    };
  }
  if (!extend && active && line >= active.startLine && line <= active.endLine) {
    if (!active.body.trim()) return removeDraft(state, active.id);
    return { ...state, activeId: null, anchor: null };
  }
  const closed = state.drafts.findLast(
    (draft) =>
      draft.body.trim() && line >= draft.startLine && line <= draft.endLine,
  );
  if (!extend && closed) {
    return { ...state, activeId: closed.id, anchor: closed.startLine };
  }
  const id = state.nextId;
  return {
    drafts: [
      ...state.drafts.filter((draft) => draft.body.trim()),
      { id, startLine: line, endLine: line, body: "" },
    ],
    activeId: id,
    anchor: line,
    nextId: id + 1,
  };
}

/** Written drafts as API line comments, ordered by line. */
export function draftLineComments(
  drafts: readonly LineDraft[],
): Schemas["AnswerLineComment"][] {
  return drafts
    .filter((draft) => draft.body.trim())
    .map((draft) => ({
      line_number: draft.startLine,
      ...(draft.endLine > draft.startLine
        ? { end_line_number: draft.endLine }
        : {}),
      body: draft.body.trim(),
    }))
    .toSorted((a, b) => a.line_number - b.line_number);
}

/** Saved, non-deleted line comments grouped by their last line. */
export function lineHistory(
  marks: readonly Schemas["MarkingResult"][],
): Record<number, LineHistoryNote[]> {
  const byLine: Record<number, LineHistoryNote[]> = {};
  for (const mark of marks) {
    for (const note of mark.line_comments ?? []) {
      if (note.deleted_at) continue;
      const endLine = note.end_line_number ?? note.line_number;
      (byLine[endLine] ??= []).push({
        id: mark.id,
        judge: mark.judge.name,
        body: note.body,
        startLine: note.line_number,
        endLine,
      });
    }
  }
  return byLine;
}

export function activeCommentCount(marking: Schemas["MarkingResult"]) {
  return marking.line_comments.filter((note) => !note.deleted_at).length;
}

export function rangeLabel(startLine: number, endLine: number) {
  return startLine === endLine
    ? `${startLine} 行目`
    : `${startLine}〜${endLine} 行目`;
}
