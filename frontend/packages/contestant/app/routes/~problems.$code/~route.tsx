import { createFileRoute } from "@tanstack/react-router";
import { fetchProblem } from "../../features/problem";
import { fetchAnswers, submitAnswer } from "../../features/answer";

export const Route = createFileRoute("/problems/$code")({
  loader: ({ context: { transport }, params: { code } }) => {
    const fetchAnswersResult = fetchAnswers(transport, code);
    const answers = fetchAnswersResult.then((r) => r.answers);
    const metadata = fetchAnswersResult.then((r) => r.metadata);
    return {
      problem: fetchProblem(transport, code),
      answers,
      metadata,
      submitAnswer: (body: string) => submitAnswer(transport, code, body),
    };
  },
});
