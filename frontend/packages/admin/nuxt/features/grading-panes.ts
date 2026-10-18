export type PaneLayout = { collapsed: boolean[]; widths: number[] };

export const defaultPaneLayout = (): PaneLayout => ({
  collapsed: [false, false, false],
  widths: [23, 55, 22],
});

/**
 * Restore the three-pane layout from storage. Falls back to the legacy
 * four-pane layout (problem, explanation, answer, marking) when no current
 * layout is saved. Throws when the stored JSON is malformed.
 */
export function restorePaneLayout(
  savedJson: string | null,
  legacyJson: string | null,
  fallback: PaneLayout,
): PaneLayout {
  const layout = {
    collapsed: [...fallback.collapsed],
    widths: [...fallback.widths],
  };
  const saved = JSON.parse(savedJson ?? "null");
  if (Array.isArray(saved?.collapsed) && saved.collapsed.length === 3)
    layout.collapsed = saved.collapsed.map(Boolean);
  if (
    Array.isArray(saved?.widths) &&
    saved.widths.length === 3 &&
    saved.widths.every(
      (n: unknown) => typeof n === "number" && n >= 10 && n <= 75,
    )
  )
    layout.widths = saved.widths;
  if (!saved) {
    const old = JSON.parse(legacyJson ?? "null");
    if (Array.isArray(old?.collapsed) && old.collapsed.length === 4)
      layout.collapsed = [
        old.collapsed[0],
        old.collapsed[2],
        old.collapsed[3],
      ].map(Boolean);
    if (
      Array.isArray(old?.widths) &&
      old.widths.length === 4 &&
      old.widths.every(
        (n: unknown) => typeof n === "number" && Number.isFinite(n),
      )
    )
      layout.widths = [
        old.widths[0],
        old.widths[1] + old.widths[2],
        old.widths[3],
      ];
  }
  return layout;
}

/**
 * Move the boundary between two adjacent panes by `delta` width units,
 * keeping both panes at least 12 units wide.
 */
export function resizePaneWidths(
  start: readonly number[],
  leftIndex: number,
  rightIndex: number,
  delta: number,
) {
  const bounded = Math.max(
    12 - start[leftIndex]!,
    Math.min(start[rightIndex]! - 12, delta),
  );
  return start.map((value, i) =>
    i === leftIndex
      ? value + bounded
      : i === rightIndex
        ? value - bounded
        : value,
  );
}
