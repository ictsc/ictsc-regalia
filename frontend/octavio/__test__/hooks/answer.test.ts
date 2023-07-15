import { renderHook } from "@testing-library/react";
import useSWR from "swr";
import { Mock, vi } from "vitest";

import useAnswers from "@/hooks/answer";
import { testAnswer } from "@/types/Answer";
import { AnswerResult, Result } from "@/types/_api";

vi.mock("swr");

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("useAnswer", () => {
  it("回答一覧が取得できる", () => {
    // setup
    const mockResult: Result<AnswerResult> = {
      code: 200,
      data: {
        answers: [testAnswer],
      },
    };

    (useSWR as Mock).mockReturnValue({
      data: mockResult,
      mutate: vi.fn(),
      isLoading: false,
    });

    // when
    const { result } = renderHook(() => useAnswers("1"));

    // then
    expect(useSWR).toBeCalledWith("problems/1/answers", expect.any(Function));
    expect(result.current.answers).toEqual(mockResult.data?.answers);
    expect(result.current.getAnswer("1")).toEqual(testAnswer);
    expect(result.current.mutate).toBeDefined();

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });

  it("取得した回答一覧が空の時 answers, getAnswer が空になる", () => {
    // setup
    const mockResult: Result<AnswerResult> = {
      code: 200,
      data: null,
    };

    (useSWR as Mock).mockReturnValue({
      data: mockResult,
      mutate: vi.fn(),
      isLoading: false,
    });

    // when
    const { result } = renderHook(() => useAnswers("1"));

    // then
    expect(useSWR).toBeCalledWith("problems/1/answers", expect.any(Function));
    expect(result.current.answers).toEqual([]);
    expect(result.current.getAnswer("1")).toBeNull();
    expect(result.current.mutate).toBeDefined();

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });

  it("ProblemId で回答が見つからなかった時 getAnswer が空になる", () => {
    // setup
    const mockResult: Result<AnswerResult> = {
      code: 200,
      data: {
        answers: [testAnswer],
      },
    };

    (useSWR as Mock).mockReturnValue({
      data: mockResult,
      mutate: vi.fn(),
      isLoading: false,
    });

    const id = "unknownId";

    // when
    const { result } = renderHook(() => useAnswers(id));

    // then
    expect(useSWR).toBeCalledWith(
      `problems/${id}/answers`,
      expect.any(Function)
    );
    expect(result.current.answers).toEqual(mockResult.data?.answers);
    expect(result.current.getAnswer(id)).toBeNull();
    expect(result.current.mutate).toBeDefined();

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });
});
