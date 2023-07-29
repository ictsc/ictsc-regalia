import { expect, test } from "@playwright/test";

import LoginPage from "./pages/login";
import UsersPage from "./pages/users";

test("画面項目が表示されること", async ({ page }) => {
  // setup
  await LoginPage.user1Login(page);

  // when
  await UsersPage.goto(page);

  // then
  await UsersPage.waitFormSelector(page);
  await UsersPage.validate(page);

  const cols = await UsersPage.Users(page);
  await expect(cols).toHaveCount(9);

  const user1 = await cols.nth(0).locator("td");
  await expect(user1.nth(0)).toHaveText("user1");
  await expect(user1.nth(1)).toHaveText("team1");
  await expect(user1.nth(2)).toHaveText("自己紹介内容1");

  const user2 = await cols.nth(1).locator("td");
  await expect(user2.nth(0)).toHaveText("user2");
  await expect(user2.nth(1)).toHaveText("team1");
  await expect(user2.nth(2)).toHaveText("");

  const user3 = await cols.nth(2).locator("td");
  await expect(user3.nth(0)).toHaveText("user3");
  await expect(user3.nth(1)).toHaveText("team1");
  await expect(user3.nth(2)).toHaveText("");

  const user4 = await cols.nth(3).locator("td");
  await expect(user4.nth(0)).toHaveText("user4");
  await expect(user4.nth(1)).toHaveText("team2");
  await expect(user4.nth(2)).toHaveText("自己紹介内容2");

  const user5 = await cols.nth(4).locator("td");
  await expect(user5.nth(0)).toHaveText("user5");
  await expect(user5.nth(1)).toHaveText("team2");
  await expect(user5.nth(2)).toHaveText("");

  const user6 = await cols.nth(5).locator("td");
  await expect(user6.nth(0)).toHaveText("user6");
  await expect(user6.nth(1)).toHaveText("team2");
  await expect(user6.nth(2)).toHaveText("");

  const user7 = await cols.nth(6).locator("td");
  await expect(user7.nth(0)).toHaveText("user7");
  await expect(user7.nth(1)).toHaveText("team3");
  await expect(user7.nth(2)).toHaveText("自己紹介内容3");

  const user8 = await cols.nth(7).locator("td");
  await expect(user8.nth(0)).toHaveText("user8");
  await expect(user8.nth(1)).toHaveText("team3");
  await expect(user8.nth(2)).toHaveText("");

  const user9 = await cols.nth(8).locator("td");
  await expect(user9.nth(0)).toHaveText("user9");
  await expect(user9.nth(1)).toHaveText("team3");
  await expect(user9.nth(2)).toHaveText("");
});
