import "@testing-library/jest-dom";
import "@testing-library/jest-dom/extend-expect";

import React from "react";

import { render, screen } from "@testing-library/react";
import { vi } from "vitest";

import Home from "@/app/page";

vi.mock("@/components/_const", () => ({
  rule: "rule",
}));
vi.mock("@/components/MarkdownPreview", () => ({
  __esModule: true,
  default: ({ content }: { content: string }) => (
    <div data-testid="markdown-preview" data-content={content} />
  ),
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

describe("Home", () => {
  test("画面が表示されることを確認する", async () => {
    // setup
    render(<Home />);

    // verify
    expect(screen.queryByTestId("common-layout")).toBeInTheDocument();
    expect(screen.queryByTestId("common-layout")).toHaveAttribute(
      "data-title",
      "ルール"
    );
    expect(screen.queryByTestId("markdown-preview")).toBeInTheDocument();
    expect(screen.queryByTestId("markdown-preview")).toHaveAttribute(
      "data-content",
      "rule"
    );
  });
});
