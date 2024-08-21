import clsx from "clsx";
import {
  MaterialSymbol,
  type MaterialSymbolType,
} from "./components/material-symbol";
import { ReactNode, useState } from "react";

function NavBar() {
  return (
    <div className="flex size-full items-center border-b-[3px] border-primary bg-surface-0">
      <span className="flex-none">ICTSC</span>
      <div className="ml-auto flex h-full items-center">
        <div className="mr-[30px]">
          <div className="flex h-[50px] items-center rounded-[10px] bg-surface-1 px-[10px] text-text">
            <MaterialSymbol icon="schedule" size={24} />
            <span className="ml-[5px] text-16">競技開始前</span>
            <span className="ml-[20px]">
              <span className="text-12">残り</span>
              <span className="ml-[5px] text-24">00 : 00 : 00</span>
            </span>
          </div>
        </div>
        <div className="flex h-full w-[140px] bg-primary [clip-path:polygon(40%_0,100%_0,100%_100%,0_100%)]">
          <button className="my-auto ml-auto mr-[30px] size-[40px]">
            <MaterialSymbol
              icon="person"
              fill
              size={40}
              className="text-surface-0"
            />
          </button>
        </div>
      </div>
    </div>
  );
}

function Sidebar({
  opened = true,
  onOpenToggleClick: handleOpenToggleClick,
}: {
  readonly opened: boolean;
  readonly onOpenToggleClick: () => void;
}) {
  return (
    <div className="flex size-full flex-col items-start bg-surface-1 text-text">
      <SideBarButton
        icon="list"
        showTitle={false}
        title={opened ? "閉じる" : "開く"}
        onClick={handleOpenToggleClick}
      />
      <SideBarButton showTitle={opened} icon="developer_guide" title="ルール" />
      <SideBarButton
        showTitle={opened}
        icon="brand_awareness"
        title="アナウンス"
      />
      <SideBarButton showTitle={opened} icon="lan" title="接続情報" />
      <SideBarButton showTitle={opened} icon="help" title="問題" />
      <SideBarButton showTitle={opened} icon="trophy" title="ランキング" />
      <SideBarButton showTitle={opened} icon="groups" title="チーム一覧" />
      <SideBarButton showTitle={opened} icon="chat" title="お問い合わせ" />
    </div>
  );
}

function SideBarButton({
  icon,
  showTitle = true,
  title,
  onClick: handleClick,
}: {
  icon: MaterialSymbolType;
  showTitle?: boolean;
  title: string;

  onClick?: React.MouseEventHandler;
}) {
  return (
    <button
      className={clsx(
        "flex flex-row items-center rounded-[10px] bg-surface-1 text-text hover:bg-surface-2 motion-safe:hover:transition-colors",
        showTitle && "w-full",
      )}
      title={title}
      onClick={handleClick}
    >
      <div className="flex size-[50px] shrink-0 basis-[50px] items-center justify-center">
        <MaterialSymbol icon={icon} size={24} />
      </div>
      {showTitle && (
        <span className="line-clamp-1 overflow-x-hidden text-left text-16">
          {title}
        </span>
      )}
    </button>
  );
}

function Layout({
  children,
  navbar,
  sidebar,
  sidebarOpened,
}: {
  readonly children?: ReactNode;
  readonly navbar: ReactNode;
  readonly sidebar: ReactNode;
  readonly sidebarOpened: boolean;
}) {
  return (
    <div
      className={clsx(
        "grid h-screen w-screen grid-rows-[70px_1fr] duration-75 motion-safe:transition-[grid-template-columns]",
        sidebarOpened ? "grid-cols-[220px_1fr]" : "grid-cols-[50px_1fr]",
      )}
    >
      <header className="sticky top-0 col-span-full row-start-1 row-end-2">
        {navbar}
      </header>
      <aside className="sticky top-[70px] col-start-1 col-end-2 row-start-2 row-end-3 h-[calc(100vh-70px)]">
        {sidebar}
      </aside>
      <main className="col-start-2 col-end-3 row-start-2 row-end-3 overflow-y-scroll">
        {children}
      </main>
    </div>
  );
}

export function App() {
  const [sidebarOpened, setSidebarOpened] = useState(true);
  return (
    <Layout
      navbar={<NavBar />}
      sidebar={
        <Sidebar
          opened={true}
          onOpenToggleClick={() => setSidebarOpened((o) => !o)}
        />
      }
      sidebarOpened={sidebarOpened}
    >
      <p>
        {`Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
            eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
            ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
            aliquip ex ea commodo consequat. Duis aute irure dolor in
            reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
            pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
            culpa qui officia deserunt mollit anim id est laborum.`.repeat(
          1000,
        )}
      </p>
    </Layout>
  );
}
