import { Button } from "@headlessui/react";
import { clsx } from "clsx";
import { MaterialSymbol, type MaterialSymbolType } from "../material-symbol";

export function NavbarView({
  collapsed = false,
  onOpenToggleClick: handleOpenToggleClick,
}: {
  readonly collapsed: boolean;
  readonly onOpenToggleClick?: () => void;
}) {
  return (
    <div className="flex size-full flex-col items-start bg-surface-1 text-text">
      <NavbarButton
        icon="list"
        showTitle={false}
        title={collapsed ? "開く" : "閉じる"}
        onClick={handleOpenToggleClick}
      />
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
      <NavbarButton showTitle={!collapsed} icon="lan" title="接続情報" />
      <NavbarButton showTitle={!collapsed} icon="help" title="問題" />
      <NavbarButton showTitle={!collapsed} icon="trophy" title="ランキング" />
      <NavbarButton showTitle={!collapsed} icon="groups" title="チーム一覧" />
      <NavbarButton showTitle={!collapsed} icon="chat" title="お問い合わせ" />
    </div>
  );
}

function NavbarButton({
  icon,
  showTitle = true,
  title,
  ...props
}: {
  icon: MaterialSymbolType;
  showTitle?: boolean;
  title: string;

  component?: "button";
  onClick?: React.MouseEventHandler;
}) {
  return (
    <Button
      className={clsx(
        "flex flex-row items-center rounded-[10px] bg-surface-1 text-text transition data-[hover]:bg-surface-2",
        showTitle && "w-full",
      )}
      title={title}
      {...props}
    >
      <div className="flex size-[50px] shrink-0 items-center justify-center">
        <MaterialSymbol icon={icon} size={24} />
      </div>
      {showTitle && (
        <span className="line-clamp-1 overflow-x-hidden text-left text-16">
          {title}
        </span>
      )}
    </Button>
  );
}
