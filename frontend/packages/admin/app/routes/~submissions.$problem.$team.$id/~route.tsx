import { createFileRoute, useRouter } from "@tanstack/react-router";
import { createClient } from "@connectrpc/connect";
import { timestampMs } from "@bufbuild/protobuf/wkt";
import {
  type MarkingResult,
  MarkService,
  ProblemService,
} from "@ictsc/proto/admin/v1";
import {
  use,
  useActionState,
  useDeferredValue,
  useId,
  useMemo,
  useOptimistic,
} from "react";
import {
  Button,
  Flex,
  Grid,
  Modal,
  NumberInput,
  Table,
  TableOfContents,
  Textarea,
  Title,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import ReactMarkdown from "react-markdown";

export const Route = createFileRoute("/submissions/$problem/$team/$id")({
  component: RouteComponent,
  loader({ context: { transport }, params }) {
    const problemClient = createClient(ProblemService, transport);
    const markClient = createClient(MarkService, transport);

    const problemCode = params.problem;
    const teamCode = parseInt(params.team, 10);
    const answerID = parseInt(params.id, 10);

    const problem = problemClient
      .getProblem({ code: problemCode })
      .then((r) => r.problem);
    const answer = markClient
      .getAnswer({
        problemCode,
        teamCode,
        id: answerID,
      })
      .then((r) => r.answer);
    const markingResults = markClient
      .listMarkingResults({})
      .then((r) =>
        r.markingResults.filter(
          (m) =>
            m.answer?.id === answerID &&
            m.answer?.problem?.code === problemCode &&
            m.answer?.team?.code === BigInt(teamCode),
        ),
      );

    const submitMark = async (score: number, comment: string) => {
      await markClient.createMarkingResult({
        markingResult: {
          answer: await answer,
          score,
          rationale: {
            body: {
              case: "descriptive",
              value: {
                comment,
              },
            },
          },
        },
      });
    };

    return {
      problem,
      answer,
      markingResults,
      submitMark,
    };
  },
});

function RouteComponent() {
  const {
    problem: problemPromise,
    answer: answerPromise,
    markingResults: markingResultsPromise,
    submitMark,
  } = Route.useLoaderData();
  const deferredProblemPromise = useDeferredValue(problemPromise);
  const deferredAnswerPromise = useDeferredValue(answerPromise);
  const deferredMarkingResults = useDeferredValue(markingResultsPromise);
  const problem = use(deferredProblemPromise);
  const answer = use(deferredAnswerPromise);
  const markingResults = use(deferredMarkingResults);

  const router = useRouter();

  const [markingResultsItems, addOptimisticMarkingResult] =
    useMarkingResults(markingResults);

  return (
    <Grid>
      <Grid.Col span={2}>
        <TableOfContents
          style={{ position: "sticky", top: 70 }}
          scrollSpyOptions={{
            selector: "h1",
          }}
          getControlProps={({ data }) => ({
            onClick: () =>
              data.getNode().scrollIntoView({ behavior: "smooth" }),
            children: data.value,
          })}
        />
      </Grid.Col>
      <Grid.Col span={10}>
        <article>
          <Title>解答</Title>
          <ReactMarkdown>{answer?.body?.body.value?.body ?? ""}</ReactMarkdown>
        </article>
        <article>
          <Title>問題解説</Title>
          <ReactMarkdown>
            {problem?.body?.body.value?.explanationMarkdown ?? ""}
          </ReactMarkdown>
        </article>
        <MarkForm
          maxScore={answer?.problem?.maxScore ?? 0}
          action={async ({ score, comment }) => {
            addOptimisticMarkingResult({ score, comment });
            await submitMark(score, comment);
            await router.invalidate();
          }}
        />
        <Flex mt="md" direction="column">
          <Title>採点履歴</Title>
          <MarkingResultsTable
            maxScore={answer?.problem?.maxScore ?? 0}
            items={markingResultsItems}
          />
        </Flex>
      </Grid.Col>
    </Grid>
  );
}

function MarkForm(props: {
  maxScore: number;
  action: (value: { score: number; comment: string }) => Promise<void>;
}) {
  const [dialogOpened, dialogControl] = useDisclosure(false);

  const [lastResult, action, isPending] = useActionState(
    async (_prev: unknown, formData: FormData) => {
      const score = Number(formData.get("score"));
      const comment = formData.get("comment") as string;

      if (isNaN(score) || score < 0 || score > props.maxScore) {
        notifications.show({
          color: "red",
          title: "得点が不正です",
          message: `0~${props.maxScore}の範囲で入力してください`,
        });
        return;
      }

      switch (formData.get("intent")) {
        case "confirm":
          dialogControl.open();
          return { confirm: true, score, comment };
        case "submit":
          await props.action({ score, comment }).catch((e) => {
            console.error(e);
            notifications.show({
              color: "red",
              title: "採点結果の送信に失敗しました",
              message: "",
            });
          });
          return;
      }
    },
    null,
  );
  const isSending = isPending && lastResult?.confirm;

  const formID = useId();

  return (
    <form id={formID} action={action}>
      <Title>採点</Title>
      <NumberInput
        mt="md"
        name="score"
        label="得点"
        description={`0~${props.maxScore}`}
        required
        disabled={isSending}
        min={0}
        max={props.maxScore}
        defaultValue={lastResult?.score}
      />
      <Textarea
        mt="md"
        name="comment"
        label="コメント"
        disabled={isSending}
        defaultValue={lastResult?.comment}
      />
      <Button
        mt="md"
        name="intent"
        value="confirm"
        type="submit"
        disabled={isSending}
      >
        送信
      </Button>

      <Modal title="確認" opened={dialogOpened} onClose={dialogControl.close}>
        <div>
          <p>得点: {lastResult?.score}</p>
          <p>コメント: {lastResult?.comment}</p>
        </div>
        <Flex gap="md">
          <Button onClick={dialogControl.close}>キャンセル</Button>
          <Button
            variant="filled"
            color="blue"
            form={formID}
            type="submit"
            name="intent"
            value="submit"
            onClick={dialogControl.close}
          >
            送信
          </Button>
        </Flex>
      </Modal>
    </form>
  );
}

type MarkingResultItem = {
  number: number;
  score: number;
  comment: string;
  createdAtMs: number;
  pending?: boolean;
};

function useMarkingResults(
  markingResults: MarkingResult[],
): [
  items: readonly MarkingResultItem[],
  optimisticAdd: (newItem: { score: number; comment: string }) => void,
] {
  const items = useMemo(() => {
    const items = markingResults.flatMap((mark) => {
      if (mark.createdAt == null) {
        return [];
      }
      return [
        {
          number: -1,
          score: mark.score,
          comment: mark.rationale?.body?.value?.comment ?? "",
          createdAtMs: timestampMs(mark.createdAt),
        },
      ];
    });
    items.sort((a, b) => b.createdAtMs - a.createdAtMs);
    return items.map((item, index) => ({
      ...item,
      number: items.length - index,
    }));
  }, [markingResults]);
  const [optimisticItems, optimisticAdd] = useOptimistic(
    items,
    (items, newItem: { score: number; comment: string }) => {
      return [
        {
          number: items.length + 1,
          createdAtMs: Date.now(),
          pending: true,
          ...newItem,
        },
        ...items,
      ];
    },
  );

  return [optimisticItems, optimisticAdd];
}

const submitTimeFormatter = new Intl.DateTimeFormat("ja-JP", {
  dateStyle: "short",
  timeStyle: "medium",
});

function MarkingResultsTable(props: {
  maxScore: number;
  items: readonly MarkingResultItem[];
}) {
  return (
    <Table>
      <Table.Thead>
        <Table.Tr>
          <Table.Th>#</Table.Th>
          <Table.Th>提出時刻</Table.Th>
          <Table.Th>得点</Table.Th>
          <Table.Th>コメント</Table.Th>
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>
        {props.items.map((item) => (
          <Table.Tr key={item.createdAtMs} opacity={item.pending ? 0.5 : 1}>
            <Table.Td>{item.number}</Table.Td>
            <Table.Td>
              {submitTimeFormatter.format(new Date(item.createdAtMs))}
            </Table.Td>
            <Table.Td>
              {item.score}/{props.maxScore}
            </Table.Td>
            <Table.Td>{item.comment}</Table.Td>
          </Table.Tr>
        ))}
      </Table.Tbody>
    </Table>
  );
}
