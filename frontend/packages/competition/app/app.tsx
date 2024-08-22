import { useReducer } from "react";
import {
  AppShellLayout,
  HeaderView,
  NavbarView,
} from "./components/app-shell/layout";

export function App() {
  const [opened, toggle] = useReducer((o) => !o, true);
  return (
    <AppShellLayout
      header={<HeaderView />}
      navbar={<NavbarView collapsed={opened} onOpenToggleClick={toggle} />}
      navbarCollapsed={!opened}
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
    </AppShellLayout>
  );
}
