import { expectData, type ApiClient } from "@ictsc/api";
import type { ScoreModel } from "./models";

export type Answer = {
  readonly id: number;
  readonly submittedAt: string;
  readonly contentCommit: string;
  readonly score?: ScoreModel;
};

export type AnswerMetadata = {
  readonly submitIntervalSeconds: number;
  readonly lastSubmittedAt: string;
};

type FetchAnswersResult = {
  readonly answers: Answer[];
  readonly metadata: AnswerMetadata;
};

export async function fetchAnswers(
  client: ApiClient,
  problemCode: string,
): Promise<FetchAnswersResult> {
  const response = expectData(
    await client.GET("/api/v1/contestant/problems/{problem_code}/answers", {
      params: { path: { problem_code: problemCode } },
    }),
  );
  return {
    answers: response.answers.map((answer) => ({
      id: answer.number,
      submittedAt: answer.submitted_at,
      contentCommit: answer.content_commit,
      score: answer.score == null ? undefined : mapScore(answer.score),
    })),
    metadata: {
      submitIntervalSeconds: response.submit_interval_seconds,
      lastSubmittedAt: response.last_submitted_at ?? "",
    },
  };
}

type FetchAnswerResult = {
  readonly answerBody: string;
  readonly submittedAtString: string;
};

export async function fetchAnswer(
  client: ApiClient,
  problemCode: string,
  answerNumber: number,
): Promise<FetchAnswerResult> {
  const response = expectData(
    await client.GET(
      "/api/v1/contestant/problems/{problem_code}/answers/{answer_number}",
      {
        params: {
          path: {
            problem_code: problemCode,
            answer_number: answerNumber,
          },
        },
      },
    ),
  );
  return {
    answerBody: response.answer.body.body,
    submittedAtString: response.answer.submitted_at,
  };
}

export async function submitAnswer(
  client: ApiClient,
  problemCode: string,
  body: string,
): Promise<void> {
  const response = expectData(
    await client.POST("/api/v1/contestant/problems/{problem_code}/answers", {
      params: { path: { problem_code: problemCode } },
      body: { body },
    }),
  );
  void response;
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
