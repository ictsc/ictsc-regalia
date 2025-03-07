import { createFileRoute } from "@tanstack/react-router";
import { createClient } from "@connectrpc/connect";
import { MarkService, ProblemService } from "@ictsc/proto/admin/v1";
import { use, useDeferredValue } from "react";
import {
  Button,
  Container,
  Grid,
  NumberInput,
  TableOfContents,
  Textarea,
  Title,
} from "@mantine/core";
import ReactMarkdown from "react-markdown";

export const Route = createFileRoute("/submissions/$problem/$team/$id")({
  component: RouteComponent,
  loader({
    context: { transport },
    params: { problem: problemCode, team: teamCode, id: answerID },
  }) {
    const problemClient = createClient(ProblemService, transport);
    const markClient = createClient(MarkService, transport);

    const problem = problemClient.getProblem({ code: problemCode });
    const answer = markClient.getAnswer({
      problemCode,
      teamCode: parseInt(teamCode, 10),
      id: parseInt(answerID, 10),
    });

    return {
      problem,
      answer,
    };
  },
});

function RouteComponent() {
  const { problem: problemPromise, answer: answerPromise } =
    Route.useLoaderData();
  const deferredProblemPromise = useDeferredValue(problemPromise);
  const deferredAnswerPromise = useDeferredValue(answerPromise);
  const { problem } = use(deferredProblemPromise);
  const { answer } = use(deferredAnswerPromise);

  return (
    <Container>
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
            <ReactMarkdown>
              {answer?.body?.body.value?.body ?? ""}
            </ReactMarkdown>
          </article>
          <article>
            <Title>問題解説</Title>
            <ReactMarkdown>
              {problem?.body?.body.value?.explanationMarkdown ?? ""}
            </ReactMarkdown>
          </article>
          <form>
            <Title>採点</Title>
            <NumberInput
              mt="md"
              name="score"
              label="得点"
              description={`0~${problem?.maxScore}`}
              min={0}
              max={problem?.maxScore}
            />
            <Textarea mt="md" name="comment" label="コメント" />
            <Button mt="md" type="submit">
              送信
            </Button>
          </form>
        </Grid.Col>
      </Grid>
    </Container>
  );
}
