import { startTransition, Suspense, use, useDeferredValue } from "react";
import { createLazyFileRoute, useRouter } from "@tanstack/react-router";
import type { ProblemDetail } from "../../features/problem";
import type { Answer } from "../../features/answer";
import { protoScoreToProps } from "../../features/score";
import * as View from "./page";
import { Deployment } from "@app/features/problem/deployment";

export const Route = createLazyFileRoute("/problems/$code")({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  const { problem, answers, metadata, submitAnswer } = Route.useLoaderData();

  const redeployable = useRedeployable(problem);
  const deferredMetadata = use(useDeferredValue(metadata));

  const deferredAnswers = useDeferredValue(answers);

  const remainRedeployCount = useRemainRedeployCount(problem);
  const latestPenalty = useLatestPenalty(problem);

  return (
    <View.Page
      onTabChange={() => {
        startTransition(() => router.load());
      }}
      redeployable={redeployable}
      remainRedeployCount={remainRedeployCount}
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
          latestPenalty={latestPenalty}
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
      deploymentList={
        <Suspense>
          <DeploymentList problem={problem} />
        </Suspense>
      }
    />
  );
}

function useRedeployable(problemPromise: Promise<ProblemDetail>) {
  const problem = use(useDeferredValue(problemPromise));
  return problem.deployment.redeployable;
}

function useRemainRedeployCount(problemPromise: Promise<ProblemDetail>) {
  const problem = use(useDeferredValue(problemPromise));
  return problem.deployment.maxRedeployment - problem.deployment.events.length;
}

function useLatestPenalty(problemPromise: Promise<ProblemDetail>) {
  const problem = use(useDeferredValue(problemPromise));
  if (problem.deployment.events.length === 0) {
    return 0;
  }
  return problem.deployment.events[0].totalPenalty;
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
    return <View.EmptyListContainer message="解答はまだありません！" />;
  }
  return (
    <View.ListContainer isPending={props.isPending}>
      {answers.map((answer) => (
        <View.SubmissionListItem
          key={answer.id}
          id={answer.id}
          submittedAt={answer.submittedAt}
          score={protoScoreToProps(problem.maxScore, answer?.score)}
        />
      ))}
    </View.ListContainer>
  );
}

function DeploymentList(props: { problem: Promise<ProblemDetail> }) {
  const problem = use(useDeferredValue(props.problem));
  const deployment: Deployment = problem.deployment;

  if (deployment.events.length === 0) {
    return <View.EmptyListContainer message="再展開はまだありません" />;
  }
  const latestEvent = deployment.events[0];
  return (
    <View.ListContainer>
      {deployment.events.map((event) => (
        <View.DeploymentListItem
          key={event.revision}
          event={event}
          maxRedeployment={deployment.maxRedeployment}
          deploymentDetail={{
            revision: event.revision,
            remainingRedeploys: deployment.maxRedeployment - event.revision,
            exceededRedeployLimit: event.revision > deployment.maxRedeployment,
            totalPenalty: event.totalPenalty,
          }}
          isDeploying={event.isDeploying}
          isLatest={event.revision === latestEvent.revision}
        />
      ))}
    </View.ListContainer>
  );
}
