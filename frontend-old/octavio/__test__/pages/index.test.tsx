import "@testing-library/jest-dom";

import React from "react";

import { render, screen } from "@testing-library/react";
import { vi } from "vitest";

import Home from "@/app/page";

vi.mock("@/components/_const", () => ({
  rule: "rule",
}));
vi.mock("@/components/markdown-preview", () => ({
  __esModule: true,
  default: ({ content }: { content: string }) => (
    <div data-testid="markdown-preview" data-content={content} />
  ),
}));

beforeEach(() => {
  // toHaveBeenCalledTimes がテストごとにリセットされるようにする
  vi.clearAllMocks();
});

describe("Home", () => {
  test("画面が表示されることを確認する", async () => {
    // when
    render(<Home />);

    // then
    expect(screen.queryByTestId("markdown-preview")).toBeInTheDocument();
    expect(screen.queryByTestId("markdown-preview")).toHaveAttribute(
      "data-content",
      "rule"
    );
  });
});
