import { withTanstack } from "./tanstack-decorator";
import "../app/index.css";

/** @type{import("@storybook/react").Preview} */
const preview = {
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
