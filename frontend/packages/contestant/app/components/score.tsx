import { clsx } from "clsx";

export function Score(props: {
  readonly maxScore: number;
  readonly score?: number;
  readonly rawScore?: number;
  readonly penalty?: number;
  readonly fullScore?: boolean;
  readonly rawFullScore?: boolean;
}) {
  return (
    <div className="flex flex-col">
      <p className="flex flex-row items-baseline gap-4 border-b border-text pb-4 pl-8 *:inline-block">
        <span
          className={clsx(
            "text-24 font-bold",
            props.fullScore && "text-primary",
            props.score == null && "px-16",
          )}
        >
          {props.score != null ? props.score : "-"}
        </span>
        <span className="-translate-y-2 text-20">/</span>
        <span className="text-14 font-bold">{props.maxScore}</span>
      </p>
      <div className="grid grid-cols-[repeat(2,auto)] grid-rows-2 place-content-end gap-4 text-14 font-bold">
        <p>素点:</p>
        <p
          className={clsx(
            "place-self-end",
            props.rawFullScore && "text-primary",
            props.rawScore == null && "px-8",
          )}
        >
          {props.rawScore != null ? props.rawScore : "-"}
        </p>
        <p>減点:</p>
        <p className={clsx("place-self-end", props.penalty == null && "px-8")}>
          {props.penalty != null ? props.penalty : "-"}
        </p>
      </div>
    </div>
  );
}
