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
exports.InContest = void 0;
var announces_index_page_1 = require("./announces.index.page");
exports.default = {
    title: "pages/announces",
    component: announces_index_page_1.AnnounceList,
};
exports.InContest = {
    args: {
        announces: __spreadArray([], Array.from({ length: 10 }, function () { return ({
            $typeName: "contestant.v1.Notice",
            slug: "hogehoge",
            title: "hogehoge",
            body: "hogehoge",
        }); }), true),
    },
};
