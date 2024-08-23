import type { Preview } from "@storybook/react";
import { withTanstack } from "./tanstack-decorator";
import "../app/index.css";

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
  },
  decorators: [withTanstack],
};

export default preview;
