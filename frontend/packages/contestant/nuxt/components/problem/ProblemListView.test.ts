import { describe, expect, it } from "vitest";
import { mount, RouterLinkStub } from "@vue/test-utils";
import ProblemListView, { type ProblemListItem } from "./ProblemListView.vue";
import type { Problem } from "~/features/models";
import { problemStatus } from "~/features/problem/status";

const section = (name: string, startAt: string, endAt: string) => ({
  name,
  startAt,
  endAt,
});
const day1 = section(
  "day1",
  "2026-09-01T00:00:00+09:00",
  "2026-09-02T00:00:00+09:00",
);
const day2 = section(
  "day2",
  "2026-09-02T00:00:00+09:00",
  "2026-09-03T00:00:00+09:00",
);

const problem = (
  code: string,
  schedule: typeof day1,
  score?: number,
  category = "ネットワーク",
): Problem => ({
  code,
  title: `問題 ${code}`,
  category,
  maxScore: 100,
  submissionableSchedules: [schedule],
  submissionStatus: { isSubmittable: schedule === day2 },
  score:
    score == null
      ? undefined
      : { markedScore: score, penalty: 0, score, maxScore: 100 },
});

const items: ProblemListItem[] = [
  problem("A01", day1, 100),
  problem("A02", day2, 40, "サーバー"),
  problem("A03", day2),
].map((p) => ({ problem: p, status: problemStatus(p), cooldownMinutes: 0 }));

function mountList(
  props: Partial<InstanceType<typeof ProblemListView>["$props"]> = {},
) {
  return mount(ProblemListView, {
    props: { items, ...props },
    global: { stubs: { NuxtLink: RouterLinkStub } },
  });
}

const rowCodes = (wrapper: ReturnType<typeof mountList>) =>
  wrapper.findAll("a.problem-row .problem-id b").map((b) => b.text());
const row = (wrapper: ReturnType<typeof mountList>, code: string) =>
  wrapper
    .findAllComponents(RouterLinkStub)
    .find((link) => link.props("to") === `/problems/${code}`)!;

describe("ProblemListView", () => {
  it("groups problems by section with open sections first", () => {
    const wrapper = mountList();
    expect(
      wrapper.findAll(".problem-day-heading").map((h) => h.text()),
    ).toEqual(["2日目の問題", "1日目の問題"]);
    expect(rowCodes(wrapper)).toEqual(["A02", "A03", "A01"]);
    expect(row(wrapper, "A01").classes()).toContain("is-answer-closed");
  });

  it("filters by status and resets the filters", async () => {
    const wrapper = mountList();
    await wrapper
      .findAll('input[type="checkbox"][value="unanswered"]')[0]!
      .setValue(true);
    expect(rowCodes(wrapper)).toEqual(["A03"]);
    await wrapper.find(".filter-reset").trigger("click");
    expect(rowCodes(wrapper)).toHaveLength(3);
  });

  it("shows cooldown and pending badges", () => {
    const wrapper = mountList({
      items: items.map((item) =>
        item.problem.code === "A02"
          ? { ...item, cooldownMinutes: 5 }
          : item.problem.code === "A03"
            ? { ...item, status: "pending" }
            : item,
      ),
    });
    expect(row(wrapper, "A02").find(".problem-cooldown-badge").text()).toBe(
      "再提出まで5分",
    );
    expect(row(wrapper, "A03").find(".problem-pending").text()).toBe("採点中");
  });

  it("shows the empty and error states", async () => {
    expect(mountList({ items: [] }).text()).toContain(
      "該当する問題はありません。",
    );
    const failed = mountList({ items: [], error: "failed" });
    await failed.find("button").trigger("click");
    expect(failed.emitted("retry")).toHaveLength(1);
  });
});
