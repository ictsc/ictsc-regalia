import { useCallback, useSyncExternalStore } from "react";
import type { Notice } from "@ictsc/proto/contestant/v1";

const STORAGE_KEY = "readAnnouncements";

function getSnapshot(): string[] {
  if (typeof window === "undefined") return [];
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored == null) return [];
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

function handleStorageEvent(e: StorageEvent): void {
  if (e.key === STORAGE_KEY) {
    cachedSlugs = getSnapshot();
    notifyListeners();
  }
}

function subscribeInternal(callback: () => void): () => void {
  if (listeners.length === 0) {
    window.addEventListener("storage", handleStorageEvent);
  }
  listeners.push(callback);
  return () => {
    listeners = listeners.filter((l) => l !== callback);
    if (listeners.length === 0) {
      window.removeEventListener("storage", handleStorageEvent);
    }
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
