import { expect, it } from "vitest";
import {
  answerLength,
  attachmentName,
  readAttachment,
  textAttachment,
  uniqueAttachmentName,
} from "./attachments";

it("counts characters without line breaks", () => {
  expect(answerLength("あい\r\nう\nx")).toBe(4);
  expect(answerLength("🙂")).toBe(1);
});

it("normalizes file names and keeps only .txt and .log", () => {
  expect(attachmentName("app (1).LOG")).toBe("app__1_.log");
  expect(attachmentName(".bashrc.txt")).toBe("bashrc.txt");
  expect(attachmentName("../../etc/passwd.txt")).toBe("_.._etc_passwd.txt");
  expect(attachmentName("run.sh")).toBeUndefined();
  expect(attachmentName("page.html")).toBeUndefined();
});

it("numbers a name that is already attached", () => {
  expect(uniqueAttachmentName("app.log", [])).toBe("app.log");
  expect(uniqueAttachmentName("app.log", ["app.log"])).toBe("app-2.log");
  expect(uniqueAttachmentName("app.log", ["app.log", "app-2.log"])).toBe(
    "app-3.log",
  );
});

it("reads UTF-8 text and rejects binary files", async () => {
  const text = await readAttachment(
    new File(["\u001b[31merror\u001b[0m\n"], "app.log"),
  );
  expect(text).toEqual({
    name: "app.log",
    size: 15,
    content: "\u001b[31merror\u001b[0m\n",
  });
  expect(
    await readAttachment(new File([new Uint8Array([0xff, 0xfe])], "a.log")),
  ).toContain("UTF-8");
  expect(await readAttachment(new File(["a\u0000b"], "a.log"))).toContain(
    "バイナリ",
  );
  expect(await readAttachment(new File([""], "a.log"))).toContain("空");
});

it("builds attachments from pasted text", () => {
  expect(textAttachment("nginx error", "line\n")).toEqual({
    name: "nginx_error.txt",
    size: 5,
    content: "line\n",
  });
  expect(textAttachment("app.log", "ログ")).toEqual({
    name: "app.log",
    size: 6,
    content: "ログ",
  });
  expect(textAttachment("a", "")).toContain("入力");
  expect(textAttachment("a", "x\u0000")).toContain("バイナリ");
});
