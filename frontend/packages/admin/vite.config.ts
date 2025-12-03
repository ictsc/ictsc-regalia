/// <reference types="vitest/config" />
import type { UserConfig } from "vite";
import { TanStackRouterVite } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";

const isTest = process.env.NODE_ENV === "test";

export default {
  server: {
    port: 3001,
    proxy: {
      "/api": {
        target: "http://localhost:8081",
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
    react(),
  ],
} satisfies UserConfig;
