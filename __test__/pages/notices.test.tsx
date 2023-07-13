import "@testing-library/jest-dom";

import { render, screen } from "@testing-library/react";
import { Mock, vi } from "vitest";

import useNotice from "@/hooks/notice";
import Notices from "@/pages/notices";
import { testNotice } from "@/types/Notice";

vi.mock("@/hooks/notice");
vi.mock("@/components/Navbar", () => ({
  __esModule: true,
  default: () => <div data-testid="navbar" />,
}));
vi.mock("@/components/NotificationCard", () => ({
  __esModule: true,
  default: () => <div data-testid="notification-card" />,
}));

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("Notices", () => {
  test("画面が表示されることを確認する", async () => {
    // setup
    (useNotice as Mock).mockReturnValue({
      notices: [testNotice],
      isLoading: false,
    });

    // when
    render(<Notices />);

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    expect(screen.getByTestId("notification-card")).toBeInTheDocument();

    // verify
    expect(useNotice).toHaveBeenCalledTimes(1);
  });

  test("通知が取得中の場合、ローディング画面が表示されることを確認する", async () => {
    // setup
    (useNotice as Mock).mockReturnValue({
      notices: [testNotice],
      isLoading: true,
    });

    // when
    render(<Notices />);

    // then
    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    expect(screen.queryByTestId("notification-card")).not.toBeInTheDocument();

    // verify
    expect(useNotice).toHaveBeenCalledTimes(1);
  });
});
