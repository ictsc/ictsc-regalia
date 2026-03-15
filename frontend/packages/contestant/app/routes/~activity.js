"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g = Object.create((typeof Iterator === "function" ? Iterator : Object).prototype);
    return g.next = verb(0), g["throw"] = verb(1), g["return"] = verb(2), typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Route = void 0;
var react_router_1 = require("@tanstack/react-router");
var activity_1 = require("@app/features/activity");
var answer_1 = require("@app/features/answer");
var react_1 = require("react");
var score_1 = require("../features/score");
var activity_page_1 = require("./activity.page");
exports.Route = (0, react_router_1.createFileRoute)("/activity")({
    component: RouteComponent,
    loader: function (_a) {
        var transport = _a.context.transport;
        return {
            activity: (0, activity_1.fetchActivity)(transport),
            fetchAnswer: function (problemCode, num) {
                return (0, answer_1.fetchAnswer)(transport, problemCode, num);
            },
        };
    },
});
function RouteComponent() {
    var _this = this;
    var _a = exports.Route.useLoaderData(), activityPromise = _a.activity, fetchAnswerFn = _a.fetchAnswer;
    var deferredActivityPromise = (0, react_1.useDeferredValue)(activityPromise);
    var activity = (0, react_1.use)(deferredActivityPromise);
    var formatDate = function (date) {
        var pad = function (num) { return num.toString().padStart(2, "0"); };
        var year = date.getFullYear();
        var month = pad(date.getMonth() + 1);
        var day = pad(date.getDate());
        var hour = pad(date.getHours());
        var minute = pad(date.getMinutes());
        var second = pad(date.getSeconds());
        return "".concat(year, "-").concat(month, "-").concat(day, "-").concat(hour, "-").concat(minute, "-").concat(second);
    };
    var downloadFile = function (filename, content) {
        var blob = new Blob([content], { type: "text/markdown" });
        var url = URL.createObjectURL(blob);
        var a = document.createElement("a");
        a.href = url;
        a.download = filename;
        a.click();
        URL.revokeObjectURL(url);
    };
    var downloadAnswer = function (problemCode, num) { return __awaiter(_this, void 0, void 0, function () {
        var _a, answerBody, submittedAtString, submittedAt, submittedAtFormattedString, filename;
        return __generator(this, function (_b) {
            switch (_b.label) {
                case 0: return [4 /*yield*/, fetchAnswerFn(problemCode, num)];
                case 1:
                    _a = _b.sent(), answerBody = _a.answerBody, submittedAtString = _a.submittedAtString;
                    submittedAt = new Date(submittedAtString);
                    submittedAtFormattedString = formatDate(submittedAt);
                    filename = "".concat(problemCode, "-").concat(submittedAtFormattedString, ".md");
                    downloadFile(filename, answerBody);
                    return [2 /*return*/];
            }
        });
    }); };
    return (<activity_page_1.ActivityPage entries={activity.map(function (entry) { return ({
            problemCode: entry.problemCode,
            problemTitle: entry.problemTitle,
            answerId: entry.answerId,
            submittedAt: entry.submittedAt,
            score: (0, score_1.protoScoreToProps)(entry.maxScore, entry.score),
            scored: entry.score != null,
            onDownload: function () {
                return (0, react_1.startTransition)(function () {
                    return downloadAnswer(entry.problemCode, entry.answerId);
                });
            },
        }); })}/>);
}
