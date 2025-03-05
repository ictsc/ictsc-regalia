import { createContext } from "react";

export const NavbarLayoutContext = createContext<{
  navbarTransitioning: boolean;
}>({
  navbarTransitioning: false,
});
