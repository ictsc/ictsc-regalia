import { useEffect, useState, type ReactNode } from "react";
import {
  differenceInDays,
  format,
  formatDuration,
  intervalToDuration,
} from "date-fns";
import type { ScheduleEntry } from "@ictsc/proto/contestant/v1";
import { timestampDate } from "@bufbuild/protobuf/wkt";
import { Logo } from "../components/logo";
import { MaterialSymbol } from "../components/material-symbol";
import { Title } from "../components/title";

export type ContestState = "in_contest" | "waiting" | "ended";

export type IndexPageProps = {
  readonly state: ContestState;
  readonly currentScheduleName?: string;
  readonly nextScheduleName?: string;
  readonly timer?: ReactNode;
  readonly entries: ScheduleEntry[];
};

export function IndexPage(props: IndexPageProps) {
  switch (props.state) {
    case "in_contest":
      return <InContest {...props} />;
    case "waiting":
      return <OutOfContest {...props} />;
    case "ended":
      return <EndOfContest {...props} />;
    default:
      return null;
  }
}

export function Timer(props: {
  readonly nowMs?: number;
  readonly endMs: number;
}) {
  const [nowState, setNow] = useState(() => Date.now());
  useEffect(() => {
    const interval = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(interval);
  }, []);
  const nowMs = props.nowMs ?? nowState;
  const days = differenceInDays(props.endMs, nowMs);
  const dur = intervalToDuration({ start: nowMs, end: props.endMs });
  return (
    <>
      <Title />
      <p
        className="flex w-[5em] items-baseline justify-end"
        title={formatDuration(dur)}
      >
        {days > 0 ? (
          `${days}日`
        ) : (
          <>
            <span className="w-[1.5em] text-center">
              {`${dur.hours ?? 0}`.padStart(2, "0")}
            </span>
            <span className="">:</span>
            <span className="w-[1.5em] text-center">
              {`${dur.minutes ?? 0}`.padStart(2, "0")}
            </span>
            <span className="">:</span>
            <span className="w-[1.5em] text-center">
              {`${dur.seconds ?? 0}`.padStart(2, "0")}
            </span>
          </>
        )}
      </p>
    </>
  );
}

function PageLayout(props: { readonly children: ReactNode }) {
  return (
    <div className="mx-40 flex h-full flex-col items-center justify-center">
      {props.children}
    </div>
  );
}

function InContest(props: IndexPageProps) {
  return (
    <PageLayout>
      <Logo width={500} />
      <span className="text-16 mt-16 underline">
        左のサイドメニューからタブを選択してください
      </span>
      <div className="rounded-16 border-primary mt-[48px] flex flex-col gap-8 border-2 p-16 *:px-8">
        <div className="flex">
          <MaterialSymbol icon="schedule" size={40} className="text-icon" />
          <div className="ml-8 flex flex-col">
            <div className="text-24 leading-[40px]">
              競技中
              {props.currentScheduleName && ` (${props.currentScheduleName})`}
            </div>
            {props.timer != null && (
              <div className="flex items-baseline">
                <div className="text-14">残り</div>
                <div className="text-32 w-[168px] text-end">{props.timer}</div>
              </div>
            )}
          </div>
        </div>
        <ScheduleTimeline entries={props.entries} />
      </div>
    </PageLayout>
  );
}

function OutOfContest(props: IndexPageProps) {
  const title = props.nextScheduleName
    ? `${props.nextScheduleName} まであと`
    : "次の競技まであと";
  return (
    <PageLayout>
      <h1 className="text-48 font-bold underline">{title}</h1>
      <div className="mt-40 flex items-center">
        <MaterialSymbol icon="schedule" size={48} className="text-icon" />
        <span className="text-48 ml-16 font-bold">{props.timer}</span>
      </div>
      <ScheduleTimeline entries={props.entries} />
    </PageLayout>
  );
}

function EndOfContest(props: IndexPageProps) {
  return (
    <PageLayout>
      <h1 className="text-48 font-bold underline">競技は終了しました</h1>
      <ScheduleTimeline entries={props.entries} />
    </PageLayout>
  );
}

const temporalStatusLabel = {
  past: "終了済のスケジュール",
  current: "現在進行中のスケジュール",
  future: "未来のスケジュール",
} as const;

const temporalColorClass = {
  past: "*:!text-disabled",
  current: "*:!text-primary",
  future: "*:!text-icon",
} as const;

function getTemporalStatus(
  entry: ScheduleEntry,
  now: Date,
): "past" | "current" | "future" {
  if (entry.endAt != null && now >= timestampDate(entry.endAt)) return "past";
  if (entry.startAt != null && now >= timestampDate(entry.startAt))
    return "current";
  return "future";
}

function formatTime(entry: ScheduleEntry): string {
  const fmt = (ts: ScheduleEntry["startAt"]) =>
    ts != null ? format(timestampDate(ts), "MM/dd HH:mm") : "";
  return `${fmt(entry.startAt)} - ${fmt(entry.endAt)}`;
}

function ScheduleTimeline(props: { readonly entries: ScheduleEntry[] }) {
  if (props.entries.length === 0) return null;

  const now = new Date();

  return (
    <div className="border-primary mt-8 flex w-full flex-col gap-4 border-t pt-8">
      {props.entries.map((entry) => {
        const status = getTemporalStatus(entry, now);
        return (
          <div
            key={entry.name}
            className={`grid grid-cols-[1em_auto_1fr] gap-x-8 ${temporalColorClass[status]}`}
          >
            <span>{status === "current" ? "▶" : ""}</span>
            <span>{entry.name}</span>
            <span
              className="font-mono tabular-nums"
              title={temporalStatusLabel[status]}
            >
              {formatTime(entry)}
            </span>
          </div>
        );
      })}
    </div>
  );
}
