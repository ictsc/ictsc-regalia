"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ScheduleProvider = ScheduleProvider;
var react_1 = require("react");
var feature_1 = require("./feature");
var use_schedule_1 = require("./use-schedule");
function ScheduleProvider(props) {
    var deferredInitialData = (0, react_1.useDeferredValue)(props.initialData);
    var _a = (0, react_1.useState)({
        data: props.initialData,
        base: props.initialData,
    }), promiseState = _a[0], setPromise = _a[1];
    var _b = (0, react_1.useTransition)(), isStatePending = _b[0], startTransision = _b[1];
    var _c = promiseState.base === props.initialData
        ? [promiseState.data, isStatePending]
        : [deferredInitialData, deferredInitialData !== props.initialData], promise = _c[0], isPending = _c[1];
    return (<>
      <use_schedule_1.ScheduleContext value={{ promise: promise, isPending: isPending }}>
        {props.children}
      </use_schedule_1.ScheduleContext>
      <react_1.Suspense>
        <ScheduleReloader schedule={promise} load={function () {
            if (isPending)
                return;
            startTransision(function () {
                setPromise({
                    data: props.loadData(),
                    base: props.initialData,
                });
            });
        }}/>
      </react_1.Suspense>
    </>);
}
function ScheduleReloader(props) {
    var onLoad = (0, react_1.useEffectEvent)(props.load);
    var schedule = (0, react_1.use)(props.schedule);
    (0, react_1.useEffect)(function () {
        var reloadAt = (0, feature_1.nextReloadAt)(schedule);
        if (reloadAt == null)
            return;
        var delayMs = reloadAt.getTime() - Date.now();
        if (delayMs <= 0) {
            onLoad();
            return;
        }
        var timer = window.setTimeout(function () {
            onLoad();
        }, delayMs);
        return function () { return window.clearTimeout(timer); };
    }, [schedule]);
    return null;
}
