import type { ComponentProps, ReactNode } from "react";
import { clsx } from "clsx";
import { Score } from "../../components/score";

export function SubmissionListContainer(props: { children?: ReactNode }) {
  return (
    <div className="size-full rounded-12 bg-surface-1 py-12">
      <div className="size-full overflow-y-auto px-12 [scrollbar-gutter:stable_both-edges]">
        {props.children}
      </div>
    </div>
  );
}

export function SubmissionList(props: {
  readonly isPending?: boolean;
  readonly children?: ReactNode;
}) {
  return (
    <ul
      className={clsx(
        "flex size-full flex-col gap-16 py-12",
        props.isPending && "opacity-75",
      )}
    >
      {props.children}
    </ul>
  );
}

export function EmptySubmissionList() {
  return (
    <div className="grid size-full place-items-center text-16 font-bold text-text">
      解答はまだありません！
    </div>
  );
}

const submissionListDateTimeFormatter = new Intl.DateTimeFormat("ja-JP", {
  dateStyle: "medium",
  timeStyle: "short",
});

export function SubmissionListItem(props: {
  readonly id: number;
  readonly submittedAt: string;
  readonly score: ComponentProps<typeof Score>;
}) {
  return (
    <li className="flex justify-between gap-8 rounded-12 bg-surface-0 p-16">
      <div className="flex flex-col">
        <h2 className="text-20 font-bold text-[#000]">#{props.id}</h2>
        <h3 className="text-12">
          提出:{" "}
          {submissionListDateTimeFormatter.format(new Date(props.submittedAt))}
        </h3>
        <p className="mt-4 text-12"></p>
      </div>
      <Score {...props.score} />
    </li>
  );
}
