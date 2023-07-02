import { expect, Page } from "@playwright/test";

import BasePage from "./page";

const IndexPage: BasePage = {
  goto: async (page: Page) => {
    await page.goto("/");
  },
  validate: async (page: Page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("ルール");
  },
  waitFormSelector: async (page: Page) => {
    await page.waitForSelector(".title-ictsc >> text=ルール");
  },
};

export default IndexPage;
