import globals from "globals";
import jsPlugin from "@eslint/js";
import prettierConfig from "eslint-config-prettier";
import importPlugin from "eslint-plugin-import";
import tseslint from "typescript-eslint";
import reactPlugin from "eslint-plugin-react";
import reactHooksPlugin from "eslint-plugin-react-hooks";
import reactRefreshPlugin from "eslint-plugin-react-refresh";
import tailwindPlugin from "eslint-plugin-tailwindcss";
import storybookPlugin from "eslint-plugin-storybook";

const jsFiles = ["**/*.{js,jsx,cjs,mjs}"];
const tsFiles = ["**/*.{ts,tsx,cts,mts}"];

/** @type {import("eslint").Linter.FlatConfig[]} */
const config = [
  { ignores: ["dist"] },
  // プラグイン
  {
    name: "@ictsc/eslint-config/plugins",
    plugins: {
      import: importPlugin,
      "@typescript-eslint": tseslint.plugin,
      react: reactPlugin,
      "react-hooks": reactHooksPlugin,
      "react-refresh": reactRefreshPlugin,
      tailwindcss: tailwindPlugin,
      storybook: storybookPlugin,
    },
  },
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
    settings: {
      ...importPlugin.configs.typescript.settings,
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
    { rules: importPlugin.configs.recommended.rules },
    {
      rules: {
        "import/no-unresolved": "off",
      },
    },
    ...tailwindPlugin.configs["flat/recommended"],
    { rules: reactPlugin.configs.recommended.rules },
    { rules: reactPlugin.configs["jsx-runtime"].rules },
    { rules: reactHooksPlugin.configs.recommended.rules },
    {
      rules: {
        "react-refresh/only-export-components": [
          "warn",
          { allowConstantExport: true },
        ],
      },
    },
  ].map(({ languageOptions: _ignore, ...config }) => ({
    ...config,
    files: [jsFiles, tsFiles],
  })),
  // TS固有のルール設定
  ...[
    ...tseslint.configs.recommendedTypeChecked,
    { rules: importPlugin.configs.typescript.rules },
    {
      rules: {
        "import/named": "off",
        "import/namespace": "off",
        "import/default": "off",
        "import/no-named-as-default-member": "off",
        "import/no-unresolved": "off",
      },
    },
  ].map(({ languageOptions: _ignore, ...config }) => ({
    ...config,
    files: [tsFiles],
  })),
  // Storybookルール設定
  {
    files: ["**/*.stories.{ts,tsx}"],
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
  // フォーマットに関するルールを無効化
  { rules: prettierConfig.rules },
];

export default config;
