import { render, screen } from "@testing-library/react";
import { vi } from "vitest";

import { Props } from "@/app/problems/__components/DescriptiveProblemEdit";
import Index from "@/app/problems/new/descriptive/page";
import "@testing-library/jest-dom";

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
    // when
    render(<Index />);

    // then
    expect(
      screen.queryByTestId("descriptive-problem-edit"),
    ).toBeInTheDocument();
  });
});
