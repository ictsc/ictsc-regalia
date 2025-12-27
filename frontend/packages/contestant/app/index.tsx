import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import { createConnectTransport } from "@connectrpc/connect-web";
import { routeTree } from "./routes.gen";

const transport = createConnectTransport({
  baseUrl: "/api",
});
const router = createRouter({
  routeTree,
  context: {
    transport,
  },
  defaultStaleTime: 1000 * 60,
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

const rootElement = document.getElementById("root");
if (rootElement != null) {
  const root = createRoot(rootElement);
  root.render(
    <StrictMode>
      <RouterProvider router={router} />
    </StrictMode>,
  );
}
