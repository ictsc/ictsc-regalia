import { expect, Page } from "@playwright/test";

const SignUpPage = {
  goto: async (page: Page, query?: string) => {
    await page.goto(`/signUp${query}`);
  },
  validate: async (page: Page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("ユーザー登録");
  },
  waitFormSelector: async (page: Page) => {
    await page.waitForSelector(".title-ictsc >> text=ユーザー登録");
  },
  UsernameInput: (page: Page) => page.locator("#username"),
  PasswordInput: (page: Page) => page.locator("#password"),
  SignUpButton: (page: Page) => page.locator("#signUpBtn"),
};

export default SignUpPage;
