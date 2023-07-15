import { renderHook } from "@testing-library/react";
import useSWR from "swr";
import { Mock, vi } from "vitest";

import useUserGroups from "@/hooks/userGroups";
import { testUserGroup, UserGroup } from "@/types/UserGroup";
import { Result } from "@/types/_api";

vi.mock("swr");

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("useUserGroups", () => {
  it("ユーザーグループ一覧を取得できる", () => {
    // setup
    const mockResult: Result<UserGroup[]> = {
      code: 200,
      data: [testUserGroup],
    };

    (useSWR as Mock).mockReturnValue({
      data: mockResult,
      isLoading: false,
    });

    // when
    const { result } = renderHook(() => useUserGroups());

    // then
    expect(useSWR).toBeCalledWith("usergroups", expect.any(Function));
    expect(result.current.userGroups).toEqual(mockResult.data);
    expect(result.current.isLoading).toEqual(false);

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });

  it("取得したユーザーグループ一覧が空の時 userGroups が Null になる", () => {
    // setup
    const mockResult: Result<UserGroup[]> = {
      code: 200,
      data: null,
    };

    (useSWR as Mock).mockReturnValue({
      data: mockResult,
      isLoading: false,
    });

    // when
    const { result } = renderHook(() => useUserGroups());

    // then
    expect(useSWR).toBeCalledWith("usergroups", expect.any(Function));
    expect(result.current.userGroups).toBeNull();
    expect(result.current.isLoading).toEqual(false);

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });
});
