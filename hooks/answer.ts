import useSWR from "swr";

import { useApi } from "@/hooks/api";
import { Answer } from "@/types/Answer";
import { Result } from "@/types/_api";

type AnswerResult = {
  answers: Answer[];
};

export const useAnswers = (id: string | null) => {
  const { apiClient } = useApi();

  const fetcher = (url: string) =>
    apiClient.get(url).json<Result<AnswerResult>>();

  const { data, mutate } = useSWR(id && `problems/${id}/answers`, fetcher);

  const getAnswer = (id: string): Answer | null => {
    return data?.data?.answers.find((answer) => answer.id === id) ?? null;
  };

  return {
    answers: data?.data?.answers ?? [],
    getAnswer,
    mutate,
  };
};
