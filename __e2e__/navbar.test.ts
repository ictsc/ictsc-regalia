import { expect, test } from "@playwright/test";

import IndexPage from "./pages/Index";
import LoginPage from "./pages/login";
import ICTSCNavbar from "./pages/navbar";
import ProblemsPage from "./pages/problems";
import ProfilePage from "./pages/profile";
import RankingPage from "./pages/ranking";
import ScoringPage from "./pages/scoring";
import TeamInfoPage from "./pages/teamInfo";
import UsersPage from "./pages/users";

test.describe("未ログイン状態", () => {
  test("ログインページに遷移できる", async ({ page }) => {
    // setup
    await IndexPage.goto(page);
    await IndexPage.validate(page);

    // when
    await ICTSCNavbar.LoginLink(page).click();

    // then
    await LoginPage.waitFormSelector(page);
    expect(page.url().split("/").pop()).toBe("login");
    await LoginPage.validate(page);
  });
});

test.describe("ログイン状態", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/login");

    await page.fill("#username", "admin");
    await page.fill("#password", "password");
    await page.click("#loginBtn");

    // ログインに成功しましたというアラートが出るまで待つ
    await page.waitForURL("/");
    // await page.waitForSelector(".alert-success");
  });

  test("ルールページに遷移できる", async ({ page }) => {
    // setup
    await IndexPage.goto(page);
    await IndexPage.validate(page);

    // when
    await ICTSCNavbar.RuleLink(page).click();

    // then
    await IndexPage.waitFormSelector(page);
    expect(page.url().split("/").pop()).toBe("");
    await IndexPage.validate(page);
  });

  test("チーム情報ページに遷移できる", async ({ page }) => {
    // setup
    await IndexPage.goto(page);
    await IndexPage.validate(page);

    // when
    await ICTSCNavbar.TeamInfoLink(page).click();

    // then
    await TeamInfoPage.waitFormSelector(page);
    expect(page.url().split("/").pop()).toBe("team_info");
    await TeamInfoPage.validate(page);
  });

  test("問題ページに遷移できる", async ({ page }) => {
    // setup
    await IndexPage.goto(page);
    await IndexPage.validate(page);

    // when
    await ICTSCNavbar.ProblemsLink(page).click();

    // then
    await ProblemsPage.waitFormSelector(page);
    expect(page.url().split("/").pop()).toBe("problems");
    await ProblemsPage.validate(page);
  });

  test("順位ページに遷移できる", async ({ page }) => {
    // setup
    await IndexPage.goto(page);
    await IndexPage.validate(page);

    // when
    await ICTSCNavbar.RankingLink(page).click();

    // then
    await RankingPage.waitFormSelector(page);
    expect(page.url().split("/").pop()).toBe("ranking");
    await RankingPage.validate(page);
  });

  test("参加者ページに遷移できる", async ({ page }) => {
    // setup
    await IndexPage.goto(page);
    await IndexPage.validate(page);

    // when
    await ICTSCNavbar.UsersLink(page).click();

    // then
    await UsersPage.waitFormSelector(page);
    expect(page.url().split("/").pop()).toBe("users");
    await UsersPage.validate(page);
  });

  test("採点ページに遷移できる", async ({ page }) => {
    // setup
    await IndexPage.goto(page);
    await IndexPage.validate(page);

    // when
    await ICTSCNavbar.ScoringLink(page).click();

    // then
    await ScoringPage.waitFormSelector(page);
    expect(page.url().split("/").pop()).toBe("scoring");
    await ScoringPage.validate(page);
  });

  test("プロフィールページに遷移できる", async ({ page }) => {
    // setup
    await IndexPage.goto(page);
    await IndexPage.validate(page);
    await ICTSCNavbar.DropdownMenu(page).click();

    // when
    await ICTSCNavbar.ProfileLink(page).click();

    // then
    await ProfilePage.waitFormSelector(page);
    expect(page.url().split("/").pop()).toBe("profile");
    await ProfilePage.validate(page);
  });

  test("ログアウトできる", async ({ page }) => {
    // setup
    // 初期は problems ページに遷移しておく indexページに設定してしまうと、同じページに遷移するためテストがうまくできない
    await ProblemsPage.goto(page);
    await ProblemsPage.validate(page);
    await ICTSCNavbar.DropdownMenu(page).click();

    // when
    await ICTSCNavbar.LogoutButton(page).click();

    // then
    await IndexPage.waitFormSelector(page);
    expect(page.url().split("/").pop()).toBe("");
    await expect(ICTSCNavbar.LoginLink(page)).toBeVisible();
  });
});
