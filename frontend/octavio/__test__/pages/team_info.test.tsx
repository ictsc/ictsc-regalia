import "@testing-library/jest-dom";

import React, { useState } from "react";

import { act, render, screen } from "@testing-library/react";
import { Mock, vi } from "vitest";

import useAuth from "@/hooks/auth";
import useRanking from "@/hooks/ranking";
import TeamInfo from "@/pages/team_info";
import { testRank } from "@/types/Rank";
import { testUser } from "@/types/User";
import { testUserGroup } from "@/types/UserGroup";

vi.mock("react");
vi.mock("@/hooks/auth");
vi.mock("@/hooks/ranking");
vi.mock("@/components/Navbar", () => ({
  __esModule: true,
  default: () => <div data-testid="navbar" />,
}));

vi.mock("@/components/HiddenInput", () => ({
  __esModule: true,
  default: ({
    value,
    isHidden,
    onClick,
  }: {
    value: string;
    isHidden: boolean;
    onClick: (e: React.MouseEvent<HTMLInputElement>) => void;
  }) => (
    <input
      data-testid="hidden-input"
      data-value={value}
      data-is-hidden={isHidden}
      onClick={onClick}
    />
  ),
}));
vi.mock("@/components/LoadingPage", () => ({
  __esModule: true,
  default: () => <div data-testid="loading" />,
}));
vi.mock("@/components/icons/Eye", () => ({
  __esModule: true,
  default: () => <div data-testid="eye" />,
}));
vi.mock("@/components/icons/EyeSlash", () => ({
  __esModule: true,
  default: () => <div data-testid="eye-slash" />,
}));

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("TeamInfo", () => {
  test("画面が表示されることを確認する", async () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([true, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });

    // when
    render(<TeamInfo />);

    screen.debug();

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    expect(
      screen.getByText(`${testUserGroup.name}@${testUserGroup.organization}`)
    ).toBeInTheDocument();
    expect(screen.getByText("256")).toBeInTheDocument(); // ランキング
    expect(screen.getByText("/1 teams")).toBeInTheDocument();
    const hiddenInputs = screen.getAllByTestId("hidden-input");
    expect(hiddenInputs[0]).toHaveAttribute(
      "data-value",
      `ssh ${testUserGroup.bastion?.bastion_user}@${testUserGroup.bastion?.bastion_host} -p ${testUserGroup.bastion?.bastion_port}`
    );
    expect(hiddenInputs[1]).toHaveAttribute(
      "data-value",
      testUserGroup.bastion?.bastion_password
    );

    // verify
    expect(useState).toHaveBeenCalledTimes(2);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("ユーザー情報が取得中の場合、ローディング画面が表示されることを確認する", () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([true, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: true,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });

    // when
    render(<TeamInfo />);

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    expect(screen.getByTestId("loading")).toBeInTheDocument();

    // verify
    expect(useState).toHaveBeenCalledTimes(2);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("ランキングが取得中の場合、ローディング画面が表示されることを確認する", () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([true, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: true,
    });

    // when
    render(<TeamInfo />);

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    expect(screen.getByTestId("loading")).toBeInTheDocument();

    // verify
    expect(useState).toHaveBeenCalledTimes(2);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("ユーザー情報とランキングが取得中の場合、ローディング画面が表示されることを確認する", () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([true, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: true,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: true,
    });

    // when
    render(<TeamInfo />);

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    expect(screen.getByTestId("loading")).toBeInTheDocument();

    // verify
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("bastion が存在しない場合、各情報が空で表示される", () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([true, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: { ...testUser, user_group: { ...testUserGroup, bastion: null } },
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });

    // when
    render(<TeamInfo />);

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    const hiddenInputs = screen.getAllByTestId("hidden-input");
    expect(hiddenInputs[0]).toHaveAttribute("data-value", "ssh @ -p ");
    expect(hiddenInputs[1]).toHaveAttribute("data-value", "");

    // verify
    expect(useState).toHaveBeenCalledTimes(2);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("bastion_user が空の場合、空で表示される", () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([true, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: {
        ...testUser,
        user_group: {
          ...testUserGroup,
          bastion: { ...testUserGroup.bastion, bastion_user: null },
        },
      },
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });

    // when
    render(<TeamInfo />);

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    const hiddenInputs = screen.getAllByTestId("hidden-input");
    expect(hiddenInputs[0]).toHaveAttribute(
      "data-value",
      `ssh @${testUserGroup.bastion?.bastion_host} -p ${testUserGroup.bastion?.bastion_port}`
    );
    expect(hiddenInputs[1]).toHaveAttribute(
      "data-value",
      testUserGroup.bastion?.bastion_password
    );

    // verify
    expect(useState).toHaveBeenCalledTimes(2);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("bastion_host が空の場合、空で表示される", () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([true, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: {
        ...testUser,
        user_group: {
          ...testUserGroup,
          bastion: { ...testUserGroup.bastion, bastion_host: null },
        },
      },
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });

    // when
    render(<TeamInfo />);

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    const hiddenInputs = screen.getAllByTestId("hidden-input");
    expect(hiddenInputs[0]).toHaveAttribute(
      "data-value",
      `ssh ${testUserGroup.bastion?.bastion_user}@ -p ${testUserGroup.bastion?.bastion_port}`
    );
    expect(hiddenInputs[1]).toHaveAttribute(
      "data-value",
      testUserGroup.bastion?.bastion_password
    );

    // verify
    expect(useState).toHaveBeenCalledTimes(2);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("bastion_port が空の場合、空で表示される", () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([true, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: {
        ...testUser,
        user_group: {
          ...testUserGroup,
          bastion: { ...testUserGroup.bastion, bastion_port: null },
        },
      },
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });

    // when
    render(<TeamInfo />);

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    const hiddenInputs = screen.getAllByTestId("hidden-input");
    expect(hiddenInputs[0]).toHaveAttribute(
      "data-value",
      `ssh ${testUserGroup.bastion?.bastion_user}@${testUserGroup.bastion?.bastion_host} -p `
    );
    expect(hiddenInputs[1]).toHaveAttribute(
      "data-value",
      testUserGroup.bastion?.bastion_password
    );

    // verify
    expect(useState).toHaveBeenCalledTimes(2);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("bastion_password が空の場合、空で表示される", () => {
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([true, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: {
        ...testUser,
        user_group: {
          ...testUserGroup,
          bastion: { ...testUserGroup.bastion, bastion_password: null },
        },
      },
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });

    // when
    render(<TeamInfo />);

    screen.debug();

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    const hiddenInputs = screen.getAllByTestId("hidden-input");
    expect(hiddenInputs[0]).toHaveAttribute(
      "data-value",
      `ssh ${testUserGroup.bastion?.bastion_user}@${testUserGroup.bastion?.bastion_host} -p ${testUserGroup.bastion?.bastion_port}`
    );
    expect(hiddenInputs[1]).toHaveAttribute("data-value", "");

    // verify
    expect(useState).toHaveBeenCalledTimes(2);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("SSHが Hidden の場合、SSHコマンドが表示されない", () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([true, vi.fn()])
      .mockReturnValueOnce([true, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });

    // when
    render(<TeamInfo />);

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    const hiddenInputs = screen.getAllByTestId("hidden-input");
    expect(hiddenInputs[0]).toHaveAttribute("data-is-hidden", "true");
    expect(screen.queryAllByTestId("eye")).toHaveLength(2);
    expect(screen.queryAllByTestId("eye-slash")).toHaveLength(0);

    // verify
    expect(useState).toHaveBeenCalledTimes(2);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("SSHが Hidden でない場合、SSHコマンドが表示される", () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([true, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });

    // when
    render(<TeamInfo />);

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    const hiddenInputs = screen.getAllByTestId("hidden-input");
    expect(hiddenInputs[0]).toHaveAttribute("data-is-hidden", "false");
    expect(screen.queryAllByTestId("eye")).toHaveLength(1);
    expect(screen.queryAllByTestId("eye-slash")).toHaveLength(1);

    // verify
    expect(useState).toHaveBeenCalledTimes(2);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("passwordが Hidden の場合、passwordが表示されない", () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([true, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });

    // when
    render(<TeamInfo />);

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    const hiddenInputs = screen.getAllByTestId("hidden-input");
    expect(hiddenInputs[1]).toHaveAttribute("data-is-hidden", "true");
    expect(screen.queryAllByTestId("eye")).toHaveLength(1);
    expect(screen.queryAllByTestId("eye-slash")).toHaveLength(1);

    // verify
    expect(useState).toHaveBeenCalledTimes(2);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("passwordが Hidden でない場合、passwordが表示されない", () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([false, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });

    // when
    render(<TeamInfo />);

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    const hiddenInputs = screen.getAllByTestId("hidden-input");
    expect(hiddenInputs[1]).toHaveAttribute("data-is-hidden", "false");
    expect(screen.queryAllByTestId("eye")).toHaveLength(0);
    expect(screen.queryAllByTestId("eye-slash")).toHaveLength(2);

    // verify
    expect(useState).toHaveBeenCalledTimes(2);
    expect(useAuth).toHaveBeenCalledTimes(1);
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("ssh をクリックすると、SSH コマンドが選択される", async () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([false, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });
    const mockSelect = vi.spyOn(HTMLInputElement.prototype, "select");
    render(<TeamInfo />);

    // when
    await act(async () => {
      screen.getAllByTestId("hidden-input")[0].click();
    });

    // then
    expect(mockSelect).toHaveBeenCalledTimes(1);
  });

  test("password をクリックすると、password が選択される", async () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([false, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });
    const mockSelect = vi.spyOn(HTMLInputElement.prototype, "select");
    render(<TeamInfo />);

    // when
    await act(async () => {
      screen.getAllByTestId("hidden-input")[1].click();
    });

    // then
    expect(mockSelect).toHaveBeenCalledTimes(1);
  });

  test("SSH の hidden が false のとき目のマークをクリックすると、Hidden の状態が true になる", async () => {
    // setup
    const mockSetIsSSHHidden = vi.fn();
    (useState as Mock)
      .mockReturnValueOnce([false, mockSetIsSSHHidden])
      .mockReturnValueOnce([false, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });
    render(<TeamInfo />);

    // when
    await act(async () => {
      screen.getAllByRole("button")[0].click();
    });

    // then
    expect(mockSetIsSSHHidden).toHaveBeenCalledTimes(1);
    expect(mockSetIsSSHHidden).toHaveBeenCalledWith(true);
  });

  test("SSH の hidden が true のとき目のマークをクリックすると、Hidden の状態が false になる", async () => {
    // setup
    const mockSetIsSSHHidden = vi.fn();
    (useState as Mock)
      .mockReturnValueOnce([true, mockSetIsSSHHidden])
      .mockReturnValueOnce([false, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });
    render(<TeamInfo />);

    // when
    await act(async () => {
      screen.getAllByRole("button")[0].click();
    });

    // then
    expect(mockSetIsSSHHidden).toHaveBeenCalledTimes(1);
    expect(mockSetIsSSHHidden).toHaveBeenCalledWith(false);
  });

  test("Password  の hidden が false のとき目のマークをクリックすると、Hidden の状態が true になる", async () => {
    // setup
    const mockSetIsPasswordHidden = vi.fn();
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([false, mockSetIsPasswordHidden]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });
    render(<TeamInfo />);

    // when
    await act(async () => {
      screen.getAllByRole("button")[2].click();
    });

    // then
    expect(mockSetIsPasswordHidden).toHaveBeenCalledTimes(1);
    expect(mockSetIsPasswordHidden).toHaveBeenCalledWith(true);
  });

  test("Password  の hidden が true のとき目のマークをクリックすると、Hidden の状態が false になる", async () => {
    // setup
    const mockSetIsPasswordHidden = vi.fn();
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([true, mockSetIsPasswordHidden]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });
    render(<TeamInfo />);

    // when
    await act(async () => {
      screen.getAllByRole("button")[2].click();
    });

    // then
    expect(mockSetIsPasswordHidden).toHaveBeenCalledTimes(1);
    expect(mockSetIsPasswordHidden).toHaveBeenCalledWith(false);
  });

  test("SSH のコピーボタンを押すと、クリップボードに SSH コマンドがコピーされる", async () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([false, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });
    global.navigator = {
      // @ts-ignore
      clipboard: {
        writeText: vi.fn(),
      },
    };
    const mockStopPropagation = vi.spyOn(Event.prototype, "stopPropagation");
    const mockWriteText = vi.spyOn(navigator.clipboard, "writeText");
    render(<TeamInfo />);

    // when
    await act(async () => {
      screen.getAllByRole("button")[1].click();
    });

    // then
    expect(mockStopPropagation).toHaveBeenCalledTimes(1);
    expect(mockWriteText).toHaveBeenCalledTimes(1);
    expect(mockWriteText).toHaveBeenCalledWith(
      `ssh ${testUserGroup.bastion?.bastion_user}@${testUserGroup.bastion?.bastion_host} -p ${testUserGroup.bastion?.bastion_port}`
    );
  });

  test("Password のコピーボタンを押すと、クリップボードに Password がコピーされる", async () => {
    // setup
    (useState as Mock)
      .mockReturnValueOnce([false, vi.fn()])
      .mockReturnValueOnce([false, vi.fn()]);
    (useAuth as Mock).mockReturnValue({
      user: testUser,
      isLoading: false,
    });
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });
    global.navigator = {
      // @ts-ignore
      clipboard: {
        writeText: vi.fn(),
      },
    };
    const mockStopPropagation = vi.spyOn(Event.prototype, "stopPropagation");
    const mockWriteText = vi.spyOn(navigator.clipboard, "writeText");
    render(<TeamInfo />);

    // when
    await act(async () => {
      screen.getAllByRole("button")[3].click();
    });

    // then
    expect(mockStopPropagation).toHaveBeenCalledTimes(1);
    expect(mockWriteText).toHaveBeenCalledTimes(1);
    expect(mockWriteText).toHaveBeenCalledWith(
      testUserGroup.bastion?.bastion_password
    );
  });
});
