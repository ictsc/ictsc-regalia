"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
var __spreadArray = (this && this.__spreadArray) || function (to, from, pack) {
    if (pack || arguments.length === 2) for (var i = 0, l = from.length, ar; i < l; i++) {
        if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.NoSchedules = exports.SingleSchedule = exports.GroupedBySchedule = exports.Default = void 0;
var protobuf_1 = require("@bufbuild/protobuf");
var wkt_1 = require("@bufbuild/protobuf/wkt");
var v1_1 = require("@ictsc/proto/contestant/v1");
var page_1 = require("./page");
exports.default = {
    title: "pages/problems",
    component: page_1.ProblemsPage,
};
var day1Am = (0, protobuf_1.create)(v1_1.ScheduleEntrySchema, {
    name: "day1-am",
    startAt: (0, wkt_1.timestampFromDate)(new Date("2026-01-01T10:00:00+09:00")),
    endAt: (0, wkt_1.timestampFromDate)(new Date("2026-01-01T12:00:00+09:00")),
});
var day1Pm = (0, protobuf_1.create)(v1_1.ScheduleEntrySchema, {
    name: "day1-pm",
    startAt: (0, wkt_1.timestampFromDate)(new Date("2026-01-01T13:00:00+09:00")),
    endAt: (0, wkt_1.timestampFromDate)(new Date("2099-12-31T23:59:59+09:00")),
});
var day2Am = (0, protobuf_1.create)(v1_1.ScheduleEntrySchema, {
    name: "day2-am",
    startAt: (0, wkt_1.timestampFromDate)(new Date("2100-01-01T10:00:00+09:00")),
    endAt: (0, wkt_1.timestampFromDate)(new Date("2100-01-01T12:00:00+09:00")),
});
var day2Pm = (0, protobuf_1.create)(v1_1.ScheduleEntrySchema, {
    name: "day2-pm",
    startAt: (0, wkt_1.timestampFromDate)(new Date("2100-01-01T13:00:00+09:00")),
    endAt: (0, wkt_1.timestampFromDate)(new Date("2100-01-01T16:00:00+09:00")),
});
var submittable = (0, protobuf_1.create)(v1_1.SubmissionStatusSchema, {
    isSubmittable: true,
});
var notSubmittable = (0, protobuf_1.create)(v1_1.SubmissionStatusSchema, {
    isSubmittable: false,
});
function makeProblem(code, title, opts) {
    var _a;
    return {
        $typeName: "contestant.v1.Problem",
        code: code,
        title: title,
        maxScore: 200,
        category: "Network",
        score: (opts === null || opts === void 0 ? void 0 : opts.score)
            ? __assign({ $typeName: "contestant.v1.Score", maxScore: 200 }, opts.score) : undefined,
        submissionStatus: opts === null || opts === void 0 ? void 0 : opts.submissionStatus,
        submissionableSchedules: (_a = opts === null || opts === void 0 ? void 0 : opts.submissionableSchedules) !== null && _a !== void 0 ? _a : [],
    };
}
exports.Default = {
    args: {
        notices: [],
        problems: __spreadArray(__spreadArray(__spreadArray(__spreadArray([], Array.from({ length: 4 }, function () { return ({
            $typeName: "contestant.v1.Problem",
            code: "ABC",
            title: "問題 ABC",
            maxScore: 200,
            category: "Network",
            submissionableSchedules: [],
        }); }), true), Array.from({ length: 5 }, function () { return ({
            $typeName: "contestant.v1.Problem",
            code: "ABC",
            title: "問題 ABC",
            maxScore: 200,
            category: "Server",
            submissionableSchedules: [],
            score: {
                $typeName: "contestant.v1.Score",
                score: 100,
                markedScore: 120,
                penalty: -20,
            },
        }); }), true), Array.from({ length: 2 }, function () { return ({
            $typeName: "contestant.v1.Problem",
            code: "ABC",
            title: "問題 ABC",
            maxScore: 200,
            category: "Network",
            submissionableSchedules: [],
            score: {
                $typeName: "contestant.v1.Score",
                score: 160,
                markedScore: 200,
                penalty: -40,
            },
        }); }), true), Array.from({ length: 5 }, function () { return ({
            $typeName: "contestant.v1.Problem",
            code: "ABC",
            title: "問題 ABC",
            maxScore: 200,
            category: "Server",
            submissionableSchedules: [],
            score: {
                $typeName: "contestant.v1.Score",
                score: 200,
                markedScore: 200,
                penalty: 0,
            },
        }); }), true),
    },
};
exports.GroupedBySchedule = {
    args: {
        notices: [],
        problems: [
            // day1-am（過去 → 提出不可）
            makeProblem("NET01", "ネットワーク基礎問題", {
                submissionStatus: notSubmittable,
                submissionableSchedules: [day1Am],
                score: { score: 200, markedScore: 200, penalty: 0 },
            }),
            makeProblem("NET02", "VLAN設定問題", {
                submissionStatus: notSubmittable,
                submissionableSchedules: [day1Am],
                score: { score: 100, markedScore: 120, penalty: -20 },
            }),
            makeProblem("SRV01", "Webサーバー構築", {
                submissionStatus: notSubmittable,
                submissionableSchedules: [day1Am],
            }),
            // day1-pm（現在 → 提出可能）
            makeProblem("SRV02", "データベース復旧", {
                submissionStatus: submittable,
                submissionableSchedules: [day1Pm],
            }),
            makeProblem("SRV03", "コンテナ運用管理", {
                submissionStatus: submittable,
                submissionableSchedules: [day1Pm],
                score: { score: 80, markedScore: 100, penalty: -20 },
            }),
            makeProblem("DNS01", "DNS権威サーバー設定", {
                submissionStatus: submittable,
                submissionableSchedules: [day1Pm],
            }),
            makeProblem("DNS02", "DNSキャッシュ問題", {
                submissionStatus: submittable,
                submissionableSchedules: [day1Pm],
            }),
            // day2-am（未来 → 提出不可）
            makeProblem("SEC01", "セキュリティ診断", {
                submissionStatus: notSubmittable,
                submissionableSchedules: [day2Am],
            }),
            makeProblem("SEC02", "ファイアウォール設定", {
                submissionStatus: notSubmittable,
                submissionableSchedules: [day2Am],
            }),
            // day2-pm（未来 → 提出不可）
            makeProblem("APP01", "ロードバランサ冗長化", {
                submissionStatus: notSubmittable,
                submissionableSchedules: [day2Pm],
            }),
            // day1-pm + day2-am（複数スケジュール）
            makeProblem("MON01", "監視システム構築", {
                submissionStatus: submittable,
                submissionableSchedules: [day1Pm, day2Am],
            }),
            // 全スケジュール
            makeProblem("ALL01", "総合演習", {
                submissionStatus: submittable,
                submissionableSchedules: [day1Am, day1Pm, day2Am, day2Pm],
            }),
        ],
    },
};
exports.SingleSchedule = {
    args: {
        notices: [],
        problems: [
            makeProblem("NET01", "ネットワーク基礎問題", {
                submissionStatus: submittable,
                submissionableSchedules: [day1Pm],
            }),
            makeProblem("NET02", "VLAN設定問題", {
                submissionStatus: submittable,
                submissionableSchedules: [day1Pm],
                score: { score: 100, markedScore: 120, penalty: -20 },
            }),
            makeProblem("SRV01", "Webサーバー構築", {
                submissionStatus: submittable,
                submissionableSchedules: [day1Pm],
            }),
        ],
    },
};
exports.NoSchedules = {
    args: {
        notices: [],
        problems: [
            makeProblem("OLD01", "終了済み：OSPF経路制御", {
                submissionStatus: notSubmittable,
                score: { score: 200, markedScore: 200, penalty: 0 },
            }),
            makeProblem("OLD02", "終了済み：BGPピアリング", {
                submissionStatus: notSubmittable,
                score: { score: 150, markedScore: 180, penalty: -30 },
            }),
        ],
    },
};
