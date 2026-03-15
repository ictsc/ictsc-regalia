"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Route = void 0;
var react_1 = require("react");
var react_router_1 = require("@tanstack/react-router");
var announces_index_page_1 = require("./announces.index.page");
var _announces_1 = require("./~announces");
exports.Route = (0, react_router_1.createFileRoute)("/announces/")({
    component: RouteComponent,
});
function RouteComponent() {
    var AnnouncePromise = _announces_1.Route.useLoaderData().announces;
    var deferredAnnouncePromise = (0, react_1.useDeferredValue)(AnnouncePromise);
    var announces = (0, react_1.use)(deferredAnnouncePromise);
    return <announces_index_page_1.AnnounceList announces={announces}/>;
}
