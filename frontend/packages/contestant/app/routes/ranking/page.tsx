import { clsx } from "clsx";
import { Button } from "@headlessui/react";
import { Link } from "@tanstack/react-router";

type RankingItemProps = {
  rank?: number
  timeStamp?: Date
  teamName: string;
  affiliation: string;
  score?: number;
  penalty?: number;
  fullScore?: boolean;
  rawFullScore?: boolean;
};

export function RankingItem(props: RankingItemProps) {
  return (
    <div className=" w-full flex flex-wrap md:flex-nowrap gap-x-40 pb-64 pl-8 items-center justify-center"> {/* ここで横幅を固定 */}
      <div>
        <div className="w-[175px] flex flex-row items-baseline justify-between ">
          {/* ランキング表示 */}
          <div className=" flex flex-row items-baseline">
            <span
              className={clsx(
                "font-bold text-32"
              )}
            >
                {props.rank != null ? props.rank : "-"}
            </span>
            <p className="text-14 font-bold">st</p>
          </div>
          {/* トータルスコア表示 */}
          <div className="flex flex-row gap-8 items-baseline">
            <span
              className={clsx(
                "text-32 font-bold",
                props.rank === 1 && "text-primary",
                props.score == null && "px-32"
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
              props.timeStamp == null && "px-12"
            )}
          >
            {props.timeStamp != null ? props.timeStamp.toLocaleString() : "-"}
          </span>
        </div>
      </div>
      <div className="flex gap-16 items-center justify-center shadow-lg min-w-[300px] max-w-[650px] w-[90%] md:w-[650px] rounded-16 px-20 py-24">
          <div className="flex flex-col flex-[1] w-0 overflow-hidden ml-64 ">
            <p className="pb-20 text-14">チーム名</p>
            <p className="text-16 truncate w-full overflow-hidden whitespace-nowrap" title={props.teamName}>{props.teamName}</p>
          </div>
          <div className="flex flex-col flex-[1] w-0 overflow-hidden">
            <p className="pb-20 text-14">所属</p>
            <p className="text- truncate w-full overflow-hidden whitespace-nowrap" title={props.affiliation}>{props.affiliation}</p>
          </div>
      </div>
    </div>
  );
}