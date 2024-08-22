import { MenuButton } from "@headlessui/react";
import { MaterialSymbol } from "../material-symbol";
import { ContestStateView } from "./contest-state";
import { AccountMenu } from "./account-menu";

export function HeaderView({
  staticAccountMenu,
}: {
  readonly staticAccountMenu?: boolean;
}) {
  return (
    <div className="flex size-full items-center border-b-[3px] border-primary bg-surface-0">
      <span className="flex-none">ICTSC</span>
      <div className="ml-auto flex h-full items-center">
        <div className="mr-[30px]">
          <ContestStateView />
        </div>
        <div className="flex h-full w-[140px] items-center justify-end bg-primary pt-[3px] [clip-path:polygon(40%_0,100%_0,100%_100%,0_100%)]">
          <AccountMenu static={staticAccountMenu}>
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
