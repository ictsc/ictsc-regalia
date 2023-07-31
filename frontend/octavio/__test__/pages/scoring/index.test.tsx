import "@testing-library/jest-dom";

import React from "react";

import { act, render, screen } from "@testing-library/react";
import { Mock, vi } from "vitest";

import Index from "@/app/(operational)/scoring/page";
import useAuth from "@/hooks/auth";
import useProblems from "@/hooks/problems";
import { testProblem } from "@/types/Problem";
import { testAdminUser, testUser } from "@/types/User";

vi.mock("next/error", () => ({
  __esModule: true,
  default: ({ statusCode }: { statusCode: number }) => (
    <div data-testid="error" data-status-code={statusCode} />
  ),
}));
vi.mock("next/router", () => require("next-router-mock"));
vi.mock("@/hooks/auth");
vi.mock("@/hooks/problems");
vi.mock("@/components/LoadingPage", () => ({
  __esModule: true,
  default: () => <div data-testid="loading" />,
}));

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("Scoring", () => {
  test("未ログインで、エラーページが表示される", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    (useProblems as Mock).mockReturnValue({
      problems: [],
      isLoading: false,
    });
    render(<Index />);

    // when
    expect(screen.getByTestId("error")).toBeInTheDocument();
    expect(screen.getByTestId("error")).toHaveAttribute(
      "data-status-code",
      "404"
    );

    // then
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblems).toHaveBeenCalledTimes(2);
  });

  test("ログイン済みで問題が取得できない場合、エラーページが表示される", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblems as Mock).mockReturnValue({
      problems: [],
      isLoading: false,
    });
    render(<Index />);

    // when
    expect(screen.getByTestId("error")).toBeInTheDocument();
    expect(screen.getByTestId("error")).toHaveAttribute(
      "data-status-code",
      "404"
    );

    // then
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblems).toHaveBeenCalledTimes(2);
  });

  test("参加者でアクセスした場合、エラーページが表示される", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testUser,
    });
    (useProblems as Mock).mockReturnValue({
      problems: [],
      isLoading: false,
    });
    render(<Index />);

    // when
    expect(screen.getByTestId("error")).toBeInTheDocument();
    expect(screen.getByTestId("error")).toHaveAttribute(
      "data-status-code",
      "404"
    );

    // then
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblems).toHaveBeenCalledTimes(2);
  });

  test("問題を取得中の場合、ローディング画面が表示される", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblems as Mock).mockReturnValue({
      problems: [],
      isLoading: true,
    });
    render(<Index />);

    // when
    expect(screen.queryByTestId("loading")).toBeInTheDocument();

    // then
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblems).toHaveBeenCalledTimes(2);
  });

  test("問題一覧が表示される", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblems as Mock).mockReturnValue({
      problems: [testProblem],
      isLoading: false,
    });
    render(<Index />);
    const tds = screen.queryAllByRole("cell");

    // when
    expect(screen.queryByText("採点")).toBeInTheDocument();
    expect(tds[1]).toHaveTextContent("-/-/-");
    expect(tds[2]).toHaveTextContent("id");
    expect(tds[3]).toHaveTextContent("XYZ");
    expect(tds[4]).toHaveTextContent("テスト問題タイトル");
    expect(tds[5]).toHaveTextContent("--- title: テスト --- #...");
    expect(tds[6]).toHaveTextContent("100");
    expect(tds[7]).toHaveTextContent("150");
    expect(tds[8]).toHaveTextContent("");
    expect(tds[9]).toHaveTextContent("自分");

    // then
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblems).toHaveBeenCalledTimes(2);
  });

  test("15分未満の問題がある場合、未採点の ~15分 に表示される", () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblems as Mock).mockReturnValue({
      problems: [{ ...testProblem, unchecked: 1 }],
      isLoading: false,
    });
    render(<Index />);
    const tds = screen.queryAllByRole("cell");

    // when
    expect(tds[1]).toHaveTextContent("1/-/-");

    // then
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblems).toHaveBeenCalledTimes(2);
  });

  test("15分以上かつ19分以下の問題がある場合、未採点の 15~19分 に表示される", () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblems as Mock).mockReturnValue({
      problems: [{ ...testProblem, unchecked_near_overdue: 1 }],
      isLoading: false,
    });
    render(<Index />);
    const tds = screen.queryAllByRole("cell");

    // when
    expect(tds[1]).toHaveTextContent("/1/-");

    // then
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblems).toHaveBeenCalledTimes(2);
  });

  test("20分以上の問題がある場合、未採点の 20分~ に表示される", () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblems as Mock).mockReturnValue({
      problems: [{ ...testProblem, unchecked_overdue: 1 }],
      isLoading: false,
    });
    render(<Index />);
    const tds = screen.queryAllByRole("cell");

    // when
    expect(tds[1]).toHaveTextContent("/-/1");

    // then
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblems).toHaveBeenCalledTimes(2);
  });

  test("20文字未満の問題文は、20文字目以降が省略されない", () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblems as Mock).mockReturnValue({
      problems: [{ ...testProblem, body: "a".repeat(19) }],
      isLoading: false,
    });
    render(<Index />);
    const tds = screen.queryAllByRole("cell");

    // when
    expect(tds[5]).toHaveTextContent("a".repeat(19));
  });

  test("20文字以上の問題文は、20文字目以降が省略される", () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblems as Mock).mockReturnValue({
      problems: [{ ...testProblem, body: "a".repeat(20) }],
      isLoading: false,
    });
    render(<Index />);
    const tds = screen.queryAllByRole("cell");

    // when
    // aaa... となる
    expect(tds[5]).toHaveTextContent(`${"a".repeat(20)}...`);

    // then
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblems).toHaveBeenCalledTimes(2);
  });

  test("問題文が Null の場合、空文字が表示される", () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblems as Mock).mockReturnValue({
      problems: [{ ...testProblem, body: null }],
      isLoading: false,
    });
    render(<Index />);
    const tds = screen.queryAllByRole("cell");

    // when
    // aaa... となる
    expect(tds[5]).toHaveTextContent("");

    // then
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblems).toHaveBeenCalledTimes(2);
  });

  test("問題作成者id が自分でない場合空文字が表示される", () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblems as Mock).mockReturnValue({
      problems: [{ ...testProblem, author_id: "other" }],
      isLoading: false,
    });
    render(<Index />);
    const tds = screen.queryAllByRole("cell");

    // when
    expect(tds[9]).toHaveTextContent("");

    // then
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblems).toHaveBeenCalledTimes(2);
  });

  test("問題をクリックした場合、問題文が表示される", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblems as Mock).mockReturnValue({
      problems: [{ ...testProblem }],
      isLoading: false,
    });
    render(<Index />);
    await act(async () => {
      await screen.queryAllByRole("row")[1].click();
    });

    // when
    expect(screen.queryByText("採点する")).toBeInTheDocument();

    // then
    expect(useAuth).toHaveBeenCalledTimes(2);
    expect(useProblems).toHaveBeenCalledTimes(4);
  });
});
