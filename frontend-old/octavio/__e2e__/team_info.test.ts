import { expect, test } from "@playwright/test";

import LoginPage from "./pages/login";
import TeamInfoPage from "./pages/teamInfo";

test("画面項目が表示されること", async ({ page }) => {
  // setup
  await LoginPage.user1Login(page);

  // when
  await TeamInfoPage.goto(page);

  // then
  await TeamInfoPage.waitFormSelector(page);
  await TeamInfoPage.validate(page);

  await expect(TeamInfoPage.GroupNameAndOrganization(page)).toHaveText(
    "team1@user-org1"
  );
  await expect(TeamInfoPage.Ranking(page)).toHaveText("1/3 teams");
  await expect(TeamInfoPage.Score(page)).toHaveText("600pt");
  await expect(TeamInfoPage.SSHInfo(page)).toHaveValue("ssh user1@host -p 22");
  await expect(TeamInfoPage.PasswordInfo(page)).toHaveValue("password");
});
