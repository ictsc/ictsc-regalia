import { describe, expect, it } from "vitest";
import { mount, RouterLinkStub } from "@vue/test-utils";
import ActivityView from "./activity/ActivityView.vue";
import NoticeListView from "./announce/NoticeListView.vue";
import WebPushSettings from "./announce/WebPushSettings.vue";
import ProfileView from "./profile/ProfileView.vue";
import RankingView from "./ranking/RankingView.vue";
import SignupView from "./signup/SignupView.vue";
import TeamDirectoryView from "./team/TeamDirectoryView.vue";
import type { Notice, Rank } from "~/features/models";

const global = {
  stubs: { NuxtLink: RouterLinkStub, MarkdownContent: true },
};

describe("NoticeListView", () => {
  const notices: Notice[] = [
    {
      slug: "a",
      title: "A",
      body: "",
      markdown: "a",
      effectiveFrom: "2026-09-01T00:00:00Z",
    },
    {
      slug: "b",
      title: "B",
      body: "",
      markdown: "b",
      effectiveFrom: "2026-09-02T00:00:00Z",
    },
  ];
  it("labels read state and reports read toggles", async () => {
    const wrapper = mount(NoticeListView, {
      props: { notices, readSlugs: ["b"] },
      slots: { settings: "<section class='settings-slot' />" },
      global,
    });
    expect(wrapper.find(".settings-slot").exists()).toBe(true);
    const buttons = wrapper.findAll(".notice-entry button");
    expect(buttons.map((b) => b.text())).toEqual(["既読にする", "未読に戻す"]);
    await buttons[1]!.trigger("click");
    expect(wrapper.emitted("toggleRead")).toEqual([["b"]]);
    await wrapper.find(".actions button").trigger("click");
    expect(wrapper.emitted("markAllRead")).toHaveLength(1);
    expect(
      wrapper.findAllComponents(RouterLinkStub).map((l) => l.props("to")),
    ).toEqual(["/announces/a", "/announces/b"]);
  });
  it("shows the empty state", () => {
    expect(
      mount(NoticeListView, {
        props: { notices: [], readSlugs: [] },
        global,
      }).text(),
    ).toContain("現在お知らせはありません。");
  });
});

describe("WebPushSettings", () => {
  const props = {
    subscribed: false,
    supported: true,
    serverEnabled: true,
    busy: false,
  };
  const button = (wrapper: ReturnType<typeof mount>) =>
    wrapper.find<HTMLButtonElement>("button");
  it("enables when supported and configured", async () => {
    const wrapper = mount(WebPushSettings, { props });
    expect(button(wrapper).text()).toBe("通知を有効にする");
    await button(wrapper).trigger("click");
    expect(wrapper.emitted("enable")).toHaveLength(1);
  });
  it("explains why notifications are unavailable", () => {
    const unsupported = mount(WebPushSettings, {
      props: { ...props, supported: false },
    });
    expect(unsupported.text()).toContain("Web Pushに対応していません");
    expect(button(unsupported).element.disabled).toBe(true);
    const disabled = mount(WebPushSettings, {
      props: { ...props, serverEnabled: false },
    });
    expect(disabled.text()).toContain("サーバー側のWeb Push設定が必要です");
    expect(button(disabled).element.disabled).toBe(true);
    const loading = mount(WebPushSettings, {
      props: { ...props, serverEnabled: undefined },
    });
    expect(loading.text()).toContain("ブラウザの通知で受け取れます");
    expect(button(loading).element.disabled).toBe(true);
  });
  it("disables an existing subscription and shows errors", async () => {
    const wrapper = mount(WebPushSettings, {
      props: { ...props, subscribed: true, error: new Error("失敗") },
    });
    expect(button(wrapper).text()).toBe("通知を停止する");
    await button(wrapper).trigger("click");
    expect(wrapper.emitted("disable")).toHaveLength(1);
    expect(wrapper.find('[role="alert"]').text()).toBe("Error: 失敗");
  });
});

describe("RankingView", () => {
  const ranking: Rank[] = [
    { rank: 1, teamCode: 2, teamName: "A", organization: "X", score: 2130 },
    { rank: 2, teamCode: 12, teamName: "B", organization: "Y", score: 1230 },
  ];
  it("highlights the own team and shows the frozen notice", () => {
    const wrapper = mount(RankingView, {
      props: {
        snapshot: { ranking, frozen: true, frozenAt: "2026-09-02T13:00Z" },
        ownTeamCode: 12,
      },
    });
    const rows = wrapper.findAll(".ranking-row:not(.ranking-heading)");
    expect(rows[1]!.classes()).toContain("is-team");
    expect(rows[1]!.find(".ranking-own-marker").text()).toBe("自チーム");
    expect(rows[0]!.find(".ranking-score").text()).toBe("2,130");
    expect(wrapper.text()).toContain("順位表は凍結中です 2026-09-02T13:00Z");
  });
  it("shows the empty state", () => {
    const wrapper = mount(RankingView, {
      props: { snapshot: { ranking: [], frozen: false } },
    });
    expect(wrapper.text()).toContain("順位はまだありません。");
  });
});

describe("ActivityView", () => {
  it("links each answer to the problem's answer history", () => {
    const wrapper = mount(ActivityView, {
      props: {
        entries: [
          {
            problemCode: "B02",
            problemTitle: "コンテナ",
            maxScore: 300,
            answerId: 3,
            submittedAt: "2026-09-02T03:00:00Z",
          },
        ],
      },
      global,
    });
    expect(wrapper.findComponent(RouterLinkStub).props("to")).toEqual({
      path: "/problems/B02",
      query: { answer: "3" },
      hash: "#answer-history",
    });
    expect(wrapper.find('[data-label="得点"]').text()).toBe("採点中 / 300");
    expect(
      mount(ActivityView, { props: { entries: [] }, global }).text(),
    ).toContain("提出はありません。");
  });
});

describe("TeamDirectoryView", () => {
  it("lists teams and members", () => {
    const wrapper = mount(TeamDirectoryView, {
      props: {
        teams: [
          {
            code: 12,
            name: "KERNEL PANIC",
            organization: "ICTSC University",
            memberLimit: 4,
            team: {
              code: 12,
              name: "KERNEL PANIC",
              organization: "ICTSC University",
              memberLimit: 4,
            },
            members: [
              { name: "alice", displayName: "Alice", selfIntroduction: "" },
            ],
          },
        ],
      },
      global,
    });
    expect(wrapper.find("h2").text()).toBe("KERNEL PANIC");
    expect(wrapper.find("h3").text()).toBe("Alice @alice");
  });
});

describe("ProfileView", () => {
  const props = {
    profile: { name: "alice", displayName: "Alice", selfIntroduction: "hi" },
    saving: false,
    message: "",
    signingOut: false,
    signOutError: "",
  };
  it("submits the edited profile", async () => {
    const wrapper = mount(ProfileView, { props });
    const inputs = wrapper.findAll("input");
    expect((inputs[0]!.element as HTMLInputElement).value).toBe("alice");
    await inputs[1]!.setValue("Alice B");
    await wrapper.find("textarea").setValue("よろしく");
    await wrapper.find("form").trigger("submit");
    expect(wrapper.emitted("save")).toEqual([
      [{ displayName: "Alice B", selfIntroduction: "よろしく" }],
    ]);
  });
  it("reports logout and shows its state", async () => {
    const wrapper = mount(ProfileView, {
      props: { ...props, signingOut: true, signOutError: "失敗" },
    });
    const logout = wrapper.find<HTMLButtonElement>(".profile-logout");
    expect(logout.text()).toBe("処理中…");
    expect(logout.element.disabled).toBe(true);
    expect(wrapper.find(".profile-account [role='alert']").text()).toBe("失敗");
    const idle = mount(ProfileView, { props });
    await idle.find(".profile-logout").trigger("click");
    expect(idle.emitted("logout")).toHaveLength(1);
  });
});

describe("SignupView", () => {
  const props = {
    initialName: "alice",
    initialDisplayName: "Alice",
    pending: false,
    message: "",
  };
  it("asks for an invitation code unless a team is assigned", async () => {
    const wrapper = mount(SignupView, { props });
    const inputs = wrapper.findAll("input");
    expect(inputs).toHaveLength(3);
    await inputs[0]!.setValue("INVITE");
    await wrapper.find("form").trigger("submit");
    expect(wrapper.emitted("submit")).toEqual([
      [{ invitationCode: "INVITE", name: "alice", displayName: "Alice" }],
    ]);
    const assigned = mount(SignupView, {
      props: { ...props, registrationTeamName: "KERNEL PANIC" },
    });
    expect(assigned.findAll("input")).toHaveLength(2);
    expect(assigned.text()).toContain("参加チーム: KERNEL PANIC");
  });
});
