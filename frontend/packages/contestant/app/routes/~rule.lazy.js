"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Route = void 0;
var react_1 = require("react");
var react_router_1 = require("@tanstack/react-router");
var title_1 = require("../components/title");
var markdown_1 = require("../components/markdown");
exports.Route = (0, react_router_1.createLazyFileRoute)("/rule")({
    component: RouteComponent,
});
function RouteComponent() {
    var rulePromise = exports.Route.useLoaderData().rule;
    var rule = (0, react_1.use)(rulePromise);
    return (<>
      <title_1.Title>ルール</title_1.Title>
      <div className="mx-40 mt-20 mb-64 max-w-screen-lg">
        <markdown_1.Typography>
          <markdown_1.Markdown>{rule.markdown}</markdown_1.Markdown>
        </markdown_1.Typography>
      </div>
    </>);
}
