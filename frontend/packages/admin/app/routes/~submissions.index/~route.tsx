import { startTransition, use, useDeferredValue, useEffect } from "react";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { createClient } from "@connectrpc/connect";
import { timestampDate } from "@bufbuild/protobuf/wkt";
import { MarkService } from "@ictsc/proto/admin/v1";
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

  const { answers } = use(deferredAnswersPromise);
  return (
    <Center>
      <Table>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>問題</Table.Th>
            <Table.Th>チーム</Table.Th>
            <Table.Th>解答ID</Table.Th>
            <Table.Th>提出時刻</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {answers.map((answer) => (
            <Table.Tr
              key={`${answer.problem?.code}-${answer.team?.code}-${answer.id}`}
              onClick={() => {
                void router.navigate({
                  to: "/submissions/$problem/$team/$id",
                  params: {
                    problem: String(answer.problem?.code ?? ""),
                    team: String(answer.team?.code ?? ""),
                    id: String(answer.id ?? ""),
                  },
                });
              }}
              style={{ cursor: "pointer" }}
            >
              <Table.Td>
                <Text size="sm" maw="15em" lineClamp={1}>
                  {answer.problem?.code}: {answer.problem?.title}
                </Text>
              </Table.Td>
              <Table.Td>
                <Text
                  size="sm"
                  maw="10em"
                  lineClamp={1}
                  title={answer.team?.name}
                >
                  {answer.team?.name}
                </Text>
              </Table.Td>
              <Table.Td>{answer.id}</Table.Td>
              <Table.Td>
                {answer.createdAt != null
                  ? submitTimeFormatter.format(timestampDate(answer.createdAt))
                  : "-"}
              </Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>
    </Center>
  );
}
