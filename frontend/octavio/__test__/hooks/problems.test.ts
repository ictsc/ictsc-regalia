import { renderHook } from "@testing-library/react";
import useSWR from "swr";
import { Mock, vi } from "vitest";

import useProblems from "@/hooks/problems";
import { Problem } from "@/types/Problem";
import { Result } from "@/types/_api";

vi.mock("swr");

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("useProblems", () => {
  it("問題一覧が取得できる", () => {
    // setup
    const mockProblemResult: Result<Problem[]> = {
      code: 200,
      data: [],
    };

    (useSWR as Mock).mockReturnValue({
      data: mockProblemResult,
      isLoading: false,
    });

    // when
    const { result } = renderHook(() => useProblems());

    // then
    expect(result.current.problems).toEqual(mockProblemResult.data);

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });

  it("問題一覧が取得できない", () => {
    // setup
    const mockProblemResult: Result<Problem[]> = {
      code: 200,
      data: null,
    };

    (useSWR as Mock).mockReturnValue({
      data: mockProblemResult,
    });

    // when
    const { result } = renderHook(() => useProblems());

    // then
    expect(result.current.problems).toEqual([]);

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });
});
