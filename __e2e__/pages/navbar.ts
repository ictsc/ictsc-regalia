import { Page } from "@playwright/test";

const ICTSCNavbar = {
  RuleLink: (page: Page) => page.locator("a >> text=ルール"),
  TeamInfoLink: (page: Page) => page.locator("a >> text=チーム情報"),
  ProblemsLink: (page: Page) => page.locator("a >> text=問題"),
  RankingLink: (page: Page) => page.locator("a >> text=順位"),
  UsersLink: (page: Page) => page.locator("a >> text=参加者"),
  ScoringLink: (page: Page) => page.locator("a >> text=採点"),
  LoginLink: (page: Page) => page.locator("a >> text=ログイン"),
  DropdownMenu: (page: Page) => page.getByText("admin", { exact: true }),
  ProfileLink: (page: Page) => page.locator("a >> text=プロフィール"),
  LogoutButton: (page: Page) => page.locator("button >> text=ログアウト"),
};

export default ICTSCNavbar;
