import { describe, expect, it } from "vitest";
import { sectionLabel } from "./section";

describe("sectionLabel", () => {
  it("labels contest days and shared sections", () => {
    expect(sectionLabel("day1")).toBe("1日目");
    expect(sectionLabel("day12")).toBe("12日目");
    expect(sectionLabel("both")).toBe("両日");
  });
  it("keeps unknown slugs", () => {
    expect(sectionLabel("extra")).toBe("extra");
    expect(sectionLabel(undefined)).toBeUndefined();
  });
});
