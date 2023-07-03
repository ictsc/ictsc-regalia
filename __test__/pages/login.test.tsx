import "@testing-library/jest-dom";

import { act, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import mockRouter from "next-router-mock";
import { Mock, vi } from "vitest";

import useAuth from "@/hooks/auth";
import Login from "@/pages/login";

vi.mock("@/hooks/auth");

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

vi.mock("next/router", () => require("next-router-mock"));

describe("Login", () => {
  test("画面が表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    render(<Login />);

    const loginButton = screen.getByRole("button");

    // verify
    expect(screen.queryByPlaceholderText("ユーザー名")).toBeInTheDocument();
    expect(screen.queryByPlaceholderText("パスワード")).toBeInTheDocument();
    expect(loginButton).toBeInTheDocument();
    expect(loginButton).not.toHaveAttribute("loading");
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
      screen.queryByText("パスワードを入力して下さい")
    ).toBeInTheDocument();
  });

  test("ユーザー名が未入力で送信した時にエラーメッセージが表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    render(<Login />);

    await userEvent.type(screen.getByTestId("password"), "password");

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

    await userEvent.type(screen.getByTestId("username"), "admin");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    expect(
      screen.queryByText("ユーザー名を入力してください")
    ).not.toBeInTheDocument();
    expect(
      screen.queryByText("パスワードを入力して下さい")
    ).toBeInTheDocument();
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

    await userEvent.type(screen.getByTestId("username"), "admin");
    await userEvent.type(screen.getByTestId("password"), "password");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    expect(mockRouter).toMatchObject({
      pathname: "/",
    });
    expect(screen.queryByText("ログインに成功しました")).toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
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

    await userEvent.type(screen.getByTestId("username"), "admin");
    await userEvent.type(screen.getByTestId("password"), "password");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    expect(screen.queryByText("ログインに失敗しました")).toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
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

    await userEvent.type(screen.getByTestId("username"), "admin");
    await userEvent.type(screen.getByTestId("password"), "password");
    const loginButton = screen.getByRole("button");

    // when
    await act(async () => {
      loginButton.click();
    });

    // then
    expect(loginButton).toHaveClass("loading");

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
    expect(signIn).toHaveBeenCalledTimes(1);
  });
});
