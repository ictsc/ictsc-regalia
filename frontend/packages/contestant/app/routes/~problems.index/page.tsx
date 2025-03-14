import type { Problem } from "@ictsc/proto/contestant/v1";
import { ProblemItem } from "./problem-item";
import { protoScoreToProps } from "../../features/score";
import { Title } from "../../components/title";

type PageProps = {
  problems: Problem[];
};

export function ProblemsPage(props: PageProps) {
  return (
    <>
      <Title>問題一覧</Title>
      <div className="mx-16 my-64 flex justify-center">
        <ul className="grid grid-flow-row grid-cols-1 gap-x-40 gap-y-24 lg:grid-cols-2">
          {props.problems.map((problem) => (
            <li key={problem.code}>
              <ProblemItem
                code={problem.code}
                title={problem.title}
                score={protoScoreToProps(problem.maxScore, problem.score)}
              />
            </li>
          ))}
        </ul>
      </div>
    </>
  );
}
