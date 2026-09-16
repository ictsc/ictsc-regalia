import { demoProblemState } from "./demo";
import type { Problem } from "../models";
export function problemStatus(
  problem: Problem,
  submitted = false,
  pending = false,
  demoMode = false,
) {
  if (pending) return "pending";
  if (!problem.score) {
    if (submitted) return "pending";
    return demoMode ? demoProblemState(problem.code).status : "unanswered";
  }
  if (problem.score.score >= problem.maxScore) return "complete";
  return problem.score.score > 0 ? "partial" : "zero";
}
export const statusLabels = {
  unanswered: "未回答",
  pending: "採点中",
  complete: "満点",
  partial: "部分点",
  zero: "0点",
};
