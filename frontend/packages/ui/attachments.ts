/** A text log attached to an answer. */
export type AttachmentFile = {
  name: string;
  size: number;
  content: string;
};

export const maxAttachments = 3;
export const maxAttachmentBytes = 1 << 20;

// Mirrors the server rule: plain file names ending in .txt or .log.
const namePattern = /^[A-Za-z0-9_-][A-Za-z0-9._-]{0,95}\.(txt|log)$/;

/** Characters counted against answer_max_length; line breaks are free. */
export function answerLength(body: string): number {
  let length = 0;
  for (const character of body) {
    if (character !== "\n" && character !== "\r") length++;
  }
  return length;
}

/**
 * Makes a browser file name acceptable to the server without changing its
 * extension, e.g. "app (1).log" becomes "app__1_.log".
 */
export function attachmentName(fileName: string): string | undefined {
  const match = /\.(txt|log)$/i.exec(fileName);
  if (!match) return undefined;
  const stem = fileName
    .slice(0, -match[0].length)
    .replace(/[^A-Za-z0-9._-]/g, "_")
    .replace(/^\.+/, "")
    .slice(0, 92);
  const name = `${stem || "log"}.${match[1]!.toLowerCase()}`;
  return namePattern.test(name) ? name : undefined;
}

/**
 * Returns name, or name with -2, -3, … before the extension when it is
 * already taken, e.g. "app.log" becomes "app-2.log".
 */
export function uniqueAttachmentName(
  name: string,
  takenNames: readonly string[],
): string {
  if (!takenNames.includes(name)) return name;
  const dot = name.lastIndexOf(".");
  for (let index = 2; ; index++) {
    const candidate = `${name.slice(0, dot)}-${index}${name.slice(dot)}`;
    if (!takenNames.includes(candidate)) return candidate;
  }
}

/**
 * Reads a file as an attachment, rejecting anything that is not UTF-8 text.
 * Returns an error message instead of throwing.
 */
export async function readAttachment(
  file: File,
): Promise<AttachmentFile | string> {
  const name = attachmentName(file.name);
  if (!name) return `${file.name}: .txt か .log のファイルだけ添付できます`;
  if (file.size === 0) return `${file.name}: 空のファイルは添付できません`;
  if (file.size > maxAttachmentBytes)
    return `${file.name}: 1MB を超えるファイルは添付できません`;
  let content: string;
  try {
    content = new TextDecoder("utf-8", { fatal: true }).decode(
      await file.arrayBuffer(),
    );
  } catch {
    return `${file.name}: UTF-8 のテキストファイルだけ添付できます`;
  }
  if (content.includes("\u0000"))
    return `${file.name}: バイナリファイルは添付できません`;
  return { name, size: file.size, content };
}

/**
 * Builds an attachment from text typed or pasted into the log editor. A name
 * without .txt or .log gets .txt. Returns an error message instead of throwing.
 */
export function textAttachment(
  fileName: string,
  content: string,
): AttachmentFile | string {
  const trimmed = fileName.trim() || "log";
  const name = attachmentName(
    /\.(txt|log)$/i.test(trimmed) ? trimmed : `${trimmed}.txt`,
  );
  if (!name) return "ファイル名を確認してください";
  if (content === "") return "ログの内容を入力してください";
  if (content.includes("\u0000")) return "バイナリは添付できません";
  const size = new TextEncoder().encode(content).length;
  if (size > maxAttachmentBytes) return "1MB を超えるログは添付できません";
  return { name, size, content };
}

/**
 * Removes what a terminal would act on instead of print: escape sequences
 * (colors, cursor moves, window titles, clipboard writes), other control
 * characters and carriage returns that could hide a line when the saved
 * log is printed with cat. Tabs and line breaks are kept.
 */
export function stripTerminalControls(text: string): string {
  return (
    text
      // OSC (e.g. title, hyperlink, clipboard) ended by BEL or ST.
      .replace(/\u001b\][^\u0007\u001b]*(?:\u0007|\u001b\\)?/g, "")
      // DCS, SOS, PM and APC strings ended by ST.
      .replace(/\u001b[PX^_][^\u001b]*(?:\u001b\\)?/g, "")
      // CSI (colors, cursor movement, erasing).
      .replace(/\u001b\[[0-?]*[ -/]*[@-~]/g, "")
      // Any other escape sequence.
      .replace(/\u001b[\s\S]?/g, "")
      // A bare CR returns to the line start and would overwrite it.
      .replace(/\r\n?/g, "\n")
      // Remaining C0 controls except tab and line feed, DEL and C1 controls.
      .replace(/[\u0000-\u0008\u000b-\u001f\u007f-\u009f]/g, "")
  );
}

/**
 * Saves an attachment from memory; nothing is fetched from the server.
 * Terminal control sequences are removed so printing the file is safe.
 */
export function downloadAttachment(attachment: AttachmentFile) {
  const url = URL.createObjectURL(
    new Blob([stripTerminalControls(attachment.content)], {
      type: "text/plain;charset=utf-8",
    }),
  );
  const link = document.createElement("a");
  link.href = url;
  link.download = attachment.name;
  link.rel = "noopener";
  link.click();
  setTimeout(() => URL.revokeObjectURL(url), 0);
}

export function formatBytes(size: number): string {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / 1024 / 1024).toFixed(2)} MB`;
}
