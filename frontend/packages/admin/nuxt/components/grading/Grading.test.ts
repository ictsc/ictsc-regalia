import { describe, expect, it } from "vitest";
import { mount, RouterLinkStub } from "@vue/test-utils";
import GradingBoard from "./GradingBoard.vue";
import GradingFilters from "./GradingFilters.vue";
import { answerKey } from "~/features/grading";
import { answers, fixedNow, myDiscordId } from "~/stories/fixtures";

const global = { stubs: { NuxtLink: RouterLinkStub } };

describe("GradingBoard", () => {
  const mountBoard = (busyKey = "") =>
    mount(GradingBoard, {
      props: { answers, myDiscordId, busyKey, now: fixedNow },
      global,
    });
  it("groups answers into columns with the viewer's first", () => {
    const wrapper = mountBoard();
    const columns = wrapper.findAll(".grading-column");
    expect(
      columns.map((c) => c.find(".grading-column-header span").text()),
    ).toEqual(["3 件", "2 件", "2 件"]);
    const first = columns[0]!.find(".grading-card");
    expect(first.find(".grading-card-meta").text()).toBe("自分の担当");
    expect(first.findComponent(RouterLinkStub).props("to")).toBe(
      "/submissions/B02/4/3",
    );
    expect(first.find(".grading-card-age").text()).toBe("提出から 2 分");
    expect(columns[2]!.text()).toContain("300 / 300 点");
  });
  it("reports status moves and assignment changes", async () => {
    const wrapper = mountBoard();
    const card = wrapper.findAll(".grading-card")[0]!;
    await card.find("select").setValue("IN_PROGRESS");
    await card.find("button").trigger("click");
    expect(wrapper.emitted("update")).toEqual([
      [answers[0], { status: "IN_PROGRESS" }],
      [answers[0], { assignment: "CLAIM_SELF" }],
    ]);
    const claimed = wrapper
      .findAll(".grading-column")[1]!
      .find(".grading-card");
    expect(claimed.find("button").text()).toBe("担当を戻す");
    expect(claimed.find(".grading-card-meta").text()).toContain("個別担当");
  });
  it("disables the card being updated", () => {
    const wrapper = mountBoard(answerKey(answers[0]!));
    const card = wrapper.findAll(".grading-card")[0]!;
    expect(card.find<HTMLSelectElement>("select").element.disabled).toBe(true);
    expect(card.find<HTMLButtonElement>("button").element.disabled).toBe(true);
  });
});

describe("GradingFilters", () => {
  it("updates and clears the filter models", async () => {
    const wrapper = mount(GradingFilters, {
      props: {
        teams: ["KERNEL PANIC"],
        problems: ["B02"],
        recent: true,
        perfect: false,
        unscored: true,
        problemCodes: ["A01", "B02"],
        teamNames: ["KERNEL PANIC"],
      },
    });
    await wrapper.findAll('input[type="checkbox"]')[0]!.setValue(false);
    expect(wrapper.emitted("update:recent")).toEqual([[false]]);
    await wrapper.find("button").trigger("click");
    expect(wrapper.emitted("update:teams")?.at(-1)).toEqual([[]]);
    expect(wrapper.emitted("update:problems")?.at(-1)).toEqual([[]]);
    expect(wrapper.emitted("update:perfect")?.at(-1)).toEqual([true]);
    expect(wrapper.emitted("update:unscored")?.at(-1)).toEqual([false]);
  });
});
