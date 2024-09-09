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
  readonly navbar: ReactNode;
  readonly navbarCollapsed: boolean;
}) {
  return (
    <div
      className={clsx(
        "grid h-screen w-screen grid-rows-[70px_1fr] duration-75 motion-safe:transition-[grid-template-columns]",
        navbarCollapsed ? "grid-cols-[50px_1fr]" : "grid-cols-[220px_1fr]",
      )}
    >
      <header className="sticky top-0 col-span-full row-start-1 row-end-2">
        {header}
      </header>
      <nav className="sticky top-[70px] col-start-1 col-end-2 row-start-2 row-end-3 h-[calc(100vh-70px)]">
        {navbar}
      </nav>
      <main className="col-start-2 col-end-3 row-start-2 row-end-3 overflow-y-auto">
        {children}
      </main>
    </div>
  );
}
