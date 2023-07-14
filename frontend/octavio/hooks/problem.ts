import matter from "gray-matter";

import useProblems from "@/hooks/problems";
import { Matter } from "@/types/Problem";

const useProblem = (code: string | null) => {
  const { problems } = useProblems();

  const problem = problems.find((prob) => prob.code === code) ?? null;

  if (problem === null) {
    return {
      matter: null,
      problem: null,
      isLoading: false,
    };
  }

  const matterResult = matter(problem.body ?? "");
  const matterData = matterResult.data as Matter;
  const newProblem = { ...problem, body: matterResult.content ?? "" };

  return {
    matter: matterData,
    problem: newProblem,
    isLoading: false,
  };
};

export default useProblem;
