import { type Transport, createClient } from "@connectrpc/connect";
import {
  AnswerService,
  type Score as ProtoScore,
} from "@ictsc/proto/contestant/v1";
import { timestampDate } from "@bufbuild/protobuf/wkt";

export type Answer = {
  readonly id: number;
  readonly submittedAt: string;
  readonly score?: ProtoScore;
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
  transport: Transport,
  problemCode: string,
): Promise<FetchAnswersResult> {
  const client = createClient(AnswerService, transport);
  const {
    answers = [],
    submitInterval,
    lastSubmittedAt,
  } = await client.listAnswers({ problemCode });
  return {
    answers: answers.map((answer) => {
      const submittedAt =
        answer.submittedAt != null
          ? timestampDate(answer.submittedAt).toISOString()
          : "";
      return {
        id: answer.id,
        submittedAt: submittedAt,
        score: answer.score,
      };
    }),
    metadata: {
      submitIntervalSeconds:
        submitInterval != null ? Number(submitInterval.seconds) : 0,
      lastSubmittedAt:
        lastSubmittedAt != null
          ? timestampDate(lastSubmittedAt).toISOString()
          : "",
    },
  };
}

type FetchAnswerResult = {
  readonly answerBody: string;
  readonly submittedAtString: string;
};

export async function fetchAnswer(
  transport: Transport,
  problemCode: string,
  answerNumber: number,
): Promise<FetchAnswerResult> {
  const client = createClient(AnswerService, transport);
  const { answer } = await client.getAnswer({ problemCode, id: answerNumber });
  const submittedAtString = answer?.submittedAt
    ? timestampDate(answer.submittedAt).toISOString()
    : "";
  return {
    answerBody: answer?.body?.body.value?.body ?? "",
    submittedAtString,
  };
}

export async function submitAnswer(
  transport: Transport,
  problemCode: string,
  body: string,
): Promise<void> {
  const client = createClient(AnswerService, transport);
  await client.submitAnswer({ problemCode, body });
}
