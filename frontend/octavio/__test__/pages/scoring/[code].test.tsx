import "@testing-library/jest-dom";

import React from "react";

import { render, screen } from "@testing-library/react";
import mockRouter from "next-router-mock";
import { useForm } from "react-hook-form";
import { Mock, vi } from "vitest";

import useAnswers from "@/hooks/answer";
import useAuth from "@/hooks/auth";
import useProblem from "@/hooks/problem";
import ScoringProblem from "@/pages/scoring/[code]";
import { Answer, testAnswer } from "@/types/Answer";
import { testProblem } from "@/types/Problem";
import { testAdminUser, testUser } from "@/types/User";

vi.mock("next/error", () => ({
  __esModule: true,
  default: ({ statusCode }: { statusCode: number }) => (
    <div data-testid="error" data-status-code={statusCode} />
  ),
}));
vi.mock("next/router", () => require("next-router-mock"));
vi.mock("react-hook-form", () => ({
  useForm: vi.fn(),
}));
vi.mock("@/hooks/auth");
vi.mock("@/hooks/problem");
vi.mock("@/hooks/answer");
vi.mock("@/components/LoadingPage", () => ({
  __esModule: true,
  default: () => <div data-testid="loading" />,
}));
vi.mock("@/components/ScoringAnswerForm", () => ({
  __esModule: true,
  default: ({ answer }: { answer: Answer }) => (
    <div data-testid="scoring-answer-form" data-key={answer.id} />
  ),
}));
vi.mock("@/layouts/BaseLayout", () => ({
  __esModule: true,
  default: ({
    children,
    title,
  }: {
    children: React.ReactNode;
    title: string;
  }) => (
    <div data-testid="base-layout" data-title={title}>
      {children}
    </div>
  ),
}));

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();

  (useForm as Mock).mockReturnValue({
    register: vi.fn(),
    watch: vi.fn(),
  });
});

describe("ScoringProblem", () => {
  test("未ログインで、エラーページが表示される", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    (useProblem as Mock).mockReturnValue({
      problem: null,
      isLoading: false,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [],
    });
    render(<ScoringProblem />);

    // when
    expect(screen.getByTestId("error")).toBeInTheDocument();
    expect(screen.getByTestId("error")).toHaveAttribute(
      "data-status-code",
      "404"
    );

    // then
    // baselayout が表示される
    expect(screen.queryByTestId("base-layout")).not.toBeInTheDocument();
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblem).toHaveBeenCalledTimes(1);
    expect(useAnswers).toHaveBeenCalledTimes(1);
  });

  test("ログイン済みで問題が取得できない場合、NotFound が表示される", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: null,
      isLoading: false,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [],
    });
    render(<ScoringProblem />);

    // when
    expect(screen.getByTestId("error")).toBeInTheDocument();
    expect(screen.getByTestId("error")).toHaveAttribute(
      "data-status-code",
      "404"
    );

    // then
    expect(screen.queryByTestId("base-layout")).not.toBeInTheDocument();
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblem).toHaveBeenCalledTimes(1);
    expect(useAnswers).toHaveBeenCalledTimes(1);
  });

  test("参加者でアクセスした場合、エラーページが表示される", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
      isLoading: false,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [],
    });
    render(<ScoringProblem />);

    // when
    expect(screen.getByTestId("error")).toBeInTheDocument();
    expect(screen.getByTestId("error")).toHaveAttribute(
      "data-status-code",
      "404"
    );

    // then
    expect(screen.queryByTestId("base-layout")).not.toBeInTheDocument();
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblem).toHaveBeenCalledTimes(1);
    expect(useAnswers).toHaveBeenCalledTimes(1);
  });

  test("isReadOnly 権限でアクセスした場合エラーページが表示される", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: { ...testAdminUser, is_read_only: true },
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
      isLoading: false,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [],
    });
    render(<ScoringProblem />);

    // when
    expect(screen.getByTestId("error")).toBeInTheDocument();
    expect(screen.getByTestId("error")).toHaveAttribute(
      "data-status-code",
      "404"
    );

    // then
    expect(screen.queryByTestId("base-layout")).not.toBeInTheDocument();
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblem).toHaveBeenCalledTimes(1);
    expect(useAnswers).toHaveBeenCalledTimes(1);
  });

  test("問題が読み込み中の時ローディング画面が表示される", async () => {
    // setup
    (useForm as Mock).mockReturnValue({
      register: vi.fn(),
      watch: vi.fn().mockReturnValue(0),
    });
    (useAuth as Mock).mockReturnValue({
      user: { ...testAdminUser, is_read_only: true },
    });
    (useProblem as Mock).mockReturnValue({
      problem: null,
      isLoading: true,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [],
    });
    render(<ScoringProblem />);

    // when
    expect(screen.queryByTestId("loading")).toBeInTheDocument();

    // then
    expect(screen.queryByTestId("base-layout")).toBeInTheDocument();
    expect(screen.queryByTestId("base-layout")).toHaveAttribute(
      "data-title",
      "採点"
    );
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useProblem).toHaveBeenCalledTimes(1);
    expect(useAnswers).toHaveBeenCalledTimes(1);
  });

  test("未採点一覧表が表示される", () => {
    // setup
    (useForm as Mock).mockReturnValue({
      register: vi.fn(),
      watch: vi.fn().mockReturnValue("0"),
    });
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [testAnswer],
    });

    // when
    render(<ScoringProblem />);

    // then
    const cells = screen.queryAllByRole("cell");
    expect(cells[0]).toHaveTextContent("");
    expect(cells[1]).toHaveTextContent("-");
    expect(cells[2]).toHaveTextContent("-");
  });

  test("15分未満の問題がある場合、未採点の ~15分 に表示される", () => {
    // setup
    (useForm as Mock).mockReturnValue({
      register: vi.fn(),
      watch: vi.fn().mockReturnValue("0"),
    });
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: { ...testProblem, unchecked: 1 },
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [testAnswer],
    });

    // when
    render(<ScoringProblem />);

    // then
    const cells = screen.queryAllByRole("cell");
    expect(cells[0]).toHaveTextContent("1");
    expect(cells[1]).toHaveTextContent("-");
    expect(cells[2]).toHaveTextContent("-");
  });

  test("15分以上かつ19分以下の問題がある場合、未採点の 15~19分 に表示される", () => {
    // setup
    (useForm as Mock).mockReturnValue({
      register: vi.fn(),
      watch: vi.fn().mockReturnValue("0"),
    });
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: { ...testProblem, unchecked_near_overdue: 1 },
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [testAnswer],
    });

    // when
    render(<ScoringProblem />);

    // then
    const cells = screen.queryAllByRole("cell");
    expect(cells[0]).toHaveTextContent("");
    expect(cells[1]).toHaveTextContent("1");
    expect(cells[2]).toHaveTextContent("-");
  });

  test("20分以上の問題がある場合、未採点の 20分~ に表示される", () => {
    // setup
    (useForm as Mock).mockReturnValue({
      register: vi.fn(),
      watch: vi.fn().mockReturnValue("0"),
    });
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: { ...testProblem, unchecked_overdue: 1 },
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [testAnswer],
    });

    // when
    render(<ScoringProblem />);

    // then
    const cells = screen.queryAllByRole("cell");
    expect(cells[0]).toHaveTextContent("");
    expect(cells[1]).toHaveTextContent("-");
    expect(cells[2]).toHaveTextContent("1");
  });

  test("「すべて」を選択している場合かつ未採点の場合採点フォームが表示される", async () => {
    // setup
    (useForm as Mock).mockReturnValue({
      register: vi.fn(),
      watch: vi.fn().mockReturnValue("0"),
    });
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [testAnswer],
    });

    // when
    render(<ScoringProblem />);

    // then
    expect(screen.queryByTestId("scoring-answer-form")).toBeInTheDocument();
  });

  test("「すべて」を選択している場合かつ採点済みの場合採点フォームが表示される", async () => {
    // setup
    (useForm as Mock).mockReturnValue({
      register: vi.fn(),
      watch: vi.fn().mockReturnValue("0"),
    });
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [{ ...testAnswer, point: 100 }],
    });

    // when
    render(<ScoringProblem />);

    // then
    expect(screen.queryByTestId("scoring-answer-form")).toBeInTheDocument();
  });

  test("「採点済みのみ」を選択している場合かつ未採点の場合採点フォームが表示される", async () => {
    // setup
    (useForm as Mock).mockReturnValue({
      register: vi.fn(),
      watch: vi.fn().mockReturnValue("1"),
    });
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [testAnswer],
    });

    // when
    render(<ScoringProblem />);

    // then
    expect(screen.queryByTestId("scoring-answer-form")).not.toBeInTheDocument();
  });

  test("「採点済みのみ」を選択している場合かつ採点済みの場合採点フォームが表示される", async () => {
    // setup
    (useForm as Mock).mockReturnValue({
      register: vi.fn(),
      watch: vi.fn().mockReturnValue("1"),
    });
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [{ ...testAnswer, point: 100 }],
    });

    // when
    render(<ScoringProblem />);

    // then
    expect(screen.queryByTestId("scoring-answer-form")).toBeInTheDocument();
  });

  test("「未採点のみ」を選択している場合かつ未採点の場合採点フォームが表示される", async () => {
    // setup
    (useForm as Mock).mockReturnValue({
      register: vi.fn(),
      watch: vi.fn().mockReturnValue("2"),
    });
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [testAnswer],
    });

    // when
    render(<ScoringProblem />);

    // then
    expect(screen.queryByTestId("scoring-answer-form")).toBeInTheDocument();
  });

  test("「未採点のみ」を選択している場合かつ採点済みの場合採点フォームが表示される", async () => {
    // setup
    (useForm as Mock).mockReturnValue({
      register: vi.fn(),
      watch: vi.fn().mockReturnValue("2"),
    });
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [{ ...testAnswer, point: 100 }],
    });

    // when
    render(<ScoringProblem />);

    // then
    expect(screen.queryByTestId("scoring-answer-form")).not.toBeInTheDocument();
  });

  test("answerId を指定した場合指定の回答が表示される", async () => {
    // setup
    mockRouter.query = { code: "abc", answer_id: "1" };
    (useForm as Mock).mockReturnValue({
      register: vi.fn(),
      watch: vi.fn().mockReturnValue(null),
    });
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [
        { ...testAnswer, id: "1" },
        { ...testAnswer, id: "2" },
        { ...testAnswer, id: "3" },
        { ...testAnswer, id: "4" },
      ],
    });

    // when
    render(<ScoringProblem />);

    // then
    expect(screen.queryByTestId("scoring-answer-form")).toHaveAttribute(
      "data-key",
      "1"
    );
  });

  test("解答フォームが正しい順番で表示される", async () => {
    // setup
    mockRouter.query = { answer_id: undefined };
    (useForm as Mock).mockReturnValue({
      register: vi.fn(),
      watch: vi.fn().mockReturnValue(null),
    });
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
    });
    (useProblem as Mock).mockReturnValue({
      problem: testProblem,
    });
    (useAnswers as Mock).mockReturnValue({
      answers: [
        { ...testAnswer, id: "3", created_at: "2021-01-03" },
        { ...testAnswer, id: "1", created_at: "2021-01-01" },
        { ...testAnswer, id: "4", created_at: "2021-01-01" },
        { ...testAnswer, id: "2", created_at: "2021-01-02" },
      ],
    });

    // when
    render(<ScoringProblem />);

    // then
    const forms = screen.queryAllByTestId("scoring-answer-form");
    expect(forms[0]).toHaveAttribute("data-key", "3");
    expect(forms[1]).toHaveAttribute("data-key", "2");
    expect(forms[2]).toHaveAttribute("data-key", "1");
    expect(forms[3]).toHaveAttribute("data-key", "4");
  });
});
