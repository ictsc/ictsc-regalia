import type { Meta, StoryObj } from "@storybook/react";
import { DiscordLoginButton } from "./discord-login-button";

export default {
  title: "DiscordLoginButton",
  component: DiscordLoginButton,
} satisfies Meta<typeof DiscordLoginButton>;

type Story = StoryObj<typeof DiscordLoginButton>;

export const All: Story = {
  render: () => (
    <div className="flex flex-col items-start gap-24">
      <DiscordLoginButton
        disabled={false}
        onClick={() => {
          document
            .getElementById("discord-login-button-onclick-text")!
            .toggleAttribute("hidden");
        }}
      />
      <DiscordLoginButton disabled={true} />
      <span id="discord-login-button-onclick-text" hidden>
        clicked!
      </span>
    </div>
  ),
};
