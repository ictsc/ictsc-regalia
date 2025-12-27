import { Fragment } from "react";
import { clsx } from "clsx";
import { Title } from "../components/title";

type Rank = {
  rank: number;
  teamName: string;
  organization: string;
  score: number;
  lastEffectiveSubmitAt?: string;
};

type RankingProps = {
  ranking: Rank[];
};

const formatter = new Intl.DateTimeFormat("ja-JP", {
  dateStyle: "short",
  timeStyle: "medium",
});

const rankPruralRule = new Intl.PluralRules("en-US", { type: "ordinal" });

export function RankingPage(props: RankingProps) {
  return (
    <>
      <Title>ランキング</Title>
      <div className="my-40 flex size-full flex-col items-center px-20">
        <div className="grid w-full max-w-screen-md grid-cols-1 gap-40 md:grid-cols-[auto_1fr]">
          {props.ranking.map((ranking, index) => {
            const { rank, score } = ranking;

            return (
              <Fragment key={index}>
                <div className="flex flex-col self-center">
                  <div className="flex gap-x-40 md:justify-between">
                    {/* ランキング表示 */}
                    <div className="flex flex-row items-baseline">
                      <span className={clsx("text-32 font-bold")}>{rank}</span>
                      <p className="text-14 font-bold">
                        {
                          {
                            zero: "",
                            one: "st",
                            two: "nd",
                            few: "rd",
                            many: "th",
                            other: "th",
                          }[rankPruralRule.select(rank)]
                        }
                      </p>
                    </div>
                    {/* トータルスコア表示 */}
                    <div className="flex flex-row items-baseline gap-8">
                      <span
                        className={clsx(
                          "text-32 font-bold",
                          rank === 1 && "text-primary",
                          score == null && "px-32",
                        )}
                      >
                        {score != null ? score : "-"}
                      </span>
                      <p className="text-14 font-bold">pt</p>
                    </div>
                  </div>
                  <div className="flex flex-row gap-4 md:justify-between">
                    <p className="text-12">Time:</p>
                    <span
                      className={clsx(
                        "text-12 font-bold",
                        ranking.lastEffectiveSubmitAt == null && "px-12",
                      )}
                    >
                      {ranking.lastEffectiveSubmitAt
                        ? formatter.format(
                            new Date(ranking.lastEffectiveSubmitAt),
                          )
                        : "-"}
                    </span>
                  </div>
                </div>
                <div className="rounded-16 flex flex-col justify-center gap-16 p-24 shadow-lg sm:flex-row sm:items-center">
                  <div className="flex flex-1 sm:justify-center">
                    <div className="flex flex-col justify-self-center overflow-hidden">
                      <p className="text-14">チーム名</p>
                      <p
                        className="text-16 mt-20 w-full truncate overflow-hidden whitespace-nowrap"
                        title={ranking.teamName}
                      >
                        {ranking.teamName}
                      </p>
                    </div>
                  </div>
                  <div className="flex flex-1 sm:justify-center">
                    <div className="flex flex-col overflow-hidden">
                      <p className="text-14">所属</p>
                      <p
                        className="text-16 mt-20 w-full truncate overflow-hidden whitespace-nowrap"
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
      </div>
    </>
  );
}
