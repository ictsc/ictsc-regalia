import { type ReactNode } from "react";
import { Link } from "@tanstack/react-router";
import { Logo } from "@app/components/logo";

export type HeaderViewProps = {
  readonly contestState?: ReactNode;
  readonly accountMenu?: ReactNode;
};

export function Header({ contestState, accountMenu }: HeaderViewProps) {
  return (
    <div className="border-primary bg-surface-0 flex size-full items-center border-b-[3px]">
      <div className="ml-4 flex-none sm:ml-16">
        <Link to="/">
          <Logo className="scale-75 sm:scale-100" height={50} />
        </Link>
      </div>
      <div className="ml-auto flex h-full items-center">
        <div className="mr-[30px]">{contestState}</div>
        <div className="bg-primary flex h-full w-[140px] items-center justify-end pt-[3px] [clip-path:polygon(40%_0,100%_0,100%_100%,0_100%)]">
          <div className="mr-[30px]">{accountMenu}</div>
        </div>
      </div>
    </div>
  );
}
