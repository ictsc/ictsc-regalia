import { type ReactNode } from "react";
import { Link } from "@tanstack/react-router";
import { Logo } from "@app/components/logo";

export type HeaderViewProps = {
  readonly contestState?: ReactNode;
  readonly accountMenu?: ReactNode;
};

export function Header({ contestState, accountMenu }: HeaderViewProps) {
  return (
    <div className="flex size-full items-center border-b-[3px] border-primary bg-surface-0">
      <div className="ml-16 flex-none">
        <Link to="/">
          <Logo height={50} />
        </Link>
      </div>
      <div className="ml-auto flex h-full items-center">
        <div className="mr-[30px]">{contestState}</div>
        <div className="flex h-full w-[140px] items-center justify-end bg-primary pt-[3px] [clip-path:polygon(40%_0,100%_0,100%_100%,0_100%)]">
          <div className="mr-[30px]">{accountMenu}</div>
        </div>
      </div>
    </div>
  );
}
