"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var protobuf_1 = require("@bufbuild/protobuf");
var wkt_1 = require("@bufbuild/protobuf/wkt");
var v1_1 = require("@ictsc/proto/contestant/v1");
var vitest_1 = require("vitest");
var feature_1 = require("./feature");
function scheduleEntry(name, startAt, endAt) {
    return (0, protobuf_1.create)(v1_1.ScheduleEntrySchema, {
        name: name,
        startAt: (0, wkt_1.timestampFromDate)(startAt),
        endAt: (0, wkt_1.timestampFromDate)(endAt),
    });
}
function schedule(data) {
    return (0, protobuf_1.create)(v1_1.ScheduleSchema, {
        hasStarted: data.current != null,
        current: data.current,
        next: data.next,
    });
}
(0, vitest_1.describe)("nextReloadAt", function () {
    (0, vitest_1.it)("returns the current schedule end while a slot is active", function () {
        var endAt = new Date("2026-03-11T10:30:00.000Z");
        (0, vitest_1.expect)((0, feature_1.nextReloadAt)(schedule({
            current: scheduleEntry("day1-am", new Date("2026-03-11T09:00:00.000Z"), endAt),
            next: scheduleEntry("day1-pm", new Date("2026-03-11T13:00:00.000Z"), new Date("2026-03-11T15:00:00.000Z")),
        }))).toEqual(endAt);
    });
    (0, vitest_1.it)("returns the next schedule start while waiting for the next slot", function () {
        var startAt = new Date("2026-03-11T13:00:00.000Z");
        (0, vitest_1.expect)((0, feature_1.nextReloadAt)(schedule({
            next: scheduleEntry("day1-pm", startAt, new Date("2026-03-11T15:00:00.000Z")),
        }))).toEqual(startAt);
    });
    (0, vitest_1.it)("returns null after all schedule slots have ended", function () {
        (0, vitest_1.expect)((0, feature_1.nextReloadAt)(schedule({}))).toBeNull();
    });
    (0, vitest_1.it)("returns null when schedule loading failed", function () {
        (0, vitest_1.expect)((0, feature_1.nextReloadAt)(null)).toBeNull();
    });
});
