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
        <header className="bg-surface-0 fixed top-0 right-0 left-0 z-10 h-(--header-height)">
          {header}
        </header>
        <nav
          className={clsx(
            "fixed top-(--header-height) bottom-0 left-0 w-(--navbar-width)",
            navbarTransitioning && "motion-safe:transition-[width]",
          )}
        >
          {navbar}
        </nav>
        <main
          className={clsx(
            "fixed top-(--header-height) right-0 bottom-0 h-(--content-height) w-(--content-width) overflow-x-clip overflow-y-auto [scroll-gutter:stable]",
            navbarTransitioning && "motion-safe:transition-[width]",
          )}
        >
          {children}
        </main>
      </div>
    </NavbarLayoutContext>
  );
}
