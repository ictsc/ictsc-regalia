import { Button, type ButtonProps } from "@headlessui/react";
import clydeIcon from "../../assets/icon_clyde_white_RGB.svg";
import { Logo } from "../components/logo";
import { Title } from "../components/title";

type Props = {
  signInURL: string;
};

export function SignInPage({ signInURL }: Props) {
  return (
    <>
      <Title>ログイン</Title>
      <div className="mx-40 flex h-full flex-col items-center justify-center gap-[90px]">
        <Logo width={500} />
        <DiscordLoginButton href={signInURL} />
      </div>
    </>
  );
}

function DiscordLoginButton(props: ButtonProps<"a">) {
  return (
    <Button
      as="a"
      className="rounded-16 text-32 bg-[#5865f2] py-[22px] ps-16 pe-[20px] shadow-md disabled:bg-[#a0a0a0] data-[hover]:bg-[#4752c4]"
      {...props}
    >
      <span className="text-surface-0 flex flex-row gap-[12px]">
        <img src={clydeIcon} width={40} height={40} alt="" />
        Discord でログイン
      </span>
    </Button>
  );
}
