import "@testing-library/jest-dom";

import React from "react";

import { render, screen } from "@testing-library/react";
import { Mock, vi } from "vitest";

import Ranking from "@/app/ranking/page";
import useRanking from "@/hooks/ranking";
import { testRank } from "@/types/Rank";

vi.mock("@/hooks/ranking");
vi.mock("@/components/LoadingPage", () => ({
  __esModule: true,
  default: () => <div data-testid="loading" />,
}));

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("Ranking", () => {
  test("画面が表示されることを確認する", async () => {
    // setup
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: false,
    });

    // when
    render(<Ranking />);

    // then
    const cells = screen.getAllByRole("cell");
    expect(cells[0]).toHaveTextContent(testRank.rank.toString());
    expect(cells[1]).toHaveTextContent(testRank.user_group.name);
    expect(cells[2]).toHaveTextContent(testRank.user_group.organization);
    expect(cells[3]).toHaveTextContent(`${testRank.point}pt`);

    // verify
    expect(useRanking).toHaveBeenCalledTimes(1);
  });

  test("ランキングが取得中の場合、ローディング画面が表示されることを確認する", async () => {
    // setup
    (useRanking as Mock).mockReturnValue({
      ranking: [testRank],
      loading: true,
    });

    // when
    render(<Ranking />);

    // then
    expect(screen.getByTestId("loading")).toBeInTheDocument();

    // verify
    expect(useRanking).toHaveBeenCalledTimes(1);
  });
});
