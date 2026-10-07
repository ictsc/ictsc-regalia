// Limit IBM Plex Sans JP to Latin glyphs so Japanese UI text falls back to Zen Maru Gothic.
const latinGlyphs = Array.from({ length: 95 }, (_, index) =>
  String.fromCharCode(index + 32),
).join("");

export const latinFontStylesheet = `https://fonts.googleapis.com/css2?family=IBM+Plex+Sans+JP:wght@400;500;600;700&text=${encodeURIComponent(latinGlyphs)}&display=swap`;

export const japaneseFontStylesheet =
  "https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500&family=Noto+Sans+JP:wght@400;500;600;700&family=Zen+Maru+Gothic:wght@700&display=swap";
