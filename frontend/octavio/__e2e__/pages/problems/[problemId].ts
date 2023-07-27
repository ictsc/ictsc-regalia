import { expect, Page } from "@playwright/test";

const ProblemPage = {
  goto: async (page: Page, code: string) => {
    await page.goto(`/problems/${code}`);
  },
  validate: async (page: Page, title: string) => {
    await expect(page.locator("h1")).toHaveText(title);
  },
  waitFormSelector: async (page: Page, title: string) => {
    await page.waitForSelector(`h1 >> text=${title}`);
  },
  Body: (page: Page) => page.locator(".problem-body"),
  MatterCols: (page: Page) => page.locator("tbody").nth(0).locator("tr"),
  AnswerCols: (page: Page) => page.locator("tbody").nth(1).locator("tr"),
  ReCreateButton: (page: Page) =>
    page.getByRole("button", { name: "再展開を行う" }),
  ReCreateAcceptButton: (page: Page) =>
    page.getByRole("button", { name: "問題の再展開を行う" }),
  AnswerTextArea: (page: Page) => page.getByRole("textbox"),
  AnswerSubmitButton: (page: Page) =>
    page.getByRole("button", { name: "提出確認" }),
  AnswerSubmitAcceptButton: (page: Page) =>
    page.getByRole("button", { name: "この内容で提出" }),
  PreviewButton: (page: Page) => page.getByRole("button", { name: "Preview" }),
  AnswerFormPreview: (page: Page) => page.locator(".answer-form-preview"),
  AnswerPreview: (page: Page) => page.locator(".answer-preview"),
  AnswerPreviewCreatedAt: (page: Page) =>
    page.locator(".answer-preview-created-at"),
  AnswerPreviewTeamInfo: (page: Page) =>
    page.locator(".answer-preview-team-info"),
};

export default ProblemPage;
