"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AppShell = AppShell;
var react_1 = require("react");
var viewer_1 = require("../features/viewer");
var schedule_1 = require("../features/schedule");
var app_shell_1 = require("../components/app-shell");
function AppShell(_a) {
    var children = _a.children, viewerPromise = _a.viewer;
    var viewer = (0, react_1.use)((0, react_1.useDeferredValue)(viewerPromise));
    var schedule = (0, schedule_1.useSchedule)()[0];
    var _b = (0, react_1.useReducer)(function (o) { return !o; }, false), collapsed = _b[0], toggle = _b[1];
    var signOutAction = (0, viewer_1.useSignOut)();
    var inContest = (0, schedule_1.hasContestStarted)(schedule);
    return (<app_shell_1.Layout header={<app_shell_1.Header accountMenu={(viewer === null || viewer === void 0 ? void 0 : viewer.type) === "contestant" && (<app_shell_1.AccountMenu name={viewer.name} onSignOut={function () {
                    (0, react_1.startTransition)(function () { return signOutAction(); });
                }}/>)}/>} navbar={(viewer === null || viewer === void 0 ? void 0 : viewer.type) === "contestant" ? (<app_shell_1.Navbar canViewProblems={inContest} canViewAnnounces={inContest} canViewActivity={true} collapsed={collapsed} onOpenToggleClick={toggle}/>) : null} navbarCollapsed={collapsed}>
      {children}
    </app_shell_1.Layout>);
}
