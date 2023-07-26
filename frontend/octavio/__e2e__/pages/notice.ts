import { expect, Page } from "@playwright/test";

const NoticesPage = {
  goto: async (page: Page) => {
    await page.goto("/notices");
  },
  validate: async (page: Page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("通知一覧");
  },
  waitFormSelector: async (page: Page) => {
    await page.waitForSelector(".title-ictsc >> text=通知一覧");
  },
};

export default NoticesPage;
