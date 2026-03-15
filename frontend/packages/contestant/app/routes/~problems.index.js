"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Route = void 0;
var react_router_1 = require("@tanstack/react-router");
var problem_1 = require("@app/features/problem");
var announce_1 = require("@app/features/announce");
var react_1 = require("react");
var wkt_1 = require("@bufbuild/protobuf/wkt");
var page_1 = require("./problems.index/page");
exports.Route = (0, react_router_1.createFileRoute)("/problems/")({
    component: RouteComponent,
    loader: function (_a) {
        var transport = _a.context.transport;
        return {
            problems: (0, problem_1.fetchProblems)(transport),
            notices: (0, announce_1.fetchNotices)(transport),
        };
    },
});
function RouteComponent() {
    var router = (0, react_router_1.useRouter)();
    var _a = exports.Route.useLoaderData(), problemsPromise = _a.problems, noticesPromise = _a.notices;
    var deferredProblemsPromise = (0, react_1.useDeferredValue)(problemsPromise);
    var problems = (0, react_1.use)(deferredProblemsPromise);
    var notices = (0, react_1.use)(noticesPromise);
    // 最も近い submittableUntil/submittableFrom に達したらリフェッチ
    (0, react_1.useEffect)(function () {
        var _a, _b;
        var earliest = null;
        for (var _i = 0, problems_1 = problems; _i < problems_1.length; _i++) {
            var p = problems_1[_i];
            var until = (_a = p.submissionStatus) === null || _a === void 0 ? void 0 : _a.submittableUntil;
            var from = (_b = p.submissionStatus) === null || _b === void 0 ? void 0 : _b.submittableFrom;
            var target = until !== null && until !== void 0 ? until : from;
            if (target == null)
                continue;
            var ms = (0, wkt_1.timestampDate)(target).getTime() - Date.now();
            if (ms <= 0)
                continue;
            if (earliest == null || ms < earliest)
                earliest = ms;
        }
        if (earliest == null)
            return;
        var timer = setTimeout(function () {
            (0, react_1.startTransition)(function () { return router.invalidate(); });
        }, earliest);
        return function () { return clearTimeout(timer); };
    }, [problems, router]);
    return <page_1.ProblemsPage problems={problems} notices={notices}/>;
}
