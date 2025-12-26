import { createContext, useContext } from "react";
import {
  Link,
  RouterProvider,
  createMemoryHistory,
  createRootRoute,
  createRoute,
  createRouter,
  useRouterState,
} from "@tanstack/react-router";

/* eslint-disable react-refresh/only-export-components */

const StoryContext = createContext(undefined);
const storyPath = "/__story__";

function NotFound() {
  const state = useRouterState();
  return (
    <div>
      <p>Simulated route not found for path: {state.location.href}</p>
      <Link to={storyPath}>Back to story</Link>
    </div>
  );
}

function RoutedStory() {
  const Story = useContext(StoryContext);
  if (Story == null) {
    throw new Error("story not found");
  }
  // eslint-disable-next-line react-hooks/static-components
  return <Story />;
}

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
