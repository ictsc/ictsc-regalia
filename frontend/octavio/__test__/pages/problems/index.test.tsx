import "@testing-library/jest-dom";

import React from "react";

import { render, screen } from "@testing-library/react";
import { useRecoilState } from "recoil";
import { Mock, vi } from "vitest";

import useNotice from "@/hooks/notice";
import useProblems from "@/hooks/problems";
import Problems from "@/pages/problems";
import { testNotice } from "@/types/Notice";
import { testProblem } from "@/types/Problem";

vi.mock("recoil");
vi.mock("@/hooks/problems");
vi.mock("@/hooks/notice");
vi.mock("next/router", () => require("next-router-mock"));
vi.mock("@/layouts/CommonLayout", () => ({
  __esModule: true,
  default: ({
    children,
    title,
  }: {
    children: React.ReactNode;
    title: string;
  }) => (
    <div data-testid="common-layout" data-title={title}>
      {children}
    </div>
  ),
}));

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("Problems", () => {
  test("画面が表示されることを確認する", async () => {
    // setup
    (useRecoilState as Mock).mockReturnValue([[], vi.fn()]);
    (useProblems as Mock).mockReturnValue({
      problems: [testProblem],
      isLoading: false,
    });
    (useNotice as Mock).mockReturnValue({
      notices: [testNotice],
      isLoading: false,
    });

    render(<Problems />);

    // verify
    expect(screen.getByTestId("common-layout")).toBeInTheDocument();
    expect(screen.getByTestId("common-layout")).toHaveAttribute(
      "data-title",
      "問題一覧"
    );
    expect(screen.queryByText("テスト通知タイトル")).toBeInTheDocument();
    expect(screen.queryByText("テスト通知本文")).toBeInTheDocument();
    expect(screen.queryByText("XYZ")).toBeInTheDocument();
    expect(screen.queryByText("テスト問題タイトル")).toBeInTheDocument();
    expect(screen.queryByText("100/100pt")).toBeInTheDocument();
    expect(useProblems).toHaveBeenCalledTimes(1);
    expect(useNotice).toHaveBeenCalledTimes(1);
  });

  test("問題一覧とお知らせ一覧が取得中の場合、ローディング画面が表示されることを確認する", async () => {
    // setup
    (useRecoilState as Mock).mockReturnValue([[], vi.fn()]);
    (useProblems as Mock).mockReturnValue({
      problems: [],
      isLoading: true,
    });
    (useNotice as Mock).mockReturnValue({
      notices: [],
      isLoading: true,
    });
    render(<Problems />);

    // then
    expect(screen.getByTestId("common-layout")).toBeInTheDocument();
    expect(screen.queryByTestId("loading")).toBeInTheDocument();

    // verify
    expect(useProblems).toHaveBeenCalledTimes(1);
    expect(useNotice).toHaveBeenCalledTimes(1);
  });

  test("問題一覧が取得中の場合、ローディング画面が表示されることを確認する", async () => {
    // setup
    (useRecoilState as Mock).mockReturnValue([[], vi.fn()]);
    (useProblems as Mock).mockReturnValue({
      problems: [],
      isLoading: true,
    });
    (useNotice as Mock).mockReturnValue({
      notices: [],
      isLoading: false,
    });
    render(<Problems />);

    // then
    expect(screen.getByTestId("common-layout")).toBeInTheDocument();
    expect(screen.queryByTestId("loading")).toBeInTheDocument();

    // verify
    expect(useProblems).toHaveBeenCalledTimes(1);
    expect(useNotice).toHaveBeenCalledTimes(1);
  });

  test("お知らせ一覧が取得中の場合、ローディング画面が表示されることを確認する", async () => {
    // setup
    (useRecoilState as Mock).mockReturnValue([[], vi.fn()]);
    (useProblems as Mock).mockReturnValue({
      problems: [],
      isLoading: false,
    });
    (useNotice as Mock).mockReturnValue({
      notices: [],
      isLoading: true,
    });
    render(<Problems />);

    // then
    expect(screen.getByTestId("common-layout")).toBeInTheDocument();
    expect(screen.queryByTestId("loading")).toBeInTheDocument();

    // verify
    expect(useProblems).toHaveBeenCalledTimes(1);
    expect(useNotice).toHaveBeenCalledTimes(1);
  });

  test("見えなくされている場合、お知らせが表示されないことを確認する", async () => {
    // setup
    (useRecoilState as Mock).mockReturnValue([[testNotice.source_id], vi.fn()]);
    (useProblems as Mock).mockReturnValue({
      problems: [testProblem],
      isLoading: false,
    });
    (useNotice as Mock).mockReturnValue({
      notices: [testNotice],
      isLoading: false,
    });
    render(<Problems />);

    // then
    expect(screen.getByTestId("common-layout")).toBeInTheDocument();
    expect(screen.queryByText("テスト通知タイトル")).not.toBeInTheDocument();
    expect(screen.queryByText("テスト通知本文")).not.toBeInTheDocument();

    // verify
    expect(useProblems).toHaveBeenCalledTimes(1);
    expect(useNotice).toHaveBeenCalledTimes(1);
  });

  test("お知らせを見えなくするボタンが動作することを確認する", async () => {
    // setup
    const onDismiss = vi.fn();
    (useRecoilState as Mock).mockReturnValue([[], onDismiss]);
    (useProblems as Mock).mockReturnValue({
      problems: [testProblem],
      isLoading: false,
    });
    (useNotice as Mock).mockReturnValue({
      notices: [testNotice],
      isLoading: false,
    });
    render(<Problems />);

    // when
    screen.getByRole("button").click();

    // then
    expect(screen.getByTestId("common-layout")).toBeInTheDocument();
    expect(onDismiss).toHaveBeenCalledWith([testNotice.source_id]);

    // verify
    expect(onDismiss).toHaveBeenCalledTimes(1);
    expect(useProblems).toHaveBeenCalledTimes(1);
    expect(useNotice).toHaveBeenCalledTimes(1);
  });

  test("お知らせを見えなくする id がすでに存在している時、正しく配列にセットされる", async () => {
    // setup
    const onDismiss = vi.fn();
    (useRecoilState as Mock).mockReturnValue([["TEST"], onDismiss]);
    (useProblems as Mock).mockReturnValue({
      problems: [testProblem],
      isLoading: false,
    });
    (useNotice as Mock).mockReturnValue({
      notices: [testNotice],
      isLoading: false,
    });
    render(<Problems />);

    // when
    screen.getByRole("button").click();

    // then
    expect(screen.getByTestId("common-layout")).toBeInTheDocument();
    expect(onDismiss).toHaveBeenCalledWith(["TEST", testNotice.source_id]);

    // verify
    expect(onDismiss).toHaveBeenCalledTimes(1);
    expect(useProblems).toHaveBeenCalledTimes(1);
    expect(useNotice).toHaveBeenCalledTimes(1);
  });

  test("shortRule が正しく表示されることを確認する", async () => {
    // setup
    vi.mock("@/components/_const", () => ({
      title: "title",
      site: "site",
      shortRule: "# ルール本文",
    }));
    (useRecoilState as Mock).mockReturnValue([[], vi.fn()]);
    (useProblems as Mock).mockReturnValue({
      problems: [testProblem],
      isLoading: false,
    });
    (useNotice as Mock).mockReturnValue({
      notices: [testNotice],
      isLoading: false,
    });
    render(<Problems />);

    // then
    expect(screen.getByTestId("common-layout")).toBeInTheDocument();
    expect(screen.queryByText("ルール本文")).toBeInTheDocument();

    // verify
    expect(useProblems).toHaveBeenCalledTimes(1);
    expect(useNotice).toHaveBeenCalledTimes(1);
  });
});
