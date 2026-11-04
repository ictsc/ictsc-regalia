import type { Problem } from "../models";
/**
 * resultsHidden は二次予選モードで最終公開前の状態。採点自体が競技終了後なので
 * 「採点中」ではなく「提出済み」として扱う。
 */
export function problemStatus(
  problem: Problem,
  submitted = false,
  pending = false,
  resultsHidden = false,
): ProblemStatus {
  if (resultsHidden) return submitted ? "submitted" : "unanswered";
  if (pending) return "pending";
  if (!problem.score) return submitted ? "pending" : "unanswered";
  if (problem.score.score >= problem.maxScore) return "complete";
  return problem.score.score > 0 ? "partial" : "zero";
}
export const statusLabels = {
  unanswered: "未回答",
  pending: "採点中",
  submitted: "提出済み",
  complete: "満点",
  partial: "部分点",
  zero: "0点",
};
export type ProblemStatus = keyof typeof statusLabels;

/** 現在のモードで取り得る状態だけを絞り込みの選択肢に出す。 */
export function statusOptions(
  resultsHidden: boolean,
): Partial<Record<ProblemStatus, string>> {
  if (resultsHidden)
    return {
      unanswered: statusLabels.unanswered,
      submitted: statusLabels.submitted,
    };
  const { submitted: _submitted, ...options } = statusLabels;
  return options;
}
