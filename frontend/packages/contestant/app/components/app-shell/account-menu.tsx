import type { ReactNode } from "react";
import {
  Menu,
  MenuItems,
  MenuItem,
  Button,
  type ButtonProps,
} from "@headlessui/react";
import { clsx } from "clsx";
import { MaterialSymbol, type MaterialSymbolType } from "../material-symbol";

export function AccountMenu({
  children,
  static: staticOpen,
}: {
  readonly children?: ReactNode;
  readonly static?: boolean;
}) {
  return (
    <Menu>
      {children}
      <MenuItems
        static={staticOpen}
        anchor={{ to: "bottom", gap: 15 }}
        transition
        className={clsx(
          "flex w-[200px] flex-col gap-[5px] rounded-[12px] bg-surface-0 py-[15px] drop-shadow",
          "transition duration-200 ease-out data-[closed]:opacity-0",
        )}
      >
        <span className="mx-[15px] text-14 text-text">ictsc</span>
        <MenuItem>
          <AccountMenuButton icon="settings">アカウント設定</AccountMenuButton>
        </MenuItem>
        <MenuItem>
          <AccountMenuButton icon="logout">ログアウト</AccountMenuButton>
        </MenuItem>
      </MenuItems>
    </Menu>
  );
}

function AccountMenuButton({
  icon,
  className,
  children,
  ...restProps
}: Omit<ButtonProps, "children"> & {
  readonly icon: MaterialSymbolType;
  readonly children?: ReactNode;
}) {
  return (
    <Button
      {...restProps}
      className={clsx(
        className,
        "flex items-center px-[15px] py-[10px] transition data-[focus]:bg-surface-1",
      )}
    >
      <MaterialSymbol icon={icon} size={20} />
      <span className="ml-[5px] text-14">{children}</span>
    </Button>
  );
}
