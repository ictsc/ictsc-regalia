import { createContext, useContext } from "react";
import {
  RouterProvider,
  createMemoryHistory,
  createRootRoute,
  createRoute,
  createRouter,
  useRouterState,
} from "@tanstack/react-router";

/* eslint-disable react-refresh/only-export-components */

const StoryContext = createContext(undefined);

function NotFound() {
  const state = useRouterState();
  return <div>Simulated route not found for path: {state.location.href}</div>;
}

function RoutedStory() {
  const Story = useContext(StoryContext);
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

/** @type{import("@storybook/react").Decorator} */
export const withTanstack = (storyFn) => (
  <StoryContext.Provider value={storyFn}>
    <RouterProvider router={router} />
  </StoryContext.Provider>
);
