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
Object.defineProperty(exports, "__esModule", { value: true });
exports.Navbar = Navbar;
var react_1 = require("react");
var clsx_1 = require("clsx");
var react_2 = require("@headlessui/react");
var react_router_1 = require("@tanstack/react-router");
var material_symbol_1 = require("../../components/material-symbol");
function Navbar(props) {
    var collapsed = props.collapsed;
    var state = (0, react_router_1.useRouterState)();
    return (<div className="bg-surface-1 text-text flex size-full flex-col items-start gap-4">
      <react_2.Button as={react_1.Fragment}>
        {function (buttonProps) { return (<button title={collapsed ? "開く" : "閉じる"} onClick={props.onOpenToggleClick} className={navbarButtonClassName(__assign({ collapsed: true }, buttonProps))}>
            <NavbarButtonInner collapsed icon="list"/>
          </button>); }}
      </react_2.Button>
      <react_2.Button as={react_1.Fragment}>
        {function (buttonProps) {
            var _a;
            return (<react_router_1.Link to="/rule" title="ルール" className={navbarButtonClassName(__assign({ collapsed: collapsed, matched: (_a = state.location.pathname) === null || _a === void 0 ? void 0 : _a.startsWith("/rule") }, buttonProps))}>
            <NavbarButtonInner collapsed={collapsed} icon="developer_guide" title="ルール"/>
          </react_router_1.Link>);
        }}
      </react_2.Button>
      <react_2.Button as={react_1.Fragment} disabled={!props.canViewAnnounces}>
        {function (buttonProps) {
            var _a;
            return (<react_router_1.Link to="/announces" title={props.canViewAnnounces ? "アナウンス" : "開催期間外です"} className={navbarButtonClassName(__assign({ collapsed: collapsed, matched: (_a = state.location.pathname) === null || _a === void 0 ? void 0 : _a.startsWith("/announces") }, buttonProps))}>
            <NavbarButtonInner collapsed={collapsed} icon="brand_awareness" title="アナウンス"/>
          </react_router_1.Link>);
        }}
      </react_2.Button>
      {/* <NavbarButton showTitle={!collapsed} icon="lan" title="接続情報" /> */}
      <react_2.Button as={react_1.Fragment} disabled={!props.canViewProblems}>
        {function (buttonProps) {
            var _a;
            return (<react_router_1.Link to="/problems" title={props.canViewProblems ? "問題" : "開催期間外です"} className={navbarButtonClassName(__assign({ collapsed: collapsed, matched: (_a = state.location.pathname) === null || _a === void 0 ? void 0 : _a.startsWith("/problems") }, buttonProps))}>
            <NavbarButtonInner collapsed={collapsed} icon="help" title="問題"/>
          </react_router_1.Link>);
        }}
      </react_2.Button>
      <react_2.Button as={react_1.Fragment}>
        {function (buttonProps) {
            var _a;
            return (<react_router_1.Link to="/teams" title="チーム一覧" className={navbarButtonClassName(__assign({ collapsed: collapsed, matched: (_a = state.location.pathname) === null || _a === void 0 ? void 0 : _a.startsWith("/teams") }, buttonProps))}>
            <NavbarButtonInner collapsed={collapsed} icon="groups" title="チーム一覧"/>
          </react_router_1.Link>);
        }}
      </react_2.Button>
      <react_2.Button as={react_1.Fragment}>
        {function (buttonProps) {
            var _a;
            return (<react_router_1.Link to="/ranking" title="ランキング" className={navbarButtonClassName(__assign({ collapsed: collapsed, matched: (_a = state.location.pathname) === null || _a === void 0 ? void 0 : _a.startsWith("/ranking") }, buttonProps))}>
            <NavbarButtonInner collapsed={collapsed} icon="trophy" title="ランキング"/>
          </react_router_1.Link>);
        }}
      </react_2.Button>
      <react_2.Button as={react_1.Fragment} disabled={!props.canViewActivity}>
        {function (buttonProps) {
            var _a;
            return (<react_router_1.Link to="/activity" title={props.canViewActivity ? "アクティビティ" : "開催期間外です"} className={navbarButtonClassName(__assign({ collapsed: collapsed, matched: (_a = state.location.pathname) === null || _a === void 0 ? void 0 : _a.startsWith("/activity") }, buttonProps))}>
            <NavbarButtonInner collapsed={collapsed} icon="history" title="アクティビティ"/>
          </react_router_1.Link>);
        }}
      </react_2.Button>
      {/* <NavbarButtonInner
          collapsed={collapsed}
          icon="groups"
          title="チーム一覧"
        /> */}
      {/* <NavbarButton showTitle={!collapsed} icon="chat" title="お問い合わせ" /> */}
    </div>);
}
function navbarButtonClassName(_a) {
    var collapsed = _a.collapsed, hover = _a.hover, active = _a.active, _b = _a.disabled, disabled = _b === void 0 ? false : _b, _c = _a.matched, matched = _c === void 0 ? false : _c;
    return (0, clsx_1.clsx)("bg-surface-1 flex flex-row items-center rounded-[10px] transition", !collapsed && "w-full", disabled ? "cursor-not-allowed opacity-75" : "text-text", (hover || matched) && "bg-surface-2", active && "opacity-75");
}
function NavbarButtonInner(props) {
    return (<>
      <div className="flex size-[50px] shrink-0 items-center justify-center">
        <material_symbol_1.MaterialSymbol icon={props.icon} size={24}/>
      </div>
      {!props.collapsed && (<span className="text-16 line-clamp-1 overflow-x-hidden text-left">
          {props.title}
        </span>)}
    </>);
}
