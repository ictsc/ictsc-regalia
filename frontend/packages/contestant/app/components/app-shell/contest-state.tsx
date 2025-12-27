import { MaterialSymbol } from "@app/components/material-symbol";

export type ContestState = "before" | "running" | "break" | "finished";
export type ContentStateViewProps = {
  readonly state: ContestState;
  readonly restDurationSeconds: number;
};

const stateMap: Record<ContestState, string> = {
  before: "競技開始前",
  running: "競技中",
  break: "休憩中",
  finished: "競技終了",
};

export const MINUTE = 60;
export const HOUR = 60 * MINUTE;
export const DAY = 24 * HOUR;

export function ContestStateView({
  state,
  restDurationSeconds,
}: ContentStateViewProps) {
  let rest = restDurationSeconds;
  const days = Math.floor(rest / DAY);
  rest %= DAY;
  const hours = Math.floor(rest / HOUR);
  rest %= HOUR;
  const minutes = Math.floor(rest / MINUTE);
  const seconds = rest % MINUTE;

  return (
    <div className="bg-surface-1 text-text flex h-[48px] w-[288px] items-center justify-between rounded-[8px] px-[8px]">
      <div className="flex items-center">
        <MaterialSymbol icon="schedule" size={24} />
        <span className="text-16 ml-[4px] line-clamp-1 w-[80px] text-clip">
          {stateMap[state]}
        </span>
      </div>
      {restDurationSeconds !== 0 && (
        <div className="flex items-baseline">
          <span className="text-12">残り</span>
          <span className="w-[128px] text-end">
            {days > 0 ? (
              <>
                <span className="text-24">{days}</span>
                <span className="text-16">日</span>
              </>
            ) : (
              <span className="text-24">
                {hours} : {minutes} : {seconds}
              </span>
            )}
          </span>
        </div>
      )}
    </div>
  );
}
