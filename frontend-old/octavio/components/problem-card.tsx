import Link from "next/link";

import clsx from "clsx";

import { preRoundMode } from "@/components/_const";
import { Problem } from "@/types/Problem";

type Props = {
  problem: Problem;
};

function ProblemCard({ problem }: Props) {
  let problemText = "";

  if (preRoundMode) {
    if (problem.is_answered) {
      problemText = "font-bold text-amber-500";
    }
  } else {
    if (
      problem.current_point >=
      (problem.solved_criterion ?? problem.current_point)
    ) {
      problemText = "font-bold text-gray-500";
    }
    if (problem.current_point === problem.point) {
      problemText = "font-bold text-amber-500";
    }
  }

  return (
    <Link
      href={`/problems/${problem.code}`}
      className="problem-card border p-4 hover:bg-base-200 hover:cursor-pointer rounded-md shadow-sm min-h-[212px] justify-between flex flex-col"
    >
      <div>
        <span className="problem-code font-bold text-2xl text-primary pr-2">
          {problem.code}
        </span>
        <span className="problem-title text-xl font-bold">{problem.title}</span>
      </div>
      <div>
        <div className={clsx("problem-point text-right", problemText)}>
          {preRoundMode
            ? `${problem.point}pt`
            : `${problem.current_point}/${problem.point}pt`}
        </div>
        <div className="font-bold text-primary">問題文へ→</div>
      </div>
    </Link>
  );
}

export default ProblemCard;
