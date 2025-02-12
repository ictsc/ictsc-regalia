import { clsx } from "clsx";
import { useState } from "react";

export type DiscordLoginButtonViewProps = {
  readonly disabled: boolean;
  readonly onClick?: () => void;
};

export function DiscordLoginButton({
  disabled,
  onClick,
}: DiscordLoginButtonViewProps) {
  const [isHovered, setIsHovered] = useState(false);
  const backgroundColor = disabled
    ? "rgba(160, 160, 160, 1)"
    : isHovered
      ? "rgba(71, 82, 196, 1)"
      : "rgba(88, 101, 242, 1)";

  return (
    <button
      disabled={disabled}
      onClick={onClick}
      style={{ backgroundColor }}
      className={clsx("rounded-[12px] pb-8 pt-8 text-16 shadow-md")}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
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
    </button>
  );
}
