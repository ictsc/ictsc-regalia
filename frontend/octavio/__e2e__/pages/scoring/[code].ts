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
};

export default ScoringProblemPage;
