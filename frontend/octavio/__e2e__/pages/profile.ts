import { expect, Page } from "@playwright/test";

const ProfilePage = {
  goto: async (page: Page) => {
    await page.goto("/profile");
  },
  validate: async (page: Page) => {
    await expect(page.locator(".title-ictsc")).toHaveText("プロフィール");
  },
  waitFormSelector: async (page: Page) => {
    await page.waitForSelector(".title-ictsc >> text=プロフィール");
  },
  DisplayNameInput: (page: Page) => page.getByRole("textbox").nth(0),
  Teams: (page: Page) => page.locator(".teams"),
  SelfIntroductionInput: (page: Page) => page.getByRole("textbox").nth(1),
  GitHubIDInput: (page: Page) => page.getByRole("textbox").nth(2),
  TwitterIDInput: (page: Page) => page.getByRole("textbox").nth(3),
  FacebookIDInput: (page: Page) => page.getByRole("textbox").nth(4),
  SubmitButton: (page: Page) => page.getByRole("button", { name: "更新" }),
};

export default ProfilePage;
