import { expect, test } from "@playwright/test";

import LoginPage from "../pages/login";
import ScoringPage from "../pages/scoring";
import ScoringProblemPage from "../pages/scoring/[code]";

test("画面項目が表示されること", async ({ page }) => {
  // setup
  await LoginPage.adminLogin(page);

  // when
  await ScoringPage.goto(page);

  // then
  await ScoringPage.validate(page);

  const problems = await ScoringPage.Problems(page);
  await expect(problems).toHaveCount(3);

  const problem1 = problems.nth(0).locator("td");
  await expect(problem1.nth(1)).toHaveText("1/1/1");
  await expect(problem1.nth(2)).toHaveText(
    "00000000-0000-4000-a000-000000000000"
  );
  await expect(problem1.nth(3)).toHaveText("abc");
  await expect(problem1.nth(4)).toHaveText("問題タイトル1");
  await expect(problem1.nth(5)).toHaveText("--- code: abc title:...");
  await expect(problem1.nth(6)).toHaveText("100");
  await expect(problem1.nth(7)).toHaveText("100");
  await expect(problem1.nth(8)).toHaveText("");
  await expect(problem1.nth(9)).toHaveText("自分");

  const problem2 = problems.nth(1).locator("td");
  await expect(problem2.nth(1)).toHaveText("1/1/1");
  await expect(problem2.nth(2)).toHaveText(
    "00000000-0000-4000-a000-000000000001"
  );
  await expect(problem2.nth(3)).toHaveText("def");
  await expect(problem2.nth(4)).toHaveText("問題タイトル2");
  await expect(problem2.nth(5)).toHaveText("問題内容2...");
  await expect(problem2.nth(6)).toHaveText("200");
  await expect(problem2.nth(7)).toHaveText("200");
  await expect(problem2.nth(8)).toHaveText("");
  await expect(problem2.nth(9)).toHaveText("自分");

  const problem3 = problems.nth(2).locator("td");
  await expect(problem3.nth(1)).toHaveText("1/1/1");
  await expect(problem3.nth(2)).toHaveText(
    "00000000-0000-4000-a000-000000000002"
  );
  await expect(problem3.nth(3)).toHaveText("ghi");
  await expect(problem3.nth(4)).toHaveText("問題タイトル3");
  await expect(problem3.nth(5)).toHaveText("問題内容3...");
  await expect(problem3.nth(6)).toHaveText("300");
  await expect(problem3.nth(7)).toHaveText("300");
  await expect(problem3.nth(8)).toHaveText("");
  await expect(problem3.nth(9)).toHaveText("自分");
});

test("採点ページへ遷移できること", async ({ page }) => {
  // setup
  await LoginPage.adminLogin(page);
  await ScoringPage.goto(page);
  const problems = await ScoringPage.Problems(page);
  const problem1 = problems.nth(0).locator("td");

  // when
  await problem1.nth(0).click();

  // then
  await ScoringProblemPage.validate(page, "問題タイトル1");
});

test("問題のプレビューができること", async ({ page }) => {
  // setup
  await LoginPage.adminLogin(page);
  await ScoringPage.goto(page);
  const problems = await ScoringPage.Problems(page);
  const problem1 = problems.nth(0).locator("td");

  // when
  await problem1.nth(1).click();

  // then
  await expect(ScoringPage.ProblemPreviewProblemInfo(page)).toHaveText(
    "問題タイトル1満点100 pt 採点基準100 pt採点する"
  );
  await expect(ScoringPage.ProblemPreviewProblemContent(page)).toHaveText(
    "問題内容1"
  );
});

// 問題のプレビューから採点ページへ遷移できること
test("問題のプレビューから採点ページへ遷移できること", async ({ page }) => {
  // setup
  await LoginPage.adminLogin(page);
  await ScoringPage.goto(page);
  const problems = await ScoringPage.Problems(page);
  const problem1 = problems.nth(0).locator("td");
  await problem1.nth(1).click();

  // when
  await ScoringPage.ProblemPreviewProblemInfo(page)
    .getByRole("link", { name: "採点する" })
    .click();

  // then
  await ScoringProblemPage.validate(page, "問題タイトル1");
});
