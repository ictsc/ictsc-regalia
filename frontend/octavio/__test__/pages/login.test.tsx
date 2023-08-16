import "@testing-library/jest-dom";

import React from "react";

import { act, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import mockRouter from "next-router-mock";
import { Mock, vi } from "vitest";

import Login from "@/app/login/page";
import useAuth from "@/hooks/auth";

vi.mock("@/hooks/auth");
vi.mock("@/components/Alerts", () => ({
  ICTSCSuccessAlert: ({
    message,
    subMessage,
  }: {
    message: string;
    subMessage: string | undefined;
  }) => (
    <div
      data-testid="success-alert"
      data-message={message}
      data-sub-message={subMessage}
    />
  ),
  ICTSCErrorAlert: ({
    message,
    subMessage,
  }: {
    message: string;
    subMessage: string | undefined;
  }) => (
    <div
      data-testid="error-alert"
      data-message={message}
      data-sub-message={subMessage}
    />
  ),
}));

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

vi.mock("next/router", () => require("next-router-mock"));
vi.mock("next/navigation", () => ({
  ...require("next-router-mock"),
}));
describe("Login", () => {
  test("画面が表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    render(<Login />);

    const loginButton = screen.getByRole("button");

    // then
    expect(screen.queryByPlaceholderText("ユーザー名")).toBeInTheDocument();
    expect(screen.queryByPlaceholderText("パスワード")).toBeInTheDocument();
    expect(loginButton).toBeInTheDocument();
    expect(loginButton).not.toHaveAttribute("loading");

    // verify
    expect(useAuth).toHaveBeenCalledTimes(1);
  });

  test("未入力で送信した時にエラーメッセージが表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    render(<Login />);

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    expect(
      screen.queryByText("ユーザー名を入力してください")
    ).toBeInTheDocument();
    expect(
      screen.queryByText("パスワードを入力してください")
    ).toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(2);
  });

  test("ユーザー名が未入力で送信した時にエラーメッセージが表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    render(<Login />);

    await userEvent.type(screen.getByPlaceholderText("パスワード"), "password");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    expect(
      screen.queryByText("ユーザー名を入力してください")
    ).toBeInTheDocument();
    expect(
      screen.queryByText("パスワードを入力して下さい")
    ).not.toBeInTheDocument();
  });

  test("パスワードが未入力で送信した時にエラーメッセージが表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    render(<Login />);

    await userEvent.type(screen.getByPlaceholderText("ユーザー名"), "admin");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    expect(
      screen.queryByText("ユーザー名を入力してください")
    ).not.toBeInTheDocument();
    expect(
      screen.queryByText("パスワードを入力してください")
    ).toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(2);
  });

  test("ログインが成功することを確認する", async () => {
    // setup
    await mockRouter.push("/");
    const signIn = vi.fn().mockResolvedValue({ code: 200 });
    (useAuth as Mock).mockReturnValue({
      user: null,
      signIn,
      mutate: vi.fn(),
    });
    render(<Login />);

    await userEvent.type(screen.getByPlaceholderText("ユーザー名"), "admin");
    await userEvent.type(screen.getByPlaceholderText("パスワード"), "password");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    expect(mockRouter).toMatchObject({
      pathname: "/",
    });
    const alert = screen.getByTestId("success-alert");
    expect(alert).toBeInTheDocument();
    expect(alert).toHaveAttribute("data-message", "ログインに成功しました");
    expect(alert).not.toHaveAttribute("data-sub-message");

    // verify
    expect(useAuth).toHaveBeenCalledTimes(2);
    expect(signIn).toHaveBeenCalledTimes(1);
  });

  test("ログインが失敗することを確認する", async () => {
    // setup
    await mockRouter.push("/");
    const signIn = vi.fn().mockResolvedValue({ code: 400 });
    (useAuth as Mock).mockReturnValue({
      user: null,
      signIn,
      mutate: vi.fn(),
    });
    render(<Login />);

    await userEvent.type(screen.getByPlaceholderText("ユーザー名"), "admin");
    await userEvent.type(screen.getByPlaceholderText("パスワード"), "password");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    const alert = screen.getByTestId("error-alert");
    expect(alert).toBeInTheDocument();
    expect(alert).toHaveAttribute("data-message", "ログインに失敗しました");
    expect(alert).not.toHaveAttribute("data-sub-message");

    // verify
    expect(useAuth).toHaveBeenCalledTimes(2);
    expect(signIn).toHaveBeenCalledTimes(1);
  });

  test("フォームが送信中の時にボタンが無効になることを確認する", async () => {
    // setup
    await mockRouter.push("/");
    // 1秒後にレスポンスを返すことで Loading 中を再現する
    const signIn = vi.fn().mockResolvedValue(
      new Promise((resolve) => {
        setTimeout(() => {
          resolve({ code: 200 });
        }, 1000);
      })
    );
    (useAuth as Mock).mockReturnValue({
      user: null,
      signIn,
      mutate: vi.fn(),
    });
    render(<Login />);

    await userEvent.type(screen.getByPlaceholderText("ユーザー名"), "admin");
    await userEvent.type(screen.getByPlaceholderText("パスワード"), "password");
    const loginButton = screen.getByRole("button");

    // when
    await act(async () => {
      loginButton.click();
    });

    // then
    expect(loginButton).toHaveClass("loading");

    // verify
    expect(useAuth).toHaveBeenCalledTimes(2);
    expect(signIn).toHaveBeenCalledTimes(1);
  });
});
