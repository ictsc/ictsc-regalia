import { describe, expect, it, vi } from "vitest";
import { mount } from "@vue/test-utils";
import GradingDetailView from "./GradingDetailView.vue";
import { defaultPaneLayout } from "~/features/grading-panes";
import {
  adminProblem,
  contentStatus,
  detailAnswer,
  markings,
  previousAnswers,
} from "~/stories/fixtures";

function mountView(props: Record<string, unknown> = {}) {
  return mount(GradingDetailView, {
    props: {
      answer: detailAnswer,
      problem: adminProblem,
      activeCommit: contentStatus.active_commit!,
      previous: previousAnswers,
      marks: markings,
      busy: false,
      workflowBusy: false,
      collapsed: defaultPaneLayout().collapsed,
      widths: defaultPaneLayout().widths,
      submitMarking: async () => true,
      ...props,
    },
    global: { stubs: { MarkdownContent: true } },
  });
}

describe("GradingDetailView", () => {
  it("compares with the latest previous answer by default", () => {
    const wrapper = mountView();
    expect(
      wrapper.find<HTMLSelectElement>(".grading-previous-picker select").element
        .value,
    ).toBe("2");
    expect(wrapper.find(".grading-diff-filebar").text()).toContain(
      "過去 #2 → 今回 #3",
    );
    expect(wrapper.find(".grading-answer-lines").classes()).toContain(
      "is-unified-diff",
    );
    expect(wrapper.find(".grading-diff-stats").text()).toMatch(/\+\d+−\d+/);
  });

  it("shows the whole answer without previous answers", async () => {
    const wrapper = mountView({ previous: [] });
    expect(wrapper.find(".grading-view-switch").exists()).toBe(false);
    expect(wrapper.findAll(".grading-line-number")).toHaveLength(
      detailAnswer.body.body.split("\n").length,
    );
  });

  it("opens line comment drafts and sends them with the marking", async () => {
    const submitMarking = vi.fn(async () => true);
    const wrapper = mountView({ previous: [], marks: [], submitMarking });
    const lineNumbers = wrapper.findAll(".grading-line-number");
    await lineNumbers[2]!.trigger("click");
    await lineNumbers[4]!.trigger("click", { shiftKey: true });
    const editor = wrapper.find(".grading-line-editor");
    expect(editor.find("label").text()).toBe("3〜5 行目へのコメント");
    await editor.find("textarea").setValue("範囲へのコメント");
    await wrapper.find('input[type="number"]').setValue("250");
    await wrapper.find(".form-grid textarea").setValue("全体");
    expect(wrapper.text()).toContain("行コメント 1 件を含めて保存します。");
    await wrapper.find(".form-grid").trigger("submit");
    await Promise.resolve();
    expect(submitMarking).toHaveBeenCalledWith({
      score: 250,
      comment: "全体",
      lineComments: [
        { line_number: 3, end_line_number: 5, body: "範囲へのコメント" },
      ],
    });
    await Promise.resolve();
    expect(wrapper.find(".grading-line-editor").exists()).toBe(false);
  });

  it("keeps drafts when the marking was not saved", async () => {
    const wrapper = mountView({
      previous: [],
      marks: [],
      submitMarking: async () => false,
    });
    await wrapper.findAll(".grading-line-number")[0]!.trigger("click");
    await wrapper.find(".grading-line-editor textarea").setValue("残す");
    await wrapper.find('input[type="number"]').setValue("10");
    await wrapper.find(".form-grid").trigger("submit");
    await Promise.resolve();
    await Promise.resolve();
    expect(wrapper.find(".grading-line-editor").exists()).toBe(true);
  });

  it("shows marking history and reports comment deletion and assignment", async () => {
    const wrapper = mountView();
    expect(wrapper.find(".grading-history-comments summary").text()).toBe(
      "行コメント 2 件 （削除済み 1 件）",
    );
    await wrapper.find(".grading-history-comment-head button").trigger("click");
    expect(wrapper.emitted("deleteLineComment")).toEqual([["mark-2", 0]]);
    await wrapper.find(".grading-detail-meta button").trigger("click");
    expect(wrapper.emitted("assignment")).toEqual([["RESET_TO_DEFAULT"]]);
  });

  it("collapses panes through the model", async () => {
    const wrapper = mountView();
    await wrapper.find(".grading-pane-header button").trigger("click");
    expect(wrapper.emitted("update:collapsed")).toEqual([
      [[true, false, false]],
    ]);
  });
});
