import { expect, Page } from "@playwright/test";

const ScoringProblemPage = {
  goto: async (page: Page, code: string) => {
    await page.goto(`/scoring/${code}`);
  },
  validate: async (page: Page, title: string) => {
    await expect(page.locator(".title-ictsc")).toHaveText(title);
  },
  waitFormSelector: async (page: Page, title: string) => {
    await page.waitForSelector(`.title-ictsc >> text=${title}`);
  },
  Body: (page: Page) => page.locator(".problem-body"),
  MatterCols: (page: Page) => page.locator("tbody").nth(0).locator("tr"),
  AnswerRows: (page: Page) =>
    page.locator("tbody").nth(1).locator("tr").nth(0).locator("td"),
  AnswerFilter: (page: Page) => page.getByRole("combobox"),
  AnswerPreview: (page: Page) => page.locator(".answer-preview"),
  AnswerPreviewCreatedAt: (page: Page) =>
    page.locator(".answer-preview-created-at"),
  AnswerPreviewTeamInfo: (page: Page) =>
    page.locator(".answer-preview-team-info"),
  AnswerInputForm: (page: Page) => page.getByRole("textbox"),
  AnswerSubmitButton: (page: Page) =>
    page.getByRole("button", { name: "採点" }),
};

export default ScoringProblemPage;
