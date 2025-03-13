import {
  startTransition,
  Suspense,
  use,
  useActionState,
  useDeferredValue,
  useOptimistic,
} from "react";
import { createLazyFileRoute, useRouter } from "@tanstack/react-router";
import type { ProblemDetail } from "../../features/problem";
import type { Answer } from "../../features/answer";
import { protoScoreToProps } from "../../features/score";
import * as View from "./page";
import { Deployment } from "@app/features/deployment";
import { DeploymentStatus } from "@ictsc/proto/contestant/v1";

export const Route = createLazyFileRoute("/problems/$code")({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  const { problem, answers, metadata, submitAnswer, deployments, deploy } =
    Route.useLoaderData();

  const redeployable = useRedeployable(problem);
  const deferredMetadata = use(useDeferredValue(metadata));

  const deferredAnswers = useDeferredValue(answers);

  const deferredDeployments = useDeferredValue(deployments);

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
            isPending={deferredAnswers !== answers}
            problemPromise={problem}
            answersPromise={deferredAnswers}
          />
        </Suspense>
      }
      deploymentList={
        <Suspense>
          <Deployments
            isPending={deferredDeployments !== deployments}
            deployments={deferredDeployments}
            deploy={deploy}
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
    <View.SubmissionListContainer>
      {answers.length === 0 ? (
        <View.EmptySubmissionList />
      ) : (
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
      )}
    </View.SubmissionListContainer>
  );
}

function Deployments(props: {
  deployments: Promise<Deployment[]>;
  isPending: boolean;
  deploy: () => Promise<void>;
}) {
  const router = useRouter();
  const [deployments, optimisticSetDeployments] = useOptimistic(
    use(props.deployments) as (Deployment & { isPending?: boolean })[],
  );
  const canRedeploy =
    (deployments?.[0].status ?? DeploymentStatus.DEPLOYED) ===
    DeploymentStatus.DEPLOYED;

  const [lastResult, action, isActionPending] = useActionState(
    async (_prev: unknown, _action: "redeploy") => {
      optimisticSetDeployments((ds) => [
        {
          isPending: true,
          revision: ds.length + 1,
          status: DeploymentStatus.DEPLOYING,
          requestedAt: new Date().toISOString(),
          allowedDeploymentCount: (ds?.[0].allowedDeploymentCount ?? 1) - 1,
          thresholdExceeded: ds?.[0].thresholdExceeded ?? false,
          penalty: ds?.[0].penalty ?? 0,
        },
        ...ds,
      ]);
      try {
        await props.deploy();
      } catch (e) {
        console.error(e);
        return "再展開に失敗しました";
      }
      await router.invalidate();
      return null;
    },
    null,
  );

  return (
    <View.Deployments
      canRedeploy={canRedeploy}
      isRedeploying={isActionPending}
      redeploy={() => startTransition(() => action("redeploy"))}
      error={lastResult}
      list={
        deployments.length === 0 ? (
          <View.EmptyDeploymentList />
        ) : (
          <View.DeploymentList isPending={props.isPending}>
            {deployments.map((deployment) => (
              <View.DeploymentItem key={deployment.revision} {...deployment} />
            ))}
          </View.DeploymentList>
        )
      }
    />
  );
}
