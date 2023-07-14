import { renderHook } from "@testing-library/react";
import matter from "gray-matter";
import { Mock, vi } from "vitest";

import useProblem from "@/hooks/problem";
import useProblems from "@/hooks/problems";
import { Problem, testProblem } from "@/types/Problem";
import { Result } from "@/types/_api";

vi.mock("@/hooks/problems");
vi.mock("gray-matter");

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("useProblem", () => {
  test("問題が取得できる", () => {
    // setup
    const mockProblemsResult: Result<Problem[]> = {
      code: 200,
      data: [testProblem],
    };

    (useProblems as Mock).mockReturnValue({
      problems: mockProblemsResult.data,
    });
    (matter as unknown as Mock).mockReturnValue({
      data: { title: "テスト", body: "テスト本文" },
    });

    // when
    renderHook(() => useProblem("XYZ"));

    // then
    expect(matter).toBeCalledWith(testProblem.body);

    // verify
    expect(useProblems).toBeCalledTimes(1);
    expect(matter).toBeCalledTimes(1);
  });

  test("問題が見つからない時", () => {
    // setup
    const mockProblemsResult: Result<Problem[]> = {
      code: 200,
      data: [testProblem],
    };

    (useProblems as Mock).mockReturnValue({
      problems: mockProblemsResult.data,
    });
    (matter as unknown as Mock).mockReturnValue({
      data: { title: "テスト" },
    });

    // when
    const { result } = renderHook(() => useProblem("unknownId"));

    // then
    expect(result.current.matter).toBeNull();
    expect(result.current.problem).toBeNull();
    expect(result.current.isLoading).toBeFalsy();

    // verify
    expect(useProblems).toBeCalledTimes(1);
    expect(matter).toBeCalledTimes(0);
  });

  test("問題の body が Null の時に matter が Null で返されるか", () => {
    // setup
    const bodyEmptyProblem = { ...testProblem, body: undefined };
    const mockProblemsResult: Result<Problem[]> = {
      code: 200,
      data: [bodyEmptyProblem],
    };
    (useProblems as Mock).mockReturnValue({
      problems: mockProblemsResult.data,
    });
    (matter as unknown as Mock).mockReturnValue({
      data: null,
    });

    // when
    const { result } = renderHook(() => useProblem("XYZ"));

    // then
    expect(matter).toBeCalledWith("");
    expect(result.current.matter).toBeNull();
    expect(result.current.problem).toEqual({ ...bodyEmptyProblem, body: "" });
    expect(result.current.isLoading).toBeFalsy();

    // verify
    expect(useProblems).toBeCalledTimes(1);
    expect(matter).toBeCalledTimes(1);
  });
});
