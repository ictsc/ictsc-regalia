"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.groupProblems = groupProblems;
var wkt_1 = require("@bufbuild/protobuf/wkt");
function getTemporalStatus(entry, now) {
    if (entry.endAt != null && now >= (0, wkt_1.timestampDate)(entry.endAt))
        return "past";
    if (entry.startAt != null && now >= (0, wkt_1.timestampDate)(entry.startAt))
        return "current";
    return "future";
}
function startAtMs(entry) {
    return entry.startAt != null ? (0, wkt_1.timestampDate)(entry.startAt).getTime() : 0;
}
function groupProblems(problems, now) {
    var _a;
    if (now === void 0) { now = new Date(); }
    var groupMap = new Map();
    for (var _i = 0, problems_1 = problems; _i < problems_1.length; _i++) {
        var problem = problems_1[_i];
        var schedules = problem.submissionableSchedules;
        var key = schedules
            .map(function (s) { return s.name; })
            .sort()
            .join(",");
        var group = groupMap.get(key);
        if (group == null) {
            group = {
                problems: [],
                entries: schedules.slice().sort(function (a, b) { return startAtMs(a) - startAtMs(b); }),
                hasSubmittableProblem: false,
            };
            groupMap.set(key, group);
        }
        if ((_a = problem.submissionStatus) === null || _a === void 0 ? void 0 : _a.isSubmittable) {
            group.hasSubmittableProblem = true;
        }
        group.problems.push(problem);
    }
    return Array.from(groupMap.entries())
        .map(function (_a) {
        var key = _a[0], _b = _a[1], problems = _b.problems, entries = _b.entries, hasSubmittableProblem = _b.hasSubmittableProblem;
        return ({
            key: key,
            schedules: entries.map(function (e) { return ({
                name: e.name,
                temporalStatus: getTemporalStatus(e, now),
            }); }),
            problems: problems.slice().sort(function (a, b) { return a.code.localeCompare(b.code); }),
            hasSubmittableProblem: hasSubmittableProblem,
            _sortKey: entries.length > 0 ? startAtMs(entries[0]) : Infinity,
        });
    })
        .sort(function (a, b) {
        if (a.hasSubmittableProblem !== b.hasSubmittableProblem) {
            return a.hasSubmittableProblem ? -1 : 1;
        }
        return a._sortKey - b._sortKey;
    });
}
