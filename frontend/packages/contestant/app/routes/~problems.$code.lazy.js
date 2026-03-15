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
exports.Route = void 0;
var react_1 = require("react");
var react_router_1 = require("@tanstack/react-router");
var v1_1 = require("@ictsc/proto/contestant/v1");
var wkt_1 = require("@bufbuild/protobuf/wkt");
var score_1 = require("../features/score");
var unread_announces_banner_1 = require("./problems.index/unread-announces-banner");
var View = require("./problems.$code/page");
exports.Route = (0, react_router_1.createLazyFileRoute)("/problems/$code")({
    component: RouteComponent,
});
function RouteComponent() {
    var router = (0, react_router_1.useRouter)();
    var _a = exports.Route.useLoaderData(), problem = _a.problem, notices = _a.notices, answers = _a.answers, metadata = _a.metadata, submitAnswer = _a.submitAnswer, deployments = _a.deployments, deploy = _a.deploy, fetchAnswer = _a.fetchAnswer;
    var redeployable = useRedeployable(problem);
    var deferredMetadata = (0, react_1.useDeferredValue)(metadata);
    var deferredAnswers = (0, react_1.useDeferredValue)(answers);
    var deferredDeployments = (0, react_1.useDeferredValue)(deployments);
    return (<View.Page onTabChange={function () {
            (0, react_1.startTransition)(function () { return router.load(); });
        }} redeployable={redeployable} content={<Content problem={problem} notices={notices}/>} submissionForm={<SubmissionForm submitAnswer={submitAnswer} problemPromise={problem} metatataPromise={deferredMetadata}/>} submissionList={<SubmissionList isPending={deferredAnswers !== answers} problemPromise={problem} answersPromise={deferredAnswers} fetchAnswer={fetchAnswer}/>} deploymentList={<Deployments isPending={deferredDeployments !== deployments} deployments={deferredDeployments} problemPromise={problem} deploy={deploy}/>}/>);
}
function useRedeployable(problemPromise) {
    var problem = (0, react_1.use)((0, react_1.useDeferredValue)(problemPromise));
    return problem.redeployable;
}
function SubmissionForm(props) {
    var _this = this;
    var _a;
    var router = (0, react_router_1.useRouter)();
    var problem = (0, react_1.use)(props.problemPromise);
    var metadata = (0, react_1.use)(props.metatataPromise);
    // Convert Timestamp to Date if present
    var submissionStatus = problem.submissionStatus
        ? {
            isSubmittable: problem.submissionStatus.isSubmittable,
            submittableUntil: problem.submissionStatus.submittableUntil
                ? (0, wkt_1.timestampDate)(problem.submissionStatus.submittableUntil)
                : undefined,
        }
        : undefined;
    // submittableUntil/submittableFrom に達したら自動リフェッチして提出状態を更新
    (0, react_1.useEffect)(function () {
        var _a;
        var until = submissionStatus === null || submissionStatus === void 0 ? void 0 : submissionStatus.submittableUntil;
        var from = ((_a = problem.submissionStatus) === null || _a === void 0 ? void 0 : _a.submittableFrom)
            ? (0, wkt_1.timestampDate)(problem.submissionStatus.submittableFrom)
            : undefined;
        var target = until !== null && until !== void 0 ? until : from;
        if (target == null)
            return;
        var ms = target.getTime() - Date.now();
        if (ms <= 0)
            return;
        var timer = setTimeout(function () {
            (0, react_1.startTransition)(function () { return router.invalidate(); });
        }, ms);
        return function () { return clearTimeout(timer); };
    }, [
        submissionStatus === null || submissionStatus === void 0 ? void 0 : submissionStatus.submittableUntil,
        (_a = problem.submissionStatus) === null || _a === void 0 ? void 0 : _a.submittableFrom,
        router,
    ]);
    return (<View.SubmissionForm action={function (body) { return __awaiter(_this, void 0, void 0, function () {
            var e_1;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        _a.trys.push([0, 2, , 3]);
                        return [4 /*yield*/, props.submitAnswer(body)];
                    case 1:
                        _a.sent();
                        return [3 /*break*/, 3];
                    case 2:
                        e_1 = _a.sent();
                        console.error(e_1);
                        return [2 /*return*/, "failure"];
                    case 3: return [4 /*yield*/, router.invalidate()];
                    case 4:
                        _a.sent();
                        return [2 /*return*/, "success"];
                }
            });
        }); }} submitInterval={metadata.submitIntervalSeconds} lastSubmittedAt={metadata.lastSubmittedAt} storageKey={"/problems/".concat(problem.code)} submissionStatus={submissionStatus}/>);
}
function Content(props) {
    var problem = (0, react_1.use)((0, react_1.useDeferredValue)(props.problem));
    var notices = (0, react_1.use)((0, react_1.useDeferredValue)(props.notices));
    return (<>
      <unread_announces_banner_1.UnreadAnnouncesBanner notices={notices}/>
      <View.Content {...problem}/>
    </>);
}
function SubmissionList(props) {
    var _this = this;
    var problem = (0, react_1.use)((0, react_1.useDeferredValue)(props.problemPromise));
    var answers = (0, react_1.use)(props.answersPromise);
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
    var downloadAnswer = function (num) { return __awaiter(_this, void 0, void 0, function () {
        var _a, answerBody, submittedAtString, submittedAt, submittedAtFormattedString, filename;
        return __generator(this, function (_b) {
            switch (_b.label) {
                case 0: return [4 /*yield*/, props.fetchAnswer(num)];
                case 1:
                    _a = _b.sent(), answerBody = _a.answerBody, submittedAtString = _a.submittedAtString;
                    submittedAt = new Date(submittedAtString);
                    submittedAtFormattedString = formatDate(submittedAt);
                    filename = "".concat(problem.code, "-").concat(submittedAtFormattedString, ".md");
                    downloadFile(filename, answerBody);
                    return [2 /*return*/];
            }
        });
    }); };
    return (<View.SubmissionListContainer>
      {answers.length === 0 ? (<View.EmptySubmissionList />) : (<View.SubmissionList isPending={props.isPending}>
          {answers.map(function (answer) { return (<View.SubmissionListItem key={answer.id} id={answer.id} submittedAt={answer.submittedAt} score={(0, score_1.protoScoreToProps)(problem.maxScore, answer === null || answer === void 0 ? void 0 : answer.score)} downloadAnswer={function () {
                    return (0, react_1.startTransition)(function () { return downloadAnswer(answer.id); });
                }}/>); })}
        </View.SubmissionList>)}
    </View.SubmissionListContainer>);
}
function Deployments(props) {
    var _this = this;
    var _a, _b;
    var router = (0, react_router_1.useRouter)();
    var _c = (0, react_1.useOptimistic)((0, react_1.use)(props.deployments)), deployments = _c[0], optimisticSetDeployments = _c[1];
    var canRedeploy = ((_b = (_a = deployments === null || deployments === void 0 ? void 0 : deployments[0]) === null || _a === void 0 ? void 0 : _a.status) !== null && _b !== void 0 ? _b : v1_1.DeploymentStatus.DEPLOYED) ===
        v1_1.DeploymentStatus.DEPLOYED;
    var problem = (0, react_1.use)((0, react_1.useDeferredValue)(props.problemPromise));
    var allowedDeploymentCount = problem.penaltyThreashold - deployments.length;
    var _d = (0, react_1.useActionState)(function (_prev, _action) { return __awaiter(_this, void 0, void 0, function () {
        var timer, e_2;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    timer = setTimeout(function () {
                        optimisticSetDeployments(function (ds) {
                            var _a, _b, _c, _d, _e, _f;
                            return __spreadArray([
                                {
                                    isPending: true,
                                    revision: ds.length + 1,
                                    status: v1_1.DeploymentStatus.DEPLOYING,
                                    requestedAt: new Date().toISOString(),
                                    allowedDeploymentCount: ((_b = (_a = ds === null || ds === void 0 ? void 0 : ds[0]) === null || _a === void 0 ? void 0 : _a.allowedDeploymentCount) !== null && _b !== void 0 ? _b : 1) - 1,
                                    thresholdExceeded: (_d = (_c = ds === null || ds === void 0 ? void 0 : ds[0]) === null || _c === void 0 ? void 0 : _c.thresholdExceeded) !== null && _d !== void 0 ? _d : false,
                                    penalty: (_f = (_e = ds === null || ds === void 0 ? void 0 : ds[0]) === null || _e === void 0 ? void 0 : _e.penalty) !== null && _f !== void 0 ? _f : 0,
                                }
                            ], ds, true);
                        });
                    }, 200);
                    _a.label = 1;
                case 1:
                    _a.trys.push([1, 3, 4, 5]);
                    return [4 /*yield*/, props.deploy()];
                case 2:
                    _a.sent();
                    return [3 /*break*/, 5];
                case 3:
                    e_2 = _a.sent();
                    console.error(e_2);
                    return [2 /*return*/, "再展開に失敗しました"];
                case 4:
                    clearTimeout(timer);
                    return [7 /*endfinally*/];
                case 5: return [4 /*yield*/, router.invalidate()];
                case 6:
                    _a.sent();
                    return [2 /*return*/, null];
            }
        });
    }); }, null), lastResult = _d[0], action = _d[1], isActionPending = _d[2];
    return (<View.Deployments canRedeploy={canRedeploy} isRedeploying={isActionPending} redeploy={function () { return (0, react_1.startTransition)(function () { return action("redeploy"); }); }} allowedDeploymentCount={allowedDeploymentCount} error={lastResult} list={deployments.length === 0 ? (<View.EmptyDeploymentList allowedDeploymentCount={allowedDeploymentCount}/>) : (<View.DeploymentList isPending={props.isPending}>
            {deployments.map(function (deployment) { return (<View.DeploymentItem key={deployment.revision} {...deployment}/>); })}
          </View.DeploymentList>)}/>);
}
