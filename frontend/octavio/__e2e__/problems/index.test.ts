import { expect, test } from "@playwright/test";

import LoginPage from "../pages/login";
import ProblemsPage from "../pages/problems";
import ProblemPage from "../pages/problems/[problemId]";

test("画面項目が表示されること", async ({ page }) => {
  // setup
  await LoginPage.login(page, "user1", "password");
  await ProblemsPage.goto(page);
  await ProblemsPage.validate(page);

  // when
  const notices = await ProblemsPage.Notices(page);
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
  const problems = await ProblemsPage.Problems(page);
  await expect(problems).toHaveCount(3);
  await expect(problems.nth(0).locator(".problem-code")).toHaveText("abc");
  await expect(problems.nth(0).locator(".problem-title")).toHaveText(
    "問題タイトル1"
  );
  await expect(problems.nth(0).locator(".problem-point")).toHaveText(
    "100/100pt"
  );
  await expect(problems.nth(1).locator(".problem-code")).toHaveText("def");
  await expect(problems.nth(1).locator(".problem-title")).toHaveText(
    "問題タイトル2"
  );
  await expect(problems.nth(1).locator(".problem-point")).toHaveText(
    "200/200pt"
  );
  await expect(problems.nth(2).locator(".problem-code")).toHaveText("ghi");
  await expect(problems.nth(2).locator(".problem-title")).toHaveText(
    "問題タイトル3"
  );
  await expect(problems.nth(2).locator(".problem-point")).toHaveText(
    "300/300pt"
  );
});

test("問題ページに遷移できる", async ({ page }) => {
  // setup
  await LoginPage.login(page, "user1", "password");
  await ProblemsPage.goto(page);
  await ProblemsPage.validate(page);

  // when
  await ProblemsPage.Problems(page).nth(0).locator(".problem-title").click();

  // then
  await ProblemPage.waitFormSelector(page, "問題タイトル1");
  expect(page.url().split("/").pop()).toBe("abc");
  await ProblemPage.validate(page, "問題タイトル1");
});
