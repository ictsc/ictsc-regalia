"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Layout = Layout;
var react_1 = require("react");
var clsx_1 = require("clsx");
var context_1 = require("./context");
function Layout(_a) {
    var children = _a.children, header = _a.header, navbar = _a.navbar, navbarCollapsed = _a.navbarCollapsed;
    var navbarEnabled = navbar != null;
    var _b = (0, react_1.useState)(false), navbarTransitioning = _b[0], setNavbarTransitioning = _b[1];
    return (<context_1.NavbarLayoutContext value={{ navbarTransitioning: navbarTransitioning }}>
      <div className={(0, clsx_1.clsx)("duration-75 [--content-height:calc(100vh-70px)] [--header-height:70px]", !navbarEnabled && "[--content-width:100vw] [--navbar-width:0]", navbarEnabled && [
            navbarCollapsed &&
                "[--content-width:calc(100vw-50px)] [--navbar-width:50px]",
            !navbarCollapsed &&
                "[--content-width:calc(100vw-220px)] [--navbar-width:220px]",
        ])} 
    // React19 からサポートされる
    // eslint-disable-next-line react/no-unknown-property
    onTransitionStart={function () { return setNavbarTransitioning(true); }} onTransitionEnd={function () { return setNavbarTransitioning(false); }}>
        <header className="bg-surface-0 fixed top-0 right-0 left-0 z-10 h-(--header-height)">
          {header}
        </header>
        <nav className={(0, clsx_1.clsx)("fixed top-(--header-height) bottom-0 left-0 w-(--navbar-width)", navbarTransitioning && "motion-safe:transition-[width]")}>
          {navbar}
        </nav>
        <main className={(0, clsx_1.clsx)("fixed top-(--header-height) right-0 bottom-0 h-(--content-height) w-(--content-width) overflow-x-clip overflow-y-auto [scroll-gutter:stable]", navbarTransitioning && "motion-safe:transition-[width]")}>
          {children}
        </main>
      </div>
    </context_1.NavbarLayoutContext>);
}
