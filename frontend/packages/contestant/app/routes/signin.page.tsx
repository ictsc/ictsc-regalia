import { Button, type ButtonProps } from "@headlessui/react";
import { Link } from "@tanstack/react-router";
import clydeIcon from "../../assets/icon_clyde_white_RGB.svg";
import { Logo } from "../components/logo";
import { Title } from "../components/title";

type Props = {
  signInURL: string;
  adminTokenAvailable: boolean;
};

export function SignInPage({ signInURL, adminTokenAvailable }: Props) {
  return (
    <>
      <Title>ログイン</Title>
      <div className="mx-40 flex h-full flex-col items-center justify-center gap-[48px] py-12">
        <Logo width={500} />
        <div className="flex flex-col items-center gap-12">
          <DiscordLoginButton href={signInURL} />
          {adminTokenAvailable ? (
            <Button
              as={Link}
              to="/signin/impersonation"
              className="text-16 text-text border-disabled rounded-16 data-hover:bg-surface-1 border px-20 py-12"
            >
              Admin としてログイン
            </Button>
          ) : null}
        </div>
      </div>
    </>
  );
}

function DiscordLoginButton(props: ButtonProps<"a">) {
  return (
    <Button
      as="a"
      className="rounded-16 text-32 bg-[#5865f2] py-[22px] ps-16 pe-[20px] shadow-md disabled:bg-[#a0a0a0] data-hover:bg-[#4752c4]"
      {...props}
    >
      <span className="text-surface-0 flex flex-row gap-[12px]">
        <img src={clydeIcon} width={40} height={40} alt="" />
        Discord でログイン
      </span>
    </Button>
  );
}
