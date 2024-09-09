import { expect, test } from "@playwright/test";

import LoginPage from "./pages/login";
import RankingPage from "./pages/ranking";

test("画面項目が表示されること", async ({ page }) => {
  // setup
  await LoginPage.adminLogin(page);
  await RankingPage.goto(page);
  await RankingPage.validate(page);

  // when
  const cols = await RankingPage.Cells(page);
  const rank1stCell = await cols.nth(0).locator("td");
  const rank2ndCell = await cols.nth(1).locator("td");
  const rank3rdCell = await cols.nth(2).locator("td");
  await expect(cols).toHaveCount(3);
  await expect(rank1stCell.nth(0)).toHaveText("1");
  await expect(rank1stCell.nth(1)).toHaveText("team1");
  await expect(rank1stCell.nth(2)).toHaveText("user-org1");
  await expect(rank1stCell.nth(3)).toHaveText("600pt");
  await expect(rank2ndCell.nth(0)).toHaveText("2");
  await expect(rank2ndCell.nth(1)).toHaveText("team2");
  await expect(rank2ndCell.nth(2)).toHaveText("user-org2");
  await expect(rank2ndCell.nth(3)).toHaveText("0pt");
  await expect(rank3rdCell.nth(0)).toHaveText("2");
  await expect(rank3rdCell.nth(1)).toHaveText("team3");
  await expect(rank3rdCell.nth(2)).toHaveText("user-org3");
  await expect(rank3rdCell.nth(3)).toHaveText("0pt");
});
