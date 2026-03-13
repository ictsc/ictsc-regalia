import { Link } from "@tanstack/react-router";
import type { Notice } from "../../features/announce";
import { useReadAnnouncements } from "../../features/announce-read-status";
import { MaterialSymbol } from "../../components/material-symbol";

type BannerProps = {
  notices: Notice[];
};

export function UnreadAnnouncesBanner({ notices }: BannerProps) {
  const { getUnreadNotices, markAllAsRead } = useReadAnnouncements();
  const unreadNotices = getUnreadNotices(notices);

  if (unreadNotices.length === 0) {
    return null;
  }

  const label = `未読の新着アナウンスが ${unreadNotices.length} 件あります`;

  return (
    <div className="flex min-w-0 items-center gap-8">
      <Link
        to="/announces"
        title={label}
        className="text-14 bg-surface-1 hover:bg-surface-2 rounded-full flex min-w-0 items-center gap-4 py-4 pr-12 pl-8 font-bold transition"
      >
        <MaterialSymbol icon="notifications" size={20} className="text-icon shrink-0" />
        <span className="truncate">{label}</span>
      </Link>
      <button
        type="button"
        title="既読にする"
        className="text-14 bg-surface-1 hover:bg-surface-2 rounded-full flex shrink-0 items-center gap-4 py-4 pr-12 pl-8 transition"
        onClick={() => markAllAsRead(notices)}
      >
        <MaterialSymbol icon="done_all" size={20} className="text-icon" />
        既読にする
      </button>
    </div>
  );
}

const defaultToggleClassName =
  "text-14 bg-surface-1 hover:bg-surface-2 rounded-full flex shrink-0 items-center gap-4 py-4 pr-12 pl-8 transition";

export function ReadToggleButton({ slug, className }: { slug: string; className?: string }) {
  const { isRead, markAsRead, markAsUnread } = useReadAnnouncements();
  const read = isRead(slug);

  return (
    <button
      type="button"
      className={className ?? defaultToggleClassName}
      onClick={() => (read ? markAsUnread(slug) : markAsRead(slug))}
    >
      <MaterialSymbol
        icon={read ? "mark_email_read" : "mark_email_unread"}
        size={20}
        className="text-icon"
      />
      {read ? "未読にする" : "既読にする"}
    </button>
  );
}
