import useSWR from 'swr'

import {useApi} from "./api";
import {Result} from "../types/_api";
import {Answer} from "../types/Answer";

type AnswerResult = {
  answers: Answer[]
}

export const useAnswers = (id: string) => {
  const {apiClient} = useApi()

  const fetcher = (url: string) => apiClient.get(url).json<Result<AnswerResult>>()

  const {data, mutate} = useSWR(`problems/${id}/answers`, fetcher)

  const getAnswer = (id: string): Answer | null => {
    return data?.data?.answers.find((answer) => answer.id === id) ?? null
  }

  return {
    answers: data?.data?.answers ?? [],
    getAnswer,
    mutate
  }
}