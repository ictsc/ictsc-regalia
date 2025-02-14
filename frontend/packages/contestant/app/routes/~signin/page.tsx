import { Logo } from "@app/components/logo";
import { Button, type ButtonProps } from "@headlessui/react";

type Props = {
  signInURL: string;
};

export function SignInPage({ signInURL }: Props) {
  return (
    <div className="flex h-full flex-col items-center justify-center gap-[90px]">
      <Logo width={500} />
      <DiscordLoginButton href={signInURL} />
    </div>
  );
}

function DiscordLoginButton(props: ButtonProps<"a">) {
  return (
    <Button
      as="a"
      className="rounded-16 bg-[#5865f2] py-[22px] pe-[20px] ps-16 text-32 shadow-md disabled:bg-[#a0a0a0] data-[hover]:bg-[#4752c4]"
      {...props}
    >
      <span className="flex flex-row gap-[12px] text-surface-0">
        <img
          src="/assets/icon_clyde_white_RGB.svg"
          width={40}
          height={40}
          alt=""
        />
        Discord でログイン
      </span>
    </Button>
  );
}
