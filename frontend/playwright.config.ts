import { defineConfig, devices } from "@playwright/test";

const reuseExistingServer = process.env.CI !== "true";

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: true,
  forbidOnly: process.env.CI === "true",
  retries: process.env.CI === "true" ? 1 : 0,
  workers: 2,
  expect: { timeout: 15000 },
  reporter:
    process.env.CI === "true"
      ? [["github"], ["html", { open: "never" }]]
      : "list",
  use: {
    trace: "retain-on-failure",
    screenshot: "only-on-failure",
  },
  projects: [
    {
      name: "contestant-chromium",
      testMatch: "contestant.spec.ts",
      use: {
        ...devices["Desktop Chrome"],
        baseURL: process.env.E2E_STATIC_ORIGIN ?? "http://127.0.0.1:3000",
      },
    },
    {
      name: "admin-chromium",
      testMatch: "admin.spec.ts",
      use: {
        ...devices["Desktop Chrome"],
        baseURL: process.env.E2E_STATIC_ORIGIN ?? "http://127.0.0.1:3001",
      },
    },
  ],
  webServer: process.env.E2E_STATIC_ORIGIN
    ? []
    : [
        {
          command: "pnpm --filter @ictsc/competition dev --host 127.0.0.1",
          url: "http://127.0.0.1:3000",
          reuseExistingServer,
          timeout: 120_000,
          stdout: "pipe",
          stderr: "pipe",
        },
        {
          command: "pnpm --filter @ictsc/admin dev --host 127.0.0.1",
          url: "http://127.0.0.1:3001/admin/",
          reuseExistingServer,
          timeout: 120_000,
          stdout: "pipe",
          stderr: "pipe",
        },
      ],
});
