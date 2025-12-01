import { createFileRoute } from "@tanstack/react-router";
import { fetchProblem } from "../features/problem";
import { fetchAnswer, fetchAnswers, submitAnswer } from "../features/answer";
import { deploy, fetchDeployments } from "../features/deployment";

export const Route = createFileRoute("/problems/$code")({
  loader: ({ context: { transport }, params: { code } }) => {
    const fetchAnswersResult = fetchAnswers(transport, code);
    const answers = fetchAnswersResult.then((r) => r.answers);
    const metadata = fetchAnswersResult.then((r) => r.metadata);
    const deployments = fetchDeployments(transport, code);
    return {
      problem: fetchProblem(transport, code),
      answers,
      metadata,
      submitAnswer: (body: string) => submitAnswer(transport, code, body),
      deployments,
      deploy: () => deploy(transport, code),
      fetchAnswer: (num: number) => fetchAnswer(transport, code, num),
    };
  },
});
