import { type ReactNode } from "react";
import { MenuButton } from "@headlessui/react";
import { Link } from "@tanstack/react-router";
import { MaterialSymbol } from "@app/components/material-symbol";
import { Logo } from "@app/components/logo";
import { AccountMenu } from "./account-menu";

export type HeaderViewProps = {
  readonly contestState?: ReactNode;
};

export function HeaderView({ contestState }: HeaderViewProps) {
  return (
    <div className="flex size-full items-center border-b-[3px] border-primary bg-surface-0">
      <div className="flex-none">
        <Link>
          <Logo height={60} />
        </Link>
      </div>
      <div className="ml-auto flex h-full items-center">
        {contestState != null && (
          <div className="mr-[30px]">{contestState}</div>
        )}
        <div className="flex h-full w-[140px] items-center justify-end bg-primary pt-[3px] [clip-path:polygon(40%_0,100%_0,100%_100%,0_100%)]">
          <AccountMenu>
            <MenuButton
              title="アカウントメニュー"
              className="mr-[30px] flex size-[50px] items-center justify-center rounded-full transition data-[hover]:bg-surface-0/50"
            >
              <MaterialSymbol
                icon="person"
                fill
                size={40}
                className="text-surface-0"
              />
            </MenuButton>
          </AccountMenu>
        </div>
      </div>
    </div>
  );
}
