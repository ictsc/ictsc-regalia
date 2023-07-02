import { Page } from "@playwright/test";

type BasePage = {
  goto(page: Page): Promise<void>;

  validate(page: Page): Promise<void>;

  waitFormSelector(page: Page): Promise<void>;
};

export default BasePage;
