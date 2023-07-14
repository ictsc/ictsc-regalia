import { expect, Page } from "@playwright/test";

const UsersPage = {
  goto: async (page: Page) => {
    await page.goto("/users");
  },
  validate: async (page: Page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("参加者一覧");
  },
  waitFormSelector: async (page: Page) => {
    await page.waitForSelector(".title-ictsc >> text=参加者");
  },
};

export default UsersPage;
