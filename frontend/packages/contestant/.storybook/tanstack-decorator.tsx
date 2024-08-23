import type { Decorator } from "@storybook/react";
import { createContext, use } from "react";
import {
  RouterProvider,
  createMemoryHistory,
  createRootRoute,
  createRoute,
  createRouter,
  useRouterState,
} from "@tanstack/react-router";

/* eslint-disable react-refresh/only-export-components */

const StoryContext = createContext<Parameters<Decorator>[0] | undefined>(
  undefined,
);

function NotFound() {
  const state = useRouterState();
  return <div>Simulated route not found for path: {state.location.href}</div>;
}

function RoutedStory() {
  const Story = use(StoryContext);
  if (Story == null) {
    throw new Error("story not found");
  }
  return <Story />;
}

const storyPath = "/__story__";
const rootRoute = createRootRoute({
  notFoundComponent: NotFound,
});
rootRoute.addChildren([
  createRoute({
    path: storyPath,
    getParentRoute: () => rootRoute,
    component: RoutedStory,
  }),
]);

const router = createRouter({
  history: createMemoryHistory({ initialEntries: [storyPath] }),
  routeTree: rootRoute,
});

export const withTanstack: Decorator = (storyFn) => (
  <StoryContext.Provider value={storyFn}>
    {/* 嘘のルートなので any で合わせる */}
    {/* eslint-disable-next-line @typescript-eslint/no-explicit-any,@typescript-eslint/no-unsafe-assignment */}
    <RouterProvider router={router as any} />
  </StoryContext.Provider>
);
