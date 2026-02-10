import { type ReactNode } from "react";
import { clsx } from "clsx";
import {
  Menu,
  MenuItems,
  MenuItem,
  Button,
  type ButtonProps,
  MenuButton,
} from "@headlessui/react";
import {
  MaterialSymbol,
  type MaterialSymbolType,
} from "@app/components/material-symbol";

export function AccountMenu(props: {
  readonly name: string;
  readonly onSignOut?: () => void;
}) {
  return (
    <Menu>
      <MenuButton
        title="アカウントメニュー"
        className="data-[hover]:bg-surface-0/50 flex size-[50px] items-center justify-center rounded-full transition"
      >
        <MaterialSymbol
          icon="person"
          fill
          size={40}
          className="text-surface-0"
        />
      </MenuButton>

      <MenuItems
        anchor={{ to: "bottom", gap: 15 }}
        transition
        className={clsx(
          "bg-surface-0 flex w-[200px] flex-col gap-[5px] rounded-[12px] py-[15px] drop-shadow",
          "transition duration-200 ease-out data-[closed]:opacity-0",
        )}
      >
        <span className="text-14 text-text mx-[15px]">{props.name}</span>
        {/* <MenuItem>
          <AccountMenuButton icon="settings">アカウント設定</AccountMenuButton>
        </MenuItem> */}
        <MenuItem>
          <AccountMenuButton icon="logout" onClick={props.onSignOut}>
            ログアウト
          </AccountMenuButton>
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
        "data-[focus]:bg-surface-1 flex items-center px-[15px] py-[10px] transition",
      )}
    >
      <MaterialSymbol icon={icon} size={20} />
      <span className="text-14 ml-[5px]">{children}</span>
    </Button>
  );
}
