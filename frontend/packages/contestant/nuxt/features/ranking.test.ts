import { describe, expect, it } from "vitest";
import type { Rank } from "./models";
import { headerRankingItems } from "./ranking";

const entry = (rank: number, teamCode: number): Rank => ({
  rank,
  teamCode,
  teamName: `team-${teamCode}`,
  organization: "ICTSC",
  score: 1000 - rank,
});

const visibleRanks = (ranking: Rank[], teamCode: number) =>
  headerRankingItems(ranking, teamCode)
    .filter((item) => item.kind === "rank")
    .map((item) => item.entry.rank);

describe("headerRankingItems", () => {
  it("単独1位では自チームと2位を表示する", () => {
    const ranking = [entry(1, 10), entry(2, 20), entry(3, 30)];
    expect(visibleRanks(ranking, 10)).toEqual([1, 2]);
  });

  it("同率1位では同率表示を加えて2位を表示する", () => {
    const ranking = [entry(1, 10), entry(1, 11), entry(1, 12), entry(2, 20)];
    expect(visibleRanks(ranking, 11)).toEqual([1, 2]);
    expect(headerRankingItems(ranking, 11)).toContainEqual({
      kind: "tie",
      rank: 1,
      count: 3,
    });
  });

  it("単独最下位では1位とひとつ上の順位を表示する", () => {
    const ranking = [entry(1, 10), entry(2, 20), entry(3, 30), entry(4, 40)];
    expect(visibleRanks(ranking, 40)).toEqual([1, 3, 4]);
  });

  it("同率最下位では同率表示を加える", () => {
    const ranking = [entry(1, 10), entry(2, 20), entry(3, 30), entry(3, 31)];
    expect(visibleRanks(ranking, 31)).toEqual([1, 2, 3]);
    expect(headerRankingItems(ranking, 31)).toContainEqual({
      kind: "tie",
      rank: 3,
      count: 2,
    });
  });

  it("中間順位では1位と現在順位の前後を表示する", () => {
    const ranking = [
      entry(1, 10),
      entry(2, 20),
      entry(3, 30),
      entry(4, 40),
      entry(5, 50),
    ];
    expect(visibleRanks(ranking, 40)).toEqual([1, 3, 4, 5]);
  });

  it("中間順位が同率の場合は代表だけを表示して同率表示を加える", () => {
    const ranking = [
      entry(1, 10),
      entry(2, 20),
      entry(3, 30),
      entry(3, 31),
      entry(3, 32),
      entry(4, 40),
    ];
    expect(visibleRanks(ranking, 31)).toEqual([1, 2, 3, 4]);
    expect(headerRankingItems(ranking, 31)).toContainEqual({
      kind: "tie",
      rank: 3,
      count: 3,
    });
  });
});
