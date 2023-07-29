import "@testing-library/jest-dom";

import React from "react";

import { act, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Mock, vi } from "vitest";

import Profile from "@/app/profile/page";
import useAuth from "@/hooks/auth";
import { testUser } from "@/types/User";

vi.mock("next/error", () => ({
  __esModule: true,
  default: ({ statusCode }: { statusCode: number }) => (
    <div data-testid="error" data-status-code={statusCode} />
  ),
}));
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
vi.mock("@/components/LoadingPage", () => ({
  __esModule: true,
  default: () => <div data-testid="loading" />,
}));
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

describe("Profile", () => {
  test("画面が表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });

    // when
    render(<Profile />);

    // then
    expect(screen.getByTestId("common-layout")).toBeInTheDocument();
    expect(screen.getByTestId("common-layout")).toHaveAttribute(
      "data-title",
      "プロフィール"
    );
    expect(screen.getByText(testUser.user_group.name)).toBeInTheDocument();
    const inputs = screen.getAllByRole("textbox");
    expect(inputs).toHaveLength(5);
    expect(inputs[0]).toHaveValue(testUser.display_name);
    expect(inputs[1]).toHaveValue(testUser.profile?.self_introduction);
    expect(inputs[2]).toHaveValue(testUser.profile?.github_id);
    expect(inputs[3]).toHaveValue(testUser.profile?.twitter_id);
    expect(inputs[4]).toHaveValue(testUser.profile?.facebook_id);

    // verify
    expect(useAuth).toHaveBeenCalledTimes(2);
  });

  test("ユーザーが取得中の場合、ローディング画面が表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: true,
    });

    // when
    render(<Profile />);

    // then
    expect(screen.getByTestId("common-layout")).toBeInTheDocument();
    expect(screen.getByTestId("loading")).toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(2);
  });

  test("ユーザーが取得できなかった場合、エラーページが表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: null,
      isLoading: false,
    });

    // when
    render(<Profile />);

    // then
    expect(screen.queryByTestId("common-layout")).not.toBeInTheDocument();
    expect(screen.getByTestId("error")).toBeInTheDocument();
    expect(screen.getByTestId("error")).toHaveAttribute(
      "data-status-code",
      "404"
    );

    // verify
    expect(useAuth).toHaveBeenCalledTimes(1);
  });

  test("表示名が未セットの場合、名前が表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: {
        ...testUser,
        display_name: null,
      },
      isLoading: false,
    });

    // when
    render(<Profile />);

    // then
    const inputs = screen.getAllByRole("textbox");
    expect(inputs[0]).toHaveValue(testUser.name);

    // verify
    expect(useAuth).toHaveBeenCalledTimes(2);
  });

  test("プロフィールが未セットの場合、空文字が表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: {
        ...testUser,
        user_profile: null,
      },
      isLoading: false,
    });

    // when
    render(<Profile />);

    // then
    const inputs = screen.getAllByRole("textbox");
    expect(inputs[1]).toHaveValue("");
    expect(inputs[2]).toHaveValue("");
    expect(inputs[3]).toHaveValue("");
    expect(inputs[4]).toHaveValue("");

    // verify
    expect(useAuth).toHaveBeenCalledTimes(2);
  });

  test("自己紹介が未セットの場合、空文字が表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: {
        ...testUser,
        user_profile: {
          ...testUser.user_profile,
          self_introduction: null,
        },
      },
      isLoading: false,
    });

    // when
    render(<Profile />);

    // then
    const inputs = screen.getAllByRole("textbox");
    expect(inputs[1]).toHaveValue("");

    // verify
    expect(useAuth).toHaveBeenCalledTimes(2);
  });

  test("GitHub ID が未セットの場合、空文字が表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: {
        ...testUser,
        user_profile: {
          ...testUser.user_profile,
          github_id: null,
        },
      },
      isLoading: false,
    });

    // when
    render(<Profile />);

    // then
    const inputs = screen.getAllByRole("textbox");
    expect(inputs[2]).toHaveValue("");
  });

  test("Twitter ID が未セットの場合、空文字が表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: {
        ...testUser,
        user_profile: {
          ...testUser.user_profile,
          twitter_id: null,
        },
      },
      isLoading: false,
    });

    // when
    render(<Profile />);

    // then
    const inputs = screen.getAllByRole("textbox");
    expect(inputs[3]).toHaveValue("");
  });

  test("Facebook ID が未セットの場合、空文字が表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: {
        ...testUser,
        user_profile: {
          ...testUser.user_profile,
          facebook_id: null,
        },
      },
      isLoading: false,
    });

    // when
    render(<Profile />);

    // then
    const inputs = screen.getAllByRole("textbox");
    expect(inputs[4]).toHaveValue("");
  });

  test("表示名を未入力で送信した時にエラーメッセージが表示されることを確認する", async () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    render(<Profile />);
    await userEvent.clear(screen.getAllByRole("textbox")[0]);

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    expect(screen.getByText("表示名は必須です")).toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(3);
  });

  test("プロフィールが更新されメッセージが表示される", async () => {
    const mockPutProfile = vi.fn().mockResolvedValue({ code: 202 });
    const mockMutate = vi.fn();
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
      putProfile: mockPutProfile,
      mutate: mockMutate,
    });
    render(<Profile />);

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    const alert = screen.getByTestId("success-alert");
    expect(alert).toBeInTheDocument();
    expect(alert).toHaveAttribute("data-message", "プロフィールを更新しました");
    expect(alert).not.toHaveAttribute("data-sub-message");
    expect(mockPutProfile).toHaveBeenCalledWith(testUser.id, {
      display_name: testUser.display_name,
      self_introduction: testUser.user_profile!.self_introduction,
      github_id: testUser.user_profile!.github_id,
      twitter_id: testUser.user_profile!.twitter_id,
      facebook_id: testUser.user_profile!.facebook_id,
    });

    // verify
    expect(useAuth).toHaveBeenCalledTimes(3);
    expect(mockPutProfile).toHaveBeenCalledTimes(1);
    expect(mockMutate).toHaveBeenCalledTimes(1);
  });

  test("プロフィール更新が失敗し、エラーメッセージが表示される", async () => {
    const mockPutProfile = vi.fn().mockResolvedValue({ code: 400 });
    const mockMutate = vi.fn();
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
      putProfile: mockPutProfile,
      mutate: mockMutate,
    });
    render(<Profile />);

    // when
    await act(async () => {
      screen.getByRole("button").click();
    });

    // then
    const alert = screen.getByTestId("error-alert");
    expect(alert).toBeInTheDocument();
    expect(alert).toHaveAttribute(
      "data-message",
      "プロフィールの更新に失敗しました"
    );
    expect(alert).not.toHaveAttribute("data-sub-message");

    // verify
    expect(useAuth).toHaveBeenCalledTimes(3);
    expect(mockPutProfile).toHaveBeenCalledTimes(1);
    expect(mockMutate).toHaveBeenCalledTimes(0);
  });
});
