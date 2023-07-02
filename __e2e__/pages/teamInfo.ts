import { expect } from "@playwright/test";

import BasePage from "./page";

const TeamInfoPage: BasePage = {
  goto: async (page) => {
    await page.goto("/team_info");
  },
  validate: async (page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("チーム情報");
  },
  waitFormSelector: async (page) => {
    await page.waitForSelector(".title-ictsc >> text=チーム情報");
  },
};

export default TeamInfoPage;
