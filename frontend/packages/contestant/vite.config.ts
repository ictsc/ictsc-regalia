import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";
import { TanStackRouterVite } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [
    tsconfigPaths(),
    TanStackRouterVite({
      routesDirectory: "app/routes",
      routeFilePrefix: "~", // Route Inclusion Style: https://tanstack.com/router/latest/docs/framework/react/guide/file-based-routing#route-inclusion-example
      generatedRouteTree: "app/routes.gen.ts",
      quoteStyle: "double",
      semicolons: true,
    }),
    react(),
  ],
});
