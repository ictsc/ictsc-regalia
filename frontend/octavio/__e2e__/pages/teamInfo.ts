import { expect, Page } from "@playwright/test";

const TeamInfoPage = {
  goto: async (page: Page) => {
    await page.goto("/team_info");
  },
  validate: async (page: Page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("チーム情報");
  },
  waitFormSelector: async (page: Page) => {
    await page.waitForSelector(".title-ictsc >> text=チーム情報");
  },
};

export default TeamInfoPage;
