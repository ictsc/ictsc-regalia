"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Route = void 0;
var react_router_1 = require("@tanstack/react-router");
var ranking_1 = require("@app/features/ranking");
var react_1 = require("react");
var ranking_page_1 = require("./ranking.page");
var wkt_1 = require("@bufbuild/protobuf/wkt");
exports.Route = (0, react_router_1.createFileRoute)("/ranking")({
    component: RouteComponent,
    loader: function (_a) {
        var transport = _a.context.transport;
        return {
            ranking: (0, ranking_1.fetchRanking)(transport),
        };
    },
});
function RouteComponent() {
    var rankingPromise = exports.Route.useLoaderData().ranking;
    var deferredRankingPromise = (0, react_1.useDeferredValue)(rankingPromise);
    var ranking = (0, react_1.use)(deferredRankingPromise);
    return (<ranking_page_1.RankingPage ranking={ranking.map(function (rank) { return ({
            rank: Number(rank.rank),
            teamName: rank.teamName,
            organization: rank.organization,
            score: Number(rank.score),
            lastEffectiveSubmitAt: rank.timestamp != null
                ? (0, wkt_1.timestampDate)(rank.timestamp).toISOString()
                : undefined,
        }); })}/>);
}
