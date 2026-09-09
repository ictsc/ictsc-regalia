import { describe, expect, it } from "vitest";
import { toCamelCase, toSnakeCase } from "./mapper";

describe("case mappers", () => {
  it("maps nested API objects in both directions", () => {
    const wire = { team_code: "A01", nested_items: [{ max_score: 100 }] };
    const view = toCamelCase(wire);
    expect(view).toEqual({
      teamCode: "A01",
      nestedItems: [{ maxScore: 100 }],
    });
    expect(toSnakeCase(view)).toEqual(wire);
  });
});
