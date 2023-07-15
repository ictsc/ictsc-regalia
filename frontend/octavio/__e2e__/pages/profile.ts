import { expect, Page } from "@playwright/test";

const ProfilePage = {
  goto: async (page: Page) => {
    await page.goto("/profile");
  },
  validate: async (page: Page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("プロフィール");
  },
  waitFormSelector: async (page: Page) => {
    await page.waitForSelector(".title-ictsc >> text=プロフィール");
  },
};

export default ProfilePage;
