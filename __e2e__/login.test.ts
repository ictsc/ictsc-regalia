import { test } from "@playwright/test";

import LoginPage from "./pages/login";

test("ログインできる", async ({ page }) => {
  // setup
  await LoginPage.goto(page);
  await LoginPage.validate(page);
  await LoginPage.UsernameInput(page).fill("admin");
  await LoginPage.PasswordInput(page).fill("password");

  // when
  await LoginPage.LoginButton(page).click();

  // then
  await page.waitForURL("/");
});
