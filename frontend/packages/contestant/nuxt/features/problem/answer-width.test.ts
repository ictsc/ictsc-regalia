import { describe, expect, it } from "vitest";
import {
  clampAnswerWidth,
  loadAnswerWidth,
  saveAnswerWidth,
} from "./answer-width";

describe("answer panel width", () => {
  it("keeps room for the problem text and a usable answer panel", () => {
    expect(clampAnswerWidth(200, 1200)).toBe(360);
    expect(clampAnswerWidth(700, 1200)).toBe(700);
    expect(clampAnswerWidth(1100, 1200)).toBe(880);
    expect(clampAnswerWidth(500, 500)).toBe(360);
  });

  it("stores the chosen width and forgets a reset", () => {
    const storage = new Map<string, string>();
    const fake = {
      getItem: (key: string) => storage.get(key) ?? null,
      setItem: (key: string, value: string) => storage.set(key, value),
      removeItem: (key: string) => storage.delete(key),
    } as unknown as Storage;
    expect(loadAnswerWidth(fake)).toBeUndefined();
    saveAnswerWidth(fake, 640);
    expect(loadAnswerWidth(fake)).toBe(640);
    saveAnswerWidth(fake, undefined);
    expect(loadAnswerWidth(fake)).toBeUndefined();
  });
});
