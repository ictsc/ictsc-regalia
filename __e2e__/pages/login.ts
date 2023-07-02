import { expect } from "@playwright/test";

import BasePage from "./page";

const LoginPage: BasePage = {
  goto: async (page) => {
    await page.goto("/login");
  },
  validate: async (page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("ログイン");
  },
  waitFormSelector: async (page) => {
    await page.waitForSelector(".title-ictsc >> text=ログイン");
  },
};

export default LoginPage;
