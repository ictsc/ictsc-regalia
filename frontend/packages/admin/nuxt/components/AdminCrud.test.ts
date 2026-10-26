import { beforeAll, describe, expect, it } from "vitest";
import { mount, RouterLinkStub } from "@vue/test-utils";
import AdminConfirmationDialog from "./confirmation/AdminConfirmationDialog.vue";
import ContentView from "./content/ContentView.vue";
import ProblemContentView from "./content/ProblemContentView.vue";
import ContestantsView from "./contestant/ContestantsView.vue";
import DeploymentsView from "./deployment/DeploymentsView.vue";
import ScoresView from "./score/ScoresView.vue";
import SettingsView from "./settings/SettingsView.vue";
import TeamDetailView from "./team/TeamDetailView.vue";
import TeamListView from "./team/TeamListView.vue";
import {
  adminProblem,
  catalog,
  contentStatus,
  contestants,
  deployments,
  invitations,
  ranking,
  scores,
  teams,
} from "~/stories/fixtures";

const global = { stubs: { NuxtLink: RouterLinkStub, MarkdownContent: true } };
const base = { busy: false, message: "" };

beforeAll(() => {
  // jsdom does not implement modal dialogs.
  HTMLDialogElement.prototype.showModal ??= function (this: HTMLDialogElement) {
    this.open = true;
  };
  HTMLDialogElement.prototype.close ??= function (this: HTMLDialogElement) {
    this.open = false;
  };
});

describe("TeamListView", () => {
  it("lists teams and submits a new team", async () => {
    const wrapper = mount(TeamListView, { props: { ...base, teams }, global });
    expect(wrapper.find(".admin-panel-heading span").text()).toBe("6 チーム");
    expect(wrapper.findAllComponents(RouterLinkStub)[0]!.props("to")).toBe(
      "/teams/2",
    );
    const inputs = wrapper.findAll("form input");
    await inputs[0]!.setValue("12");
    await inputs[1]!.setValue("NEW TEAM");
    await inputs[2]!.setValue("ICTSC");
    await wrapper.find("form").trigger("submit");
    expect(wrapper.emitted("create")?.[0]?.[0]).toEqual({
      code: 12,
      name: "NEW TEAM",
      organization: "ICTSC",
      member_limit: 4,
      color: "#A6E35F",
    });
  });
});

describe("TeamDetailView", () => {
  const props = {
    ...base,
    code: 4,
    team: teams[2]!,
    members: contestants.filter((c) => c.team.code === 4),
    invitations,
  };
  it("edits the team, invites and requests removal", async () => {
    const wrapper = mount(TeamDetailView, { props, global });
    expect(wrapper.find(".team-badge").text()).toBe("KERNEL PANIC");
    await wrapper.find('input[maxlength="255"]').setValue("RENAMED");
    expect(wrapper.find(".team-badge").text()).toBe("RENAMED");
    const forms = wrapper.findAll("form");
    await forms[0]!.trigger("submit");
    expect(wrapper.emitted("save")?.[0]?.[0]).toMatchObject({
      name: "RENAMED",
      member_limit: 4,
    });
    await wrapper
      .find('input[type="datetime-local"]')
      .setValue("2026-09-03T10:00");
    await forms[1]!.trigger("submit");
    expect(wrapper.emitted("invite")).toEqual([["2026-09-03T10:00"]]);
    await wrapper.find(".admin-danger button").trigger("click");
    expect(wrapper.emitted("remove")).toHaveLength(1);
    expect(wrapper.findAll("tbody tr")).toHaveLength(2);
  });
});

describe("ContestantsView", () => {
  it("filters contestants and requests impersonation", async () => {
    const wrapper = mount(ContestantsView, {
      props: { ...base, contestants },
      global,
    });
    expect(wrapper.find(".muted").text()).toBe("6 / 6 人を表示");
    await wrapper.find("input").setValue("carol");
    expect(wrapper.find(".muted").text()).toBe("1 / 6 人を表示");
    await wrapper.find("tbody button").trigger("click");
    expect(wrapper.emitted("impersonate")).toEqual([["carol"]]);
    await wrapper.find("input").setValue("nobody");
    expect(wrapper.text()).toContain("該当する参加者はいません。");
  });
});

describe("ScoresView", () => {
  it("shows ranking state and reports actions", async () => {
    const wrapper = mount(ScoresView, {
      props: {
        ...base,
        ranking: { ...ranking, frozen: true, frozen_at: "2026-09-02T04:00Z" },
        scores,
      },
      global,
    });
    expect(wrapper.find(".admin-badge").text()).toBe("凍結中");
    expect(wrapper.text()).toContain("凍結時刻:");
    const buttons = wrapper.findAll(".actions button");
    await buttons[0]!.trigger("click");
    await buttons[1]!.trigger("click");
    expect(wrapper.emitted("recalculate")).toHaveLength(1);
    expect(wrapper.emitted("reveal")).toHaveLength(1);
  });
});

describe("SettingsView", () => {
  it("saves the rule and shows the freeze from schedule.yaml", async () => {
    const wrapper = mount(SettingsView, {
      props: { ...base, initialMarkdown: "# rule", initialFreeze: "" },
      global,
    });
    await wrapper.find("textarea").setValue("# new rule");
    const forms = wrapper.findAll("form");
    expect(forms).toHaveLength(1);
    await forms[0]!.trigger("submit");
    expect(wrapper.emitted("saveRule")).toEqual([["# new rule"]]);
    expect(wrapper.text()).toContain("現在の凍結時刻: 未設定");
  });
});

describe("DeploymentsView", () => {
  it("filters history, shows events and requests deployments", async () => {
    const wrapper = mount(DeploymentsView, {
      props: {
        ...base,
        deployments,
        teams,
        problems: catalog,
        streamError: true,
      },
      global,
    });
    expect(wrapper.find('[role="alert"]').text()).toBe(
      "状態更新に再接続しています。",
    );
    const submit = wrapper.find<HTMLButtonElement>("form button");
    expect(submit.element.disabled).toBe(true);
    const selects = wrapper.findAll("select");
    await selects[0]!.setValue("4");
    await selects[1]!.setValue("B02");
    expect(wrapper.find(".admin-panel-heading span").text()).toBe("2 件");
    await wrapper.find("form").trigger("submit");
    expect(wrapper.emitted("create")).toEqual([[4, "B02"]]);

    const row = wrapper.findAll("tbody tr")[0]!;
    expect(row.find(".admin-badge").text()).toBe("展開中");
    await row.findAll("button")[1]!.trigger("click");
    expect(wrapper.emitted("sync")).toEqual([[4, "B02"]]);
    await row.findAll("button")[0]!.trigger("click");
    expect(wrapper.findAll("article.content-section")).toHaveLength(2);
  });
});

describe("ContentView", () => {
  it("shows status and per-list errors, and refreshes a commit", async () => {
    const wrapper = mount(ContentView, {
      props: {
        ...base,
        status: { ...contentStatus, last_error: "broken" },
        sections: { data: { sections: [] } },
        problems: { error: "failed" },
        announcements: { data: { announcements: [] } },
      },
      global,
    });
    expect(wrapper.find('[role="alert"]').text()).toBe("broken");
    expect(wrapper.text()).toContain("セクションはありません。");
    expect(wrapper.text()).toContain("failed");
    await wrapper.find("form input").setValue("  abc123  ");
    await wrapper.find("form").trigger("submit");
    expect(wrapper.emitted("refreshContent")).toEqual([["abc123"]]);
  });
});

describe("ProblemContentView", () => {
  it("shows the problem and asks to view another commit", async () => {
    const wrapper = mount(ProblemContentView, {
      props: { problem: adminProblem, initialCommit: "" },
      global,
    });
    expect(wrapper.find("h1").text()).toBe("B02: コンテナが起動を繰り返す");
    expect(wrapper.text()).toContain("デフォルト採点担当者: 未設定");
    await wrapper.find("input").setValue("1a2b3c4");
    await wrapper.find("form").trigger("submit");
    expect(wrapper.emitted("viewCommit")).toEqual([["1a2b3c4"]]);
  });
});

describe("AdminConfirmationDialog", () => {
  it("opens for a request and settles with the choice", async () => {
    const wrapper = mount(AdminConfirmationDialog, {
      props: { request: null },
      attachTo: document.body,
    });
    const dialog = wrapper.find<HTMLDialogElement>("dialog").element;
    await wrapper.setProps({
      request: { title: "削除", message: "削除しますか？", destructive: true },
    });
    await new Promise((resolve) => setTimeout(resolve));
    expect(dialog.open).toBe(true);
    expect(wrapper.find(".confirmation-dialog-eyebrow").text()).toBe(
      "取り消せない操作",
    );
    expect(wrapper.find(".button-primary").text()).toBe("実行する");
    await wrapper.find(".button-primary").trigger("click");
    await wrapper.find(".button-secondary").trigger("click");
    expect(wrapper.emitted("settle")).toEqual([[true], [false]]);
    wrapper.unmount();
  });
});
