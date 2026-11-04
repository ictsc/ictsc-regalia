/** Narrowest answer panel; matches the default width on mid-size screens. */
export const minAnswerWidth = 360;
/** Problem text keeps at least this much room next to a wide answer panel. */
export const minReadingWidth = 320;
export const answerWidthStorageKey = "regalia/answer-width";

/** Keeps a requested answer panel width between both minimums. */
export function clampAnswerWidth(width: number, gridWidth: number): number {
  const max = Math.max(minAnswerWidth, gridWidth - minReadingWidth);
  return Math.round(Math.min(Math.max(width, minAnswerWidth), max));
}

export function loadAnswerWidth(storage: Storage): number | undefined {
  const value = Number(storage.getItem(answerWidthStorageKey));
  return Number.isFinite(value) && value >= minAnswerWidth ? value : undefined;
}

export function saveAnswerWidth(storage: Storage, width: number | undefined) {
  if (width == null) storage.removeItem(answerWidthStorageKey);
  else storage.setItem(answerWidthStorageKey, String(width));
}
