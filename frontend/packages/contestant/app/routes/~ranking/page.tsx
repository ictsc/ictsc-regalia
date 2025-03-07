import { clsx } from "clsx";
import type { Rank } from "@ictsc/proto/contestant/v1";
import { Fragment } from "react";


type RankingProps = {
  ranking: Rank[];
};

export function RankingPage(props: RankingProps) {
  return (
    <div className="flex flex-col h-[800px] items-center pl-8 md:flex-nowrap">
      {props.ranking.map((ranking, index) => {
        const rank = Number(ranking.rank);
        const score =  Number(ranking.score) ;

        return (
          <Fragment key={index}>
          <div className="flex flex-wrap  items-center gap-x-40 pb-64">
            <div className="flex flex-col  ">
              <div className="flex w-[175px] gap-x-40">
                {/* ランキング表示 */}
                <div className="flex flex-row items-baseline">
                  <span className={clsx("text-32 font-bold")}>
                    {rank}
                  </span>
                  <p className="text-14 font-bold">st</p>
                </div>
                {/* トータルスコア表示 */}
                <div className="flex flex-row items-baseline gap-8">
                  <span
                    className={clsx(
                      "text-32 font-bold",
                      rank === 1&& "text-primary",
                      score == null && "px-32",
                    )}
                  >
                    {score != null ? score : "-"}
                  </span>
                  <p className="text-14 font-bold">pt</p>
                </div>
              </div>
              <div className="flex flex-row gap-4">
                <p className="text-12">Time:</p>
                <span
                  className={clsx(
                    "text-12 font-bold",
                    ranking.timestamp == null && "px-12",
                  )}
                >
                  {ranking.timestamp != null ? ranking.timestamp.toLocaleString() : "-"}
                </span>
              </div>
            </div>
            <div className="flex w-[90%] min-w-[300px] max-w-[650px] items-center justify-center gap-16 rounded-16 px-20 py-24 shadow-lg md:w-[650px]">
              <div className="ml-64 flex w-0 flex-[1] flex-col overflow-hidden">
                <p className="pb-20 text-14">チーム名</p>
                <p
                  className="w-full overflow-hidden truncate whitespace-nowrap text-16"
                  title={ranking.teamName}
                >
                  {ranking.teamName}
                </p>
              </div>
              <div className="flex w-0 flex-[1] flex-col overflow-hidden">
                <p className="pb-20 text-14">所属</p>
                <p
                  className="text- w-full overflow-hidden truncate whitespace-nowrap"
                  title={ranking.organization}
                >
                  {ranking.organization}
                </p>
              </div>
            </div>
          </div>
          </Fragment>
        );
      })}
    </div>
  );
}
