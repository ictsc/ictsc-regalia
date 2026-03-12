import type { Problem } from "@ictsc/proto/contestant/v1";

export type ProblemGroup = {
  type: "submittable" | "not-submittable";
  label: string;
  problems: Problem[];
};

export function groupProblems(problems: Problem[]): ProblemGroup[] {
  const submittable: Problem[] = [];
  const notSubmittable: Problem[] = [];

  for (const problem of problems) {
    if (problem.submissionStatus?.isSubmittable) {
      submittable.push(problem);
    } else {
      notSubmittable.push(problem);
    }
  }

  const groups: ProblemGroup[] = [];

  if (submittable.length > 0) {
    groups.push({ type: "submittable", label: "提出可能", problems: submittable });
  }

  if (notSubmittable.length > 0) {
    groups.push({ type: "not-submittable", label: "提出不可", problems: notSubmittable });
  }

  return groups;
}
