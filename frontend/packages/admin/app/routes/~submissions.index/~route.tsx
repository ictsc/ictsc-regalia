import { createFileRoute, useRouter } from "@tanstack/react-router";
import { createClient } from "@connectrpc/connect";
import { timestampMs } from "@bufbuild/protobuf/wkt";
import {
  MarkService,
  ProblemService,
  TeamService,
  type Answer,
  type Problem,
  type Team,
  type ListTeamsResponse,
  type ListProblemsResponse,
} from "@ictsc/proto/admin/v1";
import { Center, Table, Text, Button, MultiSelect } from "@mantine/core";
import {
  useState,
  useEffect,
  useMemo,
  useDeferredValue,
  startTransition,
  use,
  ReactNode,
} from "react";

interface LoaderData {
  answers: Promise<{ answers: Answer[] }>;
  teams: { label: string; value: string }[];
  problems: { label: string; value: string }[];
}

export const Route = createFileRoute("/submissions/")({
  component: RouteComponent,
  loader: async ({ context: { transport } }): Promise<LoaderData> => {
    const markClient = createClient(MarkService, transport);
    const teamClient = createClient(TeamService, transport);
    const problemClient = createClient(ProblemService, transport);

    let teamsData: { label: string; value: string }[] | null = null;
    try {
      const cachedTeams = localStorage.getItem("teamsCache");
      if (cachedTeams) {
        teamsData = JSON.parse(cachedTeams) as {
          label: string;
          value: string;
        }[];
      }
    } catch (e) {
      console.error(e);
      teamsData = null;
    }
    let problemData: { label: string; value: string }[] | null = null;
    try {
      const cachedProblems = localStorage.getItem("problemsCache");
      if (cachedProblems) {
        problemData = JSON.parse(cachedProblems) as {
          label: string;
          value: string;
        }[];
      }
    } catch (e) {
      console.error(e);
      problemData = null;
    }

    if (!teamsData) {
      const teamsResp: ListTeamsResponse = await teamClient.listTeams({});
      teamsData = teamsResp.teams.map((team: Team) => ({
        label: team.name,
        value: team.name,
      }));
      localStorage.setItem("teamsCache", JSON.stringify(teamsData));
    }
    if (!problemData) {
      const problemsResp: ListProblemsResponse =
        await problemClient.listProblems({});
      problemData = problemsResp.problems.map((problem: Problem) => ({
        label: `${problem.code}: ${problem.title}`,
        value: String(problem.code),
      }));
      localStorage.setItem("problemsCache", JSON.stringify(problemData));
    }

    return {
      answers: markClient.listAnswers({}),
      teams: teamsData,
      problems: problemData,
    };
  },
});

const submitTimeFormatter = new Intl.DateTimeFormat("ja-JP", {
  dateStyle: "short",
  timeStyle: "medium",
});

function RouteComponent() {
  const {
    answers: answersPromise,
    teams: teamsData,
    problems: problemsData,
  } = Route.useLoaderData();
  const deferredAnswersPromise = useDeferredValue(answersPromise);
  const router = useRouter();

  const [selectedProblemCodes, setSelectedProblemCodes] = usePersistentState<
    string[]
  >("selectedProblemCodes", []);
  const [selectedTeamNames, setSelectedTeamNames] = usePersistentState<
    string[]
  >("selectedTeamNames", []);
  const [filterRecent, setFilterRecent] = usePersistentState<boolean>(
    "filterRecent",
    false,
  );
  const [showPerfect, setShowPerfect] = usePersistentState<boolean>(
    "showPerfect",
    true,
  );
  const [showUnscored, setShowUnscored] = usePersistentState<boolean>(
    "showUnscored",
    false,
  );

  useEffect(() => {
    const timer = setInterval(() => {
      startTransition(() => router.invalidate());
    }, 60 * 1000);
    return () => clearInterval(timer);
  }, [router]);

  const answersResp = use(deferredAnswersPromise);
  const items = useAnswers(answersResp.answers ?? []);

  const filteredItems = filterAnswers(
    items,
    selectedProblemCodes,
    selectedTeamNames,
    filterRecent,
    showPerfect,
    showUnscored,
  );

  const teamOptions = teamsData;
  const problemOptions = problemsData;

  return (
    <>
      <FilterBar
        selectedProblemCodes={selectedProblemCodes}
        onSelectedProblemCodesChange={setSelectedProblemCodes}
        selectedTeamNames={selectedTeamNames}
        onSelectedTeamNamesChange={setSelectedTeamNames}
        filterRecent={filterRecent}
        onToggleFilterRecent={() => setFilterRecent((prev) => !prev)}
        showPerfect={showPerfect}
        onToggleShowPerfect={() => setShowPerfect((prev) => !prev)}
        showUnscored={showUnscored}
        onToggleShowUnscored={() => setShowUnscored((prev) => !prev)}
        problemOptions={problemOptions}
        teamOptions={teamOptions}
      />
      <Center>
        <AnswerTable answers={filteredItems} />
      </Center>
    </>
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

function usePersistentState<T>(
  key: string,
  initialValue: T,
): [T, React.Dispatch<React.SetStateAction<T>>] {
  const [state, setState] = useState<T>(() => {
    const stored = localStorage.getItem(key);
    if (stored !== null) {
      try {
        return JSON.parse(stored) as T;
      } catch {
        return initialValue;
      }
    }
    return initialValue;
  });

  useEffect(() => {
    localStorage.setItem(key, JSON.stringify(state));
  }, [key, state]);

  return [state, setState];
}

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
    const scoredItems: AnswerItem[] = [],
      unscoredItems: AnswerItem[] = [];
    for (const item of rawItems) {
      if (item.score != null) {
        scoredItems.push(item);
      } else {
        unscoredItems.push(item);
      }
    }
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

function filterAnswers(
  items: AnswerItem[],
  selectedProblemCodes: string[],
  selectedTeamNames: string[],
  filterRecent: boolean,
  showPerfect: boolean,
  showUnscored: boolean,
): AnswerItem[] {
  let filtered = items;
  if (selectedProblemCodes.length > 0) {
    filtered = filtered.filter((item) =>
      selectedProblemCodes.includes(item.problemCode),
    );
  }
  if (selectedTeamNames.length > 0) {
    filtered = filtered.filter((item) =>
      selectedTeamNames.includes(item.teamName),
    );
  }
  if (filterRecent) {
    const twentyMinutesAgo = Date.now() - 20 * 60 * 1000;
    filtered = filtered.filter((item) => item.submitTimeMs >= twentyMinutesAgo);
  }
  if (!showPerfect) {
    filtered = filtered.filter(
      (item) => !(item.score && item.score.total === item.score.max),
    );
  }
  if (showUnscored) {
    filtered = filtered.filter((item) => item.score == null);
  }
  return filtered;
}

function FilterBar(props: {
  selectedProblemCodes: string[];
  onSelectedProblemCodesChange: (value: string[]) => void;
  selectedTeamNames: string[];
  onSelectedTeamNamesChange: (value: string[]) => void;
  filterRecent: boolean;
  onToggleFilterRecent: () => void;
  showPerfect: boolean;
  onToggleShowPerfect: () => void;
  showUnscored: boolean;
  onToggleShowUnscored: () => void;
  problemOptions: { label: string; value: string }[];
  teamOptions: { label: string; value: string }[];
}) {
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        gap: "16px",
        marginBottom: "24px",
      }}
    >
      <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
        <MultiSelect
          placeholder="問題検索"
          data={props.problemOptions}
          searchable
          value={props.selectedProblemCodes}
          onChange={props.onSelectedProblemCodesChange}
        />
        <MultiSelect
          placeholder="チーム検索"
          data={props.teamOptions}
          searchable
          value={props.selectedTeamNames}
          onChange={props.onSelectedTeamNamesChange}
        />
      </div>
      <div style={{ display: "flex", gap: "32px" }}>
        <Button onClick={props.onToggleFilterRecent}>
          {props.filterRecent ? "直近20分以外も表示" : "直近20分のみ表示"}
        </Button>
        <Button onClick={props.onToggleShowPerfect}>
          {props.showPerfect ? "満点解答を非表示" : "満点解答を表示"}
        </Button>
        <Button onClick={props.onToggleShowUnscored}>
          {props.showUnscored ? "全採点解答表示" : "未採点のみ表示"}
        </Button>
      </div>
    </div>
  );
}
