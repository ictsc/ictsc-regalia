import {
  startTransition,
  use,
  useActionState,
  useDeferredValue,
  useEffect,
  useOptimistic,
} from "react";
import { createLazyFileRoute, useRouter } from "@tanstack/react-router";
import { DeploymentStatus } from "@ictsc/proto/contestant/v1";
import { timestampDate } from "@bufbuild/protobuf/wkt";
import type { ProblemDetail } from "../features/problem";
import type { Answer } from "../features/answer";
import { protoScoreToProps } from "../features/score";
import type { Deployment } from "../features/deployment";
import * as View from "./problems.$code/page";

export const Route = createLazyFileRoute("/problems/$code")({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  const {
    problem,
    answers,
    metadata,
    submitAnswer,
    deployments,
    deploy,
    fetchAnswer,
  } = Route.useLoaderData();

  const redeployable = useRedeployable(problem);
  const deferredMetadata = useDeferredValue(metadata);
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
        <SubmissionForm
          submitAnswer={submitAnswer}
          problemPromise={problem}
          metatataPromise={deferredMetadata}
        />
      }
      submissionList={
        <SubmissionList
          isPending={deferredAnswers !== answers}
          problemPromise={problem}
          answersPromise={deferredAnswers}
          fetchAnswer={fetchAnswer}
        />
      }
      deploymentList={
        <Deployments
          isPending={deferredDeployments !== deployments}
          deployments={deferredDeployments}
          problemPromise={problem}
          deploy={deploy}
        />
      }
    />
  );
}

function useRedeployable(problemPromise: Promise<ProblemDetail>) {
  const problem = use(useDeferredValue(problemPromise));
  return problem.redeployable;
}

function SubmissionForm(props: {
  submitAnswer: (body: string) => Promise<void>;
  problemPromise: Promise<ProblemDetail>;
  metatataPromise: Promise<{
    submitIntervalSeconds: number;
    lastSubmittedAt: string;
  }>;
}) {
  const router = useRouter();
  const problem = use(props.problemPromise);
  const metadata = use(props.metatataPromise);

  // Convert Timestamp to Date if present
  const submissionStatus = problem.submissionStatus
    ? {
        isSubmittable: problem.submissionStatus.isSubmittable,
        submittableUntil: problem.submissionStatus.submittableUntil
          ? timestampDate(problem.submissionStatus.submittableUntil)
          : undefined,
      }
    : undefined;

  // submittableUntil/submittableFrom に達したら自動リフェッチして提出状態を更新
  useEffect(() => {
    const until = submissionStatus?.submittableUntil;
    const from = problem.submissionStatus?.submittableFrom
      ? timestampDate(problem.submissionStatus.submittableFrom)
      : undefined;
    const target = until ?? from;
    if (target == null) return;
    const ms = target.getTime() - Date.now();
    if (ms <= 0) return;
    const timer = setTimeout(() => {
      startTransition(() => router.invalidate());
    }, ms);
    return () => clearTimeout(timer);
  }, [submissionStatus?.submittableUntil, problem.submissionStatus?.submittableFrom, router]);

  return (
    <View.SubmissionForm
      action={async (body) => {
        try {
          await props.submitAnswer(body);
        } catch (e) {
          console.error(e);
          return "failure";
        }
        await router.invalidate();
        return "success";
      }}
      submitInterval={metadata.submitIntervalSeconds}
      lastSubmittedAt={metadata.lastSubmittedAt}
      storageKey={`/problems/${problem.code}`}
      submissionStatus={submissionStatus}
    />
  );
}

function Content(props: { problem: Promise<ProblemDetail> }) {
  const problem = use(useDeferredValue(props.problem));
  return <View.Content {...problem} />;
}

function SubmissionList(props: {
  isPending: boolean;
  problemPromise: Promise<ProblemDetail>;
  answersPromise: Promise<Answer[]>;
  fetchAnswer: (
    num: number,
  ) => Promise<{ answerBody: string; submittedAtString: string }>;
}) {
  const problem = use(useDeferredValue(props.problemPromise));
  const answers = use(props.answersPromise);

  const formatDate = (date: Date) => {
    const pad = (num: number) => num.toString().padStart(2, "0");

    const year = date.getFullYear();
    const month = pad(date.getMonth() + 1);
    const day = pad(date.getDate());
    const hour = pad(date.getHours());
    const minute = pad(date.getMinutes());
    const second = pad(date.getSeconds());

    return `${year}-${month}-${day}-${hour}-${minute}-${second}`;
  };

  const downloadFile = (filename: string, content: string) => {
    const blob = new Blob([content], { type: "text/markdown" });
    const url = URL.createObjectURL(blob);

    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
  };

  const downloadAnswer = async (num: number) => {
    const { answerBody, submittedAtString } = await props.fetchAnswer(num);
    const submittedAt = new Date(submittedAtString);
    const submittedAtFormattedString = formatDate(submittedAt);
    const filename = `${problem.code}-${submittedAtFormattedString}.md`;

    downloadFile(filename, answerBody);
  };

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
              downloadAnswer={() =>
                startTransition(() => downloadAnswer(answer.id))
              }
            />
          ))}
        </View.SubmissionList>
      )}
    </View.SubmissionListContainer>
  );
}

function Deployments(props: {
  deployments: Promise<Deployment[]>;
  problemPromise: Promise<ProblemDetail>;
  isPending: boolean;
  deploy: () => Promise<void>;
}) {
  const router = useRouter();
  const [deployments, optimisticSetDeployments] = useOptimistic(
    use(props.deployments) as (Deployment & { isPending?: boolean })[],
  );
  const canRedeploy =
    (deployments?.[0]?.status ?? DeploymentStatus.DEPLOYED) ===
    DeploymentStatus.DEPLOYED;
  const problem = use(useDeferredValue(props.problemPromise));
  const allowedDeploymentCount = problem.penaltyThreashold - deployments.length;

  const [lastResult, action, isActionPending] = useActionState(
    async (_prev: unknown, _action: "redeploy") => {
      const timer = setTimeout(() => {
        optimisticSetDeployments((ds) => [
          {
            isPending: true,
            revision: ds.length + 1,
            status: DeploymentStatus.DEPLOYING,
            requestedAt: new Date().toISOString(),
            allowedDeploymentCount: (ds?.[0]?.allowedDeploymentCount ?? 1) - 1,
            thresholdExceeded: ds?.[0]?.thresholdExceeded ?? false,
            penalty: ds?.[0]?.penalty ?? 0,
          },
          ...ds,
        ]);
      }, 200);
      try {
        await props.deploy();
      } catch (e) {
        console.error(e);
        return "再展開に失敗しました";
      } finally {
        clearTimeout(timer);
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
      allowedDeploymentCount={allowedDeploymentCount}
      error={lastResult}
      list={
        deployments.length === 0 ? (
          <View.EmptyDeploymentList
            allowedDeploymentCount={allowedDeploymentCount}
          />
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
