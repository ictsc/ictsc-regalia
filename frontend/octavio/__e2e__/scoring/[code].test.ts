import { expect, test } from "@playwright/test";

import LoginPage from "../pages/login";
import ScoringProblemPage from "../pages/scoring/[code]";

test("画面項目が表示されること", async ({ page }) => {
  // setup
  await LoginPage.adminLogin(page);

  // when
  await ScoringProblemPage.goto(page, "abc");

  // then
  await ScoringProblemPage.validate(page, "問題タイトル1");
  await ScoringProblemPage.waitFormSelector(page, "問題タイトル1");

  const matterCols = await ScoringProblemPage.MatterCols(page);
  const matter1 = await matterCols.nth(0).locator("td");
  await expect(matterCols).toHaveCount(1);
  await expect(matter1.nth(0)).toHaveText("192.168.100.1 Copy");
  await expect(matter1.nth(1)).toHaveText("user Copy");
  await expect(matter1.nth(2)).toHaveText("password Copy");
  await expect(matter1.nth(3)).toHaveText("22");
  await expect(matter1.nth(4)).toHaveText("ssh");

  await expect(await ScoringProblemPage.Body(page)).toHaveText("問題内容1");

  const answerRows = await ScoringProblemPage.AnswerRows(page);
  await expect(answerRows.nth(0)).toHaveText("1");
  await expect(answerRows.nth(1)).toHaveText("1");
  await expect(answerRows.nth(2)).toHaveText("1");

  await expect(await ScoringProblemPage.AnswerPreviewTeamInfo(page)).toHaveText(
    "チーム: team1(user-org1)"
  );
  await expect(
    await ScoringProblemPage.AnswerPreviewCreatedAt(page)
  ).toHaveText("1990/01/01 21:34:56");
  await expect(await ScoringProblemPage.AnswerPreview(page)).toHaveText(
    "回答内容2"
  );
});

test("採点フィルタの「すべて」を選択した場合、すべての回答が表示されること", async ({
  page,
}) => {
  // setup
  await LoginPage.adminLogin(page);
  await ScoringProblemPage.goto(page, "abc");
  await ScoringProblemPage.AnswerFilter(page).selectOption("すべて");

  // when
  await expect(
    await ScoringProblemPage.AnswerPreviewTeamInfo(page).nth(0)
  ).toHaveText("チーム: team1(user-org1)");
  await expect(
    await ScoringProblemPage.AnswerPreviewCreatedAt(page).nth(0)
  ).toHaveText("1990/01/01 21:34:56");
  await expect(await ScoringProblemPage.AnswerPreview(page).nth(0)).toHaveText(
    "回答内容1"
  );
  await expect(
    await ScoringProblemPage.AnswerInputForm(page).nth(0)
  ).toHaveValue("0");
  await expect(
    await ScoringProblemPage.AnswerPreviewTeamInfo(page).nth(1)
  ).toHaveText("チーム: team1(user-org1)");
  await expect(
    await ScoringProblemPage.AnswerPreviewCreatedAt(page).nth(1)
  ).toHaveText("1990/01/01 21:34:56");
  await expect(await ScoringProblemPage.AnswerPreview(page).nth(1)).toHaveText(
    "回答内容3"
  );
  await expect(
    await ScoringProblemPage.AnswerInputForm(page).nth(1)
  ).toHaveValue("100");
  await expect(
    await ScoringProblemPage.AnswerPreviewTeamInfo(page).nth(2)
  ).toHaveText("チーム: team1(user-org1)");
  await expect(
    await ScoringProblemPage.AnswerPreviewCreatedAt(page).nth(2)
  ).toHaveText("1990/01/01 21:34:56");
  await expect(await ScoringProblemPage.AnswerPreview(page).nth(2)).toHaveText(
    "回答内容2"
  );
  await expect(
    await ScoringProblemPage.AnswerInputForm(page).nth(2)
  ).toHaveValue("");
});

test("採点フィルタの「採点済み」を選択した場合、採点済みの回答が表示されること", async ({
  page,
}) => {
  // setup
  await LoginPage.adminLogin(page);
  await ScoringProblemPage.goto(page, "abc");
  await ScoringProblemPage.AnswerFilter(page).selectOption("採点済みのみ");

  // when
  await expect(
    await ScoringProblemPage.AnswerPreviewTeamInfo(page).nth(0)
  ).toHaveText("チーム: team1(user-org1)");
  await expect(
    await ScoringProblemPage.AnswerPreviewCreatedAt(page).nth(0)
  ).toHaveText("1990/01/01 21:34:56");
  await expect(await ScoringProblemPage.AnswerPreview(page).nth(0)).toHaveText(
    "回答内容1"
  );
  await expect(
    await ScoringProblemPage.AnswerInputForm(page).nth(0)
  ).toHaveValue("0");
  await expect(
    await ScoringProblemPage.AnswerPreviewTeamInfo(page).nth(1)
  ).toHaveText("チーム: team1(user-org1)");
  await expect(
    await ScoringProblemPage.AnswerPreviewCreatedAt(page).nth(1)
  ).toHaveText("1990/01/01 21:34:56");
  await expect(await ScoringProblemPage.AnswerPreview(page).nth(1)).toHaveText(
    "回答内容3"
  );
  await expect(
    await ScoringProblemPage.AnswerInputForm(page).nth(1)
  ).toHaveValue("100");
});

test("採点ができる", async ({ page }) => {
  // setup
  await LoginPage.adminLogin(page);
  await ScoringProblemPage.goto(page, "abc");
  await ScoringProblemPage.AnswerInputForm(page).fill("100");

  // when
  await ScoringProblemPage.AnswerSubmitButton(page).click();

  // then
  await expect(await ScoringProblemPage.AnswerPreview(page)).not.toBeVisible();
});
