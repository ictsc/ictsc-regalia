"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Route = void 0;
var react_router_1 = require("@tanstack/react-router");
var schedule_1 = require("../features/schedule");
var react_1 = require("react");
exports.Route = (0, react_router_1.createFileRoute)("/problems")({
    component: RouteComponent,
});
function RouteComponent() {
    var _a = (0, schedule_1.useSchedule)(), schedule = _a[0], isPending = _a[1];
    var navigate = (0, react_router_1.useNavigate)();
    (0, react_1.useEffect)(function () {
        if (schedule == null || isPending) {
            return;
        }
        if (!(0, schedule_1.hasContestStarted)(schedule)) {
            (0, react_1.startTransition)(function () { return navigate({ to: "/" }); });
        }
    }, [schedule, isPending, navigate]);
    return <react_router_1.Outlet />;
}
