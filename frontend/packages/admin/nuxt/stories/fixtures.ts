// Fake data for admin Storybook stories. The Storybook is published, so use
// only made-up names.
import type { components } from "@ictsc/api";

type Schemas = components["schemas"];

export const fixedNow = Date.parse("2026-09-02T12:00:00+09:00");
export const adminName = "運営スタッフ";

export const contentStatus: Schemas["ContentStatus"] = {
  state: "SUCCEEDED",
  active_commit: "9f2c4e1a7b3d5c6e8f0a1b2c3d4e5f6a7b8c9d0e",
  activated_at: "2026-09-02T08:30:00+09:00",
  source_ref: "main",
  last_attempt_commit: "9f2c4e1a7b3d5c6e8f0a1b2c3d4e5f6a7b8c9d0e",
  last_attempt_at: "2026-09-02T08:30:00+09:00",
  last_error: null,
  serving_last_known_good: false,
};
