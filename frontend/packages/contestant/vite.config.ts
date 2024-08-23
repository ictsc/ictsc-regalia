import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { TanStackRouterVite } from "@tanstack/router-plugin/vite";

export default defineConfig({
  plugins: [
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
