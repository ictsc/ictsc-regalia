import clsx from "clsx";
import { IconTypes, MaterialSymbol } from "./components/MaterialSymbol";
import { ReactNode } from "react";

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

function SideBar() {
  return (
    <div className="flex size-full flex-col items-start bg-surface-1 text-text">
      <SideBarButton icon="list" showTitle={false} title="閉じる" />
      <SideBarButton icon="developer_guide" title="ルール" />
      <SideBarButton icon="brand_awareness" title="アナウンス" />
      <SideBarButton icon="lan" title="接続情報" />
      <SideBarButton icon="help" title="問題" />
      <SideBarButton icon="trophy" title="ランキング" />
      <SideBarButton icon="groups" title="チーム一覧" />
      <SideBarButton icon="chat" title="お問い合わせ" />
    </div>
  );
}

function SideBarButton({
  icon,
  showTitle = true,
  title,
}: {
  icon: IconTypes;
  showTitle?: boolean;
  title: string;
}) {
  return (
    <button
      className={clsx(
        "group flex w-full items-center text-text",
        showTitle &&
          "rounded-[10px] bg-surface-1 hover:bg-surface-2 motion-safe:hover:transition-all",
      )}
      title={title}
    >
      <div
        className={clsx(
          "flex size-[50px] items-center justify-center",
          !showTitle &&
            "rounded-[10px] bg-surface-1 group-hover:bg-surface-2 motion-safe:group-hover:transition-all",
        )}
      >
        <MaterialSymbol icon={icon} size={24} />
      </div>
      {showTitle && <span className="text-16">{title}</span>}
    </button>
  );
}

export function App() {
  return (
    <div className="grid grid-cols-[220px_1fr] grid-rows-[70px_1fr]">
      <header className="sticky top-0 col-span-full row-span-1">
        <NavBar />
      </header>
      <aside className="sticky top-[70px] col-span-1 col-start-1 row-span-1 row-start-2 h-[calc(100vh-70px)]">
        <SideBar />
      </aside>
      <main className="col-span-1 col-start-2 row-span-1 row-start-2">
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
      </main>
    </div>
  );
}
