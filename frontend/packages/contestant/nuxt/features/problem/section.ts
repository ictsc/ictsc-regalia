// "day1" -> "1日目", "both" -> "両日"; unknown slugs are shown as is.
export function sectionLabel(slug?: string) {
  if (slug === "both") return "両日";
  if (slug && /^day\d+$/.test(slug)) return `${slug.slice(3)}日目`;
  return slug;
}
