import { Fragment } from "react";
import { clsx } from "clsx";
import { Button } from "@headlessui/react";
import { Link, useRouterState } from "@tanstack/react-router";
import {
  MaterialSymbol,
  type MaterialSymbolType,
} from "../../components/material-symbol";

export function Navbar(props: {
  readonly collapsed: boolean;
  readonly canViewProblems: boolean;
  readonly canViewAnnounces: boolean;
  readonly onOpenToggleClick?: () => void;
}) {
  const { collapsed } = props;
  const state = useRouterState();
  return (
    <div className="bg-surface-1 text-text flex size-full flex-col items-start gap-4">
      <Button as={Fragment}>
        {(buttonProps) => (
          <button
            title={collapsed ? "開く" : "閉じる"}
            onClick={props.onOpenToggleClick}
            className={navbarButtonClassName({
              collapsed: true,
              ...buttonProps,
            })}
          >
            <NavbarButtonInner collapsed icon="list" />
          </button>
        )}
      </Button>
      <Button as={Fragment}>
        {(buttonProps) => (
          <Link
            to="/rule"
            title="ルール"
            className={navbarButtonClassName({
              collapsed,
              matched: state.location.pathname?.startsWith("/rule"),
              ...buttonProps,
            })}
          >
            <NavbarButtonInner
              collapsed={collapsed}
              icon="developer_guide"
              title="ルール"
            />
          </Link>
        )}
      </Button>
      <Button as={Fragment} disabled={!props.canViewAnnounces}>
        {(buttonProps) => (
          <Link
            to="/announces"
            title={props.canViewAnnounces ? "アナウンス" : "開催期間外です"}
            className={navbarButtonClassName({
              collapsed,
              matched: state.location.pathname?.startsWith("/announces"),
              ...buttonProps,
            })}
          >
            <NavbarButtonInner
              collapsed={collapsed}
              icon="brand_awareness"
              title="アナウンス"
            />
          </Link>
        )}
      </Button>
      {/* <NavbarButton showTitle={!collapsed} icon="lan" title="接続情報" /> */}
      <Button as={Fragment} disabled={!props.canViewProblems}>
        {(buttonProps) => (
          <Link
            to="/problems"
            title={props.canViewProblems ? "問題" : "開催期間外です"}
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
      <Button as={Fragment}>
        {(buttonProps) => (
          <Link
            to="/teams"
            title="チーム一覧"
            className={navbarButtonClassName({
              collapsed,
              matched: state.location.pathname?.startsWith("/teams"),
              ...buttonProps,
            })}
          >
            <NavbarButtonInner
              collapsed={collapsed}
              icon="groups"
              title="チーム一覧"
            />
          </Link>
        )}
      </Button>
      <Button as={Fragment}>
        {(buttonProps) => (
          <Link
            to="/ranking"
            title="ランキング"
            className={navbarButtonClassName({
              collapsed,
              matched: state.location.pathname?.startsWith("/ranking"),
              ...buttonProps,
            })}
          >
            <NavbarButtonInner
              collapsed={collapsed}
              icon="trophy"
              title="ランキング"
            />
          </Link>
        )}
      </Button>
      {/* <NavbarButtonInner
        collapsed={collapsed}
        icon="groups"
        title="チーム一覧"
      /> */}
      {/* <NavbarButton showTitle={!collapsed} icon="chat" title="お問い合わせ" /> */}
    </div>
  );
}

function navbarButtonClassName({
  collapsed,
  hover,
  active,
  disabled = false,
  matched = false,
}: {
  collapsed: boolean;
  hover: boolean;
  active: boolean;
  disabled?: boolean;
  matched?: boolean;
}): string {
  return clsx(
    "bg-surface-1 flex flex-row items-center rounded-[10px] transition",
    !collapsed && "w-full",
    disabled ? "cursor-not-allowed opacity-75" : "text-text",
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
        <span className="text-16 line-clamp-1 overflow-x-hidden text-left">
          {props.title}
        </span>
      )}
    </>
  );
}
