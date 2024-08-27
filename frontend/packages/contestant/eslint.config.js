import configs from "@ictsc/config/eslint";

export default [
  { ignores: ["dist"] },
  ...configs({ react: true, storybook: true }),
];
