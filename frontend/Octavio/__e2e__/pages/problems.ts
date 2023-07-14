import { expect, Page } from "@playwright/test";

const ProblemsPage = {
  goto: async (page: Page) => {
    await page.goto("/problems");
  },
  validate: async (page: Page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("問題一覧");
  },
  waitFormSelector: async (page: Page) => {
    await page.waitForSelector(".title-ictsc >> text=問題");
  },
};

export default ProblemsPage;
