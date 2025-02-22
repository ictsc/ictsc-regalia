import { type ReactNode } from "react";
import { clsx } from "clsx";

export function Layout({
  children,
  header,
  navbar,
  navbarCollapsed,
}: {
  readonly children?: ReactNode;
  readonly header: ReactNode;
  readonly navbar?: ReactNode;
  readonly navbarCollapsed: boolean;
}) {
  const navbarEnabled = navbar != null;
  return (
    <div
      className={clsx(
        "grid h-screen w-screen grid-rows-[70px_1fr] duration-75 [--content-height:calc(100vh-70px)] [--header-height:70px] motion-safe:transition-[grid-template-columns]",
        navbarEnabled
          ? navbarCollapsed
            ? "grid-cols-[50px_1fr] [--content-width:calc(100vw-50px)]"
            : "grid-cols-[220px_1fr] [--content-width:calc(100vw-220px)]"
          : "grid-cols-[0_1fr] [--content-width:100vw]",
      )}
    >
      <header className="sticky top-0 col-span-full row-start-1 row-end-2">
        {header}
      </header>
      <nav className="sticky top-[--header-height] col-start-1 col-end-2 row-start-2 row-end-3 h-[--content-height]">
        {navbar}
      </nav>
      <main className="col-start-2 col-end-3 row-start-2 row-end-3 overflow-y-auto overflow-x-clip">
        {children}
      </main>
    </div>
  );
}
