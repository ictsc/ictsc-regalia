import { describe, expect, it } from "vitest";
import { mount, RouterLinkStub } from "@vue/test-utils";
import AdminHeader from "./layout/AdminHeader.vue";
import AdminNav from "./layout/AdminNav.vue";
import OverviewView from "./overview/OverviewView.vue";

const global = { stubs: { NuxtLink: RouterLinkStub, ThemeToggle: true } };

describe("AdminHeader", () => {
  it("shows the account and reports logout", async () => {
    const wrapper = mount(AdminHeader, {
      props: { signedIn: true, displayName: "運営スタッフ" },
      global,
    });
    expect(wrapper.find(".actions > span").text()).toBe("運営スタッフ");
    await wrapper.find(".button-secondary").trigger("click");
    expect(wrapper.emitted("logout")).toHaveLength(1);
  });
  it("hides account controls when anonymous", () => {
    const wrapper = mount(AdminHeader, { props: { signedIn: false }, global });
    expect(wrapper.find(".button-secondary").exists()).toBe(false);
  });
});

describe("AdminNav", () => {
  const current = (path: string) =>
    mount(AdminNav, { props: { currentPath: path }, global })
      .findAllComponents(RouterLinkStub)
      .filter((link) => link.attributes("aria-current") === "page")
      .map((link) => link.text());
  it("marks the current section, including nested pages", () => {
    expect(current("/")).toEqual(["概要"]);
    expect(current("/teams")).toEqual(["チーム"]);
    expect(current("/teams/12")).toEqual(["チーム"]);
    expect(current("/submissions/A01/2/3")).toEqual(["採点"]);
    expect(current("/unknown")).toEqual([]);
  });
});

describe("OverviewView", () => {
  it("shows counts, failed cards and content status", () => {
    const wrapper = mount(OverviewView, {
      props: {
        cards: [
          { label: "チーム", detail: "管理", to: "/teams", count: 48 },
          { label: "回答", detail: "採点", to: "/submissions", error: "x" },
        ],
        content: {
          state: "SUCCEEDED",
          active_commit: "abc123",
          activated_at: null,
          source_ref: null,
          last_attempt_commit: null,
          last_attempt_at: null,
          last_error: null,
          serving_last_known_good: true,
        },
      },
      global,
    });
    const cards = wrapper.findAll(".admin-summary-card");
    expect(cards[0]!.find(".admin-summary-number").text()).toBe("48");
    expect(cards[1]!.find(".admin-summary-number").text()).toBe("—");
    expect(cards[1]!.find(".admin-summary-footer").text()).toBe(
      "取得できませんでした →",
    );
    const badges = wrapper.findAll(".admin-meta .admin-badge");
    expect(badges.map((b) => b.text())).toEqual([
      "配信中",
      "最終正常版を配信中",
    ]);
    expect(badges[0]!.classes()).toContain("is-active");
    expect(wrapper.find(".admin-code").text()).toBe("有効な commit: abc123");
  });
  it("falls back when the content status is unknown", () => {
    const wrapper = mount(OverviewView, { props: { cards: [] }, global });
    expect(wrapper.find(".admin-badge").text()).toBe("状態不明");
    expect(wrapper.find(".admin-code").text()).toBe("有効な commit: 未設定");
  });
});
