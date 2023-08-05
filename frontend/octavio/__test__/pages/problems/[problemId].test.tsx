import "@testing-library/jest-dom";

import React, { useState } from "react";

import { act, render, screen } from "@testing-library/react";
import { Mock, vi } from "vitest";

import ProblemPage from "@/app/problems/[problemId]/page";
import useAuth from "@/hooks/auth";
import useProblem from "@/hooks/problem";
import useReCreateInfo from "@/hooks/reCreateInfo";
import { testProblem } from "@/types/Problem";
import { testReCreateInfo } from "@/types/ReCreate";
import { testUser } from "@/types/User";

vi.mock("react");
vi.mock("next/error", () => ({
  __esModule: true,
  default: ({ statusCode }: { statusCode: number }) => (
    <div data-testid="error" data-status-code={statusCode} />
  ),
}));
vi.mock("@/components/_const", () => ({
  title: "title",
  site: "site",
  answerLimit: 20,
  recreateRule: "テスト再展開ルール",
}));
vi.mock("next/router", () => ({
  useRouter() {
    return {
      route: "/problems/abc",
      pathname: "",
      query: "MockVitest",
      asPath: "",
    };
  },
}));

vi.mock("@/app/problems/[problemId]/_components/answer-form", () => ({
  __esModule: true,
  default: () => <div data-testid="answerForm" />,
}));
vi.mock("@/app/problems/[problemId]/_components/answer-list-section", () => ({
  __esModule: true,
  default: () => <div data-testid="answerListSection" />,
}));
vi.mock("@/components/markdown-preview", () => ({
  __esModule: true,
  default: ({ content }: { content: string }) => (
    <div data-testid="markdown-preview" data-content={content} />
  ),
}));
vi.mock("@/hooks/auth");
vi.mock("@/hooks/problem");
vi.mock("@/hooks/reCreateInfo");

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("[problemId]", () => {
  test("画面が表示されることを確認する", async () => {
    // setup
    (useState as Mock).mockReturnValue([false, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
      isLoading: false,
    });
    (useReCreateInfo as Mock).mockReturnValue({
      reCreateInfo: null,
    });

    // when
    render(<ProblemPage params={{ problemId: "abc" }} />);

    // then
    expect(
      screen.queryByRole("button", { name: "再展開を行う" })
    ).not.toBeInTheDocument();
    expect(screen.queryByTestId("answerForm")).not.toBeInTheDocument();
    expect(screen.queryByTestId("answerListSection")).toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblem).toHaveBeenCalledTimes(1);
    expect(useReCreateInfo).toHaveBeenCalledTimes(1);
  });

  test("問題が取得中の場合、ローディング画面が表示されることを確認する", async () => {
    // setup
    (useState as Mock).mockReturnValue([false, vi.fn()]);
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
      isLoading: true,
    });
    (useReCreateInfo as Mock).mockReturnValue({
      reCreateInfo: null,
    });

    // when
    render(<ProblemPage params={{ problemId: "abc" }} />);

    // then
    expect(screen.queryByTestId("loading")).toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblem).toHaveBeenCalledTimes(1);
    expect(useReCreateInfo).toHaveBeenCalledTimes(1);
  });

  test("問題が取得できなかった場合、エラーページが表示されることを確認する", async () => {
    // setup
    (useState as Mock).mockReturnValue([false, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    (useProblem as Mock).mockReturnValue({
      problem: null,
      isLoading: false,
    });
    (useReCreateInfo as Mock).mockReturnValue({
      reCreateInfo: null,
    });
    render(<ProblemPage params={{ problemId: "abc" }} />);

    // then
    expect(screen.getByTestId("error")).toBeInTheDocument();
    expect(screen.getByTestId("error")).toHaveAttribute(
      "data-status-code",
      "404"
    );

    // verify
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblem).toHaveBeenCalledTimes(1);
    expect(useReCreateInfo).toHaveBeenCalledTimes(1);
  });

  test("再展開ボタンが押された時、モーダルが表示されることを確認する", async () => {
    // setup
    (useState as Mock).mockReturnValue([false, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
      isLoading: false,
    });
    (useReCreateInfo as Mock).mockReturnValue({
      reCreateInfo: testReCreateInfo,
    });
    render(<ProblemPage params={{ problemId: "abc" }} />);

    // when
    await act(async () => {
      screen.getByRole("button", { name: "再展開を行う" }).click();
    });

    // then
    expect(
      screen.queryByText("問題の再展開を行います。よろしいですか？")
    ).toBeInTheDocument();
    expect(screen.getAllByTestId("markdown-preview")[0]).toHaveAttribute(
      "data-content",
      "テスト再展開ルール"
    );

    // verify
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblem).toHaveBeenCalledTimes(1);
    expect(useReCreateInfo).toHaveBeenCalledTimes(1);
  });

  test("モーダルの閉じるボタンを推した時、正常に閉じられるか確認", async () => {
    // setup
    const mockSetIsReCreateModalOpen = vi.fn();
    (useState as Mock).mockReturnValue([false, mockSetIsReCreateModalOpen]);
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
      isLoading: false,
    });
    (useReCreateInfo as Mock).mockReturnValue({
      reCreateInfo: null,
    });
    render(<ProblemPage params={{ problemId: "abc" }} />);

    // when
    await act(async () => {
      screen.getByRole("button", { name: "閉じる" }).click();
    });

    // then
    expect(mockSetIsReCreateModalOpen).toHaveBeenCalledWith(false);

    // verify
    expect(mockSetIsReCreateModalOpen).toHaveBeenCalledTimes(1);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblem).toHaveBeenCalledTimes(1);
    expect(useReCreateInfo).toHaveBeenCalledTimes(1);
  });

  test("モーダルの再展開ボタンを推した時、正常に再展開されるか確認", async () => {
    // setup
    const mockSetIsReCreateModalOpen = vi.fn();
    (useState as Mock).mockReturnValue([false, mockSetIsReCreateModalOpen]);
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
      isLoading: false,
    });
    const mockMutate = vi.fn();
    const mockReCreate = vi.fn().mockResolvedValue({ code: 200 });
    (useReCreateInfo as Mock).mockReturnValue({
      reCreateInfo: null,
      mutate: mockMutate,
      reCreate: mockReCreate,
    });
    render(<ProblemPage params={{ problemId: "abc" }} />);

    // when
    await act(async () => {
      screen.getByRole("button", { name: "問題の再展開を行う" }).click();
    });

    // then
    expect(mockSetIsReCreateModalOpen).toHaveBeenCalledWith(false);
    expect(mockReCreate).toHaveBeenCalledWith("XYZ");

    // verify
    expect(mockSetIsReCreateModalOpen).toHaveBeenCalledTimes(1);
    expect(mockReCreate).toHaveBeenCalledTimes(1);
    expect(mockMutate).toHaveBeenCalledTimes(1);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblem).toHaveBeenCalledTimes(1);
    expect(useReCreateInfo).toHaveBeenCalledTimes(1);
  });

  test("ログインしている場合、再展開ボタンと回答フォームが表示されることを確認する", async () => {
    // setup
    (useState as Mock).mockReturnValue([false, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
      isLoading: false,
    });
    (useReCreateInfo as Mock).mockReturnValue({
      reCreateInfo: null,
    });
    render(<ProblemPage params={{ problemId: "abc" }} />);

    // then
    expect(
      screen.queryByRole("button", { name: "再展開を行う" })
    ).toBeInTheDocument();
    expect(screen.queryByTestId("answerForm")).toBeInTheDocument();
    expect(screen.queryByTestId("answerListSection")).toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblem).toHaveBeenCalledTimes(1);
    expect(useReCreateInfo).toHaveBeenCalledTimes(1);
  });

  test("isReadOnlyがtrueの場合、再展開ボタンと回答フォームが表示されないことを確認する", async () => {
    // setup
    (useState as Mock).mockReturnValue([false, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: { ...testUser, is_read_only: true },
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
      isLoading: false,
    });
    (useReCreateInfo as Mock).mockReturnValue({
      reCreateInfo: null,
    });
    render(<ProblemPage params={{ problemId: "abc" }} />);

    // then
    expect(
      screen.queryByRole("button", { name: "再展開を行う" })
    ).not.toBeInTheDocument();
    expect(screen.queryByTestId("answerForm")).not.toBeInTheDocument();
    expect(screen.queryByTestId("answerListSection")).toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblem).toHaveBeenCalledTimes(1);
    expect(useReCreateInfo).toHaveBeenCalledTimes(1);
  });
});
