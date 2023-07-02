import { expect } from "@playwright/test";

import BasePage from "./page";

const ProfilePage: BasePage = {
  goto: async (page) => {
    await page.goto("/profile");
  },
  validate: async (page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("プロフィール");
  },
  waitFormSelector: async (page) => {
    await page.waitForSelector(".title-ictsc >> text=プロフィール");
  },
};

export default ProfilePage;
