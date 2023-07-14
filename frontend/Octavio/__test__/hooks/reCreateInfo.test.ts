import { renderHook } from "@testing-library/react";
import useSWR from "swr";
import { Mock, vi } from "vitest";

import useReCreateInfo from "@/hooks/reCreateInfo";
import { GetReCreateInfo, testReCreateInfo } from "@/types/ReCreate";
import { Result } from "@/types/_api";

vi.mock("swr");

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("useReCreateInfo", () => {
  it("再展開情報が取得できる", () => {
    // setup
    const mockResult: Result<GetReCreateInfo> = {
      code: 200,
      data: testReCreateInfo,
    };

    (useSWR as Mock).mockReturnValue({
      data: mockResult,
      isLoading: false,
      mutate: vi.fn(),
    });

    // when
    const { result } = renderHook(() => useReCreateInfo("test"));

    // then
    expect(useSWR).toBeCalledWith("recreate/test", expect.any(Function), {
      refreshInterval: 30000,
    });
    expect(result.current.recreateInfo).toEqual(mockResult.data);
    expect(result.current.isLoading).toEqual(false);
    expect(result.current.mutate).toBeDefined();

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });

  it("取得した再展開情報が空の時 recreateInfo が Null になる", () => {
    // setup
    const mockResult: Result<GetReCreateInfo> = {
      code: 200,
      data: null,
    };

    (useSWR as Mock).mockReturnValue({
      data: mockResult,
      isLoading: false,
      mutate: vi.fn(),
    });

    // when
    const { result } = renderHook(() => useReCreateInfo("test"));

    // then
    expect(useSWR).toBeCalledWith("recreate/test", expect.any(Function), {
      refreshInterval: 30000,
    });
    expect(result.current.recreateInfo).toBeNull();
    expect(result.current.isLoading).toEqual(false);
    expect(result.current.mutate).toBeDefined();

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });
});
