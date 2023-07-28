import { expect, Page } from "@playwright/test";

const ScoringPage = {
  goto: async (page: Page) => {
    await page.goto("/scoring");
  },
  validate: async (page: Page) => {
    await expect(page.locator("th:nth-child(1)")).toHaveText("採点");
  },
  waitFormSelector: async (page: Page) => {
    await page.waitForSelector("th >> text=採点");
  },
  Problems: (page: Page) => page.locator("tbody").locator("tr"),
  ProblemPreviewProblemContent: (page: Page) =>
    page.locator(".problem-preview"),
  ProblemPreviewProblemInfo: (page: Page) => page.locator(".problem-info"),
};

export default ScoringPage;
