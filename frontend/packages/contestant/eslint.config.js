import globals from "globals";
import jsPlugin from "@eslint/js";
import tseslint from "typescript-eslint";
import reactHooksPlugin from "eslint-plugin-react-hooks";
import reactRefreshPlugin from "eslint-plugin-react-refresh";
import tailwindPlugin from "eslint-plugin-tailwindcss";
import storybookPlugin from "eslint-plugin-storybook";

const jsFiles = ["**/*.{js,jsx,cjs,mjs}"];
const tsFiles = ["**/*.{ts,tsx,cts,mts}"];

/** @type {import("eslint").Linter.FlatConfig[]} */
const config = [
  { ignores: ["dist"] },
  // パーサーの設定
  {
    files: [jsFiles, tsFiles],
    languageOptions: {
      ecmaVersion: "latest",
      parserOptions: { ecmaFeatures: { jsx: true } },
      globals: globals.browser,
    },
    settings: {
      react: { version: "detect" },
    },
  },
  {
    files: ["**/*.{js,jsx,mjs}"],
    languageOptions: {
      sourceType: "module",
    },
  },
  {
    files: ["**/*.cjs"],
    languageOptions: {
      sourceType: "commonjs",
    },
  },
  {
    files: [tsFiles],
    languageOptions: {
      sourceType: "module",
      parser: tseslint.parser,
      parserOptions: { projectService: true },
    },
  },
  // JS/TS共通ルール設定
  ...[
    jsPlugin.configs.recommended,
    {
      rules: {
        "no-unused-vars": [
          "error",
          {
            argsIgnorePattern: "^_",
            varsIgnorePattern: "^_",
            caughtErrorsIgnorePattern: "^_",
          },
        ],
      },
    },
    ...tailwindPlugin.configs["flat/recommended"].map(
      ({ languageOptions: _ignore, ...c }) => c,
    ),
    {
      plugins: { "react-hooks": reactHooksPlugin },
      rules: reactHooksPlugin.configs.recommended.rules,
    },
    {
      plugins: { "react-refresh": reactRefreshPlugin },
      rules: {
        "react-refresh/only-export-components": [
          "warn",
          { allowConstantExport: true },
        ],
      },
    },
  ].map((config) => ({ ...config, files: [jsFiles, tsFiles] })),
  // TS固有のルール設定
  ...[
    ...tseslint.configs.recommendedTypeChecked.map(
      ({ languageOptions: _ignore, ...c }) => c,
    ),
  ].map((config) => ({ ...config, files: [tsFiles] })),
  // Storybookルール設定
  {
    files: ["**/*.stories.{ts,tsx}"],
    plugins: { storybook: storybookPlugin, "react-hooks": reactHooksPlugin },
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
];

export default config;
