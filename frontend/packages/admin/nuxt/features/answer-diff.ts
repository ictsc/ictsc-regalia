import { diffArrays } from "diff";

export type AnswerDiffRow = {
  kind: "same" | "added" | "removed";
  previousLine: number | null;
  currentLine: number | null;
  text: string;
};

export type AnswerDiffHunk = {
  previousStart: number;
  previousCount: number;
  currentStart: number;
  currentCount: number;
  omittedBefore: number;
  rows: AnswerDiffRow[];
};

export type AnswerDiff = {
  rows: AnswerDiffRow[];
  addedLines: number;
  removedLines: number;
};

export function answerSourceLines(source: string): string[] {
  return source.replace(/\r\n?/g, "\n").split("\n");
}

export function compareAnswerSources(
  previous: string,
  current: string,
): AnswerDiff {
  const rows: AnswerDiffRow[] = [];
  let previousLine = 1;
  let currentLine = 1;
  let removed: string[] = [];
  let added: string[] = [];

  function flushChanges() {
    for (const text of removed) {
      rows.push({
        kind: "removed",
        previousLine: previousLine++,
        currentLine: null,
        text,
      });
    }
    for (const text of added) {
      rows.push({
        kind: "added",
        previousLine: null,
        currentLine: currentLine++,
        text,
      });
    }
    removed = [];
    added = [];
  }

  for (const change of diffArrays(
    answerSourceLines(previous),
    answerSourceLines(current),
  )) {
    if (change.removed) {
      removed.push(...change.value);
    } else if (change.added) {
      added.push(...change.value);
    } else {
      flushChanges();
      for (const text of change.value) {
        rows.push({
          kind: "same",
          previousLine: previousLine++,
          currentLine: currentLine++,
          text,
        });
      }
    }
  }
  flushChanges();
  return {
    rows,
    addedLines: rows.filter((row) => row.kind === "added").length,
    removedLines: rows.filter((row) => row.kind === "removed").length,
  };
}

export function buildAnswerDiffHunks(
  rows: AnswerDiffRow[],
  context = 3,
): AnswerDiffHunk[] {
  const ranges: { start: number; end: number }[] = [];
  for (let index = 0; index < rows.length; index++) {
    if (rows[index]?.kind === "same") continue;
    const start = Math.max(0, index - context);
    const end = Math.min(rows.length - 1, index + context);
    const last = ranges.at(-1);
    if (last && start <= last.end + 1) last.end = Math.max(last.end, end);
    else ranges.push({ start, end });
  }

  let consumedPrevious = 0;
  let consumedCurrent = 0;
  let cursor = 0;
  return ranges.map(({ start, end }) => {
    const omittedBefore = start - cursor;
    for (; cursor < start; cursor++) {
      if (rows[cursor]?.previousLine !== null) consumedPrevious++;
      if (rows[cursor]?.currentLine !== null) consumedCurrent++;
    }
    const hunkRows = rows.slice(start, end + 1);
    const previousCount = hunkRows.filter(
      (row) => row.previousLine !== null,
    ).length;
    const currentCount = hunkRows.filter(
      (row) => row.currentLine !== null,
    ).length;
    const hunk = {
      previousStart: consumedPrevious + (previousCount ? 1 : 0),
      previousCount,
      currentStart: consumedCurrent + (currentCount ? 1 : 0),
      currentCount,
      omittedBefore,
      rows: hunkRows,
    };
    for (; cursor <= end; cursor++) {
      if (rows[cursor]?.previousLine !== null) consumedPrevious++;
      if (rows[cursor]?.currentLine !== null) consumedCurrent++;
    }
    return hunk;
  });
}
