import { renderHook } from "@testing-library/react";
import useSWR from "swr";
import { Mock, vi } from "vitest";

import useNotice from "@/hooks/notice";
import { Notice, testNotice } from "@/types/Notice";
import { Result } from "@/types/_api";

vi.mock("swr");

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("useNotice", () => {
  it("通知一覧が取得できる", () => {
    // setup
    const mockResult: Result<Notice[]> = {
      code: 200,
      data: [testNotice],
    };

    (useSWR as Mock).mockReturnValue({
      data: mockResult,
      mutate: vi.fn(),
      isLoading: false,
    });

    // when
    const { result } = renderHook(() => useNotice());

    // then
    expect(useSWR).toBeCalledWith("notices", expect.any(Function));
    expect(result.current.notices).toEqual(mockResult.data);
    expect(result.current.mutate).toBeDefined();

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });

  it("取得した通知一覧が空の時 notices が空になる", () => {
    // setup
    const mockResult: Result<Notice[]> = {
      code: 200,
      data: null,
    };

    (useSWR as Mock).mockReturnValue({
      data: mockResult,
      mutate: vi.fn(),
      isLoading: false,
    });

    // when
    const { result } = renderHook(() => useNotice());

    // then
    expect(useSWR).toBeCalledWith("notices", expect.any(Function));
    expect(result.current.notices).toBeNull();
    expect(result.current.mutate).toBeDefined();

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });
});
