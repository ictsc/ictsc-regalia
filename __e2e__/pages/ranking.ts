import { expect } from "@playwright/test";

import BasePage from "./page";

const RankingPage: BasePage = {
  goto: async (page) => {
    await page.goto("/ranking");
  },
  validate: async (page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("ランキング");
  },
  waitFormSelector: async (page) => {
    await page.waitForSelector(".title-ictsc >> text=ランキング");
  },
};

export default RankingPage;
