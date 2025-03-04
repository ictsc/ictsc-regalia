import { clsx } from "clsx";

type RankingItemProps = {
  rank?: number;
  timeStamp?: Date;
  teamName: string;
  affiliation: string;
  score?: number;
  penalty?: number;
  fullScore?: boolean;
  rawFullScore?: boolean;
};

export function RankingItem(props: RankingItemProps) {
  return (
    <div className="flex w-full flex-wrap items-center justify-center gap-x-40 pb-64 pl-8 md:flex-nowrap">
      <div>
        <div className="flex w-[175px] flex-row items-baseline justify-between">
          {/* ランキング表示 */}
          <div className="flex flex-row items-baseline">
            <span className={clsx("text-32 font-bold")}>
              {props.rank != null ? props.rank : "-"}
            </span>
            <p className="text-14 font-bold">st</p>
          </div>
          {/* トータルスコア表示 */}
          <div className="flex flex-row items-baseline gap-8">
            <span
              className={clsx(
                "text-32 font-bold",
                props.rank === 1 && "text-primary",
                props.score == null && "px-32",
              )}
            >
              {props.score != null ? props.score : "-"}
            </span>
            <p className="text-14 font-bold">pt</p>
          </div>
        </div>
        <div className="flex flex-row gap-4">
          <p className="text-12">Time:</p>
          <span
            className={clsx(
              "text-12 font-bold",
              props.timeStamp == null && "px-12",
            )}
          >
            {props.timeStamp != null ? props.timeStamp.toLocaleString() : "-"}
          </span>
        </div>
      </div>
      <div className="flex w-[90%] min-w-[300px] max-w-[650px] items-center justify-center gap-16 rounded-16 px-20 py-24 shadow-lg md:w-[650px]">
        <div className="ml-64 flex w-0 flex-[1] flex-col overflow-hidden">
          <p className="pb-20 text-14">チーム名</p>
          <p
            className="w-full overflow-hidden truncate whitespace-nowrap text-16"
            title={props.teamName}
          >
            {props.teamName}
          </p>
        </div>
        <div className="flex w-0 flex-[1] flex-col overflow-hidden">
          <p className="pb-20 text-14">所属</p>
          <p
            className="text- w-full overflow-hidden truncate whitespace-nowrap"
            title={props.affiliation}
          >
            {props.affiliation}
          </p>
        </div>
      </div>
    </div>
  );
}
