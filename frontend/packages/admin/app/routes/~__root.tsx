import { type ReactNode, Suspense, lazy } from "react";
import {
  Link,
  Outlet,
  createRootRouteWithContext,
} from "@tanstack/react-router";
import type { Transport } from "@connectrpc/connect";
import { AppShell, Burger, Group, NavLink, Text } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";

type RouterContext = {
  transport: Transport;
};

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
});

const TanStackRouterDevtools = import.meta.env.DEV
  ? lazy(() =>
      import("@tanstack/router-devtools").then((mod) => ({
        default: mod.TanStackRouterDevtools,
      })),
    )
  : () => null;

function RootComponent() {
  return (
    <>
      <Shell>
        <Outlet />
      </Shell>
      <Suspense>
        <TanStackRouterDevtools />
      </Suspense>
    </>
  );
}

function Shell(props: { readonly children?: ReactNode }) {
  const [mobileOpened, { toggle: toggleMobile }] = useDisclosure();
  const [desktopOpened, { toggle: toggleDesktop }] = useDisclosure(true);
  return (
    <AppShell
      header={{ height: 60 }}
      navbar={{
        width: 300,
        breakpoint: "sm",
        collapsed: { mobile: !mobileOpened, desktop: !desktopOpened },
      }}
      padding="md"
    >
      <AppShell.Header>
        <Group h="100%" px="md">
          <Burger
            opened={mobileOpened}
            onClick={toggleMobile}
            hiddenFrom="sm"
            size="sm"
          />
          <Burger
            opened={desktopOpened}
            onClick={toggleDesktop}
            visibleFrom="sm"
            size="sm"
          />
          <Text size="lg">ICTSC Admin</Text>
        </Group>
      </AppShell.Header>
      <AppShell.Navbar>
        <NavLink component={Link} label="採点" to="/submissions" />
      </AppShell.Navbar>
      <AppShell.Main>{props.children}</AppShell.Main>
    </AppShell>
  );
}
