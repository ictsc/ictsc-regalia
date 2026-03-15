"use strict";
var __rest = (this && this.__rest) || function (s, e) {
    var t = {};
    for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p) && e.indexOf(p) < 0)
        t[p] = s[p];
    if (s != null && typeof Object.getOwnPropertySymbols === "function")
        for (var i = 0, p = Object.getOwnPropertySymbols(s); i < p.length; i++) {
            if (e.indexOf(p[i]) < 0 && Object.prototype.propertyIsEnumerable.call(s, p[i]))
                t[p[i]] = s[p[i]];
        }
    return t;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.AccountMenu = AccountMenu;
var clsx_1 = require("clsx");
var react_1 = require("@headlessui/react");
var material_symbol_1 = require("@app/components/material-symbol");
function AccountMenu(props) {
    return (<react_1.Menu>
      <react_1.MenuButton title="アカウントメニュー" className="data-[hover]:bg-surface-0/50 flex size-[50px] items-center justify-center rounded-full transition">
        <material_symbol_1.MaterialSymbol icon="person" fill size={40} className="text-surface-0"/>
      </react_1.MenuButton>

      <react_1.MenuItems anchor={{ to: "bottom", gap: 15 }} transition className={(0, clsx_1.clsx)("bg-surface-0 flex w-[200px] flex-col gap-[5px] rounded-[12px] py-[15px] drop-shadow", "transition duration-200 ease-out data-[closed]:opacity-0")}>
        <span className="text-14 text-text mx-[15px]">{props.name}</span>
        {/* <MenuItem>
          <AccountMenuButton icon="settings">アカウント設定</AccountMenuButton>
        </MenuItem> */}
        <react_1.MenuItem>
          <AccountMenuButton icon="logout" onClick={props.onSignOut}>
            ログアウト
          </AccountMenuButton>
        </react_1.MenuItem>
      </react_1.MenuItems>
    </react_1.Menu>);
}
function AccountMenuButton(_a) {
    var icon = _a.icon, className = _a.className, children = _a.children, restProps = __rest(_a, ["icon", "className", "children"]);
    return (<react_1.Button {...restProps} className={(0, clsx_1.clsx)(className, "data-[focus]:bg-surface-1 flex items-center px-[15px] py-[10px] transition")}>
      <material_symbol_1.MaterialSymbol icon={icon} size={20}/>
      <span className="text-14 ml-[5px]">{children}</span>
    </react_1.Button>);
}
