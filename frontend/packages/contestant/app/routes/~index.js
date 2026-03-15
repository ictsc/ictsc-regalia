"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Route = void 0;
var react_router_1 = require("@tanstack/react-router");
var index_page_1 = require("./index.page");
var schedule_1 = require("@app/features/schedule");
exports.Route = (0, react_router_1.createFileRoute)("/")({
    component: Page,
});
function Page() {
    var schedule = (0, schedule_1.useSchedule)()[0];
    var currentEntry = (0, schedule_1.getCurrentScheduleEntry)(schedule);
    var nextEntry = (0, schedule_1.getNextScheduleEntry)(schedule);
    var state;
    var timerEndMs;
    if (currentEntry != null) {
        // Currently in a schedule
        state = "in_contest";
        var endAt = (0, schedule_1.currentEndAt)(schedule);
        timerEndMs = endAt != null ? endAt.getTime() : 0;
    }
    else if (nextEntry != null) {
        // Waiting for next schedule
        state = "waiting";
        var startAt = (0, schedule_1.nextStartAt)(schedule);
        timerEndMs = startAt != null ? startAt.getTime() : 0;
    }
    else {
        // No current or next schedule - contest ended
        state = "ended";
        timerEndMs = 0;
    }
    return (<index_page_1.IndexPage state={state} currentScheduleName={currentEntry === null || currentEntry === void 0 ? void 0 : currentEntry.name} nextScheduleName={nextEntry === null || nextEntry === void 0 ? void 0 : nextEntry.name} timer={timerEndMs > 0 ? <index_page_1.Timer endMs={timerEndMs}/> : undefined}/>);
}
