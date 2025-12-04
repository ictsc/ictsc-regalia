import { useEffect, useState, type ReactNode } from "react";
import { differenceInDays, formatDuration, intervalToDuration } from "date-fns";
import { Phase } from "@ictsc/proto/contestant/v1";
import { Logo } from "../components/logo";
import { MaterialSymbol } from "../components/material-symbol";
import { Title } from "../components/title";

export type IndexPageProps = {
  readonly phase: Phase;
  readonly nextPhase?: Phase;
  readonly timer?: ReactNode;
};
export function IndexPage(props: IndexPageProps) {
  switch (props.phase) {
    case Phase.IN_CONTEST:
      return <InContest {...props} />;
    case Phase.OUT_OF_CONTEST:
    case Phase.BREAK:
      return <OutOfContest {...props} />;
    case Phase.AFTER_CONTEST:
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
  let nextPhase: string | undefined;
  switch (props.nextPhase) {
    case Phase.IN_CONTEST:
      nextPhase = "競技開始";
      break;
    case Phase.AFTER_CONTEST:
      nextPhase = "競技終了";
      break;
    case Phase.BREAK:
      nextPhase = "休憩";
      break;
  }
  return (
    <div className="mx-40 flex h-full flex-col items-center justify-center">
      <Logo width={500} />
      <span className="mt-16 text-16 underline">
        左のサイドメニューからタブを選択してください
      </span>
      <div className="mt-[48px] flex flex-col gap-8 rounded-16 border-2 border-primary p-16 *:px-8">
        <div className="flex">
          <MaterialSymbol icon="schedule" size={40} className="text-icon" />
          <div className="ml-8 flex flex-col">
            <div className="text-24 leading-[40px]">競技中</div>
            {props.timer != null && (
              <div className="flex items-baseline">
                <div className="text-14">残り</div>
                <div className="w-[168px] text-end text-32">{props.timer}</div>
              </div>
            )}
          </div>
        </div>
        {nextPhase != null && (
          <div className="flex w-full items-center border-t border-primary pt-8">
            <div className="flex size-40 items-center justify-center">
              <MaterialSymbol
                icon="arrow_forward_ios"
                size={24}
                className="text-icon"
              />
            </div>
            <div className="ml-8 mt-2 text-14">次のフェーズ: {nextPhase}</div>
          </div>
        )}
      </div>
    </div>
  );
}

function OutOfContest(props: IndexPageProps) {
  let title: string | undefined;
  switch (props.phase) {
    case Phase.OUT_OF_CONTEST:
      title = "競技開始まであと";
      break;
    case Phase.BREAK:
      title = "競技再開まであと";
      break;
  }
  return (
    <div className="mx-40 flex h-full flex-col items-center justify-center">
      <h1 className="text-48 font-bold underline">{title}</h1>
      <div className="mt-40 flex items-center">
        <MaterialSymbol icon="schedule" size={48} className="text-icon" />
        <span className="ml-16 text-48 font-bold">{props.timer}</span>
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
