"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Route = void 0;
var react_router_1 = require("@tanstack/react-router");
var teams_1 = require("@app/features/teams");
var react_1 = require("react");
var teams_page_1 = require("./teams.page");
exports.Route = (0, react_router_1.createFileRoute)("/teams")({
    component: RouteComponent,
    loader: function (_a) {
        var transport = _a.context.transport;
        return {
            teams: (0, teams_1.fetchTeams)(transport),
        };
    },
});
function RouteComponent() {
    var teamsPromise = exports.Route.useLoaderData().teams;
    var deferredTeamsPromise = (0, react_1.useDeferredValue)(teamsPromise);
    var teams = (0, react_1.use)(deferredTeamsPromise);
    return <teams_page_1.TeamsPage teamProfile={teams}/>;
}
