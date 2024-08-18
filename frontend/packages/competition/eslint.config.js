import js from "@eslint/js";
import globals from "globals";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import ts from "typescript-eslint";

const tsFileConfig = (conf) => ({
  ...conf,
  files: [...(conf.files ?? []), "**/*.{ts,tsx}"],
  languageOptions: {
    ...conf.languageOptions,
    ecmaVersion: 2023,
    globals: globals.browser,
  },
});

export default [
  { ignores: ["dist"] },
  ...[
    js.configs.recommended,
    ...ts.configs.recommended,
    {
      plugins: { "react-hooks": reactHooks },
      rules: reactHooks.configs.recommended.rules,
    },
    {
      plugins: { "react-refresh": reactRefresh },
      rules: {
        "react-refresh/only-export-components": [
          "warn",
          { allowConstantExport: true },
        ],
      },
    },
  ].map(tsFileConfig),
];
