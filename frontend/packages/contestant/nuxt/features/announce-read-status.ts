import { ref, onMounted, onUnmounted } from "vue";
import type { Notice } from "./models";
export function useReadAnnouncements() {
  const slugs = ref<string[]>([]);
  const load = () => {
    try {
      const parsed: unknown = JSON.parse(
        localStorage.getItem("readAnnouncements") ?? "[]",
      );
      slugs.value = Array.isArray(parsed)
        ? parsed.filter((s): s is string => typeof s === "string")
        : [];
    } catch {
      slugs.value = [];
    }
  };
  const save = () => {
    try {
      localStorage.setItem("readAnnouncements", JSON.stringify(slugs.value));
    } catch {
      /* Reading notices still works without storage. */
    }
  };
  onMounted(() => {
    load();
    window.addEventListener("storage", load);
  });
  onUnmounted(() => window.removeEventListener("storage", load));
  return {
    isRead: (slug: string) => slugs.value.includes(slug),
    markAsRead: (slug: string) => {
      load();
      slugs.value = [...new Set([...slugs.value, slug])];
      save();
    },
    markAsUnread: (slug: string) => {
      load();
      slugs.value = slugs.value.filter((s) => s !== slug);
      save();
    },
    markAllAsRead: (notices: Notice[]) => {
      load();
      slugs.value = [
        ...new Set([...slugs.value, ...notices.map((n) => n.slug)]),
      ];
      save();
    },
    getUnreadNotices: (notices: Notice[]) =>
      notices.filter((n) => !slugs.value.includes(n.slug)),
  };
}
