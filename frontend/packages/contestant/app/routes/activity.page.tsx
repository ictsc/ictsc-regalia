import type { ComponentProps } from "react";
import { clsx } from "clsx";
import { Title } from "../components/title";
import { Score } from "../components/score";

type ActivityEntry = {
  problemCode: string;
  problemTitle: string;
  answerId: number;
  submittedAt: string;
  score: ComponentProps<typeof Score>;
  scored: boolean;
};

type ActivityPageProps = {
  entries: ActivityEntry[];
};

const formatter = new Intl.DateTimeFormat("ja-JP", {
  dateStyle: "short",
  timeStyle: "medium",
});

export function ActivityPage(props: ActivityPageProps) {
  return (
    <>
      <Title>アクティビティ</Title>
      <div className="my-40 flex size-full flex-col items-center px-20">
        {props.entries.length === 0 ? (
          <p className="text-16 text-text mt-64">提出履歴がありません</p>
        ) : (
          <div className="flex w-full max-w-screen-md flex-col">
            {props.entries.map((entry, index) => (
              <div key={index} className="flex flex-row items-center">
                {/* 縦線+丸 */}
                <div className="flex w-12 shrink-0 flex-col items-center self-stretch">
                  {index === 0 ? (
                    <div
                      className="grow w-2"
                      style={{
                        backgroundImage:
                          "linear-gradient(to bottom, transparent, var(--color-text) 100%)",
                        opacity: 0.1,
                        maskImage:
                          "repeating-linear-gradient(to bottom, black 0 3px, transparent 3px 6px)",
                      }}
                    />
                  ) : (
                    <div className="bg-text/10 w-2 grow" />
                  )}
                  <div
                    className={clsx(
                      "size-12 shrink-0 rounded-full",
                      entry.scored ? "bg-primary" : "bg-text/20",
                    )}
                  />
                  <div
                    className={clsx(
                      "w-2 grow",
                      index < props.entries.length - 1
                        ? "bg-text/10"
                        : "bg-transparent",
                    )}
                  />
                </div>
                {/* 時刻 */}
                <div
                  className={clsx(
                    "text-12 ml-12 shrink-0 whitespace-nowrap",
                    !entry.scored && "text-text/40",
                  )}
                >
                  {formatter.format(new Date(entry.submittedAt))}
                </div>
                {/* 矢印 */}
                <div
                  className={clsx(
                    "text-16 shrink-0 px-8",
                    entry.scored ? "text-primary" : "text-text/20",
                  )}
                >
                  ←
                </div>
                {/* カード */}
                <div
                  className={clsx(
                    "rounded-16 my-12 flex min-w-0 flex-1 flex-col gap-12 p-20 shadow-lg sm:flex-row sm:items-center sm:justify-between",
                    !entry.scored && "opacity-50",
                  )}
                >
                  <div className="flex min-w-0 flex-1 flex-col gap-4">
                    <p className="text-14 truncate font-bold">
                      {entry.problemCode}: {entry.problemTitle}
                    </p>
                    <p className="text-12">提出 #{entry.answerId}</p>
                  </div>
                  {entry.scored ? (
                    <Score {...entry.score} />
                  ) : (
                    <p className="text-14 font-bold">採点中</p>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  );
}
