import path from "node:path";
import fs from "node:fs/promises";
import packageJSON from "../package.json";

const packagePath = path.dirname(
  new URL(import.meta.resolve("../package.json")).pathname,
);

const protoPath = path.join(packagePath, "proto");

await fs.rm("exports", { recursive: true }).catch(() => {});

for (const [key, def] of Object.entries(packageJSON.exports)) {
  if (typeof def !== "string") {
    continue;
  }
  const targetPath = path.resolve(packagePath, def);

  let result = `// Code generated by scripts/generate-exports.ts; DO NOT EDIT.\n`;

  const files = await fs.readdir(path.join(protoPath, key));
  for (const protoFile of files) {
    const filePath = path.join(protoPath, key, protoFile);
    const exportPath = path.relative(path.dirname(targetPath), filePath);
    result += `export * from "${exportPath}";\n`;
  }

  await fs.mkdir(path.dirname(targetPath), { recursive: true });
  console.log(path.relative(packagePath, targetPath))
  await fs.writeFile(targetPath, result, {});
}