import type { Problem } from "@ictsc/proto/contestant/v1";
import { ProblemItem } from "./problem-item";
import { protoScoreToProps } from "../../features/score";
import { Title } from "../../components/title";
import { groupProblems } from "../../features/problem/group";

type PageProps = {
  problems: Problem[];
};

export function ProblemsPage(props: PageProps) {
  const groups = groupProblems(props.problems);

  return (
    <>
      <Title>問題一覧</Title>
      <div className="mx-16 my-64 flex flex-col gap-48">
        {groups.map((group) => (
          <section
            key={group.type}
            className={group.type === "not-submittable" ? "mt-16" : undefined}
          >
            <h2 className="text-24 mb-16 font-bold">{group.label}</h2>
            <ul className="grid grid-flow-row grid-cols-1 gap-x-40 gap-y-24 lg:grid-cols-2">
              {group.problems.map((problem) => (
                <li key={problem.code}>
                  <ProblemItem
                    code={problem.code}
                    title={problem.title}
                    score={protoScoreToProps(problem.maxScore, problem.score)}
                    submissionStatus={problem.submissionStatus}
                  />
                </li>
              ))}
            </ul>
          </section>
        ))}
      </div>
    </>
  );
}
