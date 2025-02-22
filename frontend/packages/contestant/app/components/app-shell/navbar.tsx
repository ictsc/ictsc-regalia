import { Fragment } from "react";
import { clsx } from "clsx";
import { Button } from "@headlessui/react";
import { Link, useRouterState } from "@tanstack/react-router";
import {
  MaterialSymbol,
  type MaterialSymbolType,
} from "@app/components/material-symbol";

export function NavbarView({
  collapsed = false,
  onOpenToggleClick: handleOpenToggleClick,
}: {
  readonly collapsed: boolean;
  readonly onOpenToggleClick?: () => void;
}) {
  const state = useRouterState();
  return (
    <div className="flex size-full flex-col items-start gap-4 bg-surface-1 text-text">
      <Button as={Fragment}>
        {(buttonProps) => (
          <button
            title={collapsed ? "開く" : "閉じる"}
            onClick={handleOpenToggleClick}
            className={navbarButtonClassName({
              collapsed: true,
              ...buttonProps,
            })}
          >
            <NavbarButtonInner collapsed icon="list" />
          </button>
        )}
      </Button>
      <NavbarButton
        showTitle={!collapsed}
        icon="developer_guide"
        title="ルール"
      />
      <NavbarButton
        showTitle={!collapsed}
        icon="brand_awareness"
        title="アナウンス"
      />
      {/* <NavbarButton showTitle={!collapsed} icon="lan" title="接続情報" /> */}
      <Button as={Fragment}>
        {(buttonProps) => (
          <Link
            to="/problems"
            title="問題"
            className={navbarButtonClassName({
              collapsed,
              matched: state.location.pathname?.startsWith("/problems"),
              ...buttonProps,
            })}
          >
            <NavbarButtonInner collapsed={collapsed} icon="help" title="問題" />
          </Link>
        )}
      </Button>
      <NavbarButton showTitle={!collapsed} icon="trophy" title="ランキング" />
      <NavbarButton showTitle={!collapsed} icon="groups" title="チーム一覧" />
      {/* <NavbarButton showTitle={!collapsed} icon="chat" title="お問い合わせ" /> */}
    </div>
  );
}

function navbarButtonClassName({
  collapsed,
  hover,
  active,
  matched = false,
}: {
  collapsed: boolean;
  hover: boolean;
  active: boolean;
  matched?: boolean;
}): string {
  return clsx(
    "flex flex-row items-center rounded-[10px] bg-surface-1 text-text transition",
    !collapsed && "w-full",
    (hover || matched) && "bg-surface-2",
    active && "opacity-75",
  );
}

function NavbarButtonInner(props: {
  collapsed: boolean;
  icon: MaterialSymbolType;
  title?: string;
}) {
  return (
    <>
      <div className="flex size-[50px] shrink-0 items-center justify-center">
        <MaterialSymbol icon={props.icon} size={24} />
      </div>
      {!props.collapsed && (
        <span className="line-clamp-1 overflow-x-hidden text-left text-16">
          {props.title}
        </span>
      )}
    </>
  );
}

function NavbarButton({
  icon,
  showTitle = true,
  title,
}: {
  icon: MaterialSymbolType;
  showTitle?: boolean;
  title: string;
}) {
  return (
    <Button as={Fragment}>
      {(props) => (
        <button
          title={title}
          className={navbarButtonClassName({ collapsed: !showTitle, ...props })}
        >
          <NavbarButtonInner collapsed={!showTitle} icon={icon} title={title} />
        </button>
      )}
    </Button>
  );
}
