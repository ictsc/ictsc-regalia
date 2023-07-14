import { renderHook } from "@testing-library/react";
import useSWR from "swr";
import { Mock, vi } from "vitest";

import useRanking from "@/hooks/ranking";
import { testRank } from "@/types/Rank";
import { RankingResult, Result } from "@/types/_api";

vi.mock("swr");

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("useRanking", () => {
  it("ランキング一覧が取得できる", () => {
    // setup
    const mockResult: Result<RankingResult> = {
      code: 200,
      data: {
        ranking: [testRank],
      },
    };

    (useSWR as Mock).mockReturnValue({
      data: mockResult,
      isLoading: false,
    });

    // when
    const { result } = renderHook(() => useRanking());

    // then
    expect(useSWR).toBeCalledWith("ranking", expect.any(Function));
    expect(result.current.ranking).toEqual(mockResult.data?.ranking);
    expect(result.current.loading).toEqual(false);

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });

  it("取得したランキング一覧が空の時 ranking が空になる", () => {
    // setup
    const mockResult: Result<RankingResult> = {
      code: 200,
      data: null,
    };

    (useSWR as Mock).mockReturnValue({
      data: mockResult,
      isLoading: false,
    });

    // when
    const { result } = renderHook(() => useRanking());

    // then
    expect(useSWR).toBeCalledWith("ranking", expect.any(Function));
    expect(result.current.ranking).toBeNull();
    expect(result.current.loading).toEqual(false);

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });
});
