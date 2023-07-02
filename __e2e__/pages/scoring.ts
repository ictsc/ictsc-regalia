import { expect } from "@playwright/test";

import BasePage from "./page";

const ScoringPage: BasePage = {
  goto: async (page) => {
    await page.goto("/scoring");
  },
  validate: async (page) => {
    await expect(page.locator("th:nth-child(1)")).toHaveText("採点");
  },
  waitFormSelector: async (page) => {
    await page.waitForSelector("th >> text=採点");
  },
};

export default ScoringPage;
