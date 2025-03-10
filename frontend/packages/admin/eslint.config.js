import config from "@ictsc/config/eslint";

export default [
  { ignores: ["dist"] },
  ...config({ react: true, storybook: true }),
];
