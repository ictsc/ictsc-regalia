import { type ComponentProps, Fragment } from "react";
import { clsx } from "clsx";
import { Button } from "@headlessui/react";
import { Link } from "@tanstack/react-router";
import { Score } from "../../components/score";
import type { SubmissionStatus } from "@ictsc/proto/contestant/v1";

type ProblemItemProps = {
  code: string;
  title: string;

  score: ComponentProps<typeof Score>;
  submissionStatus?: SubmissionStatus;
};

export function ProblemItem(props: ProblemItemProps) {
  const isSubmittable = props.submissionStatus?.isSubmittable ?? true;

  return (
    <Button as={Fragment}>
      {({ hover, active, disabled }) => (
        <Link
          to="/problems/$code"
          params={{ code: props.code }}
          disabled={disabled}
          className={clsx(
            "rounded-16 flex w-full max-w-[512px] justify-between gap-24 px-20 py-12 transition",
            active ? "shadow-transparent" : "shadow-lg",
            props.score.rawFullScore
              ? hover
                ? "bg-surface-2"
                : "bg-disabled"
              : hover
                ? "bg-surface-1"
                : "bg-surface-0",
            !isSubmittable && "opacity-50 grayscale",
          )}
        >
          <div className="flex flex-col items-start justify-between gap-4">
            <div className="flex flex-col">
              <h3 className="text-24 text-primary font-bold">{props.code}</h3>
              <p className="text-16 line-clamp-1">{props.title}</p>
            </div>
          </div>
          <Score {...props.score} />
        </Link>
      )}
    </Button>
  );
}
