import "@testing-library/jest-dom";

import React from "react";

import { render, screen } from "@testing-library/react";
import { useRecoilState } from "recoil";
import { Mock, vi } from "vitest";

import Problems from "@/app/problems/page";
import useNotice from "@/hooks/notice";
import useProblems from "@/hooks/problems";
import { testNotice } from "@/types/Notice";
import { testProblem } from "@/types/Problem";

vi.mock("recoil");

vi.mock("@/components/markdown-preview", () => ({
  __esModule: true,
  default: ({ content }: { content: string }) => (
    <div data-testid="markdown-preview" data-content={content} />
  ),
}));
vi.mock("@/hooks/problems");
vi.mock("@/hooks/notice");
vi.mock("next/router", () => require("next-router-mock"));
vi.mock("next/navigation", () => ({
  ...require("next-router-mock"),
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
    expect(screen.queryByText("テスト通知タイトル")).toBeInTheDocument();
    expect(screen.getAllByTestId("markdown-preview")[1]).toHaveAttribute(
      "data-content",
      "テスト通知本文"
    );
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
    expect(screen.getAllByTestId("markdown-preview")[0]).toHaveAttribute(
      "data-content",
      "# ルール本文"
    );

    // verify
    expect(useProblems).toHaveBeenCalledTimes(1);
    expect(useNotice).toHaveBeenCalledTimes(1);
  });
});
