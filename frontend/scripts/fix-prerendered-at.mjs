import { readFile, writeFile } from "node:fs/promises";

const epoch = Number(process.env.SOURCE_DATE_EPOCH);
const fixed = epoch * 1000;
if (!Number.isSafeInteger(epoch) || epoch <= 0 || !Number.isSafeInteger(fixed)) {
  throw new Error("invalid SOURCE_DATE_EPOCH");
}

for (const app of ["contestant", "admin"]) {
  for (const page of ["index.html", "200.html", "404.html"]) {
    const path = `packages/${app}/.output/public/${page}`;
    const html = await readFile(path, "utf8");
    const script = /(<script\b[^>]*\bid="__NUXT_DATA__"[^>]*>)([\s\S]*?)(<\/script>)/g;
    const matches = [...html.matchAll(script)];
    if (matches.length !== 1) throw new Error(`${path}: expected one Nuxt payload`);
    const [, before, contents, after] = matches[0];
    const payload = JSON.parse(contents);
    const index = payload?.[0]?.prerenderedAt;
    const previous = payload?.[index];
    if (!Number.isInteger(index) || !Number.isSafeInteger(previous) ||
        previous < 1_000_000_000_000 || previous >= 10_000_000_000_000) {
      throw new Error(`${path}: expected Nuxt prerenderedAt epoch milliseconds`);
    }
    const oldValue = String(previous);
    const occurrence = new RegExp(`(?<![0-9])${oldValue}(?![0-9])`, "g");
    if ([...contents.matchAll(occurrence)].length !== 1) {
      throw new Error(`${path}: expected one prerenderedAt value`);
    }
    const updated = contents.replace(occurrence, String(fixed));
    await writeFile(path, html.replace(matches[0][0], before + updated + after));
  }
}
