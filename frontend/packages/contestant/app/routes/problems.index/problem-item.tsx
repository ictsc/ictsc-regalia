import { type ComponentProps, Fragment } from "react";
import { clsx } from "clsx";
import { Button } from "@headlessui/react";
import { Link } from "@tanstack/react-router";
import { Score } from "../../components/score";

type ProblemItemProps = {
  code: string;
  title: string;

  score: ComponentProps<typeof Score>;
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
            "flex w-full max-w-[512px] justify-between gap-24 rounded-16 px-20 py-12 transition",
            active ? "shadow-transparent" : "shadow-lg",
            props.score.rawFullScore
              ? hover
                ? "bg-surface-2"
                : "bg-disabled"
              : hover
                ? "bg-surface-1"
                : "bg-surface-0",
          )}
        >
          <div className="flex flex-col items-start justify-between">
            <h3 className="text-24 font-bold text-primary">{props.code}</h3>
            <p className="line-clamp-1 text-16">{props.title}</p>
          </div>
          <Score {...props.score} />
        </Link>
      )}
    </Button>
  );
}
