import "@testing-library/jest-dom";

import React from "react";

import { act, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import mockRouter from "next-router-mock";
import { Mock, vi } from "vitest";

import useAuth from "@/hooks/auth";
import SignUp from "@/pages/signUp";

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

describe("SignUp", () => {
  test("画面が表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    render(<SignUp />);

    // then
    expect(screen.queryByPlaceholderText("ユーザー名")).toBeInTheDocument();
    expect(screen.queryByPlaceholderText("パスワード")).toBeInTheDocument();
    expect(screen.queryByRole("button")).toBeInTheDocument();
    expect(screen.queryByRole("button")).not.toHaveAttribute("loading");

    // verify
    expect(useAuth).toHaveBeenCalledTimes(2);
  });

  test("未入力で送信した時にエラーメッセージが表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    render(<SignUp />);

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

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
  });

  test("ユーザー名が未入力で送信した時にエラーメッセージが表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    render(<SignUp />);

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

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
  });

  test("パスワード名が未入力で送信した時にエラーメッセージが表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    render(<SignUp />);

    await userEvent.type(screen.getByPlaceholderText("ユーザー名"), "user");

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

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
  });

  test("パスワードが8文字未満で入力した場合にエラーが表示される", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
    });
    render(<SignUp />);

    await userEvent.type(screen.getByPlaceholderText("ユーザー名"), "user");
    await userEvent.type(screen.getByPlaceholderText("パスワード"), "aaaaaaa");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    expect(
      screen.queryByText("パスワードは8文字以上である必要があります")
    ).toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
  });

  test("登録が成功した場合にホーム画面に遷移することを確認する", async () => {
    // setup
    const signUp = vi.fn().mockResolvedValue({ code: 201 });
    (useAuth as Mock).mockReturnValue({
      user: null,
      signUp,
    });
    render(<SignUp />);

    await userEvent.type(screen.getByPlaceholderText("ユーザー名"), "user");
    await userEvent.type(screen.getByPlaceholderText("パスワード"), "password");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    expect(mockRouter).toMatchObject({
      pathname: "/login",
    });
    const alert = screen.getByTestId("success-alert");
    expect(alert).toBeInTheDocument();
    expect(alert).toHaveAttribute(
      "data-message",
      "ユーザー登録に成功しました！"
    );
    expect(alert).not.toHaveAttribute("data-sub-message");

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
    expect(signUp).toHaveBeenCalledTimes(1);
  });

  test("ユーザーが既に存在する場合にエラーが表示されることを確認する", async () => {
    // setup
    const signUp = vi.fn().mockResolvedValue({
      code: 400,
      data: "Error 1062: Duplicate entry 'user' for key 'name'",
    });

    (useAuth as Mock).mockReturnValue({
      user: null,
      signUp,
    });

    render(<SignUp />);

    await userEvent.type(screen.getByPlaceholderText("ユーザー名"), "user");
    await userEvent.type(screen.getByPlaceholderText("パスワード"), "password");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    const alert = screen.getByTestId("error-alert");
    expect(alert).toBeInTheDocument();
    expect(alert).toHaveAttribute("data-message", "エラーが発生しました");
    expect(alert).toHaveAttribute(
      "data-sub-message",
      "ユーザー名が重複しています。"
    );

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
  });

  test("UserGroupID が無効な場合にエラーが表示されることを確認する", async () => {
    // setup
    const signUp = vi.fn().mockResolvedValue({
      code: 400,
      data: "Error:Field validation for 'UserGroupID' failed on the 'required' tag",
    });

    (useAuth as Mock).mockReturnValue({
      user: null,
      signUp,
    });

    render(<SignUp />);

    await userEvent.type(screen.getByPlaceholderText("ユーザー名"), "user");
    await userEvent.type(screen.getByPlaceholderText("パスワード"), "password");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    const alert = screen.getByTestId("error-alert");
    expect(alert).toBeInTheDocument();
    expect(alert).toHaveAttribute("data-message", "エラーが発生しました");
    expect(alert).toHaveAttribute(
      "data-sub-message",
      "無効なユーザーグループです。"
    );

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
  });

  test("UserGroupID の uuid 形式が無効な場合にエラーが表示されることを確認する", async () => {
    // setup
    const signUp = vi.fn().mockResolvedValue({
      code: 400,
      data: "Error:Field validation for 'UserGroupID' failed on the 'uuid' tag",
    });

    (useAuth as Mock).mockReturnValue({
      user: null,
      signUp,
    });

    render(<SignUp />);

    await userEvent.type(screen.getByPlaceholderText("ユーザー名"), "user");
    await userEvent.type(screen.getByPlaceholderText("パスワード"), "password");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    const alert = screen.getByTestId("error-alert");
    expect(alert).toBeInTheDocument();
    expect(alert).toHaveAttribute("data-message", "エラーが発生しました");
    expect(alert).toHaveAttribute(
      "data-sub-message",
      "無効なユーザーグループです。"
    );

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
  });

  test("InvitationCode が無効の場合", async () => {
    // setup
    const signUp = vi.fn().mockResolvedValue({
      code: 400,
      data: "Error:Field validation for 'InvitationCode' failed on the 'required' tag",
    });

    (useAuth as Mock).mockReturnValue({
      user: null,
      signUp,
    });

    render(<SignUp />);

    await userEvent.type(screen.getByPlaceholderText("ユーザー名"), "user");
    await userEvent.type(screen.getByPlaceholderText("パスワード"), "password");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    const alert = screen.getByTestId("error-alert");
    expect(alert).toBeInTheDocument();
    expect(alert).toHaveAttribute("data-message", "エラーが発生しました");
    expect(alert).toHaveAttribute("data-sub-message", "無効な招待コードです。");

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
  });

  test("不明のエラーの場合にエラーが表示されることを確認する", async () => {
    // setup
    const signUp = vi.fn().mockResolvedValue({
      code: 400,
      data: null,
    });

    (useAuth as Mock).mockReturnValue({
      user: null,
      signUp,
    });

    render(<SignUp />);

    await userEvent.type(screen.getByPlaceholderText("ユーザー名"), "user");
    await userEvent.type(screen.getByPlaceholderText("パスワード"), "password");

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    const alert = screen.getByTestId("error-alert");
    expect(alert).toBeInTheDocument();
    expect(alert).toHaveAttribute("data-message", "エラーが発生しました");
    expect(alert).toHaveAttribute("data-sub-message", "");

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
  });

  test("フォームが送信中の場合にボタンが無効になることを確認する", async () => {
    // setup
    await mockRouter.push("/login");
    const signUp = vi.fn().mockResolvedValue(
      new Promise((resolve) => {
        setTimeout(() => {
          resolve({ code: 200 });
        }, 1000);
      })
    );
    (useAuth as Mock).mockReturnValue({
      user: null,
      signUp,
    });
    render(<SignUp />);

    await userEvent.type(screen.getByPlaceholderText("ユーザー名"), "user");
    await userEvent.type(screen.getByPlaceholderText("パスワード"), "password");
    const button = screen.getByRole("button");

    // when
    await act(async () => {
      button.click();
    });

    // then
    expect(button).toHaveClass("loading");

    // verify
    expect(useAuth).toHaveBeenCalledTimes(4);
    expect(signUp).toHaveBeenCalledTimes(1);
  });
});
