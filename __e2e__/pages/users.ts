import { expect } from "@playwright/test";

import BasePage from "./page";

const UsersPage: BasePage = {
  goto: async (page) => {
    await page.goto("/users");
  },
  validate: async (page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("参加者一覧");
  },
  waitFormSelector: async (page) => {
    await page.waitForSelector(".title-ictsc >> text=参加者");
  },
};

export default UsersPage;
