import { useState, type ReactNode } from "react";
import { clsx } from "clsx";
import { NavbarLayoutContext } from "./context";

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
          "duration-75 [--content-height:calc(100vh-70px)] [--header-height:70px]",
          !navbarEnabled && "[--content-width:100vw] [--navbar-width:0]",
          navbarEnabled && [
            navbarCollapsed &&
              "[--content-width:calc(100vw-50px)] [--navbar-width:50px]",
            !navbarCollapsed &&
              "[--content-width:calc(100vw-220px)] [--navbar-width:220px]",
          ],
        )}
        // React19 からサポートされる
        // eslint-disable-next-line react/no-unknown-property
        onTransitionStart={() => setNavbarTransitioning(true)}
        onTransitionEnd={() => setNavbarTransitioning(false)}
      >
        <header className="fixed left-0 right-0 top-0 z-10 h-[--header-height] bg-surface-0">
          {header}
        </header>
        <nav
          className={clsx(
            "fixed bottom-0 left-0 top-[--header-height] w-[--navbar-width]",
            navbarTransitioning && "motion-safe:transition-[width]",
          )}
        >
          {navbar}
        </nav>
        <main
          className={clsx(
            "fixed bottom-0 right-0 top-[--header-height] h-[--content-height] w-[--content-width] overflow-y-auto overflow-x-clip [scroll-gutter:stable]",
            navbarTransitioning && "motion-safe:transition-[width]",
          )}
        >
          {children}
        </main>
      </div>
    </NavbarLayoutContext>
  );
}
