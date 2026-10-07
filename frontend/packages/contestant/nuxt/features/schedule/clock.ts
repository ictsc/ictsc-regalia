import type { ScheduleEntry } from "../models";

const contestDate = new Intl.DateTimeFormat("en-CA", {
  timeZone: "Asia/Tokyo",
  year: "numeric",
  month: "2-digit",
  day: "2-digit",
});

function clockText(seconds: number): string {
  return [
    Math.floor(seconds / 3600),
    Math.floor(seconds / 60) % 60,
    seconds % 60,
  ]
    .map((part) => String(part).padStart(2, "0"))
    .join(":");
}

/** Show the remaining competition time for the current day in Japan. */
export function contestClock(
  entries: readonly ScheduleEntry[],
  now: number,
): { label: string; text: string } {
  const today = contestDate.format(now);
  const todaysEntries = entries.filter(
    (entry) => contestDate.format(Date.parse(entry.startAt)) === today,
  );
  if (todaysEntries.length) {
    const remaining = todaysEntries.reduce((total, entry) => {
      const start = Date.parse(entry.startAt);
      const end = Date.parse(entry.endAt);
      return total + Math.max(0, end - Math.max(now, start));
    }, 0);
    const started = todaysEntries.some(
      (entry) => Date.parse(entry.startAt) <= now,
    );
    return {
      label:
        remaining === 0
          ? "本日終了"
          : started
            ? "本日の残り時間"
            : "本日の競技時間",
      text: clockText(Math.floor(remaining / 1000)),
    };
  }

  const next = entries
    .filter((entry) => Date.parse(entry.startAt) > now)
    .sort((a, b) => Date.parse(a.startAt) - Date.parse(b.startAt))[0];
  return next
    ? {
        label: "次の競技開始まで",
        text: clockText(Math.floor((Date.parse(next.startAt) - now) / 1000)),
      }
    : { label: "競技終了", text: "00:00:00" };
}
