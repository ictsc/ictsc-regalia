import { startTransition, Suspense, use, useDeferredValue } from "react";
import { createLazyFileRoute, useRouter } from "@tanstack/react-router";
import type { ProblemDetail } from "../../features/problem";
import type { Answer } from "../../features/answer";
import { protoScoreToProps } from "../../features/score";
import * as View from "./page";

export const Route = createLazyFileRoute("/problems/$code")({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  const { problem, answers, metadata, submitAnswer } = Route.useLoaderData();

  const redeployable = useRedeployable(problem);
  const deferredMetadata = use(useDeferredValue(metadata));

  const deferredAnswers = useDeferredValue(answers);

  return (
    <View.Page
      onTabChange={() => {
        startTransition(() => router.load());
      }}
      redeployable={redeployable}
      content={<Content problem={problem} />}
      submissionForm={
        <View.SubmissionForm
          action={async (body) => {
            try {
              await submitAnswer(body);
            } catch (e) {
              console.error(e);
              return "failure";
            }
            await router.invalidate();
            return "success";
          }}
          submitInterval={deferredMetadata.submitIntervalSeconds}
          lastSubmittedAt={deferredMetadata.lastSubmittedAt}
        />
      }
      submissionList={
        <Suspense>
          <SubmissionList
            isPending={deferredAnswers != answers}
            problemPromise={problem}
            answersPromise={deferredAnswers}
          />
        </Suspense>
      }
    />
  );
}

function useRedeployable(problemPromise: Promise<ProblemDetail>) {
  const problem = use(useDeferredValue(problemPromise));
  return problem.redeployable;
}

function Content(props: { problem: Promise<ProblemDetail> }) {
  const problem = use(useDeferredValue(props.problem));
  return <View.Content {...problem} />;
}

function SubmissionList(props: {
  isPending: boolean;
  problemPromise: Promise<ProblemDetail>;
  answersPromise: Promise<Answer[]>;
}) {
  const problem = use(useDeferredValue(props.problemPromise));
  const answers = use(props.answersPromise);
  if (answers.length === 0) {
    return <View.EmptySubmissionList />;
  }
  return (
    <View.SubmissionList isPending={props.isPending}>
      {answers.map((answer) => (
        <View.SubmissionListItem
          key={answer.id}
          id={answer.id}
          submittedAt={answer.submittedAt}
          score={protoScoreToProps(problem.maxScore, answer?.score)}
        />
      ))}
    </View.SubmissionList>
  );
}
