import { describe, expect, it } from "vitest";
import { contestClock } from "./clock";

const entries = [
  {
    name: "morning",
    startAt: "2026-09-26T09:00:00+09:00",
    endAt: "2026-09-26T12:00:00+09:00",
  },
  {
    name: "afternoon",
    startAt: "2026-09-26T13:00:00+09:00",
    endAt: "2026-09-26T17:00:00+09:00",
  },
  {
    name: "tomorrow",
    startAt: "2026-09-27T09:00:00+09:00",
    endAt: "2026-09-27T17:00:00+09:00",
  },
];

describe("contestClock", () => {
  it("shows today's full competition time before it starts", () => {
    expect(
      contestClock(entries, Date.parse("2026-09-26T08:00:00+09:00")),
    ).toEqual({
      label: "本日の競技時間",
      text: "07:00:00",
    });
  });

  it("counts only playable time remaining across today's break", () => {
    expect(
      contestClock(entries, Date.parse("2026-09-26T10:00:00+09:00")),
    ).toEqual({
      label: "本日の残り時間",
      text: "06:00:00",
    });
    expect(
      contestClock(entries, Date.parse("2026-09-26T12:30:00+09:00")),
    ).toEqual({
      label: "本日の残り時間",
      text: "04:00:00",
    });
  });

  it("shows today's end instead of tomorrow's countdown after the contest", () => {
    expect(
      contestClock(entries, Date.parse("2026-09-26T17:30:00+09:00")),
    ).toEqual({
      label: "本日終了",
      text: "00:00:00",
    });
  });

  it("uses the next start when there is no contest section today", () => {
    expect(
      contestClock(entries, Date.parse("2026-09-25T23:00:00+09:00")),
    ).toEqual({
      label: "次の競技開始まで",
      text: "10:00:00",
    });
  });
});
