import withNuxt from "./.nuxt/eslint.config.mjs";
export default withNuxt(
  {
    rules: {
      "vue/multi-word-component-names": "off",
      "vue/html-self-closing": "off",
    },
  },
  {
    // Storybook renders these with plain Vite, so they must not depend on Nuxt
    // runtime imports or call the API.
    files: ["nuxt/components/**/*.vue", "nuxt/stories/**/*.ts"],
    rules: {
      "@typescript-eslint/no-restricted-imports": [
        "error",
        {
          paths: [
            {
              name: "@ictsc/api",
              allowTypeImports: true,
              message: "Fetch data in the page and pass it as props.",
            },
          ],
          patterns: [
            {
              group: ["#*", "nuxt", "nuxt/*"],
              message: "Storybook components cannot use Nuxt imports.",
            },
          ],
        },
      ],
    },
  },
);
