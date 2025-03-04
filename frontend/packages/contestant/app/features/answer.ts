import { type Transport, createClient } from "@connectrpc/connect";
import { AnswerService } from "@ictsc/proto/contestant/v1";
import { timestampDate } from "@bufbuild/protobuf/wkt";

export type Answer = {
  readonly id: number;
  readonly submittedAt: string;
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

export async function submitAnswer(
  transport: Transport,
  problemCode: string,
  body: string,
): Promise<void> {
  const client = createClient(AnswerService, transport);
  await client.submitAnswer({ problemCode, body });
}
