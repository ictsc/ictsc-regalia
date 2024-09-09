import { expect, test } from "@playwright/test";

import LoginPage from "./pages/login";
import ProfilePage from "./pages/profile";

test("画面項目が表示されること", async ({ page }) => {
  // setup
  await LoginPage.user1Login(page);

  // when
  await ProfilePage.goto(page);

  // then
  await ProfilePage.validate(page);
  await ProfilePage.waitFormSelector(page);

  await expect(ProfilePage.DisplayNameInput(page)).toBeVisible();
  await expect(ProfilePage.Teams(page)).toHaveText("team1");
  await expect(ProfilePage.SelfIntroductionInput(page)).toBeVisible();
  await expect(ProfilePage.GitHubIDInput(page)).toBeVisible();
  await expect(ProfilePage.TwitterIDInput(page)).toBeVisible();
  await expect(ProfilePage.FacebookIDInput(page)).toBeVisible();
});

test("各項目が入力されていること", async ({ page }) => {
  // setup
  await LoginPage.user1Login(page);

  // when
  await ProfilePage.goto(page);

  // then
  await ProfilePage.validate(page);
  await ProfilePage.waitFormSelector(page);

  await expect(ProfilePage.DisplayNameInput(page)).toHaveValue("user1");
  await expect(ProfilePage.SelfIntroductionInput(page)).toHaveValue(
    "自己紹介内容1"
  );
  await expect(ProfilePage.GitHubIDInput(page)).toHaveValue("gid");
  await expect(ProfilePage.TwitterIDInput(page)).toHaveValue("tid");
  await expect(ProfilePage.FacebookIDInput(page)).toHaveValue("fid");
});

test("各項目が変更できること", async ({ page }) => {
  // setup
  await LoginPage.user1Login(page);
  await ProfilePage.goto(page);
  await ProfilePage.validate(page);
  await ProfilePage.waitFormSelector(page);
  await ProfilePage.DisplayNameInput(page).fill("user1-1");
  await ProfilePage.SelfIntroductionInput(page).fill("自己紹介内容1-1");
  await ProfilePage.GitHubIDInput(page).fill("gid-1");
  await ProfilePage.TwitterIDInput(page).fill("tid-1");
  await ProfilePage.FacebookIDInput(page).fill("fid-1");

  // when
  await ProfilePage.SubmitButton(page).click();

  // then
  await expect(page.locator(".alert-success")).toHaveText(
    "プロフィールを更新しました"
  );
});
