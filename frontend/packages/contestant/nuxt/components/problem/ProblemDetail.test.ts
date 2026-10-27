import { beforeAll, describe, expect, it } from "vitest";
import { mount, RouterLinkStub } from "@vue/test-utils";
import AnswerForm from "./AnswerForm.vue";
import AnswerHistory from "./AnswerHistory.vue";
import DeploymentHistory from "./DeploymentHistory.vue";
import ProblemHero from "./ProblemHero.vue";
import ProblemRail from "./ProblemRail.vue";
import ResetDialog from "./ResetDialog.vue";
import type { Problem } from "~/features/models";

const global = {
  stubs: { NuxtLink: RouterLinkStub, MarkdownContent: true },
};

beforeAll(() => {
  window.matchMedia ??= (query: string) =>
    ({ matches: false, media: query }) as MediaQueryList;
  // jsdom does not implement modal dialogs.
  HTMLDialogElement.prototype.showModal ??= function (this: HTMLDialogElement) {
    this.open = true;
  };
  HTMLDialogElement.prototype.close ??= function (this: HTMLDialogElement) {
    this.open = false;
    this.dispatchEvent(new Event("close"));
  };
});

const section = {
  name: "day2",
  startAt: "2026-09-02T00:00:00+09:00",
  endAt: "2026-09-03T00:00:00+09:00",
};
const problem = (code: string, score?: number): Problem => ({
  code,
  title: `問題 ${code}`,
  category: "ネットワーク",
  maxScore: 100,
  submissionableSchedules: [section],
  submissionStatus: { isSubmittable: true },
  score:
    score == null
      ? undefined
      : { markedScore: score, penalty: 0, score, maxScore: 100 },
});

describe("ProblemRail", () => {
  it("marks the current problem and shows cooldown or grading state", () => {
    const wrapper = mount(ProblemRail, {
      props: {
        currentCode: "A02",
        items: [
          { problem: problem("A01", 40), cooldownMinutes: 0, pending: false },
          { problem: problem("A02"), cooldownMinutes: 5, pending: false },
          { problem: problem("A03"), cooldownMinutes: 0, pending: true },
        ],
      },
      global,
    });
    const links = wrapper.findAllComponents(RouterLinkStub).slice(1);
    expect(links.map((link) => link.props("to"))).toEqual([
      "/problems/A01",
      "/problems/A02",
      "/problems/A03",
    ]);
    expect(links[1]!.attributes("aria-current")).toBe("page");
    expect(links[1]!.classes()).toEqual(
      expect.arrayContaining(["is-active", "has-status"]),
    );
    expect(links[1]!.find(".rail-cooldown").text()).toBe("あと5分");
    expect(links[1]!.attributes("aria-label")).toContain("再提出可能まで5分");
    expect(links[2]!.find(".rail-pending").text()).toBe("採点中");
    expect(links[0]!.attributes("aria-label")).toContain("40点 / 100点満点");
    expect(wrapper.find(".rail-stage-label").text()).toBe("2日目");
  });
});

describe("ProblemHero", () => {
  const props = {
    code: "B02",
    title: "コンテナ",
    category: "サーバー",
    maxScore: 300,
    cooldownMinutes: 0,
    open: true,
  };
  it("puts category, score and status on the meta line", () => {
    const meta = mount(ProblemHero, { props }).find(".problem-meta");
    expect(meta.findAll(":scope > span").map((part) => part.text())).toEqual([
      "サーバー",
      "—/ 300点満点",
      "回答受付中",
    ]);
  });
  it("shows whether answers are accepted", () => {
    expect(mount(ProblemHero, { props }).find(".problem-cooldown").text()).toBe(
      "回答受付中",
    );
    const closed = mount(ProblemHero, { props: { ...props, open: false } });
    expect(closed.find(".problem-cooldown").text()).toBe("回答終了");
    expect(closed.find(".problem-cooldown").classes()).toContain(
      "is-answer-closed",
    );
  });
  it("prefers the cooldown over the open state", () => {
    const wrapper = mount(ProblemHero, {
      props: { ...props, cooldownMinutes: 18, open: false },
    });
    expect(wrapper.find(".problem-cooldown").text()).toBe("再回答可能まで18分");
    expect(wrapper.find(".problem-cooldown").classes()).not.toContain(
      "is-answer-closed",
    );
  });
});

describe("AnswerForm", () => {
  const props = {
    body: "原因",
    open: true,
    sending: false,
    submitted: false,
    cooldownSeconds: 0,
    savedAt: "",
    saveError: "",
    maxLength: 10,
    allowAttachments: false,
  };
  const submitButton = (wrapper: ReturnType<typeof mount>) =>
    wrapper.find<HTMLButtonElement>(".button-submit").element;

  it("submits only after confirming the preview", async () => {
    const wrapper = mount(AnswerForm, {
      props,
      global: { stubs: { MarkdownContent: true } },
    });
    expect(submitButton(wrapper).disabled).toBe(false);
    await wrapper.find("form").trigger("submit");
    expect(wrapper.emitted("submit")).toBeUndefined();
    const dialog = wrapper.find(".answer-confirm-dialog");
    expect(dialog.text()).toContain("この内容で提出しますか？");
    await dialog.find(".button-primary").trigger("click");
    expect(wrapper.emitted("submit")).toHaveLength(1);
  });

  it("counts characters without line breaks and blocks over-long answers", async () => {
    const wrapper = mount(AnswerForm, {
      props: { ...props, body: "1234567890\n\n" },
    });
    expect(wrapper.find(".answer-length").text()).toContain("10 / 10 文字");
    expect(submitButton(wrapper).disabled).toBe(false);
    await wrapper.setProps({ body: "12345678901" });
    expect(wrapper.find(".answer-length").text()).toContain("11 / 10 文字");
    expect(wrapper.find(".answer-length").classes()).toContain("is-over");
    expect(submitButton(wrapper).disabled).toBe(true);
  });

  it("hints that Escape leaves the field after Tab indents", async () => {
    const wrapper = mount(AnswerForm, { props, attachTo: document.body });
    expect(wrapper.find(".tab-hint").exists()).toBe(false);
    await wrapper.find("textarea").trigger("keydown", { key: "Tab" });
    expect(wrapper.find(".tab-hint").text()).toContain("Esc");
    wrapper.unmount();
  });

  it("switches between the editor and the Markdown preview", async () => {
    const wrapper = mount(AnswerForm, {
      props,
      global: { stubs: { MarkdownContent: true } },
    });
    expect(wrapper.find(".answer-preview").exists()).toBe(false);
    const toggle = wrapper.find(".answer-preview-switch");
    expect(toggle.attributes("aria-checked")).toBe("false");
    await toggle.trigger("click");
    expect(toggle.attributes("aria-checked")).toBe("true");
    expect(wrapper.find(".answer-preview").exists()).toBe(true);
    expect(wrapper.find("textarea").isVisible()).toBe(false);
  });

  it("creates a log attachment from pasted text", async () => {
    const wrapper = mount(AnswerForm, {
      props: { ...props, allowAttachments: true, attachments: [] },
    });
    expect(wrapper.find(".log-composer").exists()).toBe(false);
    await wrapper.find(".answer-attach-button").trigger("click");
    const composer = wrapper.find(".log-composer");
    expect(composer.find("input").element.value).toBe("log-1.txt");
    await composer.find("textarea").setValue("Oct 26 kernel: oops\n");
    await composer.findAll("button").at(-1)!.trigger("click");
    expect(wrapper.emitted("update:attachments")?.at(-1)).toEqual([
      [{ name: "log-1.txt", size: 20, content: "Oct 26 kernel: oops\n" }],
    ]);
    expect(wrapper.find(".log-composer").exists()).toBe(false);
  });

  it("rejects a pasted log that reuses an attached name", async () => {
    const wrapper = mount(AnswerForm, {
      props: {
        ...props,
        allowAttachments: true,
        attachments: [{ name: "app.log", size: 1, content: "x" }],
      },
    });
    await wrapper.find(".answer-attach-button").trigger("click");
    const composer = wrapper.find(".log-composer");
    await composer.find("input").setValue("app.log");
    await composer.find("textarea").setValue("y");
    await composer.findAll("button").at(-1)!.trigger("click");
    expect(composer.find("[role='alert']").text()).toContain("添付済み");
    expect(wrapper.emitted("update:attachments")).toBeUndefined();
  });

  it("offers log attachments only when the problem allows them", () => {
    expect(
      mount(AnswerForm, { props }).find(".answer-attachments").exists(),
    ).toBe(false);
    expect(
      mount(AnswerForm, { props: { ...props, allowAttachments: true } })
        .find(".answer-attachments")
        .exists(),
    ).toBe(true);
  });

  it("disables submission while closed, cooling down, sending or empty", () => {
    for (const change of [
      { open: false },
      { cooldownSeconds: 30 },
      { sending: true },
      { body: "  " },
    ]) {
      expect(
        submitButton(mount(AnswerForm, { props: { ...props, ...change } }))
          .disabled,
      ).toBe(true);
    }
  });

  it("updates the body and shows status messages", async () => {
    const wrapper = mount(AnswerForm, {
      props: {
        ...props,
        submitted: true,
        cooldownSeconds: 754,
        actionError: new Error("失敗しました"),
      },
    });
    await wrapper.find("textarea").setValue("新しい回答");
    expect(wrapper.emitted("update:body")?.at(-1)).toEqual(["新しい回答"]);
    expect(wrapper.text()).toContain("回答を提出しました");
    expect(wrapper.text()).toContain("再回答可能まで 12分34秒");
    expect(wrapper.find('[role="alert"]').text()).toBe("失敗しました");
  });

  it("shows draft save state", () => {
    const saved = mount(AnswerForm, {
      props: { ...props, savedAt: "2026-09-02T03:00:00.000Z" },
    });
    expect(saved.find(".draft-save-status time").exists()).toBe(true);
    const failed = mount(AnswerForm, {
      props: { ...props, saveError: "保存できません" },
    });
    expect(failed.find(".answer-heading-row [role='alert']").text()).toBe(
      "保存できません",
    );
  });
});

describe("DeploymentHistory", () => {
  it("lists deployments and connection state", () => {
    const empty = mount(DeploymentHistory, {
      props: {
        deployments: [],
        penaltyThreshold: 2,
        deploying: true,
        streamError: true,
        resetDisabled: false,
      },
    });
    expect(empty.text()).toContain("減点なしの上限: 2回");
    expect(empty.text()).toContain("QUEUED — 再展開を要求中");
    expect(empty.text()).toContain("更新の接続を確認しています");
    expect(empty.text()).toContain("再展開履歴はありません。");

    const listed = mount(DeploymentHistory, {
      props: {
        deployments: [
          {
            revision: 2,
            status: "DEPLOYING",
            requestedAt: "2026-09-02T03:00:00Z",
            allowedDeploymentCount: 1,
            thresholdExceeded: false,
            penalty: 0,
          },
        ],
        penaltyThreshold: 2,
        deploying: false,
        streamError: false,
        resetDisabled: false,
      },
    });
    expect(listed.find("li").text()).toBe("#2 DEPLOYING / 減点 0 / 残り 1回");
  });
});

describe("AnswerHistory", () => {
  const answers = [
    { id: 2, submittedAt: "2026-09-02T03:00:00Z", contentCommit: "a" },
    {
      id: 1,
      submittedAt: "2026-09-02T02:00:00Z",
      contentCommit: "a",
      score: { markedScore: 50, penalty: 0, score: 50, maxScore: 100 },
    },
  ];
  it("links each answer and marks the selected one", () => {
    const wrapper = mount(AnswerHistory, {
      props: {
        problemCode: "B02",
        answers,
        selectedAnswerNumber: 1,
        selectedAnswerBody: "回答本文",
      },
      global,
    });
    const links = wrapper.findAllComponents(RouterLinkStub);
    expect(links[1]!.props("to")).toEqual({
      path: "/problems/B02",
      query: { answer: "1" },
      hash: "#answer-history",
    });
    expect(links[1]!.attributes("aria-current")).toBe("true");
    expect(links[0]!.attributes("aria-current")).toBeUndefined();
    expect(wrapper.findAll("tbody tr")[0]!.text()).toContain("採点中");
    expect(wrapper.find("[aria-live]").exists()).toBe(true);
  });
  it("shows the empty state", () => {
    const wrapper = mount(AnswerHistory, {
      props: { problemCode: "B02", answers: [], selectedAnswerNumber: null },
      global,
    });
    expect(wrapper.text()).toContain("提出はありません。");
    expect(wrapper.find("[aria-live]").exists()).toBe(false);
  });
});

describe("DeploymentHistory reset", () => {
  const props = {
    deployments: [],
    penaltyThreshold: 2,
    deploying: false,
    streamError: false,
  };
  it("requests a reset", async () => {
    const wrapper = mount(DeploymentHistory, {
      props: { ...props, resetDisabled: false },
    });
    await wrapper.find(".detail-reset").trigger("click");
    expect(wrapper.emitted("reset")).toHaveLength(1);
    const disabled = mount(DeploymentHistory, {
      props: { ...props, resetDisabled: true },
    });
    expect(
      disabled.find<HTMLButtonElement>(".detail-reset").element.disabled,
    ).toBe(true);
  });
});

describe("ResetDialog", () => {
  it("opens from the model and reports cancel and confirm", async () => {
    const wrapper = mount(ResetDialog, {
      props: { open: false, penaltyThreshold: 2 },
      attachTo: document.body,
    });
    const dialog = wrapper.find<HTMLDialogElement>("dialog").element;
    expect(dialog.open).toBe(false);
    await wrapper.setProps({ open: true });
    expect(dialog.open).toBe(true);
    expect(wrapper.text()).toContain("減点なしの上限: 2回");

    await wrapper.find(".button-primary").trigger("click");
    expect(wrapper.emitted("confirm")).toHaveLength(1);
    await wrapper.find(".button-secondary").trigger("click");
    expect(wrapper.emitted("update:open")?.at(-1)).toEqual([false]);
    wrapper.unmount();
  });
});
