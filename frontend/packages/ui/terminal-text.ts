/**
 * Removes what a terminal would act on instead of print: escape sequences
 * (colors, cursor moves, window titles, clipboard writes), other control
 * characters and carriage returns that could hide a line when the text is
 * printed (a downloaded log) or pasted (copied code) in a terminal, and
 * bidirectional overrides that make text read differently from what runs.
 * Tabs and line breaks are kept.
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
      // Bidirectional embeddings, overrides and isolates.
      .replace(/[\u202a-\u202e\u2066-\u2069]/g, "")
  );
}
