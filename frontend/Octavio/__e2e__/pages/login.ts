import { expect, Page } from "@playwright/test";

const LoginPage = {
  goto: async (page: Page) => {
    await page.goto("/login");
  },
  validate: async (page: Page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("ログイン");
  },
  waitFormSelector: async (page: Page) => {
    await page.waitForSelector(".title-ictsc >> text=ログイン");
  },
  login: async (page: Page) => {
    await LoginPage.goto(page);
    await LoginPage.validate(page);
    await LoginPage.UsernameInput(page).fill("admin");
    await LoginPage.PasswordInput(page).fill("password");
    await LoginPage.LoginButton(page).click();
    await page.waitForURL("/");
  },
  UsernameInput: (page: Page) => page.locator("#username"),
  PasswordInput: (page: Page) => page.locator("#password"),
  LoginButton: (page: Page) => page.locator("#loginBtn"),
};

export default LoginPage;
