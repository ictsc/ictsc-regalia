import globals from "globals";
import js from "@eslint/js";
import ts from "typescript-eslint";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import tailwind from "eslint-plugin-tailwindcss";
import storybook from "eslint-plugin-storybook";

const tsFileConfig = (conf) => ({
  ...conf,
  files: conf.files ?? ["**/*.{ts,tsx}"],
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
    ...tailwind.configs["flat/recommended"],
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
    {
      files: ["**/*.stories.{ts,tsx}"],
      plugins: { storybook, "react-hooks": reactHooks },
      rules: {
        "react-hooks/rules-of-hooks": "off",
        "storybook/await-interactions": "error",
        "storybook/context-in-play-function": "error",
        "storybook/default-exports": "error",
        "storybook/hierarchy-separator": "error",
        "storybook/no-redundant-story-name": "warn",
        "storybook/prefer-pascal-case": "warn",
        "storybook/story-exports": "error",
        "storybook/use-storybook-expect": "error",
        "storybook/use-storybook-testing-library": "error",
      },
    },
  ].map(tsFileConfig),
];
