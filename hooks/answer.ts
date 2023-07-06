import useSWR from "swr";

import useApi from "@/hooks/api";
import { Answer } from "@/types/Answer";
import { AnswerResult } from "@/types/_api";

const useAnswers = (problemId: string | null) => {
  const { client } = useApi();

  const fetcher = (url: string) => client.get<AnswerResult>(url);

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
