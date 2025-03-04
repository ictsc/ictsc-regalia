import type { Problem } from "@ictsc/proto/contestant/v1";
import { ProblemItem } from "./problem-item";

type PageProps = {
  problems: Problem[];
};

export function ProblemsPage(props: PageProps) {
  return (
    <div className="mx-16 mt-64 flex justify-center">
      <ul className="grid grid-flow-row grid-cols-1 gap-x-40 gap-y-24 lg:grid-cols-2">
        {props.problems.map((problem) => (
          <li key={problem.code}>
            <ProblemItem
              code={problem.code}
              title={problem.title}
              score={{
                maxScore: problem.maxScore,
                score: problem.score?.score,
                rawScore: problem.score?.markedScore,
                penalty: problem.score?.penalty,
                fullScore: problem.score?.score === problem.maxScore,
                rawFullScore: problem.score?.markedScore === problem.maxScore,
              }}
            />
          </li>
        ))}
      </ul>
    </div>
  );
}
