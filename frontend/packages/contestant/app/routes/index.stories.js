"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Ended = exports.Waiting = exports.InContest = void 0;
var index_page_1 = require("./index.page");
exports.default = {
    title: "pages/index",
    component: index_page_1.IndexPage,
};
exports.InContest = {
    name: "競技中",
    args: {
        state: "in_contest",
        currentScheduleName: "day1-am",
        nextScheduleName: "day1-pm",
    },
};
exports.Waiting = {
    name: "競技時間外（待機中）",
    args: {
        state: "waiting",
        nextScheduleName: "day1-am",
    },
};
exports.Ended = {
    name: "競技終了",
    args: {
        state: "ended",
    },
};
