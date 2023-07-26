import { expect, Page } from "@playwright/test";

const RankingPage = {
  goto: async (page: Page) => {
    await page.goto("/ranking");
  },
  validate: async (page: Page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("ランキング");
  },
  waitFormSelector: async (page: Page) => {
    await page.waitForSelector(".title-ictsc >> text=ランキング");
  },
  Cells: (page: Page) => page.locator("tbody >> tr"),
};
export default RankingPage;
