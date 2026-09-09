import { expectData, type ApiClient } from "@ictsc/api";
import type {
  Problem,
  ScoreModel,
  SubmissionStatus,
  ScheduleEntry,
} from "../models";

export type { Problem } from "../models";

export async function fetchProblems(client: ApiClient): Promise<Problem[]> {
  const [problemResponse, sectionResponse] = await Promise.all([
    client.GET("/api/v1/contestant/problems").then(expectData),
    client.GET("/api/v1/contestant/sections").then(expectData),
  ]);
  const schedules = new Map<string, ScheduleEntry>(
    sectionResponse.sections.map((section) => [
      section.slug,
      {
        name: section.slug,
        startAt: section.beginning,
        endAt: section.ending,
      },
    ]),
  );
  return problemResponse.problems.map((problem) => ({
    code: problem.code,
    title: problem.title,
    maxScore: problem.max_score,
    category: problem.category,
    sectionSlug: problem.section_slug,
    score: problem.score == null ? undefined : mapScore(problem.score),
    submissionableSchedules: schedules.has(problem.section_slug)
      ? [schedules.get(problem.section_slug)!]
      : [],
    submissionStatus: mapSubmissionStatus(problem.submission_status),
  }));
}

export type ProblemDetail = {
  code: string;
  title: string;
  category: string;
  maxScore: number;
  redeployable: boolean;
  penaltyThreashold: number;
  body: string;
  submissionStatus?: SubmissionStatus;
};

export async function fetchProblem(
  client: ApiClient,
  code: string,
): Promise<ProblemDetail> {
  const response = expectData(
    await client.GET("/api/v1/contestant/problems/{problem_code}", {
      params: { path: { problem_code: code } },
    }),
  );
  const problem = response.problem;
  return {
    code: problem.code,
    title: problem.title,
    category: problem.category,
    maxScore: problem.max_score,
    redeployable: problem.deployment?.redeployable ?? false,
    penaltyThreashold: problem.deployment?.penalty_threshold ?? 0,
    body: problem.body,
    submissionStatus: mapSubmissionStatus(problem.submission_status),
  };
}

function mapScore(score: {
  marked_score: number;
  penalty: number;
  score: number;
  max_score: number;
}): ScoreModel {
  return {
    markedScore: score.marked_score,
    penalty: score.penalty,
    score: score.score,
    maxScore: score.max_score,
  };
}

function mapSubmissionStatus(status: {
  is_submittable: boolean;
  submittable_from: string | null;
  submittable_until: string | null;
}): SubmissionStatus {
  return {
    isSubmittable: status.is_submittable,
    submittableFrom: status.submittable_from ?? undefined,
    submittableUntil: status.submittable_until ?? undefined,
  };
}
