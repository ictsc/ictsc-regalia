"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ScheduleContext = void 0;
exports.useSchedule = useSchedule;
var react_1 = require("react");
exports.ScheduleContext = (0, react_1.createContext)(null);
function useSchedule() {
    var ctx = (0, react_1.use)(exports.ScheduleContext);
    if (ctx == null)
        return [null, false];
    var schedule = (0, react_1.use)(ctx.promise);
    if (schedule == null)
        return [null, ctx.isPending];
    return [schedule, ctx.isPending];
}
