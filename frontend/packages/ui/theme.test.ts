import { describe, expect, it } from "vitest";
import { nextThemePreference } from "./theme";

describe("nextThemePreference", () => {
  it("switches from the system light theme to dark", () => {
    expect(nextThemePreference("system", false)).toBe("dark");
  });

  it("switches from the system dark theme to light", () => {
    expect(nextThemePreference("system", true)).toBe("light");
  });

  it("returns a manual theme to the system setting", () => {
    expect(nextThemePreference("dark", false)).toBe("system");
    expect(nextThemePreference("light", true)).toBe("system");
  });
});
