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
  Notices: (page: Page) => page.locator(".notice-card"),
  FirstDismissButton: (page: Page) => page.locator(".dismiss-btn").first(),
  Problems: (page: Page) => page.locator(".problem-card"),
  FirstNoticeLink: (page: Page) => page.locator(".notice-link").first(),
};

export default ProblemsPage;
