import useSWR from "swr";

import useApi from "@/hooks/api";
import { Answer } from "@/types/Answer";
import { Result } from "@/types/_api";

type AnswerResult = {
  answers: Answer[];
};
const useAnswers = (problemId: string | null) => {
  const { apiClient } = useApi();

  const fetcher = (url: string) =>
    apiClient.get(url).json<Result<AnswerResult>>();

  const { data, mutate } = useSWR(
    problemId && `problems/${problemId}/answers`,
    fetcher
  );

  const getAnswer = (id: string): Answer | null =>
    data?.data?.answers.find((answer: Answer) => answer.id === id) ?? null;

  return {
    answers: data?.data?.answers ?? [],
    getAnswer,
    mutate,
  };
};

export default useAnswers;
