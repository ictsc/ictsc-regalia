import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import { Mock, MockInstance, vi } from "vitest";

import ICTSCNavbar from "@/components/Navbar";
import useAuth from "@/hooks/auth";
import { testAdminUser, testUser } from "@/types/User";

vi.mock("next/router", () => ({
  useRouter() {
    return {
      route: "/",
      push: () => {},
    };
  },
}));
vi.mock("@/hooks/auth");

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("未ログイン状態 ICTSCNavBar", () => {
  test("正常に表示され未ログイン時の項目が表示されることを確認する", () => {
    // setup
    (useAuth as unknown as MockInstance).mockReturnValue({
      user: null,
      mutate: () => {},
    });

    // when
    render(<ICTSCNavbar />);

    // then
    expect(screen.queryByText("ルール")).toBeInTheDocument();
    expect(screen.queryByText("チーム情報")).not.toBeInTheDocument();
    expect(screen.queryByText("問題")).not.toBeInTheDocument();
    expect(screen.queryByText("順位")).toBeInTheDocument();
    expect(screen.queryByText("参加者")).not.toBeInTheDocument();
    expect(screen.queryByText("採点")).not.toBeInTheDocument();
    expect(screen.queryByText("ログイン")).toBeInTheDocument();
  });

  // verify
  expect(useAuth).toHaveBeenCalledTimes(0);
});

describe("参加者ログイン状態 ICTSCNavBar", () => {
  test("正常に表示され参加者ログイン時の項目が表示されることを確認する", () => {
    // setup
    (useAuth as jest.Mock).mockReturnValue({
      user: testUser,
      mutate: () => {},
    });

    // when
    render(<ICTSCNavbar />);

    // then
    expect(screen.queryByText("ルール")).toBeInTheDocument();
    expect(screen.queryByText("チーム情報")).toBeInTheDocument();
    expect(screen.queryByText("問題")).toBeInTheDocument();
    expect(screen.queryByText("順位")).toBeInTheDocument();
    expect(screen.queryByText("参加者")).toBeInTheDocument();
    expect(screen.queryByText("採点")).not.toBeInTheDocument();
    expect(screen.queryByText("ログイン")).not.toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(1);
  });

  test("ログアウトボタンを押した時にログアウト処理が実行されることを確認する", () => {
    // setup
    const logout = vi.fn().mockResolvedValue({ status: 200 });
    (useAuth as unknown as MockInstance).mockReturnValue({
      user: testUser,
      mutate: () => {},
      logout,
    });
    render(<ICTSCNavbar />);

    // when
    screen.getByText("ログアウト").click();

    // then
    expect(logout).toHaveBeenCalledTimes(1);

    // verify
    expect(useAuth).toHaveBeenCalledTimes(1);
  });
});

describe("管理者ログイン状態 ICTSCNavBar", () => {
  test("正常に表示され管理者ログイン時の項目が表示されることを確認する", () => {
    // setup
    (useAuth as Mock).mockReturnValue({
      user: testAdminUser,
      mutate: () => {},
    });

    // when
    render(<ICTSCNavbar />);

    // then
    expect(screen.queryByText("ルール")).toBeInTheDocument();
    expect(screen.queryByText("チーム情報")).toBeInTheDocument();
    expect(screen.queryByText("問題")).toBeInTheDocument();
    expect(screen.queryByText("順位")).toBeInTheDocument();
    expect(screen.queryByText("参加者")).toBeInTheDocument();
    expect(screen.queryByText("採点")).toBeInTheDocument();
    expect(screen.queryByText("ログイン")).not.toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(1);
  });
});
