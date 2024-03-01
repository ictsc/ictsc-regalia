import { useQuery } from "@connectrpc/connect-query";
import { render, screen } from "@testing-library/react";
import { Mock, vi } from "vitest";

import { Props } from "@/app/problems/__components/DescriptiveProblemEdit";
import Index from "@/app/problems/edit/[id]/page";
import "@testing-library/jest-dom";
import { Problem } from "@/proto/admin/v1/problem_pb";

vi.mock("next/navigation", () => ({
  useParams: vi.fn(() => ({ id: "abc" })),
}));

vi.mock("@connectrpc/connect-query");

vi.mock("@/app/problems/__components/DescriptiveProblemEdit", () => ({
  __esModule: true,
  default: (props: Props) => (
    <div data-testid="descriptive-problem-edit" data-content={props} />
  ),
}));

describe("ProblemEdit", () => {
  beforeEach(() => {
    // toHaveBeenCalledTimes がテストごとにリセットされるようにする
    vi.clearAllMocks();
  });

  it("画面が表示されることを確認する", async () => {
    // setups
    const problem = new Problem({
      id: "abc",
      body: {
        case: "descriptive",
        value: {
          body: "body",
          connectionInfos: [],
        },
      },
    });
    (useQuery as Mock).mockReturnValue({
      data: { problem },
    });

    // when
    render(<Index />);

    // then
    expect(useQuery).toHaveBeenCalledWith(expect.any(Object), { id: "abc" });
    expect(
      screen.queryByTestId("descriptive-problem-edit"),
    ).toBeInTheDocument();

    // verify
    expect(useQuery).toHaveBeenCalledTimes(1);
  });

  it("画面に「この問題タイプは未対応です。」と表示されることを確認する", async () => {
    // setups
    (useQuery as Mock).mockReturnValue({
      data: null,
    });

    // when
    render(<Index />);

    // then
    expect(useQuery).toHaveBeenCalledWith(expect.any(Object), { id: "abc" });
    expect(screen.queryByTestId("unsupported")).toBeInTheDocument();

    // verify
    expect(useQuery).toHaveBeenCalledTimes(1);
  });
});
