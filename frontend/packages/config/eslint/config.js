import { fixupPluginRules } from "@eslint/compat";
import globals from "globals";
import jsPlugin from "@eslint/js";
import prettierConfig from "eslint-config-prettier";
import importPlugin from "eslint-plugin-import";
import * as tseslint from "typescript-eslint";
import reactPlugin from "eslint-plugin-react";
import reacthooksPlugin from "eslint-plugin-react-hooks";
import reactRefreshPlugin from "eslint-plugin-react-refresh";
import storybookPlugin from "eslint-plugin-storybook";

const namePrefix = "@ictsc/config/eslint";
const jsFiles = ["**/*.{js,jsx,cjs,mjs}"];
const tsFiles = ["**/*.{ts,tsx,cts,mts}"];

/**
 * @param [opts]
 * @returns {import("eslint").Linter.Config[]}
 */
function config(opts = {}) {
  const { react = false, storybook = false } = opts;
  return [
    {
      name: namePrefix + "/plugins",
      plugins: {
        import: fixupPluginRules(importPlugin),
        "@typescript-eslint": tseslint.plugin,
        react: reactPlugin,
        "react-hooks": reacthooksPlugin,
        "react-refresh": reactRefreshPlugin,
        storybook: storybookPlugin,
      },
    },
    {
      name: namePrefix + "/parser",
      files: [...jsFiles, ...tsFiles],
      languageOptions: {
        ecmaVersion: "latest",
        sourceType: "module",
        parserOptions: {
          ecmaVersion: "latest",
          sourceType: "module",
          ...(react ? { ecmaFeatures: { jsx: true } } : {}),
        },
      },
      settings: {
        react: { version: "detect" },
        "import/parsers": {
          espree: [".js", ".cjs", ".mjs", ".jsx"],
        },
        "import/resolver": {
          typescript: {},
        },
      },
    },
    {
      name: namePrefix + "/parser/commonjs",
      files: ["**/*.{cjs,cts}"],
      languageOptions: {
        sourceType: "commonjs",
      },
    },
    {
      name: namePrefix + "/parser/ts",
      files: tsFiles,
      languageOptions: {
        parser: tseslint.parser,
        parserOptions: { project: true },
        globals: globals.browser,
      },
      settings: {
        ...importPlugin.configs.typescript.settings,
        "import/resolver": {
          typescript: {},
        },
      },
    },
    ...mapToFiles(
      [...jsFiles, ...tsFiles],
      [
        {
          name: namePrefix + "/js",
          rules: {
            ...jsPlugin.configs.recommended.rules,
            "no-unused-vars": [
              "error",
              {
                ignoreRestSiblings: true,
                argsIgnorePattern: "^_",
                caughtErrors: "all",
              },
            ],
          },
        },
        {
          name: namePrefix + "/import",
          rules: {
            ...importPlugin.flatConfigs.recommended.rules,
          },
        },
        ...(react
          ? [
              {
                name: namePrefix + "/react",
                rules: {
                  ...reactPlugin.configs.recommended.rules,
                  ...reactPlugin.configs["jsx-runtime"].rules,
                  ...reacthooksPlugin.configs.recommended.rules,
                  "react-refresh/only-export-components": [
                    "warn",
                    {
                      allowConstantExport: true,
                    },
                  ],
                },
              },
            ]
          : []),
      ],
    ),
    ...mapToFiles(tsFiles, [
      {
        name: namePrefix + "/ts",
        rules: {
          ...tseslint.configs.recommendedTypeChecked.reduce(
            (rules, config) => ({ ...rules, ...config.rules }),
            {},
          ),
          // TypeScript が解決するため必要ないルールを無効化する
          "import/named": "off",
          "import/namespace": "off",
          "import/default": "off",
          "import/no-named-as-default-member": "off",
          "import/no-unresolved": "off",
          "no-unused-vars": "off",
          "@typescript-eslint/no-unused-vars": [
            "error",
            {
              ignoreRestSiblings: true,
              argsIgnorePattern: "^_",
              caughtErrors: "all",
            },
          ],
        },
      },
    ]),
    ...(storybook
      ? [
          {
            name: namePrefix + "/storybook",
            files: ["**/*.stories.{js,jsx,ts,tsx}"],
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
        ]
      : []),
    { name: "eslint-config-prettier", rules: prettierConfig.rules },
  ];
}

/**
 * @param {string[]} [files]
 * @param {readonly import("eslint").Linter.Config[]} [configs]
 * @returns {import("eslint").Linter.Config[]}
 */
function mapToFiles(files, configs) {
  return configs.map((config) => ({
    ...config,
    files,
  }));
}

export default config;
