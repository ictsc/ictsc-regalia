import { mount } from "@vue/test-utils";
import { defineComponent } from "vue";
import { beforeEach, expect, it } from "vitest";
import { useReadAnnouncements } from "./announce-read-status";

beforeEach(() => localStorage.clear());

it("preserves earlier read notices when marking a detail before mount", () => {
  localStorage.setItem("readAnnouncements", JSON.stringify(["earlier"]));
  const wrapper = mount(
    defineComponent({
      setup() {
        const { markAsRead } = useReadAnnouncements();
        markAsRead("current");
        return () => null;
      },
    }),
  );
  expect(JSON.parse(localStorage.getItem("readAnnouncements")!)).toEqual([
    "earlier",
    "current",
  ]);
  wrapper.unmount();
});
