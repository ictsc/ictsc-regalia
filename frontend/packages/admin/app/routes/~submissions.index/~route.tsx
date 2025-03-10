import {
  ReactNode,
  startTransition,
  use,
  useDeferredValue,
  useEffect,
  useMemo,
} from "react";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { createClient } from "@connectrpc/connect";
import { timestampMs } from "@bufbuild/protobuf/wkt";
import { MarkService, type Answer } from "@ictsc/proto/admin/v1";
import { Center, Table, Text } from "@mantine/core";

export const Route = createFileRoute("/submissions/")({
  component: RouteComponent,
  loader: ({ context: { transport } }) => {
    const client = createClient(MarkService, transport);
    return {
      answers: client.listAnswers({}),
    };
  },
});

const submitTimeFormatter = new Intl.DateTimeFormat("ja-JP", {
  dateStyle: "short",
  timeStyle: "medium",
});

function RouteComponent() {
  const { answers: answersPromise } = Route.useLoaderData();
  const deferredAnswersPromise = useDeferredValue(answersPromise);

  const router = useRouter();
  useEffect(() => {
    const timer = setInterval(() => {
      startTransition(() => router.invalidate());
    }, 60 * 1000);
    return () => clearInterval(timer);
  });

  const answersResp = use(deferredAnswersPromise);
  const items = useAnswers(answersResp.answers ?? []);
  return (
    <Center>
      <AnswerTable answers={items} />
    </Center>
  );
}

type AnswerItem = {
  readonly key: string;
  readonly problemCode: string;
  readonly problemTitle: string;
  readonly teamCode: string;
  readonly teamName: string;
  readonly answerNumber: number;
  readonly submitTimeMs: number;
  readonly score?: {
    readonly total: number;
    readonly marked: number;
    readonly penalty: number;
    readonly max: number;
  };
};

function useAnswers(answers: readonly Answer[]): AnswerItem[] {
  const rawItems = useMemo(() => {
    return answers.map((answer) => {
      return {
        key: `${answer.problem?.code}-${answer.team?.code}-${answer.id}`,
        problemCode: String(answer.problem?.code ?? ""),
        problemTitle: answer.problem?.title ?? "",
        teamCode: String(answer.team?.code ?? 0),
        teamName: answer.team?.name ?? "",
        answerNumber: answer.id ?? 0,
        submitTimeMs:
          answer.createdAt != null ? timestampMs(answer.createdAt) : 0,
        score:
          answer.score != null
            ? {
                total: answer.score?.total ?? 0,
                marked: answer.score?.marked ?? 0,
                penalty: answer.score?.penalty ?? 0,
                max: answer.score?.max ?? 0,
              }
            : undefined,
      };
    });
  }, [answers]);

  const items = useMemo(() => {
    const scoredItems = [],
      unscoredItems = [];
    for (const item of rawItems) {
      if (item.score != null) {
        scoredItems.push(item);
      } else {
        unscoredItems.push(item);
      }
    }
    // 点数が付いていないものは提出時刻が早い順，付いているものは提出時刻が遅い順(最近採点された可能性が高い)に並べる
    unscoredItems.sort((a, b) => a.submitTimeMs - b.submitTimeMs);
    scoredItems.sort((a, b) => b.submitTimeMs - a.submitTimeMs);
    return [...unscoredItems, ...scoredItems];
  }, [rawItems]);

  return items;
}

function AnswerTable(props: { readonly answers: readonly AnswerItem[] }) {
  return (
    <Table>
      <Table.Thead>
        <Table.Tr>
          <Table.Th>問題</Table.Th>
          <Table.Th>チーム</Table.Th>
          <Table.Th>解答ID</Table.Th>
          <Table.Th>提出時刻</Table.Th>
          <Table.Th>点数</Table.Th>
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>
        {props.answers.map((item) => (
          <AnswerTr key={item.key} item={item}>
            <Table.Td>
              <Text size="sm" maw="15em" lineClamp={1}>
                {item.problemCode}: {item.problemTitle}
              </Text>
            </Table.Td>
            <Table.Td>
              <Text size="sm" maw="10em" lineClamp={1} title={item.teamName}>
                {item.teamName}
              </Text>
            </Table.Td>
            <Table.Td>{item.answerNumber}</Table.Td>
            <Table.Td>
              {submitTimeFormatter.format(new Date(item.submitTimeMs))}
            </Table.Td>
            <Table.Td>
              {item.score != null ? (
                <>
                  {item.score.total}({item.score.marked}-{item.score.penalty})/
                  {item.score.max}
                </>
              ) : (
                "-"
              )}
            </Table.Td>
          </AnswerTr>
        ))}
      </Table.Tbody>
    </Table>
  );
}

function AnswerTr(props: {
  readonly item: AnswerItem;
  readonly children?: ReactNode;
}) {
  const router = useRouter();
  return (
    <Table.Tr
      onClick={() => {
        void router.navigate({
          to: "/submissions/$problem/$team/$id",
          params: {
            problem: props.item.problemCode,
            team: props.item.teamCode,
            id: String(props.item.answerNumber),
          },
        });
      }}
      style={{ cursor: "pointer" }}
    >
      {props.children}
    </Table.Tr>
  );
}
