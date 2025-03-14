import { clsx } from "clsx";
import { Title } from "../../components/title";

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
      <div className="my-40 flex size-full flex-col items-center px-20 md:flex-nowrap">
        {props.ranking.map((ranking, index) => {
          const { rank, score } = ranking;

          return (
            <div
              key={index}
              className="flex w-full flex-col gap-x-40 pb-64 lg:flex-row lg:items-center"
            >
              <div className="flex w-[175px] flex-col">
                <div className="flex justify-between gap-x-40">
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
                <div className="flex flex-row justify-between gap-4">
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
              <div className="flex w-full min-w-[300px] max-w-[650px] flex-col justify-center gap-16 rounded-16 px-20 py-24 shadow-lg sm:flex-row sm:items-center">
                <div className="flex flex-1 flex-col justify-self-center overflow-hidden">
                  <p className="text-14">チーム名</p>
                  <p
                    className="mt-20 w-full overflow-hidden truncate whitespace-nowrap text-16"
                    title={ranking.teamName}
                  >
                    {ranking.teamName}
                  </p>
                </div>
                <div className="flex flex-1 flex-col overflow-hidden">
                  <p className="text-14">所属</p>
                  <p
                    className="mt-20 w-full overflow-hidden truncate whitespace-nowrap text-16"
                    title={ranking.organization}
                  >
                    {ranking.organization}
                  </p>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </>
  );
}
