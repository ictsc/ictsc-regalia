import { expect, test } from "@playwright/test";

import LoginPage from "../pages/login";
import ProblemPage from "../pages/problems/[problemId]";

test("画面項目が表示されること", async ({ page }) => {
  // setup
  await LoginPage.user1Login(page);

  // when
  await ProblemPage.goto(page, "abc");

  // then
  await ProblemPage.validate(page, "問題タイトル1");
  await ProblemPage.waitFormSelector(page, "問題タイトル1");
  const matterCols = await ProblemPage.MatterCols(page);
  const matter1 = await matterCols.nth(0).locator("td");
  await expect(matterCols).toHaveCount(1);
  await expect(matter1.nth(0)).toHaveText("192.168.100.1 Copy");
  await expect(matter1.nth(1)).toHaveText("user Copy");
  await expect(matter1.nth(2)).toHaveText("password Copy");
  await expect(matter1.nth(3)).toHaveText("22");
  await expect(matter1.nth(4)).toHaveText("ssh");
  await expect(await ProblemPage.Body(page)).toHaveText("問題内容1");
  const answerCols = await ProblemPage.AnswerCols(page);
  const answer1 = await answerCols.nth(0).locator("td");
  const answer2 = await answerCols.nth(1).locator("td");
  const answer3 = await answerCols.nth(2).locator("td");
  await expect(answerCols).toHaveCount(4);
  await expect(answer1.nth(0)).toHaveText("1990-01-01 21:34:56");
  await expect(answer1.nth(1)).toHaveText("abc");
  await expect(answer1.nth(2)).toHaveText("問題タイトル1");
  await expect(answer1.nth(3)).toHaveText("0 pt");
  await expect(answer1.nth(4)).toHaveText("○");
  await expect(answer2.nth(0)).toHaveText("1990-01-01 21:34:56");
  await expect(answer2.nth(1)).toHaveText("abc");
  await expect(answer2.nth(2)).toHaveText("問題タイトル1");
  await expect(answer2.nth(3)).toHaveText("100 pt");
  await expect(answer2.nth(4)).toHaveText("○");
  await expect(answer3.nth(0)).toHaveText("1990-01-01 21:34:56");
  await expect(answer3.nth(1)).toHaveText("abc");
  await expect(answer3.nth(2)).toHaveText("問題タイトル1");
  await expect(answer3.nth(3)).toHaveText("-- pt");
  await expect(answer3.nth(4)).toHaveText("採点中");
});

// TODO: 再展開ができること, 再展開はバックエンドが他のシステムと連携しているため失敗する.
// test("再展開ができること", async ({ page }) => {
//   // setup
//   await LoginPage.user1Login(page);
//   await ProblemPage.goto(page, "abc");
//
//   // when
//   await ProblemPage.ReCreateButton(page).click();
//   await ProblemPage.ReCreateAcceptButton(page).click();
//
//   // then
// });

test("投稿内容のプレビューができること", async ({ page }) => {
  // setup
  await LoginPage.user1Login(page);
  await ProblemPage.goto(page, "abc");
  await ProblemPage.validate(page, "問題タイトル1");
  await ProblemPage.AnswerTextArea(page).fill("# 回答テスト");

  // when
  await ProblemPage.PreviewButton(page).click();

  // then
  await expect(ProblemPage.AnswerFormPreview(page)).toHaveText("回答テスト");
});

test("ファイルのダウンロードができる", async ({ page }) => {
  // setup
  await LoginPage.user1Login(page);
  await ProblemPage.goto(page, "abc");
  await ProblemPage.validate(page, "問題タイトル1");
  const downloadPromise = page.waitForEvent("download");
  const answerCols = await ProblemPage.AnswerCols(page);
  const answer1 = await answerCols.nth(0).locator("td");

  // when
  await answer1.nth(6).getByRole("link", { name: "ダウンロード" }).click();

  // then
  const download = await downloadPromise;
  await expect(download.suggestedFilename()).toBe("ictsc-abc-631197296.md");
});

test("回答をプレビューできること", async ({ page }) => {
  // setup
  await LoginPage.user1Login(page);
  await ProblemPage.goto(page, "abc");
  await ProblemPage.validate(page, "問題タイトル1");
  const answerCols = await ProblemPage.AnswerCols(page);
  const answer1 = await answerCols.nth(0).locator("td");

  // when
  await answer1.nth(5).getByRole("link", { name: "投稿内容" }).click();

  // then
  await expect(ProblemPage.AnswerPreviewTeamInfo(page)).toHaveText(
    "チーム: team1(user-org1)"
  );
  await expect(ProblemPage.AnswerPreviewCreatedAt(page)).toHaveText(
    "1990-01-01 21:34:56"
  );
  await expect(ProblemPage.AnswerPreview(page)).toHaveText("回答内容1");
});

test("回答を提出できること", async ({ page }) => {
  // setup
  await LoginPage.user1Login(page);
  await ProblemPage.goto(page, "abc");
  await ProblemPage.validate(page, "問題タイトル1");
  await ProblemPage.AnswerTextArea(page).fill("# 回答テスト");
  await ProblemPage.AnswerSubmitButton(page).click();

  // when
  await ProblemPage.AnswerSubmitAcceptButton(page).click();

  // then
  await expect(page.locator(".alert-success")).toHaveText("投稿に成功しました");
});
