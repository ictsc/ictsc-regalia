import { describe, expect, it } from "vitest";
import { mount, RouterLinkStub } from "@vue/test-utils";
import CompetitionHeaderView from "./CompetitionHeaderView.vue";
import { headerRankingItems } from "~/features/ranking";
import type { Rank } from "~/features/models";

const ranking: Rank[] = [
  { rank: 1, teamCode: 2, teamName: "A", organization: "", score: 2000 },
  { rank: 2, teamCode: 12, teamName: "B", organization: "", score: 1500 },
  { rank: 3, teamCode: 13, teamName: "C", organization: "", score: 1000 },
];

function mountHeader(
  props: Partial<InstanceType<typeof CompetitionHeaderView>["$props"]> = {},
) {
  return mount(CompetitionHeaderView, {
    props: {
      compact: false,
      rankingItems: headerRankingItems(ranking, 12),
      ownTeamCode: 12,
      maxScore: 3000,
      clock: { label: "本日の残り時間", text: "02:13:48" },
      demoMode: false,
      hasUnreadNotices: false,
      teamName: "KERNEL PANIC",
      routeKey: "/problems",
      ...props,
    },
    global: { stubs: { NuxtLink: RouterLinkStub, ThemeToggle: true } },
  });
}

describe("CompetitionHeaderView", () => {
  it("shows neighboring ranks and the own team's score", () => {
    const wrapper = mountHeader();
    const ranks = wrapper.findAll(".header-neighbors > span");
    expect(ranks.map((r) => r.find("small").text())).toEqual([
      "1位",
      "2位",
      "3位",
    ]);
    const own = wrapper.find(".is-team");
    expect(own.text()).toContain("1,500");
    expect(own.text()).toContain("/ 3,000");
    expect(wrapper.find(".competition-clock").text()).toBe(
      "02:13:48本日の残り時間",
    );
  });

  it("leaves out the ranking but keeps the menu when the ranking is hidden", () => {
    const wrapper = mountHeader({ rankingHidden: true });
    expect(wrapper.find(".header-ranking").exists()).toBe(false);
    expect(wrapper.find(".header-account").exists()).toBe(true);
  });

  it("uses the simplified navigation in compact mode", () => {
    const wrapper = mountHeader({ compact: true });
    expect(wrapper.classes()).toContain("simple-header-version");
    expect(wrapper.find(".header-performance").exists()).toBe(false);
    expect(wrapper.find(".simple-header-nav").exists()).toBe(true);
  });

  it("marks unread notices and shows the team in the mobile menu", () => {
    const wrapper = mountHeader({ hasUnreadNotices: true });
    expect(wrapper.findAll(".notification-dot").length).toBeGreaterThan(0);
    expect(wrapper.find(".mobile-team-badge").text()).toBe("KERNEL PANIC");
    expect(
      mountHeader({ teamName: undefined }).find(".mobile-team-badge").exists(),
    ).toBe(false);
  });

  it("closes the mobile menu when the route changes", async () => {
    const wrapper = mountHeader();
    const menu = wrapper.find<HTMLDetailsElement>(".mobile-menu").element;
    menu.open = true;
    await wrapper.setProps({ routeKey: "/ranking" });
    expect(menu.open).toBe(false);
  });

  it("links the clock to the schedule", () => {
    const clock = mountHeader()
      .findAllComponents(RouterLinkStub)
      .find((link) => link.classes("competition-clock"));
    expect(clock?.props("to")).toBe("/schedule");
    expect(clock?.text()).toContain("02:13:48");
  });
});
