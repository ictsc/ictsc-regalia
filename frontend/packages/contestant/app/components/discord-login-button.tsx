import { Button } from "@headlessui/react";
import { clsx } from "clsx";

export type DiscordLoginButtonViewProps = {
  readonly disabled: boolean;
  readonly onClick?: () => void;
};

export function DiscordLoginButton({
  disabled,
  onClick,
}: DiscordLoginButtonViewProps) {
  return (
    <Button
      disabled={disabled}
      onClick={onClick}
      className={clsx(
        "rounded-[12px] bg-[#5865f2] pb-8 pt-8 text-16 shadow-md disabled:bg-[#a0a0a0] data-[hover]:bg-[#4752c4]",
      )}
    >
      <span
        className={clsx("flex flex-row gap-[6px] pe-8 ps-8 text-surface-0")}
      >
        <img
          src="/assets/icon_clyde_white_RGB.svg"
          width={20}
          height={20}
          alt="Discord"
        />
        Discord でログイン
      </span>
    </Button>
  );
}
