import { expect, Page } from "@playwright/test";

const TeamInfoPage = {
  goto: async (page: Page) => {
    await page.goto("/team_info");
  },
  validate: async (page: Page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("チーム情報");
  },
  waitFormSelector: async (page: Page) => {
    await page.waitForSelector(".title-ictsc >> text=チーム情報");
  },
  GroupNameAndOrganization: (page: Page) =>
    page.locator(".group-name-and-organization"),
  Ranking: (page: Page) => page.locator(".ranking"),
  Score: (page: Page) => page.locator(".score"),
  SSHInfo: (page: Page) => page.locator(".ssh-info"),
  PasswordInfo: (page: Page) => page.locator(".password-info"),
};

export default TeamInfoPage;
