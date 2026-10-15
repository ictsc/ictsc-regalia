export type Draft = { body: string; savedAt: string };
export function draftKey(name: string, team: number, problem: string) {
  return `regalia/draft/v1/${encodeURIComponent(name)}/${team}/${encodeURIComponent(problem)}`;
}
export function loadDraft(
  storage: Pick<Storage, "getItem">,
  key: string,
): Draft | null {
  const raw = storage.getItem(key);
  if (!raw) return null;
  const draft: unknown = JSON.parse(raw);
  if (
    typeof draft !== "object" ||
    !draft ||
    !("body" in draft) ||
    !("savedAt" in draft) ||
    typeof draft.body !== "string" ||
    typeof draft.savedAt !== "string"
  )
    throw new Error("保存された下書きを読み込めません");
  return { body: draft.body, savedAt: draft.savedAt };
}
export function saveDraft(
  storage: Pick<Storage, "setItem">,
  key: string,
  body: string,
): Draft {
  const draft = { body, savedAt: new Date().toISOString() };
  storage.setItem(key, JSON.stringify(draft));
  return draft;
}
