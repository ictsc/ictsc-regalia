import { createContext, useState, type ReactNode } from "react";
import { clsx } from "clsx";

export const NavbarLayoutContext = createContext<{
  navbarTransitioning: boolean;
}>({
  navbarTransitioning: false,
});

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
  const [navbarTransitioning, setNavbarTransitioning] = useState(false);
  return (
    <NavbarLayoutContext value={{ navbarTransitioning }}>
      <div
        className={clsx(
          "grid h-screen w-screen grid-rows-[70px_1fr] duration-75 [--content-height:calc(100vh-70px)] [--header-height:70px] motion-safe:transition-[grid-template-columns]",
          !navbarEnabled && "grid-cols-[0_1fr] [--content-width:100vw]",
          navbarEnabled && [
            navbarCollapsed &&
              "grid-cols-[50px_1fr] [--content-width:calc(100vw-50px)]",
            !navbarCollapsed &&
              "grid-cols-[220px_1fr] [--content-width:calc(100vw-220px)]",
          ],
        )}
        // React19 からサポートされる
        // eslint-disable-next-line react/no-unknown-property
        onTransitionStart={() => setNavbarTransitioning(true)}
        onTransitionEnd={() => setNavbarTransitioning(false)}
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
    </NavbarLayoutContext>
  );
}
