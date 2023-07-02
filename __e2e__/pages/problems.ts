import { expect } from "@playwright/test";

import BasePage from "./page";

const ProblemsPage: BasePage = {
  goto: async (page) => {
    await page.goto("/problems");
  },
  validate: async (page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("問題一覧");
  },
  waitFormSelector: async (page) => {
    await page.waitForSelector(".title-ictsc >> text=問題");
  },
};

export default ProblemsPage;
