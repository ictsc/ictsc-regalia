import { renderHook } from "@testing-library/react";
import expect from "expect";
import useSWR from "swr";

import useAuth from "@/hooks/auth";
import { testUser } from "@/types/User";
import { AuthSelfResult, Result } from "@/types/_api";

jest.mock("swr");

describe("useAuth", () => {
  it("ユーザーが取得できる", () => {
    // setup
    const mockAuthResult: Result<AuthSelfResult> = {
      code: 200,
      data: {
        user: testUser,
      },
    };

    (useSWR as jest.Mock).mockReturnValue({
      data: mockAuthResult,
      mutate: jest.fn(),
      isLoading: false,
    });

    // when
    const { result } = renderHook(() => useAuth());

    // then
    expect(result.current.user).toEqual(mockAuthResult.data?.user);
    expect(result.current.isLoading).toEqual(false);
    expect(result.current.mutate).toBeDefined();

    // verify
    expect(useSWR).toBeCalledTimes(1);
  });

  it("ユーザーが取得できない", () => {
    // setup
    const mockAuthResult: Result<AuthSelfResult> = {
      code: 200,
      data: null,
    };

    (useSWR as jest.Mock).mockReturnValue({
      data: mockAuthResult,
      mutate: jest.fn(),
      isLoading: false,
    });

    // when
    const { result } = renderHook(() => useAuth());

    // then
    expect(result.current.user).toBeNull();
    expect(result.current.isLoading).toEqual(false);
    expect(result.current.mutate).toBeDefined();

    // verify
    expect(useSWR).toBeCalledTimes(2);
  });
});
