import { useCallback, useSyncExternalStore } from "react";
import type { Notice } from "@ictsc/proto/contestant/v1";

const STORAGE_KEY = "readAnnouncements";

function getSnapshot(): string[] {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored == null) return [];
  try {
    return JSON.parse(stored) as string[];
  } catch {
    return [];
  }
}

let cachedSlugs: string[] = getSnapshot();

function getSnapshotCached(): string[] {
  return cachedSlugs;
}

function setSlugs(slugs: string[]): void {
  cachedSlugs = slugs;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(slugs));
}

let listeners: Array<() => void> = [];

function subscribeInternal(callback: () => void): () => void {
  listeners.push(callback);
  const handler = (e: StorageEvent) => {
    if (e.key === STORAGE_KEY) {
      cachedSlugs = getSnapshot();
      callback();
    }
  };
  window.addEventListener("storage", handler);
  return () => {
    listeners = listeners.filter((l) => l !== callback);
    window.removeEventListener("storage", handler);
  };
}

function notifyListeners(): void {
  for (const listener of listeners) {
    listener();
  }
}

export function useReadAnnouncements() {
  const readSlugs = useSyncExternalStore(
    subscribeInternal,
    getSnapshotCached,
    () => [],
  );

  const markAsRead = useCallback((slug: string) => {
    const current = getSnapshot();
    if (!current.includes(slug)) {
      const next = [...current, slug];
      setSlugs(next);
      notifyListeners();
    }
  }, []);

  const markAllAsRead = useCallback((notices: Notice[]) => {
    const current = getSnapshot();
    const currentSet = new Set(current);
    const newSlugs = notices
      .map((n) => n.slug)
      .filter((s) => !currentSet.has(s));
    if (newSlugs.length === 0) return;
    setSlugs([...current, ...newSlugs]);
    notifyListeners();
  }, []);

  const markAsUnread = useCallback((slug: string) => {
    const current = getSnapshot();
    if (current.includes(slug)) {
      const next = current.filter((s) => s !== slug);
      setSlugs(next);
      notifyListeners();
    }
  }, []);

  const isRead = useCallback(
    (slug: string): boolean => {
      return readSlugs.includes(slug);
    },
    [readSlugs],
  );

  const getUnreadNotices = useCallback(
    (notices: Notice[]): Notice[] => {
      const readSet = new Set(readSlugs);
      return notices.filter((n) => !readSet.has(n.slug));
    },
    [readSlugs],
  );

  return { markAsRead, markAsUnread, markAllAsRead, isRead, getUnreadNotices };
}
