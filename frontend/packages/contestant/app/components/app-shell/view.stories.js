"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Default = void 0;
var react_1 = require("react");
var actions_1 = require("storybook/actions");
var layout_1 = require("./layout");
var header_1 = require("./header");
var navbar_1 = require("./navbar");
var contest_state_1 = require("./contest-state");
var account_menu_1 = require("./account-menu");
var signOutAction = (0, actions_1.action)("sign-out");
function AppShell() {
    var _a = (0, react_1.useReducer)(function (o) { return !o; }, false), collapsed = _a[0], toggle = _a[1];
    return (<layout_1.Layout header={<header_1.Header contestState={<contest_state_1.ContestStateView state="before" restDurationSeconds={73850}/>} accountMenu={<account_menu_1.AccountMenu name="Alice" onSignOut={signOutAction}/>}/>} navbar={<navbar_1.Navbar canViewProblems canViewAnnounces canViewActivity collapsed={collapsed} onOpenToggleClick={toggle}/>} navbarCollapsed={collapsed}>
      <h1>Main</h1>
      <p>
        あのイーハトーヴォのすきとおった風、夏でも底に冷たさをもつ青いそら、うつくしい森で飾られたモリーオ市、郊外のぎらぎらひかる草の波。
      </p>
    </layout_1.Layout>);
}
exports.default = {
    title: "AppShell",
    component: AppShell,
};
exports.Default = {
    render: function () { return <AppShell />; },
};
