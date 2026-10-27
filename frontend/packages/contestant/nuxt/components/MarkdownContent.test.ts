import { beforeAll, expect, it, vi } from "vitest";
import { flushPromises, mount } from "@vue/test-utils";
import MarkdownContent from "@ictsc/ui/components/MarkdownContent.vue";

beforeAll(() => {
  // jsdom does not implement modal dialogs.
  HTMLDialogElement.prototype.showModal ??= function (this: HTMLDialogElement) {
    this.open = true;
  };
  HTMLDialogElement.prototype.close ??= function (this: HTMLDialogElement) {
    this.open = false;
    this.dispatchEvent(new Event("close"));
  };
});

async function render(source: string) {
  const wrapper = mount(MarkdownContent, {
    props: { source },
    attachTo: document.body,
  });
  // Rendering is async (syntax highlighting), then the copy buttons are added.
  await vi.waitFor(() => {
    expect(wrapper.find(".markdown").element.innerHTML).not.toBe("");
  });
  await flushPromises();
  return wrapper;
}

it("copies the text of a code block", async () => {
  const writeText = vi.fn().mockResolvedValue(undefined);
  Object.assign(navigator, { clipboard: { writeText } });
  const wrapper = await render(
    "```bash\n$ ip route\ndefault via 192.0.2.1\n```",
  );
  const button = wrapper.find(".code-block .code-copy");
  expect(button.attributes("aria-label")).toBe("コードをコピー");
  expect(button.find("svg").exists()).toBe(true);
  await button.trigger("click");
  await flushPromises();
  expect(writeText).toHaveBeenCalledWith("$ ip route\ndefault via 192.0.2.1");
  expect(button.attributes("aria-label")).toBe("コピーしました");
  wrapper.unmount();
});

it("opens an image in the viewer and closes it again", async () => {
  const wrapper = await render("![構成図](https://example.com/net.png)");
  const dialog = wrapper.find<HTMLDialogElement>("dialog.image-viewer");
  expect(dialog.element.open).toBe(false);
  await wrapper.find(".markdown img").trigger("click");
  await flushPromises();
  expect(dialog.element.open).toBe(true);
  expect(dialog.find("img").attributes("alt")).toBe("構成図");
  await dialog.find(".image-viewer-close").trigger("click");
  expect(dialog.element.open).toBe(false);
  wrapper.unmount();
});
