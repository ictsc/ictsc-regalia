import assert from "node:assert/strict";
import { execFileSync } from "node:child_process";
import { mkdtempSync, mkdirSync, readFileSync, writeFileSync, rmSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { fileURLToPath } from "node:url";
import test from "node:test";

const script = fileURLToPath(new URL("./fix-prerendered-at.mjs", import.meta.url));

test("replaces only the Nuxt prerenderedAt payload value", () => {
  const directory = mkdtempSync(join(tmpdir(), "regalia-prerender-"));
  try {
    for (const app of ["contestant", "admin"]) {
      const output = join(directory, "packages", app, ".output", "public");
      mkdirSync(output, { recursive: true });
      for (const page of ["index.html", "200.html", "404.html"]) {
        writeFileSync(join(output, page),
          '<html><script id="__NUXT_DATA__" type="application/json">' +
          '[{"prerenderedAt":1},1793350732000]</script>' +
          '<span>1793350732000</span></html>');
      }
    }
    execFileSync(process.execPath, [script], {
      cwd: directory,
      env: { ...process.env, SOURCE_DATE_EPOCH: "1793350700" },
    });
    const result = readFileSync(join(directory, "packages/contestant/.output/public/index.html"), "utf8");
    assert.match(result, /\[{"prerenderedAt":1},1793350700000\]/);
    assert.match(result, /<span>1793350732000<\/span>/);
  } finally {
    rmSync(directory, { recursive: true, force: true });
  }
});
