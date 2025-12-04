/// <reference types="vitest/config" />
import type { UserConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";
import { TanStackRouterVite } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";
// import { visualizer } from "rollup-plugin-visualizer"

const isTest = process.env.NODE_ENV === "test";

export default {
  server: {
    port: 3000,
    strictPort: true,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
  plugins: [
    !isTest &&
      TanStackRouterVite({
        routesDirectory: "app/routes",
        routeFilePrefix: "~", // Route Inclusion Style: https://tanstack.com/router/latest/docs/framework/react/guide/file-based-routing#route-inclusion-example
        generatedRouteTree: "app/routes.gen.ts",
        quoteStyle: "double",
        semicolons: true,
      }),
    tsconfigPaths(),
    react(),
  ],
  build: {
    rollupOptions: {
      plugins: [
        // visualizer({ open: true }),
      ],
      output: {
        manualChunks(id) {
          if (id.includes("node_modules") && id.includes("react")) {
            return "react";
          }
        },
      },
    },
  },
} satisfies UserConfig;
