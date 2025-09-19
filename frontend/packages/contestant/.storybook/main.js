/** @type { import('@storybook/react-vite').StorybookConfig } */
const config = {
  stories: ["../app/**/*.mdx", "../app/**/*.stories.@(js|jsx|mjs|ts|tsx)"],
  framework: {
    name: "@storybook/react-vite",
    options: {},
  },
};

export default config;
