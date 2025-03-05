import { startTransition, Suspense, use, useDeferredValue } from "react";
import { createLazyFileRoute, useRouter } from "@tanstack/react-router";
import type { ProblemDetail } from "../../features/problem";
import type { Answer } from "../../features/answer";
import * as View from "./page";

export const Route = createLazyFileRoute("/problems/$code")({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  const { problem, answers, submitAnswer } = Route.useLoaderData();

  const redeployable = useRedeployable(problem);
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
        />
      }
      submissionList={
        <Suspense>
          <SubmissionList answersPromise={answers} />
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

function SubmissionList(props: { answersPromise: Promise<Answer[]> }) {
  const answers = use(useDeferredValue(props.answersPromise));
  if (answers.length === 0) {
    return <View.EmptySubmissionList />;
  }
  return (
    <View.SubmissionList>
      {answers.map((answer) => (
        <View.SubmissionListItem
          key={answer.id}
          id={answer.id}
          submittedAt={answer.submittedAt}
          score={{ maxScore: 100 }}
        />
      ))}
    </View.SubmissionList>
  );
}
