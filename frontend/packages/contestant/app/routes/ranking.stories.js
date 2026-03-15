"use strict";
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
exports.Default = void 0;
var ranking_page_1 = require("./ranking.page");
exports.default = {
    title: "pages/ranking",
    component: ranking_page_1.RankingPage,
};
exports.Default = {
    args: {
        ranking: __spreadArray(__spreadArray([], Array.from({ length: 2 }, function () { return ({
            rank: 1,
            teamName: "チーム名なまえがわからない",
            organization: "testTeam",
            score: 8888,
            lastEffectiveSubmitAt: "2025-03-04T12:00:00Z",
        }); }), true), Array.from({ length: 4 }, function () { return ({
            rank: 2,
            teamName: "チーム名なまえがわからない",
            organization: "testTeam",
            score: 8888,
            lastEffectiveSubmitAt: "2025-03-04T12:00:00Z",
        }); }), true),
    },
};
