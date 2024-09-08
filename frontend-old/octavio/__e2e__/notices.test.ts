import { expect, test } from "@playwright/test";

import LoginPage from "./pages/login";
import NoticesPage from "./pages/notice";

test("画面項目が表示されること", async ({ page }) => {
  // setup
  await LoginPage.user1Login(page);

  // when
  await NoticesPage.goto(page);

  // then
  await NoticesPage.validate(page);
  const notices = await NoticesPage.Notices(page);
  await expect(notices).toHaveCount(3);
  await expect(notices.nth(0).locator(".notice-title")).toHaveText(
    "通知タイトル3"
  );
  await expect(notices.nth(0).locator(".notice-body")).toHaveText(
    "通知メッセージ3"
  );
  await expect(notices.nth(1).locator(".notice-title")).toHaveText(
    "通知タイトル2"
  );
  await expect(notices.nth(1).locator(".notice-body")).toHaveText(
    "通知メッセージ2"
  );
  await expect(notices.nth(2).locator(".notice-title")).toHaveText(
    "通知タイトル1"
  );
  await expect(notices.nth(2).locator(".notice-body")).toHaveText(
    "通知メッセージ1"
  );
});
