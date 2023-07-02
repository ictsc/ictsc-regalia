import { test } from "@playwright/test";

import IndexPage from "./pages/Index";

test("画面項目が表示されること", async ({ page }) => {
  await IndexPage.goto(page);
  await IndexPage.validate(page);
});
