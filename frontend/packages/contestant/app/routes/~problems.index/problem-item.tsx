import { Fragment } from "react";
import { clsx } from "clsx";
import { Button } from "@headlessui/react";
import { Link } from "@tanstack/react-router";

type ProblemItemProps = {
  code: string;
  title: string;

  maxScore: number;
  score?: number;
  rawScore?: number;
  penalty?: number;

  fullScore?: boolean;
  rawFullScore?: boolean;
};

export function ProblemItem(props: ProblemItemProps) {
  return (
    <Button as={Fragment}>
      {({ hover, active, disabled }) => (
        <Link
          to="/problems/$code"
          params={{ code: props.code }}
          disabled={disabled}
          className={clsx(
            "flex w-full justify-between gap-24 rounded-16 px-20 py-12 transition",
            active ? "shadow-transparent" : "shadow-lg",
            props.rawFullScore && !hover && "bg-disabled",
            props.rawFullScore && hover && "bg-surface-2",
            !props.rawFullScore && !hover && "bg-surface-0",
            !props.rawFullScore && hover && "bg-surface-1",
          )}
        >
          <div className="flex flex-col items-start justify-between">
            <h3 className="text-24 font-bold text-primary">{props.code}</h3>
            <p className="line-clamp-1 text-16">{props.title}</p>
          </div>
          <div className="flex flex-col gap-4">
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
              <span className="translate-y-[-2px] text-[20px]">/</span>
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
              <p
                className={clsx(
                  "place-self-end",
                  props.penalty == null && "px-8",
                )}
              >
                {props.penalty != null ? props.penalty : "-"}
              </p>
            </div>
          </div>
        </Link>
      )}
    </Button>
  );
}
