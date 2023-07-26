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
};

export default ProblemPage;
