import { test } from "@playwright/test";

import SignUpPage from "./pages/signUp";

test("参加者登録できる", async ({ page }) => {
  // setup
  await SignUpPage.goto(
    page,
    "?user_group_id=00000000-0000-4000-a000-000000000001&invitation_code=test-invitation-code"
  );
  await SignUpPage.validate(page);
  await SignUpPage.UsernameInput(page).fill("testUser");
  await SignUpPage.PasswordInput(page).fill("password");

  // when
  await SignUpPage.SignUpButton(page).click();

  // then
  await page.waitForURL("/login");
});
