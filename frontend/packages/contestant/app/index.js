"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var react_1 = require("react");
var client_1 = require("react-dom/client");
var react_router_1 = require("@tanstack/react-router");
var connect_web_1 = require("@connectrpc/connect-web");
var routes_gen_1 = require("./routes.gen");
var transport = (0, connect_web_1.createConnectTransport)({
    baseUrl: "/api",
});
var router = (0, react_router_1.createRouter)({
    routeTree: routes_gen_1.routeTree,
    context: {
        transport: transport,
    },
    defaultStaleTime: 1000 * 60,
});
var rootElement = document.getElementById("root");
if (rootElement != null) {
    var root = (0, client_1.createRoot)(rootElement);
    root.render(<react_1.StrictMode>
      <react_router_1.RouterProvider router={router}/>
    </react_1.StrictMode>);
}
