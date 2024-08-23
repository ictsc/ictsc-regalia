import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { TanStackRouterVite } from "@tanstack/router-plugin/vite";

export default defineConfig({
  plugins: [
    TanStackRouterVite({
      routesDirectory: "app/routes",
      generatedRouteTree: "app/routes.gen.ts",
      quoteStyle: "double",
      semicolons: true,
    }),
    react(),
  ],
});
