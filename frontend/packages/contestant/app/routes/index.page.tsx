import { useEffect, useState, type ReactNode } from "react";
import { differenceInDays, formatDuration, intervalToDuration } from "date-fns";
import { Logo } from "../components/logo";
import { MaterialSymbol } from "../components/material-symbol";
import { Title } from "../components/title";

export type ContestState = "in_contest" | "waiting" | "ended";

export type IndexPageProps = {
  readonly state: ContestState;
  readonly currentScheduleName?: string;
  readonly nextScheduleName?: string;
  readonly timer?: ReactNode;
};

export function IndexPage(props: IndexPageProps) {
  switch (props.state) {
    case "in_contest":
      return <InContest {...props} />;
    case "waiting":
      return <OutOfContest {...props} />;
    case "ended":
      return <EndOfContest />;
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

function InContest(props: IndexPageProps) {
  return (
    <div className="mx-40 flex h-full flex-col items-center justify-center">
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
        {props.nextScheduleName != null && (
          <div className="border-primary flex w-full items-center border-t pt-8">
            <div className="flex size-40 items-center justify-center">
              <MaterialSymbol
                icon="arrow_forward_ios"
                size={24}
                className="text-icon"
              />
            </div>
            <div className="text-14 mt-2 ml-8">
              次のスケジュール: {props.nextScheduleName}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

function OutOfContest(props: IndexPageProps) {
  const title = props.nextScheduleName
    ? `${props.nextScheduleName} まであと`
    : "次の競技まであと";
  return (
    <div className="mx-40 flex h-full flex-col items-center justify-center">
      <h1 className="text-48 font-bold underline">{title}</h1>
      <div className="mt-40 flex items-center">
        <MaterialSymbol icon="schedule" size={48} className="text-icon" />
        <span className="text-48 ml-16 font-bold">{props.timer}</span>
      </div>
    </div>
  );
}

function EndOfContest() {
  return (
    <div className="mx-40 flex h-full flex-col items-center justify-center">
      <h1 className="text-48 font-bold underline">競技は終了しました</h1>
    </div>
  );
}
