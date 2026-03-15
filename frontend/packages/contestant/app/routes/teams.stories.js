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
var teams_page_1 = require("./teams.page");
exports.default = {
    title: "pages/teams",
    component: teams_page_1.TeamsPage,
};
exports.Default = {
    args: {
        teamProfile: __spreadArray(__spreadArray([], Array.from({ length: 1 }, function () { return ({
            $typeName: "contestant.v1.TeamProfile",
            name: "チーム名がわからない",
            organization: "テスト",
            members: [
                {
                    name: "huge",
                    displayName: "ふが",
                    selfIntroduction: "よろしくお願いします！",
                },
                {
                    name: "hoge",
                    displayName: "ほげ",
                    selfIntroduction: "がんばります！",
                },
            ],
        }); }), true), Array.from({ length: 2 }, function () { return ({
            $typeName: "contestant.v1.TeamProfile",
            name: "チーム名",
            organization: "testTeam",
            members: [
                {
                    name: "huge",
                    displayName: "ふが",
                    selfIntroduction: "よろしくお願いします！",
                },
                {
                    name: "hoge",
                    displayName: "ほげ",
                    selfIntroduction: "がんばります！",
                },
            ],
        }); }), true),
    },
};
