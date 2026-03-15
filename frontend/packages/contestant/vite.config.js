"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var vite_1 = require("@tailwindcss/vite");
var vite_tsconfig_paths_1 = require("vite-tsconfig-paths");
var vite_2 = require("@tanstack/router-plugin/vite");
var plugin_react_1 = require("@vitejs/plugin-react");
// import { visualizer } from "rollup-plugin-visualizer"
var isTest = process.env.NODE_ENV === "test";
exports.default = {
    server: {
        port: 3000,
        strictPort: true,
        proxy: {
            "/api": {
                target: "http://localhost:8080",
                changeOrigin: true,
                rewrite: function (path) { return path.replace(/^\/api/, ""); },
            },
            "/auth": {
                target: "http://localhost:8080",
                changeOrigin: true,
            },
        },
    },
    plugins: [
        (0, vite_1.default)(),
        !isTest &&
            (0, vite_2.TanStackRouterVite)({
                routesDirectory: "app/routes",
                routeFilePrefix: "~", // Route Inclusion Style: https://tanstack.com/router/latest/docs/framework/react/guide/file-based-routing#route-inclusion-example
                generatedRouteTree: "app/routes.gen.ts",
                quoteStyle: "double",
                semicolons: true,
            }),
        (0, vite_tsconfig_paths_1.default)(),
        (0, plugin_react_1.default)(),
    ],
    build: {
        rollupOptions: {
            plugins: [
            // visualizer({ open: true }),
            ],
            output: {
                manualChunks: function (id) {
                    if (id.includes("node_modules") && id.includes("react")) {
                        return "react";
                    }
                },
            },
        },
    },
};
